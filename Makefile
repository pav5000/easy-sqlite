BINDIR=$(CURDIR)/bin
GOVER=$(shell go version | perl -nle '/(go\d\S+)/; print $$1;')
SMARTIMPORTS=${BINDIR}/smartimports_${GOVER}
LINTVER=v1.59.1
LINTBIN=${BINDIR}/lint_${GOVER}_${LINTVER}

precommit: format lint

run-example:
	go run github.com/pav5000/easy-sqlite/cmd/example

bindir:
	mkdir -p ${BINDIR}

lint: install-lint
	${LINTBIN} run

format: install-smartimports
	${SMARTIMPORTS}

install-lint: bindir
	test -f ${LINTBIN} || \
		(GOBIN=${BINDIR} go install github.com/golangci/golangci-lint/cmd/golangci-lint@${LINTVER} && \
		mv ${BINDIR}/golangci-lint ${LINTBIN})

install-smartimports: bindir
	test -f ${SMARTIMPORTS} || \
		(GOBIN=${BINDIR} go install github.com/pav5000/smartimports/cmd/smartimports@latest && \
		mv ${BINDIR}/smartimports ${SMARTIMPORTS})
