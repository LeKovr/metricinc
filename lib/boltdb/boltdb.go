/*
Package boltdb is an interface to embedded key/value database boltdb.

Этот пакет не используется напрямую в других пакетах, только в main.
Остальные пакеты (lib/grpcapi) работают через интерфейс lib/iface/kvstore.

*/
package boltdb

import (
	"encoding/binary"
	"errors"

	"github.com/boltdb/bolt"
	"github.com/gogo/protobuf/proto"
	"github.com/LeKovr/metricinc/counter/setup"
	"github.com/LeKovr/metricinc/lib/iface/logger"
	pb "github.com/LeKovr/metricinc/lib/proto/counter"
)

// Config is a program flags group used in constructor
type Config struct {
	File        string `long:"db_file" default:"base.db" description:"Bolt database file"`
	Bucket      string `long:"db_bucket" default:"counter" description:"Bucket name"`
	NumberKey   string `long:"db_number_key" default:"number" description:"Key name for current number"`
	SettingsKey string `long:"db_settings_key" default:"config" description:"Key name for settings data"`
}

// Store holds config fields as []byte and refs to logger and boltdb
type Store struct {
	Bucket      []byte
	NumberKey   []byte
	SettingsKey []byte
	db          *bolt.DB
	log         logger.Entry
}

// NewStore creates a boltdb object.
func NewStore(log logger.Entry, cfg Config) (*Store, error) {

	// Open bolt database
	// TODO: move 0644 to Config
	db, err := bolt.Open(cfg.File, 0644, nil)
	if err != nil {
		return nil, err
	}

	log.WithField("config", cfg).Debug("Create store")
	s := Store{
		Bucket:      []byte(cfg.Bucket),
		NumberKey:   []byte(cfg.NumberKey),
		SettingsKey: []byte(cfg.SettingsKey),
		db:          db,
		log:         log,
	}
	return &s, nil

}

// GetSettings reads settings from database.
// Bucket: s.Bucket, Field key: s.SettingsKey
func (s *Store) GetSettings() (*setup.Settings, error) {

	var sets *setup.Settings

	// Read-only transaction
	err := s.db.View(func(tx *bolt.Tx) error {
		var b *bolt.Bucket
		if b = tx.Bucket(s.Bucket); b == nil { // no such bucket
			s.log.Debugf("Bucket does not exists")
			return nil
		}

		var v []byte
		if v = b.Get(s.SettingsKey); v == nil {
			s.log.Debugf("Settings data does not exists")
			return nil
		}

		loaded := setup.Settings{}
		if err := UnmarshalSettings(v, &loaded); err != nil {
			s.log.Debugf("Error unmarshalling settings data: %+v", err)
			return err
		}

		// No error, load sets
		sets = &loaded
		return nil
	})
	return sets, err
}

// SetSettings saves settings to database.
// Bucket: s.Bucket, Field key: s.SettingsKey
func (s *Store) SetSettings(sets *setup.Settings) error {

	// Write transaction
	err := s.db.Update(func(tx *bolt.Tx) error {
		s.log.WithField("bucket", string(s.Bucket)).WithField("settings", *sets).Debug("Open bucket for settings")

		bucket, err := tx.CreateBucketIfNotExists(s.Bucket)
		if err != nil {
			return err
		}

		if v, err := MarshalSettings(sets); err != nil {
			return err
		} else if err := bucket.Put(s.SettingsKey, v); err != nil {
			return err
		}
		return nil
	})
	return err
}

// GetNumber reads number from database.
// Bucket: s.Bucket, Field key: s.NumberKey
func (s *Store) GetNumber() (*int64, error) {

	var num int64

	// Read-only transaction
	err := s.db.View(func(tx *bolt.Tx) error {
		var b *bolt.Bucket
		if b = tx.Bucket(s.Bucket); b == nil { // no such bucket
			return nil
		}

		if v := b.Get(s.NumberKey); v != nil {
			var n int64
			var bytes int
			if n, bytes = binary.Varint(v); bytes <= 0 {
				if bytes == 0 {
					return errors.New("read buffer too small") // is it possible?
				}
				// -bytes = total bytes read
				return errors.New(string(s.NumberKey) + " is not int64") // db is broken
			}
			num = n
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &num, nil
}

// SetNumber saves number to database.
// Bucket: s.Bucket, Field key: s.NumberKey
func (s *Store) SetNumber(number *int64) error {

	// Write transaction
	err := s.db.Update(func(tx *bolt.Tx) error {
		s.log.WithField("bucket", string(s.Bucket)).WithField("number", *number).Debug("Open bucket for number")
		bucket, err := tx.CreateBucketIfNotExists(s.Bucket)
		if err != nil {
			return err
		}
		buf := make([]byte, binary.MaxVarintLen64)
		_ = binary.PutVarint(buf, *number) // returns the number of bytes written
		if err = bucket.Put(s.NumberKey, buf); err != nil {
			return err
		}
		return nil
	})

	return err
}

// Close connection to kv database.
func (s *Store) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

// MarshalSettings encodes a settings to binary format.
func MarshalSettings(sr *setup.Settings) ([]byte, error) {
	return proto.Marshal(&pb.Settings{
		Step:  *proto.Int64(int64(sr.Step)),
		Limit: *proto.Int64(int64(sr.Limit)),
	})
}

// UnmarshalSettings decodes a settings from binary format.
func UnmarshalSettings(data []byte, d *setup.Settings) error {
	var buf pb.Settings
	if err := proto.Unmarshal(data, &buf); err != nil {
		return err
	}

	d.Step = buf.GetStep()
	d.Limit = buf.GetLimit()

	return nil
}
