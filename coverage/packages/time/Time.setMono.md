# Independent review of `time.Time.setMono`

**Verdict: FAITHFUL**

Reviewed author commit `b70dd68` against Go 1.27.0 `src/time/time.go`. Already marked times replace only ext. At in-range extended seconds the source ORs the packed marker and seconds into wall, preserving raw spare bits. Out-of-range extended values leave the whole state unchanged. The test detects replacing the old wall bits on one literal case.

Pinned Go and Lean matched 4,860 fixed and 60,000 sampled rows at seeds 17, 2718 and 65537. Seven independent literal rows matched both. A location-clearing mutation failed on six rows. Native checks preserve the exact pointer, including nil and distinct locations with identical names and offsets. Replay with `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/time/clock-mutation/review` after building `godiff-time-clock-mutation`.

The scope covers raw wall, ext and argument words on nonnil receivers. Lean uses a location identity token.

Evidence SHA-256: `ebe371438ffb4cb9f1627263ca021631a3796a94a10a14b9557faade7ee740e6`
