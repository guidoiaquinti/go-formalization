# Independent review of `math.archAtanh`

**Verdict: FAITHFUL**

Evidence SHA-256: `a8ca6698c66c3c942b4867dd15b8dff92a91128e3f0f0e6a960a64b4230bf6eb`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `ff8cd98f1d4a4252fec860ea6e8c48e3894d9fb5711cb3ebe9628c668db7ed0d`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `50f4cd10a85342ffc1d3ae39aa476acdaea1ca71a3c7382cb6a14c95d408292b`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `f6cb275996e37c2a0bd1cedda46f257e531f36c3f9b6a3c258673ede42ce2ab0`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `4cbc14a0e3e4920d6e702508c399b9467135b8c5f674bb9331d14c26114d11e5`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `a030a4d9489acfd2ba34de53f86ac8780062bd9b35d8a9414eca59a52b4b8901`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `1fcce4a071ff539f0c0c9b482ec1006431de868e2860535bcd6d36a669663679`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `13f8c414badd950f942610a10a7200d839c678738cef98c3ac114f4a46b5150c`.

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
