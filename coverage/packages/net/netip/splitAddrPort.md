# Independent review of `net/netip.splitAddrPort`

**Verdict: FAITHFUL**

Evidence SHA-256: `261ae4d74d71a7c932ba6e417a7d0a585ca21e78b31a42a36ec807f2afdc6ee9`

I reviewed evidence commit `7f43d9f` against Go 1.27.0 `src/net/netip/netip.go`. The source selects the final colon, checks raw host before port, then strips a complete bracket pair. The fixed control uses a documentation-prefix IPv6 address and verifies the rightmost split, bracket balance, host-before-port precedence and the shape-only empty bracket case.

Pinned native Go and Lean agreed on 1,227 fixed rows and 20,000 sampled rows each for seeds 17, 2718 and 65537 on darwin/arm64 with `go run ./difftest/net/netip/parse-fields`. Eight independent rows passed with `go run ./difftest/net/netip/parse-fields/review` and detected seven wrong outcomes across the group. The scoped kernel audit passed. Field parsing uses valid subrange indices and four-byte destinations. The split model covers structure. Later address and port validation and rendered IPv4 diagnostics are outside scope. [Independent control](../../../../difftest/net/netip/parse-fields/review/main.go)
