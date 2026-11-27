# Independent review of `math.pone`

**Verdict: FAITHFUL**

Evidence SHA-256: `e9fc19fb2ce60019e2770047db0c9b14efd067c66d22a23cee6d0ebfa49f98c5`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `ecc2ba4d32094a50fa95168cd5761349ef0ee64068451403bafd678e7897764b`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `f0579727541c77d09bce1f6b43fd330926b196b5b9bbfe17d12ee0cb5a388ee2`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `1ab965a08996afcf7a30664697cdd0e79a56d4202652d2d31e3f26139ad1b511`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `23b9a1ea14a29af94c36bf397d4493e52ef318b909ef0ce1d69216036352a547`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `b3739ac6bd02836d9a564bb9d21d5625e45d81ed57387423237e94afda14289f`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `48c1a4d501a21b39ce8169f1332b290636ab473c4f9201f1b068c3693c349bbf`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `191343339041075fe183548f5f6fdc2293d73487996dd54aa2d406ef10f3a6dc`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `4ba446e54c6fd3ed61031a49808c0fe4ef45d2e8bacb540193f06b51271d07ae`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `a513208ea7d0b7429c200cdb828840af9d1d7b7700838bc8af7675635d20f85e`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `08d5fa668f71ec291653d6355bd2ebc4a6a20404394f5d663e78269db2a00bc2`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `9b6f9bc539e0ba0be647fbcda7311c14756b9440d0de10cb21f9189d474a683d`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `1444c5d3bba671127378807f6a1b25dd9611102462240027df632af73f9832c6`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `b1130558824b636bf17863fc6c8fbc7a7807d826c9cad6ffd7ad288169dbb90c`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `165eba2e51a475b1b5658d34bd6882ac4ccf0358cf1366ce0c06f9de25fecfb7`.

I matched the model to pinned Go 1.27.0's private `pone` branch ordering,
decimal coefficient tables, reciprocal square, leading denominator
constant, rounded ratio and final addition. The panic class covers
negative values, values below two, and NaN. Positive infinity returns one.

The independent control evaluated the pinned decimal coefficients with
explicit rounding and contraction instructions. It passed 44 exact native
darwin/arm64 rows, including eight panics, and 132 Lean observations across
three profiles. Omitting denominator contraction, omitting Horner
contraction, and selecting the wrong coefficient band each changed exact
results on 8, 8 and 4 cases respectively. Native amd64 execution remains
open.

The author's direct comparison passed 66 boundaries and 60,000 seeded
observations using 17, 2718 and 65537, including 10,475 panics. Its independent
numerical reference differed by at most one binary64 step.
