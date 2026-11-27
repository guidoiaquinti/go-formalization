# Independent review of `slices.SortedStableFunc`

**Verdict: FAITHFUL**

Evidence SHA-256: `4a029ab723cf77aae99f4ffb2d732ba141c6c2bbacb272df085b35247a8ade99`

I checked evidence commit `68ee474` against Go 1.27.0. The public wrapper collects the finite producer before sorting. The model checks producer failure before invoking the comparator, returns nil for empty collection, and preserves stable order with a consistent comparator.

I ran `go run ./difftest/slices/public-stable` and reproduced 3,287 boundaries and 20,000 cases each at seeds 17, 2718 and 65537. The independent `go run ./difftest/slices/final-trio/review` checked four literal stable-sort Lean rows and matching pinned native observations. The comparator checks the completed producer acknowledgement count on each call. A producer panic after two yields must escape without any comparison. Tagged equal keys keep their original order, and an empty producer returns nil without calling a nil comparator. The scope covers finite producer and comparator scripts and excludes allocation and concurrent or reentrant use. [Independent control](../../../difftest/slices/final-trio/review/main.go)
