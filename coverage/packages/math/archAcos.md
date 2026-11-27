# Independent review of `math.archAcos`

**Verdict: FAITHFUL**

Evidence SHA-256: `ef36d24011409d51b8633cb0a7cec5210184c4332d1eb13d3ec23af4b4ed9227`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `ec370e2b8bba6978a1abe8dc9b3de5b4852be05928b05991fbeaa6f9338e1c33`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `dea52491d012c62b557e6962f9c35a4cd442d0179ac835f235aa8be6b85b34c7`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `bf787ce48dcf8842c1db2dd5106fdd5806bb9a1e0fc9075ef642a0e4e117502f`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `1eab7430eda7fcefb7aa41a9d911f1c3e6bc031fa6d0a622ea292b75f1ffdcf3`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `c39ccd8fa7c96b8125b38c12d3a5582b1b674ee029779891d578ef3e364b49c8`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `8b97b3fa910dbc3d500e72bef9fd3e9d2b13cd397b8e31120483933d87bd5c7a`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `f353296dcaaa1f2e0c1cd9a8dd6383266a9b0cf6589d2a1050288ce2e6b2aa59`.

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
