# Independent review of `path/filepath.walkDir`

**Verdict: FAITHFUL**

Evidence SHA-256: `dcb5ff134ebc75ea2bfe66df8f52eea56085c79229fa55023b5e9b68e9f44e08`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `063962de7f9068934e254eac7cc8bc84e8cf76e097cdbd78c4e6137e68a94efb`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `d41c28303f72f02af648f317add97cd08304f2ae0cafc22e13bbb33d46f82377`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `3b2cbf7ac3ba3745347f55a4c65f46b901a4fb2cb884c311774c3505242a29c2`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `75dbf5aea21d456cba922d125a0c3154bcfafd03f7de326bdb5a90b472ddceff`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `30af23d6bb03fe6d6b6fa8b10169cc92ed4a6617954148360e2d10f53552f1a4`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `bf00c70e3ea854340f147afb432f5687df72ac6658213f7a5e6556bb7b611300`.

The pinned Go 1.27.0 routine reports a directory before reading children. Directory SkipDir becomes a successful prune, while file SkipDir is returned to its parent and stops that parent's remaining siblings. SkipAll and ordinary errors propagate by identity, including errors that wrap SkipDir. An unsuccessful ReadDir causes a second callback with the same directory entry and any supplied partial children are visited only if that callback returns nil. The model represents the ordered filesystem observations and preserves these decisions under the supplied depth bound. [Model](../../../../lean/GoLean/Path/Filepath/WalkDir.lean) · [Comparison](../../../../difftest/path/filepath/walk-dir-private/check.go)

The pinned native overlay matched 4,096 exhaustive callback tables and 60,000 sampled tables on the fixed author tree at seeds 17, 2718 and 65537. Four additional author model rows exercise partial children after directory errors. The independent control passed nine Lean fixtures and five native cases on a separate tree. It induced a real ReadDir failure by moving a directory during the first callback and restored it during the second callback, confirming the pinned implementation did not rescan. It also checked sibling continuation, SkipAll, and wrapped sentinel identity. Native partial reads, native Windows and Plan 9 are outside local evidence. [Reviewer control](../../../../difftest/path/filepath/walk-dir-private/review/main.go)
