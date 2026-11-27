# Independent review of `index/suffixarray.readInt`

**Verdict: FAITHFUL**

The pinned code calls ReadFull on the ten-byte view and then decodes the original visible slice. An immediate EOF with a visible scratch beginning in byte 1 returns −1 and EOF. Returning zero immediately after the read error changes the decoded value.

The pinned private-entry overlay and Lean comparison passed 5,944 boundaries and 60,000 cases across seeds 17, 2718, and 65537 on darwin/arm64. The three independently specified fixtures passed on pinned Go and Lean and distinguished three altered outcomes. Scope follows the registry, including finite bounded reads, fixed writer replies, legal slice geometry, and no callback mutation.

Run `GOTOOLCHAIN=go1.27.0 go run ./difftest/index/suffixarray/integer-io/review`.

Evidence SHA-256: `8d804a7a07d7a7299b6f6050ed72408c44346e714f25e1fb9a6af856a7a6a919`
