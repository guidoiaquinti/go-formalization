# Independent review of `strings.Title`

**Verdict: FAITHFUL**

Evidence SHA-256: `f1b02391fbb303b88030b16a14dfbaf2d09b0c2601ed812ec54221e95764bb74`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `d06e3a92af3ccc78acd88327f1d8bf81118dec447886da3b6e01d547fd14f62e`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `2f9388481bf90d29e69a66294d7cf7da15cde3604af0b8b94e5bd54ebe5a0204`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `76bc3ba89cfef4b07d32a2f6a99644110f601527e2d1e77372963c78fe871c4e`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `803b5496caa06ef0f00a607f988049c6556d1fecf84cb33618f565f1d22f47b2`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `5cb7aae86a943e3cd661d45df5fb5d4a0583d53d42b3e2910401229e388e2d02`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `50ee831ee93b300c85482c74c116d978551133b9aad20d6d714866886fe05914`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `91217210e46a21d4d0bf2268d5f0f7a4a99329d2791327228fb40b1351fa5629`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `e7a10495e9c26500d7d374a127613e53f46604c305898cf37f300dce2777464f`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `eef57201eed240809ea1a8170a5365ce5f7e587177ccd129613cbb68f43c9b10`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `e8d842ba09651add3b0f9a8c20e2bb32bb1d63709be52f2e760c5caf70800887`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `19dcce092e536a18b0d77dbbc63b253ce44b967bea5390a467b3d709e8153fba`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `87fb4d44fea8b5c0ce1f9aaaae4ba3408aa4b4f0e56cd90f37511961cfa57484`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `36364cbd3171c11e5727b9f4a524768d4a821cabc4957beede9870eb2fe84998`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `a3e14a8fdb19b9cf015db35617522b4ee3a1874a7eeea9bee853f91b9c4bf4bf`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `217788b684238a647f68130bec3de8a559b54c8ca1c403361d8af072ab831c74`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `ce9d276cb1758797dad5a2459222f22ea6f1dc21ef1510d2e364019f352ad672`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `379a24bc027fd54deed2678625d57edeb4085b99c2bde229f6bedee269e05574`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `366686ec5efa5b8645579a71fd05e0016bb2d0a4f76fedb0c45eff727c0ef001`.

Reviewed a6fdec0 against Go 1.27.0. The pinned closure starts with a space and remembers each original decoded rune before processing the next. Title mapping only occurs when this previous rune meets the pinned strings.isSeparator rule. The model preserves ASCII letters, digits, and underscore as nonseparators, and Unicode punctuation without whitespace status does not cause capitalization. The output retains original valid unchanged bytes and uses replacement rune encoding when input bytes are malformed.

The independent native and Lean control checks 12 fixed inputs. It distinguishes Unicode title mapping from uppercase on the digraph `ǳ`, ASCII apostrophe and hyphen from non-ASCII punctuation, Unicode whitespace from a zero-width nonseparator, and malformed bytes from a case boundary. It includes NUL, underscores, digits, a valid replacement rune, and consecutive truncated UTF-8. The control detects whole-string ToTitle, whole-string ToUpper, and identity alternatives.

The independent replay also passed the registered 1,080 boundaries and 20,000 samples per seed for 17, 2718 and 65537. The claim concerns output bytes and excludes storage allocation and unsafe backing identity.

Run the reviewer control from the repository root after building the Lean executable.

```sh
GOTOOLCHAIN=go1.27.0 go run difftest/strings/title/review/main.go
```
