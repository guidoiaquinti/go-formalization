# Independent review of `bytes.Buffer.readSlice`

**Verdict: FAITHFUL**

Evidence SHA-256: `4c85434e2decdb774074b4a2688bb8ebe6dc99f4c0b685c7fb8ec682f81f35fe`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `339d5b02357f28df1ad505d21535d0551f84cf3603aa6bcf3392d5c3b3c9d657`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `77101ff2ee4113a4f4c05ebee8e777768544c68339a6b762c83d67ebe8d79425`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `32fbb5e770f51adfdbd8cdcda50119cc9cdfe73931856c387a7b65c135bbf6ba`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `00d79dbc6e62dc5c1a8430f49c44872252f533742ff7644c60d2d5e82137494c`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `3d0df81f8492803f3fbea1aba0ea929004b5ede62c1176da756332277633d3c2`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `0b603737a2df6a97532abb848319dd6f648f9b03f151114085ffadb7cb3403f7`.

I traced Go 1.27.0 `src/bytes/buffer.go`. The helper searches the unread suffix for the first raw delimiter. A hit includes that byte and returns nil error. A miss returns all unread bytes and the exact `io.EOF` sentinel. The returned slice retains capacity to the end of the original backing array, and the cursor advances without resetting the buffer. The source sets the ordinary read flag even when the result has length zero. The native bridge also mutates each nonzero-capacity return briefly to observe its backing alias.

The author checked 11,522 boundaries and 60,000 sampled cases using seeds 17, 2718 and 65537. Seven independent literal cases passed the private darwin/arm64 overlay and Lean. They cover nil receiver panic, nil versus nonnil empty results, an exhausted nonnil view with spare capacity, repeated raw delimiter bytes, a delimiter only in the consumed prefix or spare space, and a hit in the last stored byte. Scratch Lean mutations that excluded the delimiter, clipped the capacity, or retained the old read flag failed two, four, and six cases respectively.

The scope assumes legal buffer geometry. Native execution on other platforms was not part of this review.
