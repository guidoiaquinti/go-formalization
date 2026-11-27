# Independent review of `time.dateToAbsDays`

**Verdict: FAITHFUL**

Reviewed author commit `7732011` against Go 1.27.0 `src/time/time.go`. The March epoch shifts January and February to the preceding absolute year, uses unsigned intermediate wrap and accepts an unrestricted native day offset. Independent literal dates were derived by Gregorian ordinal counts. Day-zero coercion failed one row.

Pinned Go and Lean matched 15,552 fixed rows and 60,000 sampled rows at seeds 17, 2718 and 65537. Six independent literal fixtures passed. Replay with `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/time/date-components/review` after building `godiff-time-date-components`.

The month helpers use civil months one through twelve, and normalization uses a positive base. Native replay used signed 64-bit machine fields. The direct date helper models the uint64 day word and its modular arithmetic.

Evidence SHA-256: `9242ec12c919fd46f2942b89c0604a71ceb8dc2a19c34f1bd0233f69bb1e43d9`
