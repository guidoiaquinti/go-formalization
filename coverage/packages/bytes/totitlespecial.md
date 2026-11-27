# Independent review of `bytes.ToTitleSpecial`

**Verdict: FAITHFUL**

Evidence SHA-256: `ab0c1702508d60522d5452e9a3e693cea7abcef7738df88697ad17d50ed7320a`

Reviewed author commit `6762eff` against Go 1.27.0 `src/bytes/bytes.go:755–757`, `bytes.Map`, and `src/unicode/letter.go` range search and conversion. Tests an independent title delta, Turkish fallback behavior, and dropping malformed-input RuneError runes after an explicit negative mapping. An overlapping range fixture selects the exact midpoint hit. A custom alternating range exercises signed range bounds and int32 intermediate wrap. Native instrumentation also checks nonnil results, input-sized capacity when the output fits, and backing isolation.

Ten independent fixtures passed on all three pinned native entries and the compiled Lean executable. Using unsigned range comparisons changed one fixture, encoding negative mappings as RuneError changed three, and selecting uppercase for titlecase changed four. Replay using `GOTOOLCHAIN=go1.27.0 go run ./difftest/bytes/special-case/review` from the repository root.

The author compared 1,276 boundaries and 60,000 samples at seeds 17, 2718, and 65537. Allocation failure and exact allocator growth capacity remain outside scope.
