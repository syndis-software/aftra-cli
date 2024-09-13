mkfile_dir:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

OPENAPI_HOST ?= https://app.aftra.io

init:
	echo "Make sure you have go/bin in your PATH"
	go get .
	go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.15.0
	go generate ./...

build:
	go build ./...

download-openapi-spec:
	echo "OPENAPI_HOST is ${OPENAPI_HOST}"
	curl ${OPENAPI_HOST}/api/openapi.json > $(mkfile_dir)/scripts/openapi.json

update-openapi-subset-spec:
	python $(mkfile_dir)/scripts/subset_maker.py $(mkfile_dir)/scripts/openapi.json > $(mkfile_dir)/scripts/openapi-subset.json
	rm $(mkfile_dir)/scripts/openapi.json

generate:
	oapi-codegen --package=openapi -generate=types,client -o $(mkfile_dir)/pkg/openapi/openapi.gen.go $(mkfile_dir)/scripts/openapi-subset.json

upgrade: download-openapi-spec update-openapi-subset-spec generate

release:
	git tag $(RELEASE)
	git push --tags