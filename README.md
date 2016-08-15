# go-multihash

[![](https://img.shields.io/badge/made%20by-Protocol%20Labs-blue.svg?style=flat-square)](http://ipn.io)
[![](https://img.shields.io/badge/project-multiformats-blue.svg?style=flat-square)](http://github.com/multiformats/multiformats)
[![](https://img.shields.io/badge/freenode-%23ipfs-blue.svg?style=flat-square)](http://webchat.freenode.net/?channels=%23ipfs)
![travis](https://travis-ci.org/multiformats/go-multihash.svg)

> [multihash](//github.com/multiformats/multihash) implementation in Go.

## Table of Contents

- [Install](#install)
- [Usage](#usage)
- [Maintainers](#maintainers)
- [Contribute](#contribute)
- [License](#license)

## Install

```sh
go get multiformats/multihash
```

## Usage

```go
package main

import (
  "encoding/hex"
  "fmt"
  "github.com/multiformats/go-multihash"
)

func main() {
  // ignores errors for simplicity.
  // don't do that at home.

  buf, _ := hex.DecodeString("0beec7b5ea3f0fdbc95d0dd47f3c5bc275da8a33")
  mhbuf, _ := multihash.EncodeName(buf, "sha1");
  mhhex := hex.EncodeToString(mhbuf)
  fmt.Printf("hex: %v\n", mhhex);

  o, _ := multihash.Decode(mhbuf);
  mhhex = hex.EncodeToString(o.Digest);
  fmt.Printf("obj: %v 0x%x %d %s\n", o.Name, o.Code, o.Length, mhhex);
}
```

Run [test/foo.go](test/foo.go)

```
> cd test/
> go build
> ./test
hex: 11140beec7b5ea3f0fdbc95d0dd47f3c5bc275da8a33
obj: sha1 0x11 20 0beec7b5ea3f0fdbc95d0dd47f3c5bc275da8a33
```

## Maintainers

Captain: [@Kubuxu](https://github.com/Kubuxu).

## Contribute

Contributions welcome. Please check out [the issues](https://github.com/multiformats/go-multihash/issues).

Check out our [contributing document](https://github.com/multiformats/multiformats/blob/master/contributing.md) for more information on how we work, and about contributing in general. Please be aware that all interactions related to multiformats are subject to the IPFS [Code of Conduct](https://github.com/ipfs/community/blob/master/code-of-conduct.md).

Small note: If editing the Readme, please conform to the [standard-readme](https://github.com/RichardLitt/standard-readme) specification.

## License

[MIT](LICENSE)
