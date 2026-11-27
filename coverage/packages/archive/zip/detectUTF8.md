# Independent review of `archive/zip.detectUTF8`

**Verdict: FAITHFUL**

Evidence SHA-256: `eed94e964f2ad4f83fcbb680ea38dc14d2b156c1dbe8e76c1b4f28919a8cd1ee`

Go 1.27.0 accepts valid Unicode and sets the UTF-8 requirement outside the ASCII interval 0x20–0x7d or at backslash. It distinguishes a literal replacement rune from a malformed width-one decode. A malformed suffix clears both flags even after backslash set the requirement. The pinned author replay matched 125,802 inputs, with samples at seeds 17, 2718 and 65537. [Model](../../../lean/GoLean/Archive/Zip/DetectUTF8.lean) · [Comparison](../../../difftest/archive/zip/detect-utf8/check.go)

The independent command `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/archive/zip/detect-utf8/review` passed 11 native and Lean cases and rejected nine changed outcomes. Cases bracket both ASCII cutoffs, cover literal U+FFFD and U+10FFFF, reject surrogate bytes, and exercise malformed bytes after backslash. The scope covers finite raw byte strings and the two detector flags. [Reviewer control](../../../difftest/archive/zip/detect-utf8/review/main.go)
