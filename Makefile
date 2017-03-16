gx:
	go get github.com/whyrusleeping/gx
	go get github.com/whyrusleeping/gx-go

covertools:
	go get github.com/mattn/goveralls
	go get golang.org/x/tools/cmd/cover

deps: gx covertools
	gx --verbose install --global
	gx-go rewrite

sharness: deps
	git clone https://github.com/multiformats/multihash.git $(PWD)/sharness
	cd $(PWD)/multihash && go build -v . && ls && chmod +x ./multihash
	export MULTIHASH_BIN="$(PWD)/multihash/multihash" && export TEST_EXPENSE=1 && make -j1 -C $(PWD)/sharness/tests/sharness

publish:
	gx-go rewrite --undo

