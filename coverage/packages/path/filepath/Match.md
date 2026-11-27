# Independent review of `path/filepath.Match`

**Verdict: FAITHFUL**

Evidence SHA-256: `742beb61258fa6f9e83655458d72a1db827f6ce35ed47e0a3036e5bd3fbd3954`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `6b2ca08f1d1829ff6e62e68e99d1fbec29754a9dff0cfc563a460f54ce5122f8`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `9c366f0bc47ca3a346d8baa600447fe997ce3198feda2b42ec24be9a2d83e606`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `3dac5df6f4b8b4b786cb07c7359ffa52ad5f1b777bff805c2df80a49881a1597`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `fad020d4e3cf15d79459f606669e23e5fca46563a4c93084d938ec3d8e4cb3ac`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `6edc19bad211b164a203508da12dba109b6821f1b846c9873c6794f6b8d470b0`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `c72d909c9c0b65e4e03bce2737b5e0a2e1171173f5dc7ec6c96fa6b7cf6ffa5c`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `057b539b8261421269790de64784f2c6d47ad2eaedc61a3400c919c754e14291`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `595d631a63c8fed6626308b32f6cec74026a9e5c7226b8a6d400357425a05b4b`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `63e14c708eb7d30080fdc24a3d23ccbe675fe2980b12c8d515fa4839d73221ad`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `c2a3f31aadec57863dcc619f7a346d0997d2b93c5e413d7637f10de7883a080b`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `5c99944fdd4d603e04d5908cc8b8dfeabefad99f7b7eefcaccdee388d561e8c4`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `1e59d1371a06e4c39214750da52a8de292a82addc3030c98207e54360a8a9c7d`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `8ef614d4e3c7a83ae6cd7087a0cf1a8a084179633005817886c7bebf5447f578`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `47224012302c33b936923cb7b10e11d6cf0a578a74fdbcd762fc4b75c9fd0234`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `62868f877216baaa046f87d8ad4c636a0f41919b13a0103f4751e11c89c427bf`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `4ccde2d60c121ddbd263462210f118d959c1663980f86cc03b990336d44385b8`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `86cd96af3a3de4a2610e877572e4a51bcbf44c852b6611aa6daa22e3bacde9f5`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `3ff3f8ad0bfe3f734ac5934f627f42a81667faaf1502af306ef2eab0d9f34d81`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `de5b890b5cc20f35447bee5e367e7381b2ed2e12fb9ac4809c72849804a8ed58`.

Reviewed Go 1.27.0 src/path/filepath/match.go and evidence 2365ece. The model follows scanChunk and matchChunk at the current name position. An intermediate match commits the first successful chunk position. A terminal match also requires full name consumption, and a preceding star can keep trying positions until the platform separator. The trailing-star branch checks the remaining name directly. The matchChunk helper reports syntax errors in the current chunk even after a literal mismatch. A failed earlier chunk prevents the caller from reaching a later malformed chunk.

The control checks 22 fixed calls to pinned native Match, comparing ErrBadPattern by identity, and 44 model rows over the two platform profiles. The cases distinguish eager full-pattern validation, prefix-only terminal acceptance, separator-blind star skipping, and global backtracking. In particular, the Unix input pattern `*[a/]*b` and name `a/b` fails because the first chunk accepts `a` and the second star cannot skip the slash. Globally revisiting the first chunk would accept the slash via its character class. Separate witnesses exercise escaping, question mark versus a native separator, class matching of a separator, and malformed UTF-8 in operands.

The differential checks 532 boundary cases, 60,000 samples across seeds 17, 2718, and 65537, and 20 fixed cases per modeled profile. Native Windows execution remains outside the verified scope.

Run the reviewer control from the repository root after building the Lean executable.

```sh
GOTOOLCHAIN=go1.27.0 go run difftest/path/filepath/match/review/main.go
```
