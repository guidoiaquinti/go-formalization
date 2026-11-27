# Independent review of `math.archTanh`

**Verdict: FAITHFUL**

Evidence SHA-256: `8e8d2bc53807897bd3309309bd4914baa6507befbe45ac8cfc69ded2143c6301`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `f03aa28d800082fb7202848033d1ee4bb2551ce1c0305cbf887a29ff8ef79fdf`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `dcba5c5ac7d73e973c4c5eadcd6a1e2c5219a33637a30f3d72118ac348e6b17e`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `629cdcd769d99f70bee497fe283c55c835683a6315d5b3f4c67dbd06ebceb1bd`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `42977b5b6231458880c55f75a4d18c61f8824e98229b73fbd79b4099d01781ae`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `bdcf9fa569b01667d2aee946ffc754e081efd03cc2ba44ff02b56f22ad42abc8`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `09ea93468ed3aab994c39c31788a934bacc990473db822237ad2c3cce0434aef`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `64c6f1cc1dd3365a4d4a0c9d17e69ee1a00f55e1afcdeb070f8a185b74785dcc`.

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
