# Independent review of `path/filepath.join`

**Verdict: FAITHFUL**

Evidence SHA-256: `25e4a5d01b3b6c1a890f97894db95cfb5b7b03a0dbf06ab590df359b107df803`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `ef3f1b3f9792ae16486c124f7a0779f83eaf1bc29f1c4bcf5189e310f91b83b2`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `ebd1555d945f9484a5e78cde2779525471e2a9fb6083e45994ed5bc6c467d266`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `fd16c4ecbb1c29044a0d75dd516f2c1cac13de0ca5de8ec6e970e63d4e9df426`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `9338bc06bdecd8ed0b55f262ef64d0ff0b92b78b331b6e4de852ce93d5de7138`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `5815ae2e7999430721fd1f00a522368946d28e3be64cc4b77b8a0fc7f1ed8be0`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `a18af2a23fd63a5a7abbc023a629d859565363e775207e4bbb94201bd71d38a8`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `fae100b706c8f07fbdf0ce98ea7e44f65469886c342e002fdf8b6e0c452eb0e6`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `2eef22760a60cbf93cd7a18b4d1f0e58219508e4d87ebaeb66b2a004c3de2c30`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `01904837d8356a7fd473d5ae70a0b23fcdfa21679cd25b5090526b6136c8f554`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `5ce10188b6e7baf5f2e001000da5c8845ee945632d2f0ec767baf8b41d2c416c`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `858bf58dc153afe693583f8137ff0dd24c8c2b8737f2c67095bc4b8aeb84bb4e`.

Reviewed e5a73b3 against Go 1.27.0. The private implementations match the model shared with the public wrapper. Unix and Plan 9 select the first nonempty element, combine the suffix with slash separators, and lexically clean the combined path. Windows strips leading separator bytes after an existing separator, retains relative concatenation after a drive colon, and guards a separator-only root before a complete double-question-mark component. Windows Clean preserves its pinned behavior when stale allocated bytes trigger a prefix guard.

The runnable control calls the actual private function in pinned Go and checks 13 fixed native results and 26 Lean profile rows. It covers true UNC input versus adjacent separators, drive-relative and drive-rooted forms, special device-prefix handling, empty lists, malformed bytes, and inherited stale-colon cleanup. Inserting a separator unconditionally after a drive colon and omitting the complete device guard both disagree with fixed witnesses.

The portable extracted Windows Join and Clean helpers agreed with the private Lean executor on 60,012 cases, including 12 fixed cases and 20,000 generated element lists for each seed 17, 2718, and 65537. The independent author replay also passed 528 native boundaries, 60,000 native samples, and 16 fixed Windows profile cases. Native Windows and Plan 9 execution remain outside the local matrix.

Run the reviewer control from the repository root after building the Lean executable.

```sh
GOTOOLCHAIN=go1.27.0 go run difftest/path/filepath/join-private/review/main.go
```
