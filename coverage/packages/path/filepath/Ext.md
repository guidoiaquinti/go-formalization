# Independent review of `path/filepath.Ext`

**Verdict: FAITHFUL**

Evidence SHA-256: `dafc8415ce2266c6238a700228b386f2047e4f5262ee8399c582853dc129e7f4`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `1fefda81af0287fcec09d1e106ee7500c34c42bb90e49ec629d5a49f144ed607`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `1f0e385bff08670dc471c009bb03792aa4ed59ebadc2f1be04c7c673d1d96c73`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `0b0525e198b57aaa96d1bdd8b9068496eab5a237e9c8f21da46e0915d79a6a5a`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `6dae25771e3d5b03aa524054ba3135542a608cdb5ebe8cd15943671d26ce7e72`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `5ea1614758b4225482d8c89cbdec9291e81257de159f8332f9c1b9214c7fccc4`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `0cf3878bc2b7fb5e553410d03f62881e2e9a581885b5b14f0804735dbfacfee2`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `b278f1fc525d7732201ac74fbdd45a9021f1e76cd5f4969ae1dc93c392d726ee`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `a4bf5b32d5fcc9ac432d6efdb1ec38bd64414d697301546d48259edcd44926cc`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `95b77d5a79fcd3f2ebf3cefbc9f4f68ef6945e8bcb0b191cd400c1370f803b93`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `d5e975ffd676dafc60d5adbfc237fcad6e31b746291c9435867360c1814962e9`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `0fca3ced46f86622799c4a228b5761e5b2b5b5527b864a11e2ea888c9eec6fe7`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `7a0e9b4ebd94739024310d49af0b340432323c6ffea618c29cddf66eb55e4ddf`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `9c4360461a30deb85ac869149b45d7b27643edf9fd779f16e419b775c7e6d908`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `f9f2bcf6c96af28c7e144f866704d45debb6d81d40f8a889b3c55059a74e23f7`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `2f00b8e111905d716b78c2ce6987f59227b8bc1349d859535063e649c8f8ffb4`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `b693ea70f329d45e7c18d058df1b10783d08d76eca78557a544f992c34d56a81`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `3271199c074fe5d4c3e1ff7c49c4bcc59cf1439721d8b4e6679b5db6f82eb304`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `98b33b9521b10e91e042a9cca97dd682072d8246f3859bea22f32d9918509830`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `427e5420f514a2214defe92df7606d92b25405d8a381be817be74a4fa6b7059b`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `3900d72c4ee109d97a8ba5f4cab31018641d6d678ee9ad424faaabce149fe047`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `611025c825ad0b7c2f20dd7f2d976da5aea6b019c403d3435e2430a96a3ea7a8`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `3ab1882baaa2812f11c1edcb20d3cf4b3cd5a18962534ece81858e14a262a205`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `ebb11abbbc161e1830ce287b236cd5ef0043326bbe03b0433270405712069703`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `543c9928fef59ec9e40587f77d3fcace9b03f115d1722ac479ac4e38dad62614`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `51747306a7992325792dac6d1ddca2205a080a31148dcb4441745f3b6ee6d95b`.

The pinned wrapper delegates to the backward byte scan in `filepathlite.Ext`. The scan stops at the first recognized separator encountered from the end and returns the original suffix if it first encounters a dot. A leading dot in the final component is an extension. Slash stops the scan in both profiles, and a backslash stops it in the Windows profile. Lean proves the empty and trailing-separator cases and the trailing-dot case. [Pinned wrapper](https://github.com/golang/go/blob/go1.27.0/src/path/filepath/path.go#L138-L140) · [Pinned implementation](https://github.com/golang/go/blob/go1.27.0/src/internal/filepathlite/path.go#L215-L222) · [Lean model](../../../../lean/GoLean/Path/Filepath/Ext.lean)

Pinned Go and Lean passed 524 boundaries and 60,000 samples using seeds 17, 2718, and 65537. The public Go entry point ran on darwin/arm64. The Windows profile was checked by the compiled Lean model and a separate last-byte-index oracle. Native Windows execution is unverified locally. [Differential comparison](../../../../difftest/path/filepath/ext/check.go)

All eleven fixed native and Lean profile fixtures passed. Counting only slash as a separator in the Windows profile failed three. Searching for a dot across directory boundaries in the Unix profile failed three. Other fixtures cover leading dots, the final dot of a multi-dot component, malformed bytes, and NUL. [Reviewer control](../../../../difftest/path/filepath/ext/review/main.go)

```sh
make verify TARGET=path/filepath
GOTOOLCHAIN=go1.27.0 GOWORK=off CGO_ENABLED=0 go run difftest/path/filepath/ext/review/main.go
```
