# Independent review of `path/filepath.SplitList`

**Verdict: FAITHFUL**

Evidence SHA-256: `be5953502c271fa3c3a2fc3419ba9d5c2e8dece5e40e9e52d232bd72c316eed1`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `c76d4e46d86266fce399a1c5f5beb2d7ef6ff7d6cdf3c0b0a08ea8767a8fa742`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `e618d8b8131cdbb5ec5ffdbd56d1776368bf781821a74614c08eef0242d7d26b`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `1b706e94d47b4e9a06792bcbb9ed5b96cb96b7101642fdf4dfb037c0f2880801`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `74d7b36dc7272f27dc1db2e4e3f5d31e2697104a2035cfd0d9a9d92af2b2a721`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `6de8c18c77c772c6c3808cbc295843c6f4dd83a1e11bfa7e9b6981a6335b18fa`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `9dd85689a713fa04dbf8788af057494c9b60ca8688d7958b82931c605072295b`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `cf4b11f9a75d63717a538be9533d3c80656831bdbf19345f6e2f3010122e1658`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `871e8a7729ae62ee7a46f5521bca59e83dbcd2bd81c768c260dcbe30efd071b6`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `b7f722ee52d8741d9db46d194bec10a1a124dcc97d252e280fdeb462a1c0c69c`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `18209b7446563a9efeb46049fa4c7443177cae1af442cfaec2387da62a1819fa`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `257fded419eefbfb1745effd44772fdd736c43619c8789fd251e6a6ab1a0b949`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `31a666e3c4bc594cff555ae27d20fea7cd1e93acd27ec702021c09c0d6cce04a`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `a29e2dcf2ed82ee6f24793e26b29f5b9d2afd7d2681830061471ce945f6caa40`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `5f810da22e3372c060879344ac2285841396da4d48c4b95c82b34b41c2ef7f87`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `fcfc713ee7101006ea834416abc31c48d5276bc7e56f2cb4add236979b5ef030`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `f11af362b6f4546e965735087d18ebf05b0fe9e1c8cbe842e19b3990cd995fc9`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `a3f12e234e7c9ce88c212224a0dda09e288d668cb4db6ef4961740f42c4cdcd6`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `b20767fab6d1dc6c1de48cacc96fe129aebe8449ad1e9188aeb8377dbc809d5e`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `ec006dbd36b28223985831858be595cc947de0e9ffb356a65ab517940ab7ec00`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `ca4cb496efac6f83709163e66bd3a63537ae7b1af224099c59bb983458c30ad2`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `072fb510fb8f9199711fb2069298a8b0e1078ae11deb51e34f9dbfa41f6aec4f`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `8f1e1ac246406b5301963b9e324a795e3a45a20d8f8a347d1317b5fe6206894d`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `7cfa469d6cf7c936f3e1bb9da4f6c1ff3febe047f92264aec4d838322158b3cd`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `ee5de8c0ed5a9dafdd6434f5b5b8287d7d3cb0dfb69d727cb3bc62227e3bb6cd`.

Pinned Go 1.27.0 returns a nonnil empty list for empty input. The
Unix and Plan 9 sources otherwise split on the configured colon or
NUL delimiter and retain empty fields. Windows toggles quote parity
at each double quote, splits only at unquoted semicolons and removes
quotes from every field. The Lean parser combines splitting and
removal without losing their order because removing a quote has no
other byte effects. Unmatched quotes continue to suppress delimiters
through the end. Backslashes do not escape quotes.

The author comparison checks 524 boundary inputs and 60,000 samples
with seeds 17, 2718 and 65537. The independent control checked eleven
host Go and three-profile Lean fixtures. It asserts nonnil empty
output, both quote parities, unmatched and doubled quotes, leading
and trailing NUL, malformed UTF-8, and delimiter bytes found in
other profiles. Ignoring Windows quote parity failed four fixtures.
Using colon for Plan 9 failed five. Native tests executed on
darwin/arm64; Windows and Plan 9 were checked by source inspection,
the explicit expected fields and Lean, with native runs open.

Run the installed control from the repository root with
`GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/path/filepath/split-list/review`
after building the model executable.
