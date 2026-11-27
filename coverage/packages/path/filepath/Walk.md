# Independent review of `path/filepath.Walk`

**Verdict: FAITHFUL**

Evidence SHA-256: `6375b52e793bb67f4b3a528ccfdf2108e6eeb282a4ad32ef1433bb5ecf859ee6`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `2a4dde631027dd96e142fe949cf59af98cd0008dc35cbd09d17ee81cc16e4d42`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `58495f7063399213f04efe924afd0f566122c66db118ef2cbc81b14cca893972`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `a61ec07c12b9e054f5edc84fb6def6bae6d8e1e851d05314bd599ff8fff02d7e`.

Reviewed b6a0ffb against Go 1.27.0. The wrapper calls Lstat on the root, passes lookup failure and nil FileInfo to the callback, and delegates a successful lookup to the private walker. The public result normalizes exact SkipDir and SkipAll sentinels. A wrapped sentinel or another callback error is returned unchanged.

The independent control passed six fixed model results and seven native sequences. It distinguishes SkipDir from SkipAll in a nested walk, checks missing-root callback data and exact returned error identity, and removes a child during the parent callback to exercise the live Lstat failure branch. The child was already listed and reported with nil FileInfo, and SkipDir allowed the final root-level sibling. Listing after the parent callback would omit the removed child. Treating a wrapped sentinel as the exact sentinel would suppress an expected error.

The full replay passed 4,096 exhaustive callback maps on the fixed tree, 60,000 seeded repeats under 17, 2718 and 65537, four native missing-root outcomes and twelve fixed model rows. The model takes supplied filesystem observations and a sufficient traversal depth. Native Windows and Plan 9 remain outside the local test matrix.

Run the reviewer control from the repository root after building the Lean executable.

```sh
GOTOOLCHAIN=go1.27.0 go run difftest/path/filepath/walk/review/main.go
```
