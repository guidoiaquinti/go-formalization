# Independent review of `net/netip.uint128.bitsSetFrom`

**Verdict: FAITHFUL**

Evidence SHA-256: `e907a4f72ba0997ab9ac41d812985105a3b697634e6f00e08021aaf80da0d646`

I checked evidence commit `ab8244f` against Go 1.27.0 `net/netip/uint128.go`. The source ORs in the complement of the high-prefix mask. Bit zero sets all bits, bit 128 is an identity, and bit 129 faults. The fixture rejects preserving the entire original word at bit zero.

I ran `go run ./difftest/net/netip/uint128` on darwin/arm64 and reproduced 1,796 boundaries and 20,000 cases per seed 17, 2718 and 65537. The compiled Lean and direct pinned private Go observations agreed with the arbitrary-precision reference. Five separate literal fixtures passed both engines with `go run ./difftest/net/netip/uint128/review`, with eleven fixed mutation checks across the group. The scope records panic classification and direct field-address behavior where applicable. Exact runtime fault text and concurrent unsafe access are outside scope. [Independent control](../../../../difftest/net/netip/uint128/review/main.go)
