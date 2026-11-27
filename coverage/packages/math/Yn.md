# Independent review of `math.Yn`

**Verdict: FAITHFUL**

Evidence SHA-256: `cea8295dc8ae1514462a1f492544e58d5d5821bdd12a2faeef52f7652acedd0a`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `255066260395cddebba7b00d2c800ef7a64d7f393b1f8676aa32c001a79ff012`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `2a516ea40b5d0b63f8a798ea40b23cadc6faf4d3b36d787a18bf8690b6e18321`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `2538a59e8c585becf4032bfae75adc8e06e3629297f6730be12a07b880364d9c`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `ae96e8517a673534ef6159de6ea4a1090737b01fe534919d93d86b89b725264f`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `93a79cf7a323b42807592190797c7ca69949568512b967df811c90ea7467d439`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `1bb0d259c467faa6e2a70792148480511a711274827dc76ba8581b64e8b36766`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `a2947d764a23e6f46b0eec5f92fa8ddcd7c125c23d0a03b7cd82ef68d22bbb70`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `c9db723c4ffe2b7baea4b0e784deb3e80d8a54101c81389a0bd87bdcd9d0d373`.

I checked the pinned Go 1.27.0 special values, signed bounded orders,
delegation at orders zero and one, huge-argument phase selection, and
rounded forward recurrence. The claim covers orders -64 through 64.
Machine-integer behavior outside that range remains open.

The independent control matched 45 exact darwin/arm64 observations and
135 Lean observations across three profiles. It checked 96 observations
against a separately implemented convergent series and order recurrence,
and checked signed tiny-input overflow and the two adjacent words at the
order-one overflow boundary. Native amd64 remains unverified.

The author's 1,916 boundaries and 60,000 seeded cases passed with seeds
17, 2718 and 65537. The independent high-precision reference found ten
finite-versus-infinite rounding differences immediately at order-specific
overflow boundaries. Lean matched pinned Go exactly on those rows.
