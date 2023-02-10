CURDIR := $(shell pwd)
BINDIR ?= ${CURDIR}/.bin

PROTOC_GEN_GO_V = v1.28.1
PROTOC_GEN_GO_GRPC_V = v1.2.0

# Tooling and dependencies
INSTALL_TOOL = test -f ${BINDIR}/`echo ${TOOLURL} | perl -ne 'if (/[\w.\/-]+\/([\w-]+)@latest/) { print $$1; exit }'` || \
	GOBIN=${BINDIR} go install ${TOOLURL}

.bindir:
	test -d ${BINDIR} || mkdir -p ${BINDIR}

.install-protoc-gen-go: TOOLURL=google.golang.org/protobuf/cmd/protoc-gen-go@${PROTOC_GEN_GO_V}
.install-protoc-gen-go:
	$(INSTALL_TOOL)

.install-protoc-gen-go-grpc: TOOLURL=google.golang.org/grpc/cmd/protoc-gen-go-grpc@${PROTOC_GEN_GO_GRPC_V}
.install-protoc-gen-go-grpc:
	$(INSTALL_TOOL)

deps: .install-protoc-gen-go .install-protoc-gen-go-grpc

# Proto generation
PROTOC = protoc --plugin=${BINDIR}/protoc-gen-go \
	--go_opt=module=github.com/renbou/loggo --go_out=. $$PROTO_ARGS ${PROTO_PATH}
PROTOC_GRPC = PROTO_ARGS="\
	--plugin=${BINDIR}/protoc-gen-go-grpc \
	--go-grpc_opt=module=github.com/renbou/loggo --go-grpc_out=."; ${PROTOC}

.generate-storage-proto: PROTO_PATH=internal/storage/proto/models.proto
.generate-storage-proto:
	$(PROTOC)

.generate-api-proto: PROTO_PATH=api/pigeoneer.proto
.generate-api-proto:
	$(PROTOC_GRPC)

generate: .generate-storage-proto .generate-api-proto

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
