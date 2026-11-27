# Independent review of `math.archSin`

**Verdict: FAITHFUL**

Evidence SHA-256: `76bae9028de4c002b0331b298b53be787d52d711029868c0ff043597a96a3046`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `6499dc1ebd4e910fff3eab1609c83b31bcc5e1d5c7f2aaaa9b3c3ac995fe5836`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `2e205acd27768d923dc7bb34e2a1a4369a9bfd109cdf0b1ad67169ae776e0280`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `041ff55c1cd8177870dde768217a9107798bb9d4cb082c8eaa5fd1fb0cc93247`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `c33c223844aeed6dff3ae43a4b9f6580e446fb3f8547793265952dfa6ec362b6`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `d103cc0dc7c16318a89fdd9306364d36543bd7075831c83aa995dc1b08b75a5d`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `d90fb19ff94b0d8a5bf6576bd94ecdde6ea4461ef23b76f954058d46a6719aa2`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `af725e9500fa0954e22493246b82667b281a2922be989f603a7522d858fcfa54`.

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
