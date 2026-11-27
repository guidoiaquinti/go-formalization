# Independent review of `math.archLog2`

**Verdict: FAITHFUL**

Evidence SHA-256: `2ce27d4aaf77794e5b29fd4bc86e623cff1cb22e0ff1db48d70ade5bcc69d2b1`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `7299f286860652759355826a8d5e421cac9c1e19cc7b522399fc8ef9639072b2`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `dacc5478ffd1e03d71d4d6cc51fce98b572cdcea53e16a0a96bfe44b559f6f38`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `e4707836d6ce9a092f747f5e706077055771f4150a554b04d1eddce6468f43a9`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `8017d5d565d782b5d6db1cc9b4f93f18a95e03b0c44c59cb0393393ed7ec7100`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `5326fd6a58ba943270c551de424202ac42e836141c570bad9e2fe10656c075c2`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `72a60f5cf4e246e90faab255bd5ec0c5fcd5cf812f2b7ae2b4d3c10129c638f2`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `c036ceaa92958ef785441c020dfec379ba92eae4eaabd23f22088e258ec7715e`.

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
