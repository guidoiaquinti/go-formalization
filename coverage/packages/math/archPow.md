# Independent review of `math.archPow`

**Verdict: FAITHFUL**

Evidence SHA-256: `a3c0a8ee5107c85455ecfaa0f1c9884d1157a5b64203b50fa9fdaa0da15b674b`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `5a3b4e0b9e97dd336b0f28bfd5f1fe7efc2cdb3c62096b8a1574602c01483a02`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `eabc895c44b1a5b3ced95631e28cfe71085663b2a94c694f2ac6682d84ade4d3`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `f7d0a56cadb424fa0ce215ea44a4b3c7affc9b381fb6abbf9630689793be73e9`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `a85f0ac577fdfff64263b85d0b031fa4f6b963debd3d04b555030385fac7aa1b`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `f06b009c0c115df24a417b754adba086f02fe5b8b9baf197d7658fe842952289`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `6fe8b2e9aa7030d73b303e24fb0624713e6105557c022055e6f6422466e105e9`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `1bbf0a872ccd247d084b0f5bd6ae5225a91a462c62b902bf55aa27d848f14d47`.

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
