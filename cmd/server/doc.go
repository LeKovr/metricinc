// Copyright 2017 Alexey Kovrizhkin <lekovr@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

/*
gRPC counter service

This service offers counter methods, described in lib/grpcapi package

Usage

Run from command line:

 server [OPTIONS]

Application Options:
 --listen=          Addr and port which server listens (default: :50051)

Logging Options:
 --log_level=       Log level [warn|info|debug] (default: debug)
 --log_stdout       Log to STDOUT without color and timestamps

API Options:
 --init_number=     Initial number (default: 0)
 --init_step=       Increment step (default: 1)
 --init_limit=      Increment loop limit (default: 100)
 --store_strict     Do not ignore store errors

Storage Options:
 --db_file=         Bolt database file (default: base.db)
 --db_bucket=       Bucket name (default: counter)
 --db_number_key=   Key name for current number (default: number)
 --db_settings_key= Key name for settings data (default: config)

Help Options:
 -h, --help         Show help message

Examples

Run with defaults:

 $ ./server
 DEBU[0000] Create logger                                 config="{debug false}"
 INFO[0000] Counter server v0.1
 DEBU[0000] Create listener                               addr=":50051"
 DEBU[0000] Create store                                  config="{base.db counter number config}"
 DEBU[0000] Create API                                    config="{0 1 100 false}"
 DEBU[0000] Load Settings from db                         settings="{3 20}"
 DEBU[0000] Load Number from db                           number=3
 INFO[0000] API created

 ^C
 INFO[0001] Got signal: interrupt
 WARN[0001] Final state                                   number=3 settings="{3 20}"

Run without debug logs:

 $ ./server --log_level info
 INFO[0000] Counter server v0.1
 INFO[0000] API created

 ^C
 INFO[0003] Got signal: interrupt
 WARN[0003] Final state                                   number=3 settings="{3 20}"

*/
package main
