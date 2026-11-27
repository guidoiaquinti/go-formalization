# Independent review of `net/netip.parseIPv4`

**Verdict: FAITHFUL**

Evidence SHA-256: `4b492abe8d244854e1d2432605002db2204e44fda7f7528cc74f73f1c7a0b0dc`

I reviewed evidence commit `7f43d9f` against Go 1.27.0 `src/net/netip/netip.go`. Pinned Go supplies a zero-initialized four-byte destination and wraps successful fields in a zone-free IPv4 value. On failure the native check verifies a zero Addr and the original diagnostic. The control covers a documentation-range IPv4 address and a decimal octet above 255.

Pinned native Go and Lean agreed on 1,227 fixed rows and 20,000 sampled rows each for seeds 17, 2718 and 65537 on darwin/arm64 with `go run ./difftest/net/netip/parse-fields`. Eight independent rows passed with `go run ./difftest/net/netip/parse-fields/review` and detected seven wrong outcomes across the group. The scoped kernel audit passed. Field parsing uses valid subrange indices and four-byte destinations. The split model covers structure. Later address and port validation and rendered IPv4 diagnostics are outside scope. [Independent control](../../../../difftest/net/netip/parse-fields/review/main.go)
