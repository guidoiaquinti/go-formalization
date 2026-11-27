# Independent review of `path/filepath.EvalSymlinks`

**Verdict: FAITHFUL**

Evidence SHA-256: `e83b5158140b2ccdc1226365a00d35b10e9288b920669d3b3f0c15bd0e669755`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `dbb2f83d9031b16018edba026d25595c8e4e1bb66fb3e7e282f041787e60f812`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `cc72e1bed25e2b5f9084fa1ecf753ac5c925ebdde0e91fb605a772d481076a6d`.

Reviewed c999fb0 against Go 1.27.0. The public entry delegates the input path to the platform implementation. Unix returns the walker result directly. Windows normalizes only a successful walk and preserves a failure from either stage. Plan 9 checks a nonempty path with Lstat, converts an error containing the case-sensitive phrase `not a directory` to the platform sentinel, and lexically cleans a successful path.

The independent control passed six fixed model rows and five native filesystem scenarios. Literal native results cover a link relative to its containing directory, a two-link chain, an absolute target, a broken link, and the exact cycle error message. Treating the target as relative to the process directory would break the nested relative link. Normalizing after a failed Windows walk would discard its supplied error. The two Plan 9 messages distinguish case-sensitive forwarding from case-insensitive matching.

The full replay passed 270 native boundaries and 60,000 seeded paths using 17, 2718 and 65537, plus nine fixed model rows. The claim covers public wrapper dispatch with supplied platform delegate outcomes. The native comparisons exercise Unix link traversal. Native Windows and Plan 9 remain outside the local matrix.

Run the reviewer control from the repository root after building the Lean executable.

```sh
GOTOOLCHAIN=go1.27.0 go run difftest/path/filepath/eval-symlinks/review/main.go
```
