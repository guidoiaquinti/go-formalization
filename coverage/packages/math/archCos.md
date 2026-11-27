# Independent review of `math.archCos`

**Verdict: FAITHFUL**

Evidence SHA-256: `31994be98fd729de6ee38d928187abf49ec1b2f1456793e455f1d62426178962`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `1f6e261395066860f36f9ab75142fe6f898798ba8fc0dc5fa0741878845bf439`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `93bfef1a5d75bff6c42598c1377677c75e201f08647e26799ffbc4d2cdbf5f13`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `3296e3f2854dd344535937d935f3113579481972df5ec85c77c0f290a1c5d006`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `b6cf09480bcac4a8b7339608f55d17905bd092fed7876f8e0ca9f568e27c1bad`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `f5ee563c29a192ff9c983cf9c7ce655f8afe6e3ff784fc9548949ad7b52b4144`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `1002cfad456813708909b7c5f877ef965d2459b44c1ad24e53d208b462279aee`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `6f6810403f665614b379aa6861ee2edbf370cfd3bb258f08b7f72648c2876534`.

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
