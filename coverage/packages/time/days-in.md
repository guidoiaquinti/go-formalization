# Independent review of `time.daysIn`

**Verdict: FAITHFUL**

Reviewed author commit `7732011` against Go 1.27.0 `src/time/time.go`. February follows the Gregorian leap rule. Other months use an alternation inverted starting in August. Returning 28 for every February failed on one row.

Pinned Go and Lean matched 15,552 fixed rows and 60,000 sampled rows at seeds 17, 2718 and 65537. Six independent literal fixtures passed. Replay with `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/time/date-components/review` after building `godiff-time-date-components`.

The month helpers use civil months one through twelve, and normalization uses a positive base. Native replay used signed 64-bit machine fields. The direct date helper models the uint64 day word and its modular arithmetic.

Evidence SHA-256: `067ce02f812a1c525dcc916187fb82790815ad0d507b9a65493bf045b43b4251`
