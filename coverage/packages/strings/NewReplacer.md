# Independent review of `strings.NewReplacer`

**Verdict: FAITHFUL**

Evidence SHA-256: `8453b35fb2933e28e1d1dc2413797ba644c694c82069bbacb7a4170fb3b11ac7`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `cecdb66403ca9fb869864d50e34548cbea2e2e62771f0f6fde1799a1fadacd61`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `b05b6a32642a2cc087a28a39f17a7eae7d693702d2d6ac7b648dde46356a339d`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `bcb551da062257dccf53f1e9ad2b959208aec8fbb7fbf8bb1b705c767c4b9423`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `f465f36d5c52bd821ab6d054c4ac184a012d08f40959ed8b45c07bc0dcfcf511`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `73401572194fa5217f4ef75c4efd6a161ebd6e84bd1892932163581786554765`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `3c56ada08bb069561a0fdb7b464b62d9b61cdc3c5257e2284c9cfa2077281c94`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `9be5331935d5e034cf8009d57b5da18b038fd341bc024ed210aee5d96959c230`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `e332b1bbecbcccd89ae61fa3f81f9236af122bbf257659b6e2cd989f86b94370`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `540fe199463159a44efde3e340324357929c72fa8f489fea79b3a8b5a6a76707`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `2b98c64fefb0677bb2bc27704c2842add217296a0cac1499359e80300e53c980`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `2d674fde163ce91954260f28eedf96e65c72a6cd51952af69c4098bc61727a67`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `332f619404b31d95f977cf2482952b733b9cee1f3839099b73a00961510274f0`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `c8c3bdde1ca2f7f5eba868dccdecb1ca1a817fbfa30b3546fb1406c0d4136714`.

Reviewed edd2c47 against Go 1.27.0. The pinned constructor checks argument-count parity, then creates a new Replacer whose pending argument slots come from appending to a nil slice. The returned Once is zero and the implementation is initially nil. Both empty argument-slice forms produce nil pending storage. Nonempty valid input preserves order and exact raw bytes while detaching the slots from the caller.

The control checks eight fixed pinned Go and Lean observations. The native observation mutates the original slots after construction, then checks the retained copy, exact argument bytes, nilness, and lazy state. Cases cover odd panic text, both empty forms, empty strings, duplicate keys, and malformed bytes. Slot aliasing, eager construction, and retaining nonnil empty storage would each disagree with an explicit observation.

The differential checks 778 boundary cases and 60,000 samples under seeds 17, 2718, and 65537. The claim covers initial successful constructor state and odd-count panic. It excludes backing capacity, concurrent mutation, unsafe string mutation, subsequent lazy building, and allocator failure.

Run the reviewer control from the repository root after building the Lean executable.

```sh
GOTOOLCHAIN=go1.27.0 go run difftest/strings/new-replacer/review/main.go
```
