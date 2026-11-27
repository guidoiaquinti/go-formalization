# Independent review of `internal/bytealg.Compare`

**Verdict: FAITHFUL**

Evidence SHA-256: `9eb36912ee812419ebfce4bfa2cf7fd571e69f3f92ddf03a8a79913bd6cdff9e`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `6f56948e51b5af44c2c4ed91dcc63628ee9c2bd24c5ff43fca53afc7c91fdb8b`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `6dfdf6281a981654593deb3a65fb31cb3fa2a94fa42896fce95f25472fa47dd9`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `405ebb1f0f8a98403a2432d47190f50d15e8c31341fcbcb85ebc109ffc0db6fc`.

Reviewed 85d325f against Go 1.27.0. The Darwin arm64 slice entry extracts data pointers and slice lengths before entering the shared comparison body. The model returns the sign of unsigned byte lexicographic ordering.

The reviewer control checks ten fixed inputs through all four native entry points and Lean. It covers nil versus nonnil empty, high-byte ordering, embedded NUL, partial UTF-8, strict prefixes, first mismatch before a differing length, alias fast paths, and offsets around multiword boundaries. Signed-byte ordering reverses the 7F/80 witness. Checking length first reverses the first-mismatch witness. Comparing beyond the supplied view would distinguish equal inputs with differently poisoned backing.

The full independent replay passed 67,940 boundaries and 60,000 seeded samples using 17, 2718 and 65537. The verdict covers exact scalar results on the exercised Darwin arm64 profile. Native memory access order, timing, and unexecuted architectures remain outside scope.

Run the shared reviewer control from the repository root after building the Lean executable.

```sh
GOTOOLCHAIN=go1.27.0 go run difftest/internal/bytealg/compare/review/main.go
```
