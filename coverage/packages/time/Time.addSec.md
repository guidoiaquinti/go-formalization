# Independent review of `time.Time.addSec`

**Verdict: FAITHFUL**

Reviewed author commit `b70dd68` against Go 1.27.0 `src/time/time.go`. Packed addition keeps the marker and original ext while the candidate lies in range, otherwise it strips the marker and recovers extended wall seconds before adding. The two negative boundary rows distinguish exactly representable MinInt64 from overflow, which saturates at -MaxInt64. Replacing either result or retaining the stripped marker fails one literal case each.

Pinned Go and Lean matched 4,860 fixed and 60,000 sampled rows at seeds 17, 2718 and 65537. Seven independent literal rows matched both. A location-clearing mutation failed on six rows. Native checks preserve the exact pointer, including nil and distinct locations with identical names and offsets. Replay with `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/time/clock-mutation/review` after building `godiff-time-clock-mutation`.

The scope covers raw wall, ext and argument words on nonnil receivers. Lean uses a location identity token.

Evidence SHA-256: `d0e6e75fdacfa469e6a17dd82e1503d8990aea1c13fda068bcd2fe8aef7f1310`
