CURDIR := $(shell pwd)
BINDIR ?= ${CURDIR}/.bin

# Tooling
.bindir:
	test -d ${BINDIR} || mkdir -p ${BINDIR}

.install-antlr: .bindir
	test -f ${BINDIR}/antlr.jar || curl "https://www.antlr.org/download/antlr-4.11.1-complete.jar" -o ${BINDIR}/antlr.jar

tools: .install-antlr

# Proto generation
PROTOC = protoc --go_opt=module=github.com/renbou/loggo --go_out=. $$PROTO_ARGS ${PROTO_PATH}
PROTOC_GRPC = PROTO_ARGS="--go-grpc_opt=module=github.com/renbou/loggo --go-grpc_out=."; ${PROTOC}

.generate-storage-proto: PROTO_PATH=internal/storage/proto/models.proto
.generate-storage-proto:
	$(PROTOC)

.generate-api-proto: PROTO_PATH=api/pigeoneer.proto api/telemetry.proto
.generate-api-proto:
	$(PROTOC_GRPC)

.generate-web-proto:
	cd front && npm run generate

.generate-filter-parser:
	java -Xmx500M -cp "${BINDIR}/antlr.jar:$$CLASSPATH" org.antlr.v4.Tool \
	  -Dlanguage=JavaScript -no-listener \
	  front/src/lib/filters/Filters.g4; \
	rm -f front/src/lib/filters/*.interp front/src/lib/filters/*.tokens
	
generate: .generate-storage-proto .generate-api-proto .generate-web-proto .generate-filter-parser

# Testing
COVER_FILE := $(shell mktemp)
COVER_TMP_FILE := $(shell mktemp)
COVER_HTML_REPORT ?= .cover.html
COVER_FUNC_REPORT ?= .cover.func

test:
	go test -count=1 -race -v ./...

cover:
	go test -count=1 -race -covermode=atomic -coverpkg=./... -coverprofile ${COVER_TMP_FILE} -v ./...
	cat ${COVER_TMP_FILE}   | \
		grep -v ".*\.pb\.go" | \
		grep -v ".*_mock\.go" | \
		cat - > ${COVER_FILE}
	if [ -n "${COVER_HTML}" ]; then \
		go tool cover -html ${COVER_FILE} -o ${COVER_HTML_REPORT}; \
	else \
		go tool cover -func ${COVER_FILE} -o ${COVER_FUNC_REPORT}; \
	fi
	rm -f ${COVER_TMP_FILE} ${COVER_FILE}

# Building
.build-front:
	cd front && npm run build

build: .bindir .build-front
	\cp -r front/dist internal/web && \
	CGO_ENABLED=0 ${BUILD_FLAGS} go build -o ${BINDIR}/loggo ./cmd/loggo
