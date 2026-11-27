# Independent review of `path/filepath.unixAbs`

**Verdict: FAITHFUL**

Evidence SHA-256: `121a46093f2669aa4884887f57e6d7163718f57a135cca65bd2a246194693c9c`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `94189486db6094be412e7f091deebc3b53b7dcf8e6f418e2d67ee3482aeec4fa`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `a2064701c007ca943f363c5c40ee2896b9f89ca765e527300b98ef8fba386322`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `39974a2dd2ceed99ab9f4bd61c0e453c34aa7775494355857f85041e208de52e`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `0d008a859af65b6c14d82be108b463dd0a402b367d97e86f5a95cda66a579cc1`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `a5678a76e439753394c1d28c576b5c29dfeb5dff7e2e70cb1602dffccd6c0db8`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `02ead49a129ff89f63d98198985420f21f3ac0ac5c12dc92f0a982cd9cc47c54`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `4567d17fd1afd463bf82e74102055ae1abe3993a883f21535b5ca1781862eabf`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `fe93a96e3a7a5d876f2e1940c62824f562fa90be6ecd96ba4a8ef02533012b91`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `b02dd7795127e0253cf91001613e93fe407938b140748c5b3adfd1a92f3a1192`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `1d2f4cd7a85267ab4a8894f0c6a9a5f93f9af6c6c94e6105eb62d8675c2840fc`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `88669992264fe1fdc784311b7fb0d381054073fb7e04acfa647e37a1ab599e52`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `f0fd23a8a5f4b08e002dc10e37037cc13dd0db81e57195c94388d0646b7f3857`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `36381692933be4156f9fc84f446b96535c85356923e4b939189e1753cead96a6`.

I traced Go 1.27.0 `src/path/filepath/path.go` through the private helper. An absolute path returns the platform's cleaned path immediately. Relative paths obtain the working directory and either join the two inputs or return the exact error with an empty path. The Lean representation supplies the working-directory outcome as a value and branches before evaluating it. The Windows profile is a model of the private helper with Windows path primitives, while the public Windows absolute-path API uses a separate platform body.

The author checked 524 native boundary paths and 60,000 sampled paths using seeds 17, 2718 and 65537, plus twelve fixed Lean profiles and error observations. Six further literal Lean fixtures verified raw malformed bytes, exact error bytes containing NUL, joining, Windows volume-relative input, and Plan 9 device paths. A Darwin-only bounded test subprocess entered 360 nested temporary path components and cleared `PWD` to cause a real `os.Getwd` failure. The pinned private helper returned the same error type and value for empty and malformed relative inputs, and cleaned a malformed absolute input successfully. Removing the working directory in a separate probe did not cause `Getwd` to fail on this host, so the control uses the depth limit.

Working-directory resolution is supplied to the Lean model. Native Windows, Plan 9 and other systems were not exercised by this review.
