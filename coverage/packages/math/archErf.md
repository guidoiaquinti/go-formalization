# Independent review of `math.archErf`

**Verdict: FAITHFUL**

Evidence SHA-256: `8db91df994f4b883f8933e90614168a77213775ffb12dcd6a826b415eeece51d`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `1fe41439b134e5452d4b7eda5c76f3d7614a2821a85048eacfadb9c19c90fbec`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `d3ef82d99687ca028ccb5ad92268657a405b78b680bf344214310242380c7f7d`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `4fc1282d659eacbbe1dc73a29cd5bc06e651dcdde17073ebd4e09deeea6e03d7`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `5a6221b5bed337c9830b4204f9544259b66ca001bf6fb1fc2f61de9ddf735b10`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `5d498457e223fa79460d46b538e73ab3f0b2cac142b505b6bc2b6c0fc51a2f81`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `c7072a92e27ef0478115569f37d17b66d5d00de04bdb0937cdd74c0b2a8f7610`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `d1124a371d0dff38c7d00d2f6e553f4de5a05339c3a89f33db318e9d24d2eae8`.

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
