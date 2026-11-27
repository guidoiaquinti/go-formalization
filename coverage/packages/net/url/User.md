# Independent review of `net/url.User`

**Verdict: FAITHFUL**

Evidence SHA-256: `daf35456ec8632966e5e2cd2b6309bafdb2b3d11ebfbf7bd03419f1fe6044405`

Constructs a nonnil Userinfo, retaining the raw username and clearing both password fields. A native bridge reads the constructor state separately from the accessors. Pinned native Go and Lean agreed on all 126,043 author rows, including 60,000 samples at seeds 17, 2718 and 65537. [Model](../../../lean/GoLean/Net/Url/Userinfo.lean) · [Comparison](../../../difftest/net/url/userinfo/check.go)

Seven independent native and Lean cases passed and rejected seven altered observations in the shared evidence. The registered scope uses feasible raw strings, with constructor-reachable, nil and synthetic receiver states as applicable. [Reviewer control](../../../difftest/net/url/userinfo/review/main.go)
