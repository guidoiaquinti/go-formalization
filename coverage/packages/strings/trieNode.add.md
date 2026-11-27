# Independent review of `strings.trieNode.add`

**Verdict: FAITHFUL**

Evidence SHA-256: `b1817d0e16f4374032ca56017a262acc299a998df4905982640666eaf1b2246f`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `d7a27fbcebbbc2d0248284783433c9f8489e883dfb835f057dc6c04dfc2e1680`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `ac372f5124a0099b5fde06f9a8dce619c09c02602032dc056dc41e4c510d8b8c`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `d2fc9aac4a33a6d5f5c553f7e921cb7c4a372373825bb7652126455dd0edb47f`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `9f54bd96e0b7c14b4ce0f8288195cb8c0f10253b6f39a5bfbdc8b3e5e6462d06`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `0cb43058602eab03d1689190aa97135df0155c872311860850cc4a255e1f2cbd`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `a6d7b9183cddb51e6934525e764a39f82f057807875e74be7a71a1bff85cc1e3`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `63864507710b31f19d30abdc2e0a0abf3a395ee89243b1dacc59e84632bf1649`.

Reviewed 11e8f0d against Go 1.27.0. Empty keys install a terminal only when it is absent. A matching compressed edge advances by the entire shared prefix; divergence changes the edge into a table or splits it at the common bytes. An inserted key that ends inside an edge creates a terminal at the split and preserves the old continuation. Expanding a leaf or inserting a missing table child preserves other terminals and branches. Every recursive descent consumes at least one input byte.

The independent control checks nine exact native and Lean snapshots. It covers duplicate priority even when the new priority is higher, a key ending inside a compressed edge, partial and zero-byte divergence, extension of an existing terminal, insertion into a populated table, and raw high and NUL bytes. Overwriting a duplicate or swallowing the new terminal at the split changes a fixed snapshot.

The full replay passed 4,842 boundaries and 60,000 seeded cases, with seeds 17, 2718 and 65537. The pinned native checker starts with a root table and a complete byte map. The claim covers valid compressed tries and positive native priorities. Nil nodes, malformed states and pointer identity remain outside scope.

Run the reviewer control from the repository root after building the Lean executable.

```sh
GOTOOLCHAIN=go1.27.0 go run difftest/strings/trie-node-add/review/main.go
```
