# Makefile for building CoreDNS
GITCOMMIT:=$(shell git describe --dirty --always)
BINARY:=coredns
SYSTEM:=
CHECKS:=check
BUILDOPTS:=-v
GOPATH?=$(HOME)/go
MAKEPWD:=$(dir $(realpath $(firstword $(MAKEFILE_LIST))))
CGO_ENABLED?=0

.PHONY: all
all: coredns

.PHONY: coredns
coredns: $(CHECKS)
	CGO_ENABLED=$(CGO_ENABLED) $(SYSTEM) go build $(BUILDOPTS) -ldflags="-s -w -X github.com/coredns/coredns/coremain.GitCommit=$(GITCOMMIT)" -o $(BINARY)

.PHONY: check
check: core/plugin/zplugin.go core/dnsserver/zdirectives.go

core/plugin/zplugin.go core/dnsserver/zdirectives.go: plugin.cfg
	go generate coredns.go
	go get

.PHONY: gen
gen:
	go generate coredns.go
	go get

.PHONY: pb
pb:
	$(MAKE) -C pb

.PHONY: clean
clean:
	go clean
	rm -f coredns

.PHONY: signed
signed:
	rm -f ./plugin/bloomsec_nsec5/testdata/db.miek.nl.signed
	rm -f ./plugin/bloomsec_nsec5/testdata2/db.miek.nl.signed
	rm -f ./plugin/bloomsec_nsec5/testdata3/db.miek.nl.signed
	rm -f ./plugin/bloomsec_nsec5/testdata4/db.miek.nl.signed
	rm -f ./plugin/bloomsec/testdata/db.miek.nl.signed
	rm -f ./plugin/bloomsec/testdata2/db.miek.nl.signed
	rm -f ./plugin/bloomsec/testdata3/db.miek.nl.signed
	rm -f ./plugin/sign_nsec5/testdata/db.miek.nl.signed
	rm -f ./plugin/sign_nsec5/testdata2/db.miek.nl.signed
	rm -f ./plugin/sign_nsec5/testdata3/db.miek.nl.signed
	rm -f ./plugin/sign_nsec5/testdata4/db.miek.nl.signed
	rm -f ./plugin/sign/testdata/db.miek.nl.signed

.PHONY: testdata_ns
testdata_ns:
	./coredns -conf ./memusage/Corefiles_ns/$(COREFILE) -dns.port 1054

.PHONY: testdata_za
testdata_za:
	rm -f ./memusage/testdata/db.miek.nl.signed
	./coredns -conf ./memusage/Corefiles_za/$(COREFILE) -dns.port 1054

.PHONY: run_eval
run_eval: 
	rm -f ./evaluation/testdata/zonedata/db.miek.nl.signed
	./coredns -conf ./evaluation/Corefiles/Corefile_$(version) -dns.port 1054