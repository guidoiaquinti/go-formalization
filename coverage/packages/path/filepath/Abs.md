# Independent review of `path/filepath.Abs`

**Verdict: FAITHFUL**

Evidence SHA-256: `63a7485a2c192b8b48b714b9fe8e50bacc47190e05ae67a7a290e8141d40e7ca`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `03b18da5b6c26d0a66fa5529f7caa37ca5fe622b7a3c1e5608026cf3b45b017e`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `310431eafda8ecf9d6c056534f56acce8f5a25e28aced9f3485af48c699992b9`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `2eb55d247f06713feae433e90380f661e59a48ca850ecb32fcc03e09e3b7d30d`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `1a424b7971f96d0900f2bc17846332a7f65efde03fcb12a3f3cb96941cf3ce73`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `b5633911faba26d883610ebc52be3fa7e62080e7ee87e3743708a47d8ebf8818`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `ad52cfdf89149a608919842e0ade72aa2db79b2547fd84bbae9d13b0510a68df`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `4e942fcb470a89a8d4bc615bcf4389afd1cf6b63f8d825d91de3b6b095ab1cd3`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `d5f3730ef5a165dc11ded612a9d774ed37d0a7167029fbe41ee0ca4512fe2570`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `93d869ee12e169f8f775c6c404fa697887d3050f0c790f002498c6813c8d5340`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `f5aa16f7ad263c433832be384e313c1a3adfca05f8e66d3b188828aa93b34f8c`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `8d01bbdb278258dfde139642f21fe5d259cb8c7f25484093f420c3a0f3db59c4`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `e5c41390a4f70f039ccc7e5bbb09b5fcaed43cba1b849cc3b9f727f11d2b7437`.

The public Go 1.27.0 declaration directly delegates to the selected platform routine. Unix and Plan 9 clean already absolute input without consulting Getwd. Relative input propagates a failed Getwd or joins its successful path with the original raw input. Plan 9 recognizes a leading device marker as absolute. Windows passes the uncleaned input to FullPath, substituting a dot for empty input, and cleans only the successful returned path. The model preserves operating-system failure values parametrically. [Model](../../../../lean/GoLean/Path/Filepath/Abs.lean) · [Comparison](../../../../difftest/path/filepath/abs/check.go)

Pinned darwin/arm64 matched 522 boundaries and 60,000 samples at seeds 17, 2718 and 65537. The author also ran 10 fixed profile observations. The independent control passed 15 Lean observations, including two native calls in a subprocess with a deleted working directory. The local Getwd succeeded in that probe, and the native relative and absolute results matched. Synthetic failures distinguish an unnecessary Getwd on absolute Unix input, bypassing FullPath for absolute Windows input, returning uncleaned Windows output, and sharing Unix absolute classification with Plan 9. Native Windows and Plan 9 and a native Getwd failure remain outside the observed evidence. [Reviewer control](../../../../difftest/path/filepath/abs/review/main.go)
