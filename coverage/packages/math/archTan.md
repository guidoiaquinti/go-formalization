# Independent review of `math.archTan`

**Verdict: FAITHFUL**

Evidence SHA-256: `2ad1c8575a280fd5bf16842750d0e0903bccca6487e20d06ed75c7e615e221cb`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `6a1442a3a6ffbf682f4797c82b6936ece404408856db293c2228bd97d29b3f8a`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `880908df323ef9a8278e88963e75ce08168380e3dbe8a37dd224ac7bf4ec8ee1`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `2e3d35ce586a9854af645c8ce1299b4619977e29b7fa9c87993a0658b32bd913`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `85588ca3a78b9bede2d74b9ccd9508754bb74a8add0e50fa90c827a6b29dfe70`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `32d84bf9a317e52a9dda1db221813582d819cbd21db9db94cff1846c7cb1c344`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `9a077836e1965726635ffbd95e8d697a023629fc1464821c4377bc13c8df6d41`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `f45ae33d3b05e99514c96ba686c5dce95e2d9d78af10afb98f44ead5433d3f13`.

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
