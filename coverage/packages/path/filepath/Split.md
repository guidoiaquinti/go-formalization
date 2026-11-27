# Independent review of `path/filepath.Split`

**Verdict: FAITHFUL**

Evidence SHA-256: `32eb0b696390f540d03c647cb83c0af08e7382d4a09b46f8f02530a23650062a`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `5f00a644edbdaeb0b9611f0a035887dfa3f95eb46e76e886004f204e1d76e096`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `21bfb2c60cc90ea3cbd08ec0bb19812117869ee86162d2ad253687202ee393e1`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `e14665d717c2bdf0d831a22e7b3b779c9a351491b4d9d4caff591c26cecc7a01`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `80abc5c0bdc73486b701b515c909e6291fe8a6e76d3c268716070fedb72059f0`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `e4e3504a228525df786403bc32e74614cd3d7056d1f31cb284fef22283bda1ef`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `426e9a304fc3cc08cf3f6617982755097cfa3594ea00a71d89b806d0b42daf20`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `fb6a2b143fd4869aa6c1d16184c63bd7df433f01d91ea3759a4133c2c3f381d1`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `19c55d503fc0c0da9cbc957629a5a0e4c336f1f8a97801a37b6d485131e4cc0e`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `0e0ebaef4f440ae2c2d71ec1355b96c9ba38d3a1e7dffed0506c1555e9f57eae`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `e124a09119f253977c3b3256126032e75f852c065bcc76481e7ad7ab67b2937a`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `5d769b9330da172ede1aaa1d997a678d355b26909a7940efe522de057bc361c1`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `2801aa8f9f402bc533c87d7a34ebf9b7949c2f5ca4092aa3e7e66e913de2c831`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `9724f254715de4e72fb7f78caf7ffb0dfd0fde104bdc91407ca6ef832cc5d672`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `e2538db68d4cbe754c33f610de946bbb85a5e7f208ce416668d2f457341a7466`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `e8a09526ff58adad789cc130d48575719b703e4ded3d374bac19b16caed0d0c9`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `32f4e08da6599e1b99a15029eaa5302db78a2995c66db4daca4153bcf0c5ff55`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `49eba0ba43edb5776a8b9fa8c50ca3dd64be6762f8fcdb732987ebe3a493cf50`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `0dcee4e0e33208b9ecbbe6196290f691224091bc1ef51e68f0d6b563ef88f025`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `3b87dd8847cd39b8a39bf87e5ebf5f7cc333d4f12de3b351422b8f1ebbec7daa`.

Reviewed a9c2e53 against Go 1.27.0's public wrapper and internal/filepathlite.Split. The pinned source scans backward from the end of the original string until reaching the volume or a recognized separator. Taking the greater of the volume length and the last recognized separator's end gives the same boundary. The Windows VolumeName helper normalizes slash bytes in its result, which preserves the volume length. Split uses that length to slice the original path bytes. The result retains the separator, duplicate separators, dot and parent components, NUL and malformed UTF-8.

The independent control checks 15 pinned native calls and 30 explicit Lean profile outputs. Cases cover drive-only and drive-relative paths, a slash byte before a drive colon, mixed-separator UNC, UNC and device volumes with parent components, a case-folded extended UNC prefix, and raw malformed bytes. The volume-blind last-separator alternative fails seven Windows fixtures. Cleaning before splitting changes the preserved bytes in `a//.././b`. Every fixture independently checks that the two expected components reconstruct the original bytes.

The author differential checks 528 boundaries and 60,000 samples across seeds 17, 2718 and 65537, plus 16 explicit Windows model rows. Native Windows is outside the locally verified scope. The claim concerns bytes of the two outputs and excludes backing-storage identity.

Run the reviewer control from the repository root after building the Lean executable.

```sh
GOTOOLCHAIN=go1.27.0 go run difftest/path/filepath/split/review/main.go
```
