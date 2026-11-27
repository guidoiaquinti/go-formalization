# Independent review of `net/netip.appendDecimal`

**Verdict: FAITHFUL**

Evidence SHA-256: `c040732ad04e872e5edc453b24be4d01b41ebb0904e71b2ae370e4d3c92d434f`

I checked evidence commit `d3b9e83` against Go 1.27.0 `net/netip/netip.go`. The decimal helper takes a uint8 and emits each digit as a separate append, with the hundreds and tens gates at 100 and 10. With two spare bytes and value 255, the old backing receives ASCII 2 and 5 even though the result detaches on its last digit. The model preserves raw digit contents and agrees on whether the returned slice aliases the original backing.

I ran `go run ./difftest/net/netip/text-digits` to replay 327,683 rows per declaration on pinned native darwin/arm64 and compiled Lean. The corpus exhausts every uint16 and spare capacity from zero through four, with the low byte used for decimal, plus three nil fixtures. Five independent literal fixtures also passed with `go run ./difftest/net/netip/text-digits/review`. They distinguish single-stage and digit-by-digit writes and check exact fit, zero prefix length, nonnil empty backing and nil backing. Allocator-dependent growth capacity and resource failures are outside scope. [Independent control](../../../../difftest/net/netip/text-digits/review/main.go)
