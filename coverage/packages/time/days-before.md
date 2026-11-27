# Independent review of `time.daysBefore`

**Verdict: FAITHFUL**

Reviewed author commit `7732011` against Go 1.27.0 `src/time/time.go`. Returns the nonleap prefix count for each civil month. Omitting the post-February correction failed on three rows.

Pinned Go and Lean matched 15,552 fixed rows and 60,000 sampled rows at seeds 17, 2718 and 65537. Six independent literal fixtures passed. Replay with `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/time/date-components/review` after building `godiff-time-date-components`.

The month helpers use civil months one through twelve, and normalization uses a positive base. Native replay used signed 64-bit machine fields. The direct date helper models the uint64 day word and its modular arithmetic.

Evidence SHA-256: `a7ee8686f55f4516d19304bb8eaa56fa702f71772be2c27773585464232318ee`
