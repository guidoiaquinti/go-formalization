# Independent review of `strings.makeGenericReplacer`

**Verdict: FAITHFUL**

Evidence SHA-256: `89740fad57de1675dd8c91965ba779e407c07a89ecdb98680ae0cbbc5ce72c31`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `2188e921c3c23ad6019a9bd492e291224d704f2dcf04d00aafa7721a422cd5cd`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `bc595213a0cbad76c1617e640d3686dae6a121224c2c09506e2a82e4710dad5a`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `acddfbf93aef3cfc6bd2bd55e6dd8d877309e1ef9d75945f294127aa184b9317`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `d7fbd2907264ab16f18066e0d820d741884ff74f9fdccfccf4eb388a65597a66`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `bad7281f95b7bf8bda955d831fec2d44814f5d45f1170722d6c09c73461ab822`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `e6bb7eb45a21b5f7b4ca3330e053b1d692c0d963102993f99a0b04d72072e173`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `eca5b96a17a9ffb9a6cf6731703f0fb245c2cb3fd8b10ec72cd3c440e31a509b`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `4d4a4251efb01ea0d2e8a698a610d34b547b43b4fa78caa54d6a98f84d00cdb3`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `feaa1704ec563d27bffd784f762c082c2dff39402174c508ece646de50ad5801`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `34d167e6d6d4d0028be230ea7db1bc589c11fc71c8afa129ba40b46c1fa0a5ce`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `41b1025a1d5657f4c231187e38586aa3799de77f2ca84e89a065e970049b99cb`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `5e8d52cff17f0fd008abe553f4728211b2dcc9ecb3929f296e0031861310cfb4`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `f7adad4ad74695aeadf7a36da9e2494b96734d326f601ad4518446e7fb11626d`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `0a032b262ca6d0b1dd7fd6f8653bb9c6ae932c225af6ade59baf9fc1dbbfa4f4`.

Reviewed 2e37f75 against Go 1.27.0's makeGenericReplacer and trieNode.add. Only key bytes contribute to the dense mapping, ordered by byte value. Absent bytes use the table-size sentinel, and when all 256 bytes are present every index is occupied. The root always uses a table. Each nonterminal unary stretch becomes a compressed edge. An intermediate terminal stays at its own node even when there is just one child. The first occurrence of each duplicate key owns its priority and exact raw value, including an empty value.

The runnable control checks eight complete pinned Go and Lean snapshots. It fixes the alphabet independently of the construction and checks all 256 mapping entries, root/table/edge labels, each terminal value and priority, and the order of occupied children. Cases include an empty root, duplicate empty and nonempty keys, a fully compressed edge, terminals inside a shared prefix, branching, malformed raw bytes, and all 256 bytes supplied in reverse order. Insertion-order mapping, last-duplicate priority, and uncompressed-edge alternatives each disagree with a fixed snapshot.

The differential checks 269 boundaries and 60,000 samples with seeds 17, 2718 and 65537. The verdict covers successful construction from even old/new lists and resulting topology. It excludes pointer identity, intermediate insertion states, and odd lists.

Run the reviewer control from the repository root after building the Lean executable.

```sh
GOTOOLCHAIN=go1.27.0 go run difftest/strings/make-generic-replacer/review/main.go
```
