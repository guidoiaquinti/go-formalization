# Independent review of `bytes.Buffer.WriteByte`

**Verdict: FAITHFUL**

Evidence SHA-256: `f1a99a773da57f20701f7b6ef7936e3e5e9520f6eaadcd5bea814cf4fd0da465`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `f1caa706841800c12d576355cd4ad14ad191d6530d0dfd8f1bb16e1dae3a0995`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `fec25e1d38a77c51e9e8864c5e26988eea0d82db38c8f96d2a5616a03a267938`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `f4013bcf35ab8d8161a1ff0d545847f95606c0121984d3238470cc15d674d1b5`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `9d06526b5d7b23f086f4b00503e3cecf610aceee1a8811e54f848edb140e1701`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `d447a76e09b2807eb5be7594c67a825fbfe760a9f9b3edf4a2abadee0d40ed22`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `9c308b55d813e9f6023faa1fd84e77bf909208ebe360fc498b3fb4180c38efe4`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `c3768b91ee47d1a5dd93bc8cfcbad974541a3f6355f6d90239497790245482ff`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `82db6d777363dd6fbda7ac578683e8b1b26f4ae8d7dc19b9b1abbe30c7f632f2`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `76dd09c1b113f8a05dc76e728cacdd4f3ad06aa8d405b0fd73d1d33c03f164b2`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `1100a533fed733b93277fdbce8dc0337d15bf292fc6e8f7a13afa3e478ec49a4`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `3c9a88aba2b767b96ef221584e1d0e76e075f260669884085d5191e25197ad45`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `9ddf59da1fa8ed31c6c9740716523b64b4d9e5b8052745f5f1763c18cc4cde09`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `8795fe5bc302abce4a6c64e2e75a6b302e47504cb1705217ecdcfcb94e9eeb7c`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `d8931dd3544d1ed46da41a77d25dbc60dd7815e7a6b820457745f4f9c597f660`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `13a3080fd7d182c63b5a60b9516c6391514b29506bd7ae0a15dca8a57d071dd1`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `223a98e6946ed69573ababfef96973e09b79f90338e153dad785ba2c2e027b86`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `7253dddca82afc263e2ad984d9548bf61e202a4e8a18f3ecaebbd26b8ee439c7`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. The existing capacity helper became public without changes to its definition. Previous digest `c6efd33848bca954088f2e1f7bf88f08bcaa74eeb068bc97e80a8c687b99bd79`.

The pinned Go 1.27.0 method clears the last-read flag before its first reslice attempt. The model preserves that priority, the exhausted-buffer reset after a failed reslice, the inclusive half-capacity compaction threshold, and the different first allocations for nil storage and a nonnil empty slice. The reallocation size request simplifies to the greater of original capacity plus one and twice the capacity after dropping the consumed prefix. The model includes the pinned byte allocation size classes through a 512-byte request and preserves untouched spare bytes when it reuses storage. [Model](../../../lean/GoLean/Bytes/BufferWriteByte.lean)

The direct overlay and Lean matched all 17,417 boundary rows and 60,000 sampled rows from seeds 17, 2718 and 65537. Ten independent fixtures matched both paths and covered the 512-byte request, nil receiver, nil and nonnil empty storage, reslice before reset, and even and odd compaction thresholds. Two public checks held live aliases through compaction and reallocation and checked that WriteByte invalidated the previous read. The strict half-capacity threshold, always-compact policy and conflated empty-slice cases each failed a fixture. Allocation failure, requests above 512 bytes and invalid state geometry remain outside scope. [Comparison](../../../difftest/bytes/buffer-write-byte/check.go) · [Reviewer control](../../../difftest/bytes/buffer-write-byte/review/main.go)
