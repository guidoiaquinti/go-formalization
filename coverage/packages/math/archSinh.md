# Independent review of `math.archSinh`

**Verdict: FAITHFUL**

Evidence SHA-256: `1982f4bd6ddc05246c629fef624a922ecc497daba7b4f1b10db55583f0e34151`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `1f68b2dd868d6ccadd6e027510b86e0f9984f73a6ba8bdc03c1cb788104b393c`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `c5e1dbb6654e06a19a03aa70e68fff5f3251dc9aecfdb608747ac874c0641496`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `23708c4737d37620840ddd16b55d1ff8e9585c0453bcaa3b6480a7a5bb2fc4ce`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `bbc0d4ad115cfbd4e4399cf2342e85823efd7fed486be699020da28673f3812d`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `9d35b6a2bbc109fca4efc6156454a1f99b47514272f57919dabbb0890eb0ee97`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `0f23e2232f20808060b3fc89b3c9986439ad0ac3be30fe2a79695522ece9aef8`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `e5e6871d8086fa0e8c52f91561cead6dd750d6b12a87182b7fad9767ca12c778`.

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
