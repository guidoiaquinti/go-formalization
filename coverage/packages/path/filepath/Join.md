# Independent review of `path/filepath.Join`

**Verdict: FAITHFUL**

Evidence SHA-256: `7656a24094ddc0ffa1e403e808aefe7d3724902f9729956874ad96428f75315c`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `442b3845bac005c0588103fff2efcc66ec1bb2ae0eab14121b867858f1a9dd1e`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `8ccd4aceebcc516e9d1395fbdd1dc087348c475b3863150bb48540477a2d085d`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `3468c67ca4692f680557364b848c9136eb85c983fb23da6ae60a0cb3c990656f`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `b36fa54ec68b0a0ee17d09c1bd304c7705a7ccd93d8a09ef0431af38d59f5f45`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `d8e83351e55deb3a1ff7035e5aa8c40935b981bc6e80f294bb2d348a426cfa6b`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `9933e1c9964c2d798a3ffa535328889119b828169ffe57a550df99a4df3ecb68`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `48c4eefdc1c4140ec9dc968777a31a05908497e157ad41ea2e177c3d20bcb327`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `7247eab41bc80a555ef0e6a59ad041c7654f7fa1cc55956ccdccab8d12e70e00`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `081674eb36451863e25b83ef6234020bf1ec66319e2256c8ee4b1f4372529b43`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `af2800bb9928eb6f3f4c08e0d7e0c5b41256d53efa9936acdf52077920304abc`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `56a810992f1d768571035e53cd23817314cd2ad9620fd5c618ac49c8044221ff`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `8b0c400ed53d8feb1afcd6a3d34b50ec7f628b67897174f8dfbdfdfcfe448fed`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `6f1bf11cf6443cc89267a26fb85a2d8c7c8c9edc3f2de5a38cbf4290927605b8`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `1528b8032c6e8c1faec004c9bea626ff577af581bf81b4eaabfa6f6407e09869`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `7ed8d295a8fd08f50117f7dedb133a5dcc7671e1add1ebb623a78095d1ee0dbf`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `da83ce97b2d54693a79c6f50987b69df7878dede11f7bff7949870117a90a10b`.

Reviewed 1ee46b3 against Go 1.27.0. The non-Windows source discards leading empty elements, joins the remaining strings with slash, and cleans only a nonempty joined path. The Windows source preserves the first nonempty raw element, strips incoming separators when the accumulated path already ends in a separator, leaves drive-relative elements adjacent to a colon, and adds a guard for a root followed by a complete double-question-mark component. The final Clean step retains the pinned guard behavior triggered by allocated backing bytes.

The independent runnable control checked 13 pinned native calls and 26 model rows over the two profiles. It covers empty inputs, colon-relative and colon-rooted paths, UNC construction from an existing UNC prefix, duplicate separator removal, root device boundaries, inherited stale-colon behavior, and malformed raw bytes. An unconditional separator after the drive fails `Join("C:", "leaf")`. Omitting the root device guard fails a root before the literal double-question-mark component.

I also compared the Windows profile with mechanically extracted Go 1.27.0 Windows Join and Clean helpers using portable separator and slash-conversion functions. All 60,012 comparisons passed, including 12 fixed cases and 20,000 generated lists per seed for 17, 2718 and 65537. The author comparison separately passed 528 native boundaries, 60,000 seeded native samples, and 16 Windows model fixtures. Native Windows execution remains outside the verified scope. The claim concerns lexical output bytes and excludes filesystem and allocation behavior.

Run the reviewer control from the repository root after building the Lean executable.

```sh
GOTOOLCHAIN=go1.27.0 go run difftest/path/filepath/join/review/main.go
```
