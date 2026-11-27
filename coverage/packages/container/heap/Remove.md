# Independent review of `container/heap.Remove`

**Verdict: FAITHFUL**

Evidence SHA-256: `eafb45312caa33ba8afa65d54dcc179dd312de5e5a7db3526979f81c446c1ccb`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `c9fd562a02311db77c8ce847f131d87fa0022300aa288bed3cf80f89f68ccd47`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `1b9bde42e0894a045785253210d901da9f95a75248f55c4702053e9212014138`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `609afab519d4a3fbb134e00ea61d5c51c041d2fbb78efe9bd3abfe9dd7735b97`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `a0eed6ef3e1357064a4e59714a6a956f120a6c183f92f2f7d866b55cdb4061ed`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `2f8b71a2b6de1a88a37fce8576f70c4260a0937026509ba61d612557b9dc2eb3`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `b853fce495d8fd86f195920eafa482112a4d905c1f1214149f6b1505954b4d0c`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `dbdcc0822ef33822c2da2efa8c849fafa12f099a8a101c08c6c2b50e9acdfa11`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `bc7749958a082aea789e15e06165f25446d444a17f8216fe7f01709b7e0b4f6d`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `414948677ae690dd78136bf3d1fed34eaacc0b162edca48f6313f66d1c135ee2`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `3403ca2f9e0974b0a8078ded48849ef6f7bc2b33b0290cea9867aa9aee4fd364`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `8c94c716655485c4c168a96f613144e9b4085088a39fd964937fcfe5736071ee`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `8e9652021a3494780175b97913653714a4a11e3005491eda8ed1699224a82385`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `c84b7580c3a44ad1ca22c634ea0338e8aab3d1ec8337f7c6719a60f8ce3fc895`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `34b91deeb523f7c6251dc346c3f916715474dd123228c83cad1e55a37268f707`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `fc058c958223e259b28451bd439439576f989ff264262717c486a8295f648425`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `ebe4950cbcf00278f58187d99c5d49de1dfbf314ad05653b5a6b45d1898456e4`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `e76a8c76a64eed4ef0bd8ba9de5763c137fa7b083296b44ce4e49c8d33086691`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `7244fd1eb8286ab297d3beb38ee8649feb82d95922f7194437a2f999905280e8`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `0b7e167a4b49727cfcd99d98d0c0aaa77d7ba0c2f59358b0b8d147b961c732b3`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `ffd73189a2a5d8c457c3bf8d8ced90a73b733d72671bed10c43066d0c32d0a4c`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `637d07a0439fdf462e95144adc292f5f1e86592ec401714cea26a8a530825a1f`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `ad5d416f390d0f192f8b0bdfb76acdce1ed2437223058eb529ee6152c2b3bfd0`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `ac5a0850e0d1cbd3bf1e571e40ec30b18e8ac40f3ffa53338282f808ec065937`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `8af5efd28fc73f623a47c61ae3c2542bebfbc9f08623c2c9fb7fffe51585894c`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `09d18e8444f9e386d733785b81493d942d6f395f1c2bcfa2a63978a2394ce0dd`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `44150f967ba67f788dfd7965cc1cdb9500dc1fd0414d1f5921cb44b06ac2a9c8`.

Go 1.27.0 obtains the former last index with one Len call. If the
requested index is already last, the next callback is Pop. Otherwise
Remove first swaps that index with the last, repairs only the reduced
heap, then returns the result from Pop. The repair tries down first and
uses up only when down made no swap. The Lean model reuses Fix with
the reduced length and drops its synthetic initial Len event, so it
preserves the pinned callback ordering and arguments. Arbitrary
callback panics, invalid indices and backing storage outside the trace
are out of scope.

The author check passed 8,984 boundaries and 60,000 scripted cases
with seeds 17, 2718 and 65537. The independent control passed six
fixed tagged-heap cases on darwin/arm64. It checked the real removed
entry, survivor identities, heap ordering and exact ordered callback
trace, and replayed the observed comparison choices through Lean.
The fixtures include a singleton, a nontrivial last-index shortcut,
a former child excluded by the reduced length, an upward repair, a
downward repair and duplicate-key identity. Forcing the last-index
swap, using the old full length, and dropping upward fallback failed
two, one and three cases, respectively.

Run the installed control from the repository root with
`GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/container/heap/remove/review`
after building the model executable.
