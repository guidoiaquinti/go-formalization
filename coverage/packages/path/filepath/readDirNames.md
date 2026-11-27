# Independent review of `path/filepath.readDirNames`

**Verdict: FAITHFUL**

Evidence SHA-256: `8a9aa0a1f30442edf348a6a33f0dd0cf2ecb38b467f01460c857a77abddd7470`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `29ff91e688ec429a0bd4a78026314431284c01fd6c1f3dc375c8066e903f09ce`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `5624ee6efb1fc8652594e493cbfa5a17f264d0f55c15d926a4676b723d19f3c4`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `4d1bcc46a71a48a78cf94d00c201dfe54c10f8972ce5e75a082885c5f8b0131a`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `c8cebc5a97b6a970a57ea30da6ca066029815a8a81676596abe40450640e7b9a`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `c5779d68e7aa7065a0cfe54969e1b4574278e03dff1b1cd497088a184c3006be`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `4ebfe87fd34ced350a41ea0bcf13fd1c1182b1ab6ff7744767efbe3dda498a57`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `27a0778a155492c3820706577904670aeac6a6a62819ddf7832d71aa3edc06f1`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `ea04c6dc0912fbb75bb107834b3bb83e6b28c9120a83519e0e8d3b8f41babfd9`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `b8d35d7e1e00215df8bfacd979988f7f595c2792beae9f55577ec9761ce3f0d8`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `cb69f663a3cfc025dcc0240f33e699b05e9f18c04242e733bbb1ccc80f063640`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `aa4dbda853f8fb622c758037873346f144377d279fd246fdc58ffe0713c0eda4`.

Reviewed aeba703 against Go 1.27.0. Open failure returns immediately with nil names. Once open succeeds, the source calls Readdirnames with -1 and closes the handle before examining the read error. It preserves the read error and discards even partial names in that branch. Close errors do not change either a successful read or the read error. The successful path sorts by raw string byte order and retains duplicates in the supplied observation model.

The independent control checks five fixed model outcomes and reruns four native fixtures through the private Go method. The fixtures cover an empty directory, a directory populated with distinctly named entries, a missing path, and a regular file. The model rows distinguish close-error precedence, discarding partial results, empty success and byte-order sorting without deduplication. Returning the close error in the successful-read row or sorting case-insensitively in the mixed ASCII row would change the recorded output.

The author replay agreed on all 60,260 modeled cases, including 260 boundaries and 20,000 each under seeds 17, 2718 and 65537. The native fixtures exercise actual open and read outcomes. Injected close errors and arbitrary raw names are model-only; the replay serializes operation count while the model equations expose operation order.

Run the reviewer control from the repository root after building the Lean executable.

```sh
GOTOOLCHAIN=go1.27.0 go run difftest/path/filepath/read-dir-names/review/main.go
```
