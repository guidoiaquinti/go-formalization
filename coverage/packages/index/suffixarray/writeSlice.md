# Independent review of `index/suffixarray.writeSlice`

**Verdict: FAITHFUL**

The pinned loop reserves ten header bytes and admits an entry only with ten visible bytes remaining. A nineteen-byte view with capacity twenty and one entry leaves the count zero and offers the ten-byte header with retained padding. Admitting the entry based on its actual one-byte width changes the count, frame, and backing array.

The pinned private-entry overlay and Lean comparison passed 5,944 boundaries and 60,000 cases across seeds 17, 2718, and 65537 on darwin/arm64. The three independently specified fixtures passed on pinned Go and Lean and distinguished three altered outcomes. Scope follows the registry, including finite bounded reads, fixed writer replies, legal slice geometry, and no callback mutation.

Run `GOTOOLCHAIN=go1.27.0 go run ./difftest/index/suffixarray/integer-io/review`.

Evidence SHA-256: `7be134bea1aff2c396e82942237560a8dc32fa86cf4ef1e751f9d81a8cc96840`
