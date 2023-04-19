#/usr/bin/env sh

docker run --rm -v `pwd`/authz:/authz openpolicyagent/opa:latest test --coverage /authz > opa-coverage.json
go run ../main.go opa-coverage.json > codecov-coverage.json
