CURDIR := $(shell pwd)
BINDIR ?= ${CURDIR}/.bin

# Tooling
INSTALL_TOOL = test -f ${BINDIR}/`echo ${TOOLURL} | perl -ne 'if (/[\w.\/-]+\/([\w-]+)@latest/) { print $$1; exit }'` || \
	GOBIN=${BINDIR} go install ${TOOLURL}

.bindir:
	test -d ${BINDIR} || mkdir -p ${BINDIR}

# Proto generation
PROTOC=protoc --go_opt=module=github.com/renbou/obzerva --go_out=. ${PROTO_PATH}

.generate-storage-proto: PROTO_PATH=internal/logs/storage/proto/models.proto
.generate-storage-proto:
	$(PROTOC)

generate: .generate-storage-proto

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
		cat - > ${COVER_FILE}
	if [ -n "${COVER_HTML}" ]; then \
		go tool cover -html ${COVER_FILE} -o ${COVER_HTML_REPORT}; \
	else \
		go tool cover -func ${COVER_FILE} -o ${COVER_FUNC_REPORT}; \
	fi
	rm -f ${COVER_TMP_FILE} ${COVER_FILE}
