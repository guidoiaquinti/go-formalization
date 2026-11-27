# Independent review of `path/filepath.Glob`

**Verdict: FAITHFUL**

Evidence SHA-256: `3f4c8dc94ceb71154d9b7fa9d9165b42fe8a4f243edd3b7ae3037b6e3b9100b9`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `0719d897db4b9c756fbdf7d864218689db9980ba6299d0e9278ecd85420c924c`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `d64b5a1db841f8a40fd2459ab9965e89bcd351e1773af9d907f3316d1dea27f4`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `9144b556df4f12df6a9d8744fbe4d66fe75998b7017111662cd657e632733dfd`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `36cc14b63da77e04574f0632e11129431aa61d78a3a39e067be9d5c681405f38`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `744a5007d089ed982c6a12cae9c62445da698c3f05e04f774fc6a684dc5cc255`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `8751d9fe4a6edf0de4120f14c2109cd087e1fe49cb4eaba3ea6b6d6147292593`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `4fa00301316d7edc856d64766d8694737a93051e8d2e436a75418f1f73909534`.

I traced the pinned Go 1.27.0 public call through its private depth-zero invocation. The driver supplies depth zero to the source-shaped Lean recursion. The literal theorem applies when the pattern is valid, contains no metacharacters, and Lstat succeeds. The fixed fixture adapter preserves native lookup order for intermediate directories even when a subsequent parent component would erase them lexically. The model takes filesystem observations as inputs.

The author comparison passed 276 boundary observations and 60,000 seeded observations on darwin/arm64 using seeds 17, 2718 and 65537. The separate control passed 17 literal expected outcomes through the public Go function and Lean executable. It covers two wildcard directory levels, sorted accumulation, literal raw path spelling, nil on failed lookups, nonexistent and nondirectory intermediate components, a traversable parent path, and exact ErrBadPattern. Returning cleaned literal input, dropping matches requiring two wildcard directory levels, or returning a nonnil empty slice would fail two, two, and eight controls, respectively. I also replayed the previous directory adapter that lacked intermediate traversal validation. It failed both wildcard lookup controls with nonexistent or nondirectory intermediate components while preserving the reachable case.

The evidence is scoped to supplied filesystem outcomes and a fixed native tree. Native Windows and Plan 9 execution remain unverified.
