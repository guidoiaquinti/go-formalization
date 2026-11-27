# Independent review of `time.isLeap`

**Verdict: FAITHFUL**

Reviewed author commit `7732011` against Go 1.27.0 `src/time/time.go`. The pinned masks after the divisibility-by-25 test give the Gregorian rule on signed years. The fixtures include negative and positive centuries. Using divisibility by four alone failed one row.

Pinned Go and Lean matched 15,552 fixed rows and 60,000 sampled rows at seeds 17, 2718 and 65537. Six independent literal fixtures passed. Replay with `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/time/date-components/review` after building `godiff-time-date-components`.

The month helpers use civil months one through twelve, and normalization uses a positive base. Native replay used signed 64-bit machine fields. The direct date helper models the uint64 day word and its modular arithmetic.

Evidence SHA-256: `78e0132756c514d3838762d344efe6d34dff50e200ee4f72d4789ee5f6e02866`
