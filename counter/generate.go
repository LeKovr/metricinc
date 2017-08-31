// protoc --gofast_out=. counter.proto

//go : generate protoc --go_out=plugins=grpc:. counter.proto
package count
