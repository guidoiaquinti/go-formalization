# Independent review of `math.archAsinh`

**Verdict: FAITHFUL**

Evidence SHA-256: `942385d087e79b938a32d6ed0bc573fe565485551474d00f2a1dd35e37c96a36`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `7ec754e26b01a88225c452215735898529d269d9e81d9996c059ec650069fdaf`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `b2af1b64c88a3221bd1b7c07656ac7959ad820acb80d5ca57bc513de77a0e672`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `6d1a69b6bc44cb009be12a65f3cb336dbec52863a7354c65175ee382f773ec9a`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `f5b70ff9cf675b1a1eeaec53fe118659346a5de2a2064cfc66d16f368f522733`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `bd23ebc4d514cfb39f78a2ec518a62cd1530136855af52b1410ba75140eb0222`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `8ca25d475e1032f283f00cc8d997f6be5b078665bad4bac6c55a1b23b76c2620`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `56a750e45b456a71b13c2eb759272e714494945844e9ea36986434da6d5456dc`.

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
