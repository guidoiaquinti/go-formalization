# Independent review of `path/filepath.evalSymlinks`

**Verdict: FAITHFUL**

Evidence SHA-256: `a13385d70d9a2f54abe7ae6ae399b45dd7dbf3d0277f90c504d2590ea3bc7a4f`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `8cc0f73e42330b80032bfc4d7e308c7b3528bbe0c73052df47c61dcf79cc8567`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `565a613903f7a5343277f2425109ffd22df407c0d8ca409cc660dbe00c3c79ad`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `992f10d1c129af26331b8314afac96a71a41c57b1cc9e830b6d8c2397ed42977`.

I checked the three pinned Go 1.27.0 platform wrappers. Unix directly returns the walker result. Windows propagates a failed walk and passes a successful walk's exact path to normalization, returning the normalization result. Plan 9 skips Lstat on the empty input and otherwise cleans a successful path or propagates the observed error. It substitutes the ENOTDIR marker only when the exact lowercase byte substring appears in the error text. The model's observed error carries opaque identity and message bytes.

The author checked 268 boundary paths and 60,000 seeded paths using seeds 17, 2718 and 65537 against native walker delegation, plus nine fixed Lean wrapper results. Independent control passed six additional inputs through the pinned darwin/arm64 private wrapper and direct walker and nine literal modeled cases across the three profiles. The cases cover conflicting delegate outcomes, byte-exact error propagation, Windows normalization after success, Plan 9 skipping an error on empty input, and matching the error substring through surrounding NUL or malformed bytes. Scratch Lean mutations that skipped normalization, called Lstat on empty input, or substituted every Plan 9 failure failed two, one, and one controls.

The review covers delegation and supplied outcomes. Actual symbolic-link traversal, normalization, Lstat implementation, and native Windows and Plan 9 execution remain outside the scope.
