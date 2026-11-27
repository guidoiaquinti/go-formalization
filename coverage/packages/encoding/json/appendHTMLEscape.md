# Independent review of `encoding/json.appendHTMLEscape`

**Verdict: FAITHFUL**

Evidence SHA-256: `32663f712b0a065f504a07c1254b68ff59bea99ec7d3e3efa4aa9b6a31b97179`

The private helper appends an ordinary span and its replacement as separate chunks. A fitting span remains in the original backing even when the next six-byte replacement forces growth. The controls distinguish failing an entire replacement chunk from retaining its partial bytes. Nil input and unchanged malformed bytes are preserved. The pinned default Go 1.27.0 profile passed 67,056 author cases, including 60,000 samples at seeds 17, 2718 and 65537. [Model](../../../lean/GoLean/Encoding/Json/HtmlEscape.lean) · [Comparison](../../../difftest/encoding/json/html-escape/check.go)

The independent command `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/encoding/json/html-escape/review` matched eight native and Lean fixtures and rejected six altered outcomes. The fixtures distinguish exact UTF-8 separator byte sequences and preserve pre-escaped, ordinary and malformed bytes. Source and destination are disjoint. The alternate legacy build and allocator behavior remain outside review scope. [Reviewer control](../../../difftest/encoding/json/html-escape/review/main.go)
