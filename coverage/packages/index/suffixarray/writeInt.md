# Independent review of `index/suffixarray.writeInt`

**Verdict: FAITHFUL**

The pinned writer receives all ten bytes through the scratch alias, even when the signed encoding consumes one byte. With a visible length of one and capacity twelve, encoding −1 keeps nine old padding bytes in the offered frame and returns the fixed writer error. A mutant that offers only the encoded byte changes the observed write.

The pinned private-entry overlay and Lean comparison passed 5,944 boundaries and 60,000 cases across seeds 17, 2718, and 65537 on darwin/arm64. The three independently specified fixtures passed on pinned Go and Lean and distinguished three altered outcomes. Scope follows the registry, including finite bounded reads, fixed writer replies, legal slice geometry, and no callback mutation.

Run `GOTOOLCHAIN=go1.27.0 go run ./difftest/index/suffixarray/integer-io/review`.

Evidence SHA-256: `407176cc80b3819722c2d24edd33928ac1788ad57315ffcb5a60c9eb753d0a82`
