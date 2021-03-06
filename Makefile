VERSION:=$(shell git describe --tags)
COMMIT:=$(shell git rev-parse HEAD)
BUILT:=$(shell date +%FT%T%z)
BASE_PKG:=github.com/setheck/smartthings-exporter
IMAGE:=setheck/smartthings-exporter

LDFLAGS=-ldflags "-w -s -X main.Version=${VERSION} \
				        -X main.Built=${BUILT} \
				        -X main.Commit=${COMMIT}"

test:
	go test ./... -cover

build:
	go build ${LDFLAGS} .

dbuild:
	# *Note, docker file calls `make build`
	docker build . -t ${IMAGE}:latest
	docker run --rm ${IMAGE}:latest -version

dbuild-arm:
	docker version
	docker build --platform linux/arm/v7 --tag ${IMAGE}:latest .

dpush:
	docker push ${IMAGE}:latest

drun: dbuild
	docker run --rm -p 9119:9119 ${IMAGE}:latest

tag: MAJOR=0
tag: MINOR=0
tag: PATCH=4
tag:
	git tag "${MAJOR}.${MINOR}.${PATCH}"
	git push origin --tags

deploy: clean dbuild
	docker tag ${IMAGE}:latest ${IMAGE}:${VERSION}
	docker push ${IMAGE}:latest
	docker push ${IMAGE}:${VERSION}

clean:
	rm -rf smartthings-exporter

.PHONY: test build dbuild clean tag