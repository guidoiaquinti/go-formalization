# Independent review of `bytes.Buffer.ReadBytes`

**Verdict: FAITHFUL**

Evidence SHA-256: `29190063740d6b2c019932d76a585fd4ce79795d605f653dcc30f1c03b3dac90`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `69db2b2202870765e9e097f611cc632b69f9e8b697a3dcdeb0ccd569bd0e8aba`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `756a300dcd9c2db9153ef32d86315e569d567c6fd19e7c37941747aedb86ea50`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `4657ea042546ac8365ee9b6bbf779b836a68f76c7a90a212f2ac5fc04c2d402c`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `d10f8049daa3748a67bfb9395e202491897abc15e2ee779f198d442eeeb0c049`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `ab026b3ca95cbc2599a2343c2e7b7eb68602a94943e5a97e151f8f194aae0e34`.

Go 1.27.0 reads through the first delimiter of the unread suffix, or returns that entire suffix with io.EOF. Both paths move the cursor and set the ordinary read flag, including an empty EOF result. The outer append copies nonempty bytes into new storage and returns nil for an empty result. The model reproduces the copied bytes, rounded capacity and zero spare bytes through the stated successful 512-byte allocation bound. [Model](../../../lean/GoLean/Bytes/BufferReadBytes.lean)

Pinned Go and Lean matched all 12,034 boundaries and 60,000 samples using seeds 17, 2718 and 65537. Seven independent Go and Lean fixtures matched, including a delimiter before the cursor, one only in spare storage, an exact last-byte hit, and a nine-byte result rounded to capacity 16. The public control mutated the result without affecting its source, reused the buffer without affecting the result, and confirmed that an empty EOF reinstalls the ordinary read flag. Larger allocations and allocation failure remain outside scope. [Comparison](../../../difftest/bytes/buffer-read-bytes/check.go) · [Reviewer control](../../../difftest/bytes/buffer-read-bytes/review/main.go)
