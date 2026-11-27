# Independent review of `path/filepath.scanChunk`

**Verdict: FAITHFUL**

Evidence SHA-256: `36dafc0c04a1fcfa68259ebef8d603dd2aca396677a4d439fafbda4ece7debb1`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `44646c2588c09d937698b116c2d0405e4c2c9cc29282236ed70017fedd57590d`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `b429d50fd452cbee2a08cfc4f5b828212cb598ac3030b3de4b2454f2bcbc5831`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `d6e643ed97a047cea5cc6f938eb4e2c597c46b7c45c7b2e5b1a9319be1a1f143`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `369173f5823e2b91f5e6ccaec7dda9f577ae083e3e7c432ef685832df1e90926`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `3b40452f184e1fb1368fae52ba8d0a41aa09ef2b32af869864af2d99a27ea260`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `d9ccbc8d94c235f1f4c1d2b136ad4c94564222d92706b97fb1679e921e5f2f54`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `8cacbecad34e8616fef017df330a05d18cef713fb667e8f5b9a1a82a630b5a0a`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `9a12ab9f5334cf4612e02cbe60b4b66482d16d88566c59bf1be51a79a87755c7`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `98ab661c2854c29d026a2164b808fc58a6149bd931ed576d60419ea6b998aa66`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `3a81cd573758eca9774c5613e16f495325fca7699dc2aeb5eca54b5039267919`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `65fecf8fc296e749400c3dc8db218b6a62ad2c5e7d2b3fc6835d65868e7a3e40`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `b3fdbe6923d78c825ea50a8cf182aba6bbdc5fed14ef8bd6c3a5abf24201e97f`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `3079974ea8aecc4f2e0b540ab06ff9f797fc467fcfb4a53d25b61104148f2b36`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `075a45ae166a48663a7662acef2b1d78f9acafe43861e4f501791e2c144abc57`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `ca7d070ef7371d049652a09ed1831d21f6e484ab6a743434ad529d5d6893e945`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `157581e1a9bc377c488429407aa27d262d7430526370fdf20fa39c7447aee640`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `59ef2fbb5726a9609a304f1193c18dfc28b8db6e36957f901d2e22d79c498e7f`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `ef7148ff38167264dc8c75dca1eb906a7a757a9fdcb59400bb8875be64c3e43a`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `4fe6193c71806e73eb89d25a636fe4eb7605140e821f9483ced68e1482579bbe`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `5653c259d63b89ff9a234b1d1ff3a330c8628fa44a2d8ac9f2fe59ee2ac62a67`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `87bebbc2a84e2a1639ce55f1336116c7d695012c098cf7014e1385ec33eb1f2f`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `ffe6f638ce2f5362d88b2d97675ec44e90a06a4111dea16c7ed7bd50e9df44d8`.

The pinned scanner removes leading stars before scanning a chunk. It tracks a Boolean bracket state and returns at the first active star outside that state. A backslash skips exactly the next byte on non-Windows targets. On Windows it is ordinary for this scanner. Malformed ranges, trailing backslashes and malformed UTF-8 stay in the returned raw substrings, with validation deferred to the caller.

The pinned native Go and Lean comparison passed 524 boundaries and 60,000 sampled strings from seeds 17, 2718 and 65537. The independent control passed 13 direct darwin/arm64 cases and 26 modeled platform observations. Ignoring bracket state failed three fixtures, enabling Windows escapes failed three, and treating bracket openings as nested failed one. Native Windows execution remains unverified locally.

```sh
GOTOOLCHAIN=go1.27.0 GOWORK=off CGO_ENABLED=0 go run ./difftest/path/filepath/scan-chunk/review
```

[Model](../../../../lean/GoLean/Path/Filepath/ScanChunk.lean) · [Comparison](../../../../difftest/path/filepath/scan-chunk/check.go) · [Independent control](../../../../difftest/path/filepath/scan-chunk/review/main.go)
