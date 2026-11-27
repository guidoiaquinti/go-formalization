# Independent review of `math.archCbrt`

**Verdict: FAITHFUL**

Evidence SHA-256: `5be3f4b7abbc4a4ccaebf7e86de18f7b180bd0fbeab8c70c1663854321e9c14a`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `90db0b85d66bbbdb52f57c7c9101d985e2ed25d6d944a0054cbbd0009e5afc20`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `c3f189e874beff63ba21b646a65937544af61d7de0f868418f9935a0b062db5b`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `3e0634188f5732f37a3454f0dfe8e450187584f09ab5da7d797abd0ca54a971f`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `70499e0cc6fd6439a070f7143c3b688bdbcbf182e29a08bb97fbb2c5b9076924`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `f4bba0d3e0f989be50c5f9a76a70b8a744974cc8094314d51a98329c0dd3a45a`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `81ea8ad8bab330818ca45e609b375183d319b671120be93825199e8f3cb2ec14`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `b00b554d4aeeec6de1c649f0df5ce46b0eb5b2881c87a4668f33a9fb3b5a7bdf`.

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
