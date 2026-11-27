# Independent review of `math.archCosh`

**Verdict: FAITHFUL**

Evidence SHA-256: `bdfd70b716b70c1a0fde4e120151aa2b248419c211329851f95a2045bb43a0ba`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `c3ad0227ea224ca7f4eb9c60cf04dc8ee4e57bf8ce6b50a71c281ae55c36e999`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `038ca5cba5deca402ec16d8f81343a4655e9062f89f0a9c8c42fcf8510f2318a`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `fcf650b39a2c313d47c216f6016794a392e6263dfd75792408d606d7442ac585`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `f345eb73725b97ed00617de2cbe51a4bf1274760d1f290bace29d0ed9978c075`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `519a30a43cdacc1de060d2fcefcee5d1c5b9731da2b83fd805c183f58b9f8da6`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `5b77e16366f51a6942f1376a11816213d87dc481c0c223d1786f6205b3d44839`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `ab7938b78b01e45fdc62aafcc7d3bd3dddc2a570b3f2605346e68d0094cfa83f`.

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
