mockgen:
	go generate ./...

build:
	scripts/build.sh

integration_tests:
	go test -tags integration -p 1 ./...

generate_migration_time:
	@date +"%Y%m%d%H%M%S"
