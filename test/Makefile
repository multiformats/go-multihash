all: deps

deps:
	cd sharness && make deps

clean:
	cd sharness && make clean

test: test_expensive

test_expensive:
	cd sharness && make TEST_EXPENSIVE=1

test_cheap:
	cd sharness && make

.PHONY: all clean
