# Independent review of `math.archFrexp`

**Verdict: FAITHFUL**

Evidence SHA-256: `4447fccc9b21f0ce9aee13d22851ea17ec428b46a84085a2bf37b0b7f5d0e1e7`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `579850f2279cdc363f50d075679a55a8fcb4bc5207a4784f625119d671e37c47`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `e37e58802ea854bb7d23d4be1f5922f6dcde4c074c38eb785349f5272574093d`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `e36ce1b9bfddd7b40ee669aacdfd6b20f76dfc030e2b9fb188fe7de98c98a62f`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `59c1dede7da66d4a5401aae0f29a907fe3b67f22b6b3162da3d680a120587206`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `b708775aa49678cc744d3439a04e6fafe161fd89cab9d80e8deeae29bbfbb8ce`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `4e244fd2265d643bb3b74a9c2406c1390df62595462797298df12690ae495086`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `6ec1f5383ab990b9618b25f71bb73c5f2fde58dabaa96b8e8b75eb54df5cf2fa`.

I checked the pinned Go 1.27.0 declaration, its availability constant and
the non-s390x build constraint. The function body unconditionally panics
with the Go string `not implemented`. The shared independent control
parses the actual pinned source and checks five operand rows plus the
excluded s390x profile for each of the 25 declarations.

The author invoked the exact private symbol with a same-package Go
bridge and checked the panic value and availability flag on darwin/arm64.
Each declaration received operand boundary cases and 60,000 seeded
cases using 17, 2718 and 65537. All 1,503,048 comparisons passed across
the batch. Native non-arm64 ports and the s390x implementations remain
outside the reviewed observations.
