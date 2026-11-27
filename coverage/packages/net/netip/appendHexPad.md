# Independent review of `net/netip.appendHexPad`

**Verdict: FAITHFUL**

Evidence SHA-256: `2475f740f05203370d8721af11f055784e45a98f046e9247732eb792d6e9d1ec`

I checked evidence commit `d3b9e83` against Go 1.27.0 `net/netip/netip.go`. The padded hexadecimal helper passes all four digits to a single append. With two available bytes its abff result detaches and leaves both old spare bytes intact. At exact fit it updates all four positions, and input 10 checks the leading zero padding. The model preserves raw digit contents and agrees on whether the returned slice aliases the original backing.

I ran `go run ./difftest/net/netip/text-digits` to replay 327,683 rows per declaration on pinned native darwin/arm64 and compiled Lean. The corpus exhausts every uint16 and spare capacity from zero through four, with the low byte used for decimal, plus three nil fixtures. Five independent literal fixtures also passed with `go run ./difftest/net/netip/text-digits/review`. They distinguish single-stage and digit-by-digit writes and check exact fit, zero prefix length, nonnil empty backing and nil backing. Allocator-dependent growth capacity and resource failures are outside scope. [Independent control](../../../../difftest/net/netip/text-digits/review/main.go)
