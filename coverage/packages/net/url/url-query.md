# Independent review of `net/url.URL.Query`

**Verdict: FAITHFUL**

Evidence SHA-256: `cac2d65d72569782911c6c3f541cb814bf5b1a86357b44159237165c9d3aef11`

Reviewed author commit `de11364` against pinned Go 1.27.0 `src/net/url/url.go:1164–1167`. Returns the successfully parsed portion of RawQuery, dropping the error. A nil receiver faults, while an over-limit RawQuery returns the empty nonnil map.

Eight fixed fixtures passed against pinned native Go and compiled Lean. Discarding parsed values whenever an error occurred changed two fixtures. Replay using `GOTOOLCHAIN=go1.27.0 go run ./difftest/net/url/query/review` from the repository root.

Author evidence passed 2,463 boundaries and 60,000 sampled cases at seeds 17, 2718, and 65537. The model takes the effective GODEBUG value and native signed width as inputs. Process setting races, external backing alias changes, and allocation failure remain outside scope.
