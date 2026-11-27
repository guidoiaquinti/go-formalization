# Independent review of `path/filepath.getEsc`

**Verdict: FAITHFUL**

Evidence SHA-256: `61fa898af3dc59cce7449faf0f1353adaaeb427dc126a00081493ac31f557350`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `9acab719ae5a9c21cfa59e67362f8517a5b70c318e4163fda4a0b7b711cc0602`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `2a74f0bf4fd352d29e0b6a2e037cf340a6160f3b76edbf73f7fe50ddf1dc6556`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `ff55e8e08b9b7e6bc57ebe83463afecdbe033a2c4b52b93c0427ac4624e52053`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `2963b6c2c7e805f017323ae641a9997d6d54212d1eb8ae25f239b9c7fc9be1e7`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `2e53a9db61393fe18225f2b0e301f24a28d14e7508d625a19e1f8edb5521d104`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `700e271ef7f209e24e0f76e69e157289fe30a2ddc7704cfc43f3175517ffd9eb`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `a5853d0b3900da3f202292e65a4b065c5a5065c32173f73a55e61c9fbfd4914e`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `63e6e20cca0a0700d7ba37fe56432c037bccd8dde211ae2584415d38e98c7b01`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `e6139bfe0bbf87581855c61aa70f9aa9046c45bbe6808f2e21189db52f413652`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `9d7ee20411931a9bfa293f12a9206cf95e9fc811a30af2621bb62148b272f758`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `de9c644f0e1d191b20f5c3ae90f23a32cfd16f45330c9eb9133687f2a00105e0`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `f36981ddda42e625e57b1910a70f3f44d074d7a058964e33743e385da4d998e4`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `1f463c46990aadf9ce4c6ebe1c0c72b76726506f74766911d9cd431af19f013b`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `0ee2f97f664dc4c33c83a5dead85f7948d5fee8898b539bbc74c587255d007fe`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `e7bce33952e89b1b5c997b7f49e0ab76afd079b3d706fb13b57814fc051d5dc0`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `a504a95ca723325c21279975f1b023543a7ac4bb40ad0114031f1f97c970a84e`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `93f2d9d2f71c59598ffd68d40852f43a57c5077610b9bda5a17e1c2e5df1afb4`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `82cbae9af58a2190e621dc3b3023e1457fc074309427197b198badf663ffa1b9`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `98c780dc310d92bf4ffcc6169d9ea335efbcb0be0dad94316f78f1856895dd87`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `095f4173b6b3bffaae2a81e98d5db3c37a142173e092b18dc67e3b93a65b245e`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `df4c758664a4c112963cdd934d687befbd19d660039360a56bd7159b3c0efd55`.

The pinned helper rejects empty chunks and leading unescaped hyphen or right bracket. It consumes a leading backslash as an escape outside Windows, then decodes one rune. Malformed UTF-8 and a missing remainder return the singleton `ErrBadPattern` while preserving the decoded rune and remainder. A literal encoded RuneError can succeed when bytes follow it. Lean proves the early exits, escape distinction, and error remainder behavior. [Pinned source](https://github.com/golang/go/blob/go1.27.0/src/path/filepath/match.go#L200-L221) · [Lean model](../../../../lean/GoLean/Path/Filepath/GetEsc.lean)

The pinned private overlay on darwin/arm64, separate profile oracle, and compiled Lean model agreed on 780 boundaries and 60,000 samples with seeds 17, 2718, and 65537. The Windows escape profile was simulated. Native Windows execution is unverified locally. [Differential comparison](../../../../difftest/path/filepath/get-esc/check.go)

Pinned Go passed eleven fixed fixtures, including truncated escapes, valid and malformed RuneError, and a valid rune followed by a raw invalid byte. Discarding rune and remainder on error failed three fixtures on darwin. Rejecting a literal RuneError failed one. Replacing an error with the same text and a different identity failed seven. [Reviewer control](../../../../difftest/path/filepath/get-esc/review/main.go)

```sh
make verify TARGET=path/filepath
GOTOOLCHAIN=go1.27.0 GOWORK=off CGO_ENABLED=0 go run difftest/path/filepath/get-esc/review/main.go
```
