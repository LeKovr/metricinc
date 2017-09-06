/*
Client for gRPC counter service

Usage

Run from command line:
 client [OPTIONS] <get | inc | set>

Application Options:
  --connect= Addr and port which server listens (default: :50051)

Help Options:
  -h, --help     Show help message

Available commands:
  get  Get current values
  inc  Increment counter
  set  Set both counter settings (step and limit) using defaults if not given (see 'set -h')

[set command options]
  --step=  Increment step (default: 1)
  --limit= Increment loop limit (default: 100)

Program returns all counter data as json.

Examples

Values from clean new service:
 $ ./client get
 {"number":0,"step":1,"limit":100}

Increment counter:
 $ ./client inc
 {"number":1,"step":1,"limit":100}

Set new step:
 $ ./client set --step 2
 {"number":1,"step":2,"limit":100}

Increment by new step
 $ ./client inc
 {"number":3,"step":2,"limit":100}


*/
package main
