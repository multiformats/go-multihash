#!/bin/sh
#
# Copyright (c) 2015 Christian Couder
# MIT Licensed; see the LICENSE file in this repository.
#

test_description="blake3 tests"

. lib/test-lib.sh

test_expect_success "'multihash -a=blake3 -e=hex' succeeds" '
	echo "Hash me!" >hash_me.txt &&
	multihash -a=blake3 -e=hex hash_me.txt >/dev/null
'

test_hasher() {
        INPUT_LEN=$1
        EXPECTED=$2

        test_expect_success "'multihash -a=blake3 -e=hex' output is correct (input len=$INPUT_LEN)" '
        go run ../t0030-lib/data.go $INPUT_LEN | multihash -a=blake3 -e=hex >actual
        echo $EXPECTED >expected
        test_cmp expected actual
'
}

# Copied from official test vector, see:
#
# https://github.com/BLAKE3-team/BLAKE3/blob/080b3330159a19407dddb407dc917925ac40c4d3/test_vectors/test_vectors.json

test_hasher 0 "1e20af1349b9f5f9a1a6a0404dea36dcc9499bcb25c9adc112b7cc9a93cae41f3262"
test_hasher 1 "1e202d3adedff11b61f14c886e35afa036736dcd87a74d27b5c1510225d0f592e213"
test_hasher 2 "1e207b7015bb92cf0b318037702a6cdd81dee41224f734684c2c122cd6359cb1ee63"
test_hasher 3 "1e20e1be4d7a8ab5560aa4199eea339849ba8e293d55ca0a81006726d184519e647f"
test_hasher 4 "1e20f30f5ab28fe047904037f77b6da4fea1e27241c5d132638d8bedce9d40494f32"
test_hasher 5 "1e20b40b44dfd97e7a84a996a91af8b85188c66c126940ba7aad2e7ae6b385402aa2"
test_hasher 6 "1e2006c4e8ffb6872fad96f9aaca5eee1553eb62aed0ad7198cef42e87f6a616c844"
test_hasher 7 "1e203f8770f387faad08faa9d8414e9f449ac68e6ff0417f673f602a646a891419fe"
test_hasher 8 "1e202351207d04fc16ade43ccab08600939c7c1fa70a5c0aaca76063d04c3228eaeb"
test_hasher 63 "1e20e9bc37a594daad83be9470df7f7b3798297c3d834ce80ba85d6e207627b7db7b"
test_hasher 64 "1e204eed7141ea4a5cd4b788606bd23f46e212af9cacebacdc7d1f4c6dc7f2511b98"
test_hasher 65 "1e20de1e5fa0be70df6d2be8fffd0e99ceaa8eb6e8c93a63f2d8d1c30ecb6b263dee"
test_hasher 127 "1e20d81293fda863f008c09e92fc382a81f5a0b4a1251cba1634016a0f86a6bd640d"
test_hasher 128 "1e20f17e570564b26578c33bb7f44643f539624b05df1a76c81f30acd548c44b45ef"
test_hasher 129 "1e20683aaae9f3c5ba37eaaf072aed0f9e30bac0865137bae68b1fde4ca2aebdcb12"
test_hasher 1023 "1e2010108970eeda3eb932baac1428c7a2163b0e924c9a9e25b35bba72b28f70bd11"
test_hasher 1024 "1e2042214739f095a406f3fc83deb889744ac00df831c10daa55189b5d121c855af7"
test_hasher 1025 "1e20d00278ae47eb27b34faecf67b4fe263f82d5412916c1ffd97c8cb7fb814b8444"
test_hasher 2048 "1e20e776b6028c7cd22a4d0ba182a8bf62205d2ef576467e838ed6f2529b85fba24a"
test_hasher 2049 "1e205f4d72f40d7a5f82b15ca2b2e44b1de3c2ef86c426c95c1af0b6879522563030"
test_hasher 3072 "1e20b98cb0ff3623be03326b373de6b9095218513e64f1ee2edd2525c7ad1e5cffd2"
test_hasher 3073 "1e207124b49501012f81cc7f11ca069ec9226cecb8a2c850cfe644e327d22d3e1cd3"
test_hasher 4096 "1e20015094013f57a5277b59d8475c0501042c0b642e531b0a1c8f58d2163229e969"
test_hasher 4097 "1e209b4052b38f1c5fc8b1f9ff7ac7b27cd242487b3d890d15c96a1c25b8aa0fb995"
test_hasher 5120 "1e209cadc15fed8b5d854562b26a9536d9707cadeda9b143978f319ab34230535833"
test_hasher 5121 "1e20628bd2cb2004694adaab7bbd778a25df25c47b9d4155a55f8fbd79f2fe154cff"
test_hasher 6144 "1e203e2e5b74e048f3add6d21faab3f83aa44d3b2278afb83b80b3c35164ebeca205"
test_hasher 6145 "1e20f1323a8631446cc50536a9f705ee5cb619424d46887f3c376c695b70e0f0507f"
test_hasher 7168 "1e2061da957ec2499a95d6b8023e2b0e604ec7f6b50e80a9678b89d2628e99ada77a"
test_hasher 7169 "1e20a003fc7a51754a9b3c7fae0367ab3d782dccf28855a03d435f8cfe74605e7817"
test_hasher 8192 "1e20aae792484c8efe4f19e2ca7d371d8c467ffb10748d8a5a1ae579948f718a2a63"
test_hasher 8193 "1e20bab6c09cb8ce8cf459261398d2e7aef35700bf488116ceb94a36d0f5f1b7bc3b"
test_hasher 16384 "1e20f875d6646de28985646f34ee13be9a576fd515f76b5b0a26bb324735041ddde4"
test_hasher 31744 "1e2062b6960e1a44bcc1eb1a611a8d6235b6b4b78f32e7abc4fb4c6cdcce94895c47"
test_hasher 102400 "1e20bc3e3d41a1146b069abffad3c0d44860cf664390afce4d9661f7902e7943e085"

test_done
