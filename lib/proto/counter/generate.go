// This file used by go generate

//go:generate protoc --gogo_out=plugins=grpc:. counter.proto

package counter
