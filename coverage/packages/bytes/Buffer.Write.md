# Independent review of `bytes.Buffer.Write`

**Verdict: FAITHFUL**

Evidence SHA-256: `48ab50f8befc864469e7751cbf6e8354d70828e47d25c0ea5ef3908c7e689cca`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `dd34df260f0d75c25a9c95446ac24d3290a34f96933fabbe18a8f4115c2b0617`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `654214c250bea662562637e850004b6434296f178d35e1ea3183affd73acd682`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `407e7d46f4dbcd49404c39ecb9a0beeab061e457109b6cf63d9e0714ec3132e3`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `cadc108e32dfab38448f20e2a108e2d8229f22425b4cf61aad5a1174e293156c`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `db6c3d6b65f12171d934ffbb36149fd7eb18fe90d24ff33cedfca0fe538b7a7c`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `43040c5df4c2fa73ceb6f4239b136806a4e56286404ee5bd5e0e949f57425148`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `666acd9ace4b65ae070272d24eac9596d42277e569c2adc1e8fb10d9afa9a667`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `0d7a74c4cf66666e3bbd2d4f65775f602cdfba5214b1170a2689c83333311be8`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `c4de683f7a6f8be3e9a846950b1c1a82c601181b5c399c24da25033f70b94c44`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `e3a241368ff7e6c63e35117ea47f7ec4c437e5f736942cc3f78d9cbf23b23c8f`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `ce319fc9c685c43e04df300314b7782a46176fb38d3e5ca456be905c936fbdff`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `b02fd92d49b1374222fa6d2e48df78fbab49ae2c79f7ef45f15e68af4db12e19`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `12d7e2a7b1ce1c3aaecba9cd9aad2fd2256d7bcc8bfb98bdf6ff2b98c71908b2`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `58fa9bbc8704960392e1158b190ac497dce3f4a6e504276f38cb1fa51ca6d149`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `475ed02e5dfbd4a9cc057c7a9257417b29c3b2c88d434a28bb1389afc4864cc9`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `25e7da5df3b8e280dfa002abeac809f630f80912eb57868c3970ab3995bd73c8`.

I traced Go 1.27.0 `src/bytes/buffer.go` through `Write`, `tryGrowByReslice`, `grow`, and `growSlice`. The model clears the read flag before reserving space and preserves the consumed prefix on a successful fast reslice, including a zero-length write. When growth is necessary, a fully consumed buffer resets before trying the reslice again. A partially consumed buffer slides only when unread length plus write width is at most half capacity. This move changes any overlapping source on the old backing before the final copy. An allocation saves the unread bytes on new backing and reads the source from the old backing. Within the recorded successful requests, `max(length + count, 2 * (capacity - offset))` and the pinned small-object capacity classes give the observed capacity. The runner checks exact stored and spare bytes, offset, read flag, nilness, reuse and byte count. The same-package native bridge also checks that a successful call returns nil error.

The author comparison passed 21,714 boundary rows and 60,000 seeded rows using 17, 2718 and 65537, including detached and aliased inputs. The independent control passed nine fixed fixtures on native darwin/arm64 and Lean. It constructs live aliases using public `bytes.Buffer` operations, checks `UnreadByte` invalidation and backing contents, and supplies literal model outcomes. The compaction fixture rejects reading a pre-move snapshot, making the half-capacity threshold strict, and retaining the prior read flag. I ran each of those mutations separately against the fixture.

Aliasing behavior and spare bytes are observations of the pinned implementation. The claim covers legal buffer states and successful allocations through a raw request of 512 bytes. Allocation errors, `ErrTooLarge`, and native execution on other platforms remain outside this review.
