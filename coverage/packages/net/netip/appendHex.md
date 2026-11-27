# Independent review of `net/netip.appendHex`

**Verdict: FAITHFUL**

Evidence SHA-256: `84d1289c4ae039135178060005d946904d7061608a0944e5cdb54204b80ae07f`

I checked evidence commit `d3b9e83` against Go 1.27.0 `net/netip/netip.go`. The hexadecimal helper gates leading digits at successive powers of sixteen and appends one lowercase digit at a time. Its abff witness writes a and b into old backing with two available bytes and detaches before either f. The model preserves raw digit contents and agrees on whether the returned slice aliases the original backing.

I ran `go run ./difftest/net/netip/text-digits` to replay 327,683 rows per declaration on pinned native darwin/arm64 and compiled Lean. The corpus exhausts every uint16 and spare capacity from zero through four, with the low byte used for decimal, plus three nil fixtures. Five independent literal fixtures also passed with `go run ./difftest/net/netip/text-digits/review`. They distinguish single-stage and digit-by-digit writes and check exact fit, zero prefix length, nonnil empty backing and nil backing. Allocator-dependent growth capacity and resource failures are outside scope. [Independent control](../../../../difftest/net/netip/text-digits/review/main.go)
