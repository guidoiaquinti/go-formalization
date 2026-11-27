# Independent review of `net/url.UserPassword`

**Verdict: FAITHFUL**

Evidence SHA-256: `74d35443914ed70238d01c0e45c33f734397e0c25fbf79bddf53a31e6a75de38`

Constructs nonnil Userinfo with the exact original username and password bytes and sets the flag even for an empty password. Pinned native Go and Lean agreed on all 126,043 author rows, including 60,000 samples at seeds 17, 2718 and 65537. [Model](../../../lean/GoLean/Net/Url/Userinfo.lean) · [Comparison](../../../difftest/net/url/userinfo/check.go)

Seven independent native and Lean cases passed and rejected seven altered observations in the shared evidence. The registered scope uses feasible raw strings, with constructor-reachable, nil and synthetic receiver states as applicable. [Reviewer control](../../../difftest/net/url/userinfo/review/main.go)
