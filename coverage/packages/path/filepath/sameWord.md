# Independent review of `path/filepath.sameWord`

**Verdict: FAITHFUL**

Evidence SHA-256: `4687afe2a8370a17de38e7d95139078f12d6ef71472609bb769bf47263d1141a`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `3a3e1c26c250fa79bd392ba97a49f0ef284b08df5e868c7d1d5a51b19df18425`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `a94640654658d8cf141dd2f83a373a65d8d20fd776026eb130779373d47af0af`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `2891a556b19cf0c966b5dcd298022168c6c4ec00746c6b0264a04013a8da52a0`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `427a6e6ff03d26467bd0d814f13eccba270cf7a6c3b6badc5afdb0c8e009e8ad`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `6865fbd593062d703e25f8d5930ee254f3e12a30fcb540898fbe4bef88a2bba7`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `0f3060543d1e511590149bfaf4b0992a7cc184d1989f50f273e27b53bf025d2e`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `51fc9b7d383840d630c9a661a38d129a1afcc9d2e4586b3ac305c352ed29c803`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `dc10598d8e53547b814fd26baa6da7c9389489427286c888fccfc532fbe30588`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `50393e6f6ebf6540a11d7770e1bb23374de38e209a3338aa13a278c53f0d4a0d`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `ed9e287a97f36dfb45a03e0a424027e39253dd29e98f748f152f84a383f43eed`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `38f564bb1e94b3a87fff2ee2f88af19e4c3ecce5db5669ec4ad2c6dfa2cd3b92`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `f81ccd20e079e763c76cc8e6d4bc9a9968c9a804bb97068e96c23d6c3b91c853`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `5f506e06148223d3b477aedd55d1ee2794bbc0984d143fa93c0b9d5b14bb10cf`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `9f10d02437cd254a1eda63285ae36746c88dc87ecb5cd49c28ff5adcb813bc3f`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `c2b4abd4aba9ecf72e6ec199b1f3db992616299b0b95d188411bde2b16a87567`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `4055593215ea9d135994da3dd5eda586df492c7d863033e29db97b551c8de408`.

Pinned Go 1.27.0 uses byte equality on Unix and Plan 9, and delegates
to strings.EqualFold on Windows. The Lean wrapper selects the same
byte equality or checked Unicode simple-fold function. It keeps raw
separator and component bytes and does not apply cleaning, canonical
composition or multicharacter folding.

The author comparison checks 528 boundaries and 60,000 samples with
seeds 17, 2718 and 65537. The independent control passed 16 fixed
direct private Go and dual-profile Lean fixtures. Cases include long
s, Greek final sigma, Kelvin sign, German sharp s, unequal composed
and decomposed strings, separator and dot-component differences,
embedded NUL and both one- and two-byte malformed strings. Raw equality
in the Windows profile failed seven, lowercase equality failed two,
and slash replacement before comparison failed one. Native Go ran
on darwin/arm64. Windows and Plan 9 were checked by source inspection,
fixed semantics and Lean; native Windows and Plan 9 runs remain open.

Run the installed control from the repository root with
`GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/path/filepath/same-word/review`
after building the model executable.
