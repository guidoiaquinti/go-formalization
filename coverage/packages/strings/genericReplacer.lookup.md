# Independent review of `strings.genericReplacer.lookup`

**Verdict: FAITHFUL**

Evidence SHA-256: `868a032298d6821cc5a6ef41af2fdedd7b26fc8c08f3f7e3719e9a63184a14f8`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `1ae1ea41103875609f45be6b59301d1132d91923988ae35e1c8d860a2ca711e1`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `16ddae02b4e6e4ef362e64bf73b7562407cb4da620f4ae34d1442919ab2c37dc`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `b8f13b42b8be249f7215673f6f897876d0b65e0ade361a6a62ff4cae19ec1584`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `4e3b6c4f217d66a84583d0438acaf517184e825b5fd2656e4105ce4df87ebaec`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `c028d2e8f32ae4093ec3da5c16ac3a4c68fb5f4d66db648eba6f9abdc6539201`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `46dd4cc86f0c46cd81bf7809d1fa6e4672d6704d23a3ff2ecc14b0ac4e2ed59c`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `c0f484cafafa7e1fea79c20c0916c220cfcd67096ece4138a25bfbb11e9e6cc7`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `32edb739cd06aa442e2a4b764e85303085971292ed624818a5f5f3c9264eff24`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `90e548b85bcac423b6b17995f8b0d470005a8471cfd0f1fe01eac6f021d97f2e`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `29a11a3d30971097a374d3b0396131bfa1c0e02004064748b5f02f4f580ab6c0`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `8ae9c8449ffcbf79473de57d157c4f41b386a2b6d0b4e648d2e12848662859a1`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `f4eb69987abfb35f62bde7b5ce85953f82ab7df11fe6edcfc7a9840c9382da9e`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `a0eb1c023c40884d66823de57c2fba85b888c4ebcd296e99e28e7d30f9399f34`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `7da825f7645b04315928a4a08fcf2b0378c3b132289f2525cbbfe11e85acb8c9`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `2ad0dadf36f1b74458efbd267b635920b5787fcf8a267dfa994b476bfc6a04ef`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `1aa454acdd8d35df416769f69a16dbcc3aff5251cdc5fcb23bb94e9ce537c04f`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `6032f0f6d5c47e86cc3661be35254da1d199c8a2dd764efca79309e837397c11`.

Reviewed 8b5c68e against Go 1.27.0. The constructor assigns decreasing priority in pair order and preserves the first duplicate terminal. Lookup follows a single path through byte tables and compressed prefix edges, retaining the greatest priority encountered. The candidate suffix model traverses the same raw byte sequence. It checks terminals before exhausting the input and skips only the root terminal when ignoreRoot is set. Returned key length counts original raw bytes, and an empty replacement is distinguishable from a miss.

The control checks 14 fixed Lean and pinned native lookups through the actual private method on constructor-built tries. It covers short and long keys in both priority orders, duplicate keys, empty replacement, root selection and skipping, input ending within a compressed edge, an absent table byte, and the 256-entry table with malformed key bytes. The longest-first and root-always-eligible alternatives each fail two witnesses.

The author replay checks 832 boundaries and 60,000 samples using seeds 17, 2718 and 65537. The scope covers exact return tuples on constructor-built tries from even old/new inputs. Mutated tries, nil receivers and allocation behavior are outside the claim.

Run the reviewer control from the repository root after building the Lean executable.

```sh
GOTOOLCHAIN=go1.27.0 go run difftest/strings/generic-replacer-lookup/review/main.go
```
