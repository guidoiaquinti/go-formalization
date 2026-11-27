# Independent review of `strings.ToUpperSpecial`

**Verdict: FAITHFUL**

Evidence SHA-256: `48e211fc2765f93aa3d6f84af649a379c1aa5a856f7fc7ef3d2112b9e68f31dd`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `249c7694b1e0a0ad4f0446cb87c40d489acb8e356d5e43c1a201650506be9460`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `f0317700082faccf6c314ed6f3d2f5354f575ca3b41c526c2762fb455ffc36e2`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `cbf877e2e2d06305c54c66ec27f9d376a45eaaaad99082b8ecf2f1416a28c2f7`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `08e4e883af7537cfdfa665c722f10a5b7dca1172df5d741e23630497675f05a1`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `5cf40c0ac79e99a048c53652dc3b8df51618998bd57633b297b7d6f422f6566a`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `4851d5fe611d60c3426c783c542b211270ddc0bbabcfa56978fd5328e8315c02`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `95ced3c546d632a7b6648f0176306a2396047bbf4ca6ff240434a5674406c575`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `4d357fc5e1c1215c8846a12f8dc5d74726e5ccadd163f0f7c968572b70a8fb58`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `9fedf00122048bc48e66914097de055a9e53b211340be5863c6e02500712593a`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `9464a2ab33f042a81e5642d34e3345a2f9b9d7a6d8989ed6f921cf2097a39eed`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `74303684d52c7d01a842f1c128e4e28584f6122ddb39ff8096f2746857951968`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `9f863c1869b9116d5949b5a5eadb20189a5be3fdb690973bfcf68ffd7f4092fa`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `f1eadee5349adc15f5626415add50097829db79e0d5215d66db81f9661156698`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `69a8db0daf1e662411f0d2bee59f8a4e0504ba8597c6171c75f5355db8ff2173`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `8971fd30c8b5ccda7f1e14fa6823300c85350fc7b95d8aa629677c79a903f624`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `d1b32149320d2aefd0893546215beac2c01dd962019d7fb1e49d0b5f2d1e5af0`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `2e8104c3938e55508f4bede517010c445306d458522a41e7052cf1adb63c2fc1`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `412c18e003c7d993aa48d500e4b4152c8b607d402322b5170cab575461b22791`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `e4e5e4c1ee6c7eefd72a7ea9663000ac3eeaf5d3cbaa62a95ccfe41c802595e7`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `05121860acc3f67e3ea52c744d828a1351fd80a8daa538315d6a09feb473989c`.

Go 1.27.0 routes the supplied SpecialCase through ToUpper and then through Map. The model uses the upper delta selected by the pinned half-open midpoint traversal. Inclusive signed 32-bit comparisons preserve traversal through overlapping or unsorted data. A found entry suppresses fallback even when its upper delta is zero. The alternating marker computes the even offset from the range's low endpoint with signed 32-bit wraparound. Misses use the pinned Unicode uppercase mapping.

The independent control invokes the pinned public Go function and the Lean executable on nine fixed-output witnesses. It checks midpoint and zero-delta identity, an unsorted miss, signed endpoint conversion and wrapping, negative drop, invalid surrogate output, Turkish case and malformed bytes adjacent to a literal replacement rune. The control rejects first matching range, fallback on an identity hit, and using the lower delta on the overlap witness.

The author differential replay checks 1,248 boundaries and 20,000 samples for each of seeds 17, 2718, and 65537. The verdict covers exact output bytes over representable strings and supplied SpecialCase tables. It excludes allocation and unsafe backing storage observations.

Run the reviewer control from the repository root after building the Lean executable.

```sh
GOTOOLCHAIN=go1.27.0 go run difftest/strings/to-upper-special/review/main.go
```
