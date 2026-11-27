# Independent review of `internal/bytealg.abigen_runtime_cmpstring`

**Verdict: FAITHFUL**

Evidence SHA-256: `e22f3aa0787e2f24a96ecceaabfab0d670c4ee95cf51e022c050d42db0428b47`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `ef7a4c6d11fc3a793140d51f1dd7f4cfe5cf7ed8bf1ed40fbabc3a678727f2c7`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `a74e530f3b57a2bec6b7f04ea95639e51bd314d4b4e92043c46f35e5e4e4a4a1`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `cd3a618f53bc5f9c8c04e3fce0b854f6f66ee3779cbbdb71a79986a420383e9d`.

Reviewed 85d325f against Go 1.27.0. The linkname entry targets runtime.cmpstring directly. The Darwin arm64 string entry supplies raw string data and byte lengths to the common comparison body.

The reviewer control checks ten fixed inputs through all four native entry points and Lean. It covers nil versus nonnil empty, high-byte ordering, embedded NUL, partial UTF-8, strict prefixes, first mismatch before a differing length, alias fast paths, and offsets around multiword boundaries. Signed-byte ordering reverses the 7F/80 witness. Checking length first reverses the first-mismatch witness. Comparing beyond the supplied view would distinguish equal inputs with differently poisoned backing.

The full independent replay passed 67,940 boundaries and 60,000 seeded samples using 17, 2718 and 65537. The verdict covers exact scalar results on the exercised Darwin arm64 profile. Native memory access order, timing, and unexecuted architectures remain outside scope.

Run the shared reviewer control from the repository root after building the Lean executable.

```sh
GOTOOLCHAIN=go1.27.0 go run difftest/internal/bytealg/compare/review/main.go
```
