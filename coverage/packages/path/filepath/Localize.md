# Independent review of `path/filepath.Localize`

**Verdict: FAITHFUL**

Evidence SHA-256: `304a8f875e440cbb5d73fbead0745eff9bc4bd18671546617ce2d0f698ad25e2`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `7d73d072a23d3c093d34ab496274eac7d4a2186fb0f7e3b3950a9711d869f8e1`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `689f277bafebe6d6b9e22183983041ff5fc41ec849f02df0874c72726d0ab306`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `a8e6b8b27810475493bea82ff0f043ec4156a4c0d674e8b690b35036157dc0b8`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `c64a4f859470247e1847d3f4125d35fed6c22b0c5382bda07b0ad2c80c4d89d2`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `3c4840cb3a78776e2114572d6158500387bd79f3d0d4506bf5f4bb8a286631fe`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `4f6a33a6c8c83fba9e28692f74b163579a2ed7f778ba272eb7495cab19b652e1`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `9536f0388fa430fddcac73cb049f3739d12cf1f1d13494eda5f5ea659bdd63cb`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `722fc2193c3041fc6bdfdfab6855c10af0226a5e13159bcbd976fbb27159d726`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `7a970344eaaac07842bed729f3d7a8191e1106415bc8f6aef51ed4988687a774`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `dfbac2a3ab92c8d4890b6d2f7b05eeaf53a312c94a98e444c83b2762cb737b51`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `30c1dbf971c6843ab7a96cb27622f4d578d7adfc66e1375834fb18e651e829b9`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `6b32dc670dfe108abed3e5a657a43c49bed43e539ee65df68f0617c1739e25ef`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `c9313b60c63c101e894ca4b0591a500f355ddc014329d1a268008a86ded41c2a`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `42db3056f428aa72a6899697b2944efbf57c45783cee167271cee1129bd8cefe`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `091c84a0efb55fc942087e53d600e6f31337b28586f43793e36065eec2fe6cb5`.

Go 1.27.0 first requires valid UTF-8 and nonempty slash-delimited elements other than dot or dot-dot, with a special case for the whole path `.`. All profiles then reject NUL. Plan 9 additionally rejects a leading `#`, while Windows rejects colon, backslash and reserved device components before converting slash to backslash. Using the Windows locality predicate for individual validated elements preserves the reserved-name decision because valid elements already satisfy its other checks. [Model](../../../../lean/GoLean/Path/Filepath/Localize.lean)

Pinned darwin/arm64 matched 533 boundaries and 60,000 seeded cases at seeds 17, 2718 and 65537. The author also compared 21 fixed profile cases. The independent control passed 13 native and 30 Lean observations, including malformed UTF-8, actual RuneError and noncharacters, NUL accepted by ValidPath but rejected by Localize, lexical dot validation, Windows raw device-predicate arguments, and Plan 9 leading versus interior device markers. Incorrect variants that skip the NUL check, substitute IsLocal for ValidPath, omit Windows conversion, ignore or alter the device callback, or mishandle the leading Plan 9 marker are distinguished. Native Windows and Plan 9 execution, OS device-name detection and error object identity remain outside local evidence. [Comparison](../../../../difftest/path/filepath/localize/check.go) · [Reviewer control](../../../../difftest/path/filepath/localize/review/main.go)
