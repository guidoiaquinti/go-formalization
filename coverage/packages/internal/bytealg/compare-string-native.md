# Independent review of `internal/bytealg.CompareString`

**Verdict: FAITHFUL**

Evidence SHA-256: `bc53137b3103d2e1962b95cdb3a75cba8bf9e80604dc8450fd32d11b5684b61b`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `72f0bd73e7d6e992dc65f39a450ab87f87b3c465e9cfff349102e51242a9e461`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `ca7639d9da66fe4e455d070d71ff84b7d90a6cd5137e84bfa545f6390a70a64f`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `b8a58a00b72aae7dec4120ac38da06321d91a5cbb750e10c3d33e9062185407b`.

Reviewed 85d325f against Go 1.27.0. The pinned Go function forwards to runtime.cmpstring. Its return value is the sign of unsigned lexicographic comparison over the raw string bytes.

The reviewer control checks ten fixed inputs through all four native entry points and Lean. It covers nil versus nonnil empty, high-byte ordering, embedded NUL, partial UTF-8, strict prefixes, first mismatch before a differing length, alias fast paths, and offsets around multiword boundaries. Signed-byte ordering reverses the 7F/80 witness. Checking length first reverses the first-mismatch witness. Comparing beyond the supplied view would distinguish equal inputs with differently poisoned backing.

The full independent replay passed 67,940 boundaries and 60,000 seeded samples using 17, 2718 and 65537. The verdict covers exact scalar results on the exercised Darwin arm64 profile. Native memory access order, timing, and unexecuted architectures remain outside scope.

Run the shared reviewer control from the repository root after building the Lean executable.

```sh
GOTOOLCHAIN=go1.27.0 go run difftest/internal/bytealg/compare/review/main.go
```
