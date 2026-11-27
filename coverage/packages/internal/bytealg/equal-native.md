# Independent review of `internal/bytealg.Equal`

**Verdict: FAITHFUL**

Evidence SHA-256: `b2094356f0559891f46d833ae52c73357dd5afa7d6e6be0eb96f46c1d03e4ef1`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `d7463944b35b0759594a481f2b648d2d9009af50d9752658290bbba0aa30e3c5`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `6b68133c3b98eecf08dd169e2e27ed1881afbef8b47dc40afd4faa7a64728700`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `7ac54a8c89419598f6635ae0264802cbcea612a04728b82464e4e72c72268bb0`.

Reviewed 85d325f against Go 1.27.0. The pinned Go function compares string conversions of the slices, which tests byte sequences and length. Nil and allocated empty slices therefore compare equal.

The reviewer control checks ten fixed inputs through all four native entry points and Lean. It covers nil versus nonnil empty, high-byte ordering, embedded NUL, partial UTF-8, strict prefixes, first mismatch before a differing length, alias fast paths, and offsets around multiword boundaries. Signed-byte ordering reverses the 7F/80 witness. Checking length first reverses the first-mismatch witness. Comparing beyond the supplied view would distinguish equal inputs with differently poisoned backing.

The full independent replay passed 67,940 boundaries and 60,000 seeded samples using 17, 2718 and 65537. The verdict covers exact scalar results on the exercised Darwin arm64 profile. Native memory access order, timing, and unexecuted architectures remain outside scope.

Run the shared reviewer control from the repository root after building the Lean executable.

```sh
GOTOOLCHAIN=go1.27.0 go run difftest/internal/bytealg/compare/review/main.go
```
