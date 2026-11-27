# Independent review of `slices.Chunk`

**Verdict: FAITHFUL**

Evidence SHA-256: `eee4ed1396b3914dc383c65c903bb825b69314c7cb25fbf9fae9363f613c9cbe`

I checked evidence commit `68ee474` against Go 1.27.0. The size guard runs before returning the closure. Each invocation starts at zero, clips each yielded view to its remaining length and stops when the callback returns false. The model snapshots the contents passed to the callback, applies the modeled write to the shared backing, and starts the second invocation from that updated backing.

I ran `go run ./difftest/slices/chunk` and reproduced 5,402 boundaries and 20,000 cases each at seeds 17, 2718 and 65537. The independent `go run ./difftest/slices/final-trio/review` checked three literal Lean rows and matching pinned Go observations. Writing 44 through the first chunk and stopping must produce a restarted first chunk of [44,6]. A mutation that resumes from the second chunk loses this observation. Native address and append checks distinguish the original element alias from new append storage, and verify the untouched next cell and spare storage. The scope covers finite sequential callbacks and the registered first-cell mutation. [Independent control](../../../difftest/slices/final-trio/review/main.go)
