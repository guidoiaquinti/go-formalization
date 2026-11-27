# Independent review of `net/url.splitHostPort`

**Verdict: FAITHFUL**

Evidence SHA-256: `a082ffea065d4454c079b000902ed1d6cb580bb0311fea2e0d9c35710f791bc3`

Uses the last colon with a valid suffix. The controls distinguish invalid suffix retention, earlier versus last delimiter, and stripping brackets after a successful split. Pinned native Go 1.27.0 and Lean agreed on 126,608 author inputs, including 60,000 samples at seeds 17, 2718 and 65537. [Model](../../../lean/GoLean/Net/Url/Lexical.lean) · [Comparison](../../../difftest/net/url/lexical/check.go)

The independent replay, `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/net/url/lexical/review`, passed 11 native and Lean fixtures and rejected 12 altered outcomes across the group. The scope uses feasible raw byte strings and excludes allocation and concurrency. [Reviewer control](../../../difftest/net/url/lexical/review/main.go)
