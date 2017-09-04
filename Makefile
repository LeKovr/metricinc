SHELL      = /bin/bash

pb:
	go generate ./lib/proto/...

doc:
	@echo "Open http://localhost:6060/pkg/lekovr/exam"
	godoc -http=:6060

cov-status:
	for f in counter lib/* ; do [[ $$f == lib/iface ]] || [[ $$f == lib/proto ]] || \
  { pushd $$f ; go test -coverprofile=coverage.out ; popd ; } ; done

# go get github.com/golang/mock/mockgen
mocks: lib/mock_kvstore/mock_kvstore.go

lib/mock_kvstore/mock_kvstore.go:
	mockgen -source=lib/iface/kvstore/kvstore.go > $@

#go tool cover -html=coverage.out
#	go test ./lib/... -coverprofile=coverage.out