# Independent review of `net/url.parseQuery`

**Verdict: FAITHFUL**

Evidence SHA-256: `cc671ec96cf92bab5f9d25b372e907270542a7b52ad03127d97c4e98b207775c`

Reviewed author commit `de11364` against pinned Go 1.27.0 `src/net/url/url.go:979–1011`. The raw ampersand count plus one determines the limit before any mutation. Duplicate keys append to prior values in order. The first malformed escape is kept until a semicolon error overrides it. The fixtures check nil-map panic on a reached assignment.

Eight fixed fixtures passed against pinned native Go and compiled Lean. Counting only ampersands changed one fixture, and retaining an earlier escape over a semicolon changed one. Replay using `GOTOOLCHAIN=go1.27.0 go run ./difftest/net/url/query/review` from the repository root.

Author evidence passed 2,463 boundaries and 60,000 sampled cases at seeds 17, 2718, and 65537. The model takes the effective GODEBUG value and native signed width as inputs. Process setting races, external backing alias changes, and allocation failure remain outside scope.
