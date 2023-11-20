
help: # show all commands
	@egrep -h '\s#\s' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?# "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

test: # run tests with output to xml file
	mkdir -p tmp
	go mod download
	go run gotest.tools/gotestsum@latest --junitfile tmp/report.xml --format testname -- --coverprofile=tmp/cover.out --covermode count ./...
	go run github.com/boumenot/gocover-cobertura@latest < tmp/cover.out > tmp/coverage.xml

build: # build executables
	go build -o tmp/meet cmd/meet.go
