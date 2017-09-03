package boltdb

import (
	"encoding/binary"
	"errors"

	"github.com/boltdb/bolt"
	"github.com/gogo/protobuf/proto"
	pb "lekovr/exam/lib/proto/counter"
	"lekovr/exam/lib/struct/logger"
	"lekovr/exam/counter/setup"
)

type Config struct {
	File        string `long:"db_file" default:"base.db" description:"Bolt database file"`
	NumberKey   string `long:"db_number_key" default:"number" description:"Show verbose debug information"`
	SettingsKey string `long:"db_settings_key" default:"config" description:"Show verbose debug information"`
	Bucket      string `long:"db_bucket" default:"counter" description:"Show verbose debug information"`
}

type Store struct {
	Bucket      []byte
	NumberKey   []byte
	SettingsKey []byte
	db          *bolt.DB
	log         logger.Entry
}

func NewStore(log logger.Entry, cfg Config) (*Store, error) {

	db, err := bolt.Open(cfg.File, 0644, nil)
	if err != nil {
		return nil, err
	}

	log.Debugf("Got config: %+v", cfg)
	s := Store{
		Bucket:      []byte(cfg.Bucket),
		NumberKey:   []byte(cfg.NumberKey),
		SettingsKey: []byte(cfg.SettingsKey),
		db:          db,
		log:         log,
	}
	return &s, nil

}

func (s *Store) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

// MarshalDial encodes a dial to binary format.
func MarshalSettings(sr *setup.Settings) ([]byte, error) {
	return proto.Marshal(&pb.SettingsRequest{
		Step:  *proto.Int64(int64(sr.Step)),
		Limit: *proto.Int64(int64(sr.Limit)),
	})
}

func UnmarshalSettings(data []byte, d *setup.Settings) error {
	var buf pb.SettingsRequest
	if err := proto.Unmarshal(data, &buf); err != nil {
		return err
	}

	d.Step = buf.GetStep()
	d.Limit = buf.GetLimit()

	return nil
}

func (s *Store) GetSettings() (*setup.Settings, error) {

	var sets *setup.Settings
	err := s.db.View(func(tx *bolt.Tx) error {

		if b := tx.Bucket(s.Bucket); b == nil { // no such bucket
			s.log.Debugf("Bucket does not exists")
			return nil
		} else {
			if v := b.Get(s.SettingsKey); v == nil {
				s.log.Debugf("Settings data does not exists")
				return nil
			} else {
				// не закрываем if, чтобы иметь доступ к v
				loaded := setup.Settings{}
				if err := UnmarshalSettings(v, &loaded); err != nil {
					s.log.Debugf("Error unmarshalling settings data: %+v", err)
					return err
				}
				sets = &loaded
			}
			// No error. sets loaded
			return nil
		}
	})
	return sets, err
}

func (s *Store) SetSettings(sets *setup.Settings) error {

	err := s.db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(s.Bucket))
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

func (s *Store) GetNumber() (*int64, error) {

	var num int64
	err := s.db.View(func(tx *bolt.Tx) error {
		if b := tx.Bucket(s.Bucket); b == nil { // no such bucket
			return nil
		} else {

			if v := b.Get(s.NumberKey); v != nil {
				var nn int64
				var bytes int
				if nn, bytes = binary.Varint(v); bytes <= 0 {
					// значение длиннее 8и байт или буфер короче 8. Наш случай - первый
					if bytes == 0 {
						return errors.New("read buffer too small") //
					}
					// -bytes - сколько байт прочитано
					return errors.New(string(s.NumberKey) + " is not int64") //
				} else {
					num = nn
				}
			}
		}
		return nil
	})
	if err == nil {
		return &num, nil
	}
	return nil, err
}

func (s *Store) SetNumber(number *int64) error {
	s.log.Debugf("SetNumber: %d", *number)

	err := s.db.Update(func(tx *bolt.Tx) error {
		s.log.Debugf("Open bucket %s for %d", string(s.Bucket), *number)
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
