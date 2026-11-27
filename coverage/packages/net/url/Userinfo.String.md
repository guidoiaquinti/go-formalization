# Independent review of `net/url.Userinfo.String`

**Verdict: FAITHFUL**

Evidence SHA-256: `bda4e73762e5482d519fdb84136b49751cd5d853a52516a1e5ebb0efe8f5166a`

Adds the colon only when the flag is set, including when the password is empty. The byte controls distinguish escaping reserved bytes and invalid UTF-8, preserving literal plus, and representing space as %20. Pinned native Go and Lean agreed on all 126,043 author rows, including 60,000 samples at seeds 17, 2718 and 65537. [Model](../../../lean/GoLean/Net/Url/Userinfo.lean) · [Comparison](../../../difftest/net/url/userinfo/check.go)

Seven independent native and Lean cases passed and rejected seven altered observations in the shared evidence. The registered scope uses feasible raw strings, with constructor-reachable, nil and synthetic receiver states as applicable. [Reviewer control](../../../difftest/net/url/userinfo/review/main.go)
