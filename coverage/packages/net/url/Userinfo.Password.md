# Independent review of `net/url.Userinfo.Password`

**Verdict: FAITHFUL**

Evidence SHA-256: `ec235dd3b4606f29a07681a407e3932e1d0e3ac3af621e4e532fbb59f2aa5a60`

Returns the raw stored password and presence flag independently, including a synthetic nonempty password with the flag cleared. A nil receiver returns empty and false. Pinned native Go and Lean agreed on all 126,043 author rows, including 60,000 samples at seeds 17, 2718 and 65537. [Model](../../../lean/GoLean/Net/Url/Userinfo.lean) · [Comparison](../../../difftest/net/url/userinfo/check.go)

Seven independent native and Lean cases passed and rejected seven altered observations in the shared evidence. The registered scope uses feasible raw strings, with constructor-reachable, nil and synthetic receiver states as applicable. [Reviewer control](../../../difftest/net/url/userinfo/review/main.go)
