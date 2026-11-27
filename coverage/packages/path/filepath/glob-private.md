# Independent review of `path/filepath.glob`

**Verdict: FAITHFUL**

Evidence SHA-256: `355aad08f5bbfd25e9f26d1ee6c53af34e8fc4f539b17c638605976a1a701768`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `60af0254bbbb6460ad3abb8498ec18c99fa1f222063cf47b7dd14537e811e58c`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `11989d2baa1a92c3dced61640039ccf9bf3b0525a452f6385dc347ef19d4314f`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `fe421e2082313bfbe5f7886dd9610c7986ac22c2a01410e107dcc6d7565656ca`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `18bbda7afabcc4e2748b8e33464170ce15ef47f4784d387c9f94de6039fb6af7`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `44b9da6161aa97ba1bf1d52b48c193e3b03f64d02d4247032e81b43508054f7b`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `dce2e6c98f69eb07744349fbb0b601420ab5111144732d5bc42596a397cd3e15`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `b78862e6dcffe8afd0904616100f124e744c50b6b502f6ac3f433254683a2152`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `ac265ad164e2a5d3c2f19a8be020979fffe8090d6007aa9eb867a60f89c9e41c`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `f480f1606802c9d92c6574d62a4a34660fc92d7de46256943b7c2723a66d4bf4`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `30ef540dbc9f83392920d1fbe46a31b47a63373bef3a8fc1409564bd31c11296`.

Reviewed cee927b against Go 1.27.0. The private helper returns the incoming collection on stat errors, non-directory paths, or open errors. Once open, it ignores the read error and processes the available names, whether the read completed or returned partial names. It sorts names by raw bytes before matching. The first Match error returns the current collection and ErrBadPattern. The retained collection can include matches supplied by the caller before this directory was visited. The deferred close has no effect on the return value.

The runnable control checks eight fixed model outcomes and reruns seven native fixtures. It distinguishes eager pattern validation from checking only after reaching a name, including a malformed later chunk that a name mismatch can avoid. It checks preserving a supplied collection on syntax error, preserving nil versus allocated empty collections, and appending sorted joined results. An eager validation alternative disagrees with open failure, empty directory, and an early name mismatch.

The complete replay passed 264 boundaries and 60,000 simulated outcomes using seeds 17, 2718 and 65537. Native tests exercise real directories, missing paths and regular files, and compare ErrBadPattern by identity. Partial reads and open failure are represented by supplied observations. Native partial-read or close-error injection, native Windows and Plan 9, and backing allocation remain outside the verified scope.

Run the reviewer control from the repository root after building the Lean executable.

```sh
GOTOOLCHAIN=go1.27.0 go run difftest/path/filepath/glob-private/review/main.go
```
