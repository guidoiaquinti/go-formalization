# Independent review of `net/netip.uint128.halves`

**Verdict: FAITHFUL**

Evidence SHA-256: `2aa3b08391ae95b80ca82fe60afba1c1b9d51843fe08aeaa428eef5df2056972`

I checked evidence commit `ab8244f` against Go 1.27.0 `net/netip/uint128.go`. The Go overlay checks exact high and low field addresses and rejects swapped or aliased references. Each write is observed before or after the other as appropriate, and the nil receiver panics.

I ran `go run ./difftest/net/netip/uint128` on darwin/arm64 and reproduced 1,796 boundaries and 20,000 cases per seed 17, 2718 and 65537. The compiled Lean and direct pinned private Go observations agreed with the arbitrary-precision reference. Five separate literal fixtures passed both engines with `go run ./difftest/net/netip/uint128/review`, with eleven fixed mutation checks across the group. The scope records panic classification and direct field-address behavior where applicable. Exact runtime fault text and concurrent unsafe access are outside scope. [Independent control](../../../../difftest/net/netip/uint128/review/main.go)
