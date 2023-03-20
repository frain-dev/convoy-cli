mockgen:
	go generate ./...

build:
	scripts/build.sh

integration_tests:
	go test -tags integration -p 1 ./...

install:
	go build -o convoy-cli ./cmd
	mv ./convoy-cli ${GOPATH}/bin

generate_migration_time:
	@date +"%Y%m%d%H%M%S"
