# Independent review of `slices.SortStableFunc`

**Verdict: FAITHFUL**

Evidence SHA-256: `83d90dd2181f405b30f9d50dec087b9dd894913f8cf310f2482a5beb383952b2`

I checked evidence commit `68ee474` against Go 1.27.0. The public wrapper sorts the caller's slice with the pinned stable helper. The model retains the post-panic backing and exact comparison trace for the registered scripts.

I ran `go run ./difftest/slices/public-stable` and reproduced 3,287 boundaries and 20,000 cases each at seeds 17, 2718 and 65537. The independent `go run ./difftest/slices/final-trio/review` checked four literal stable-sort Lean rows and matching pinned native observations. The tagged [513,256,512] fixture must sort to [256,513,512], with 513 still before 512. Exchanging equal keys reverses that tie. A panic on the second comparison observes the same first completed exchange, which rules out rolling back to the original backing. The scope covers the tagged int64 corpus and fixed scripted comparators. [Independent control](../../../difftest/slices/final-trio/review/main.go)
