# Independent review of `encoding/json.HTMLEscape`

**Verdict: FAITHFUL**

Evidence SHA-256: `0e793114b3bef83e591600ef33d2c4f2805821da55408ae084b0c3589bac6860`

The public function appends to the existing visible prefix and faults on a nil buffer. It emits lowercase escape sequences for the targeted ASCII bytes and UTF-8 separators, copying other bytes unchanged. The review does not model buffer offset or internal storage effects. The pinned default Go 1.27.0 profile passed 67,056 author cases, including 60,000 samples at seeds 17, 2718 and 65537. [Model](../../../lean/GoLean/Encoding/Json/HtmlEscape.lean) · [Comparison](../../../difftest/encoding/json/html-escape/check.go)

The independent command `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/encoding/json/html-escape/review` matched eight native and Lean fixtures and rejected six altered outcomes. The fixtures distinguish exact UTF-8 separator byte sequences and preserve pre-escaped, ordinary and malformed bytes. Source and destination are disjoint. The alternate legacy build and allocator behavior remain outside review scope. [Reviewer control](../../../difftest/encoding/json/html-escape/review/main.go)
