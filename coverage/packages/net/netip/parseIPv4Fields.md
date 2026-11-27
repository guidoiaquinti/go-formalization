# Independent review of `net/netip.parseIPv4Fields`

**Verdict: FAITHFUL**

Evidence SHA-256: `a56bff8ca8863a961f7af9ba21ec90cc10b601219275c89f0fd89036190610b1`

I reviewed evidence commit `7f43d9f` against Go 1.27.0 `src/net/netip/netip.go`. The source parses the selected byte subrange but retains the original input in the structured error. The control verifies two completed field writes, untouched remaining bytes and an offending suffix clipped to the selected range.

Pinned native Go and Lean agreed on 1,227 fixed rows and 20,000 sampled rows each for seeds 17, 2718 and 65537 on darwin/arm64 with `go run ./difftest/net/netip/parse-fields`. Eight independent rows passed with `go run ./difftest/net/netip/parse-fields/review` and detected seven wrong outcomes across the group. The scoped kernel audit passed. Field parsing uses valid subrange indices and four-byte destinations. The split model covers structure. Later address and port validation and rendered IPv4 diagnostics are outside scope. [Independent control](../../../../difftest/net/netip/parse-fields/review/main.go)
