# Independent review of `bytes.NewBuffer`

**Verdict: FAITHFUL**

Evidence SHA-256: `c041d983d3cf978d6980168309f26e2bda6911c75fb6d452f5ca79aa270acd27`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `6ea62738e05280d8b10bf61833c2c04b12933d909858b7f99598d921e4fe192d`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `dddf4dc5a44db436934ce31a6eec446fd0b254bee4dfd48ffdc7aa0f294ff1d1`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `70af6f54f77ff37ca991de072514d4fafd24f8ab5663d2915557eeeafd80bc1f`.

Pinned Go 1.27.0 stores the original slice descriptor in a fresh Buffer, leaving the cursor and previous-read flag at zero. The model preserves the supplied view's stored bytes, spare capacity, spare bytes and nilness. The Go comparison checks a saved address from before ownership transfer and observes the adopted storage through the Buffer. All 4,489 boundaries and 60,000 samples at seeds 17, 2718 and 65537 matched. [Model](../../../lean/GoLean/Bytes/NewBuffer.lean) · [Comparison](../../../difftest/bytes/new-buffer/check.go)

Six independent Go and Lean fixtures passed, distinguishing nil and nonnil empty input, capacity truncation, offset subviews and malformed bytes. An additional native first write reused an adopted empty view's positive capacity and preserved the untouched spare byte. The fixtures distinguish collapse of empty nilness, losing the supplied spare capacity and zeroing its bytes. They cover successful construction with legal slice geometry. [Reviewer control](../../../difftest/bytes/new-buffer/review/main.go)
