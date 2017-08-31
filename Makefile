SHELL      = /bin/bash

pb:
	go generate ./counter


#	protoc  -I counter counter.proto  --go_out=plugins=grpc:counter

#	protoc -I /home/src/go/src -I counter counter.proto  --go_out=plugins=grpc:counter


#protoc -I routeguide/ routeguide/route_guide.proto --go_out=plugins=grpc:routeguide