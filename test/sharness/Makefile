# Run tests
#
# Copyright (c) 2014 Christian Couder
# MIT Licensed; see the LICENSE file in this repository.
#

# NOTE: run with TEST_VERBOSE=1 for verbose sharness tests.

T = $(sort $(wildcard t[0-9][0-9][0-9][0-9]-*.sh))
BINS = bin/multihash
SHARNESS = lib/sharness/sharness.sh
MULTIHASH_CMD = ../../multihash
MULTIHASH_ROOT = ../../

all: clean deps $(T) aggregate

clean:
	@echo "*** $@ ***"
	-rm -rf test-results
	-rm -rf $(BINS)

$(T):
	@echo "*** $@ ***"
	./$@

aggregate:
	@echo "*** $@ ***"
	lib/test-aggregate-results.sh

deps: $(SHARNESS) $(BINS)

$(SHARNESS):
	@echo "*** installing $@ ***"
	lib/install-sharness.sh

bin/multihash: $(MULTIHASH_ROOT)/**/*.go
	go build -o bin/multihash $(MULTIHASH_CMD)

bins: $(BINS)

.PHONY: all clean $(T) aggregate
