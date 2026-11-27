# Independent review of `encoding/binary.littleEndian.AppendUint32`

**Verdict: FAITHFUL**

Evidence SHA-256: `986cc8aa1b87133830cd2bb2820b594a6f43139ee6795358a71f27de00e7bd28`

I reviewed evidence commit `9c2b71b` against the pinned Go 1.27.0 `encoding/binary/binary.go` body for this method. The source passes 4 bytes to Go's built-in append, in little-endian order, using shifted byte conversions of the unsigned input. The model obtains those same bytes from the checked endian writer and appends them to the visible prefix. Exactly 4 spare slots are enough for reuse. In that branch the model changes the same original slots and preserves the rest of the backing. With fewer slots, the model retains the old backing unchanged and emits the new visible contents. The specialized Lean theorems establish the reuse criterion and original-backing preservation in the detached branch.

I independently replayed 65,649 boundaries and 20,000 cases per seed 17, 2718, and 65537 against pinned native Go, the separate byte-position and copying oracle, and Lean. I also ran the [independent reviewer control](../../../../difftest/encoding/binary/endian-append/review/main.go) with asymmetric literal bytes, destination start at offset one, and capacity one below, exactly at, and one above the threshold. It checked exact contents and full old backing against Go and Lean and tested native returned-view alias writes. Native nil and nonnil empty destinations also produced the expected bytes. The literal layouts distinguish reversed order, and the adjacent capacity witnesses distinguish either one-slot threshold error.

The reviewed observations are visible contents, original backing bytes, and backing reuse for valid destination slices and feasible appends. Growth capacity and allocation failure behavior are outside the claim.
