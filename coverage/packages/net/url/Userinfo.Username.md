# Independent review of `net/url.Userinfo.Username`

**Verdict: FAITHFUL**

Evidence SHA-256: `1db7bed152de22140d1089140007778917b013847a0d1fc135a8a743e688edc2`

Reads the stored username bytes independently of the password and its presence flag. A nil receiver returns the empty string. Pinned native Go and Lean agreed on all 126,043 author rows, including 60,000 samples at seeds 17, 2718 and 65537. [Model](../../../lean/GoLean/Net/Url/Userinfo.lean) · [Comparison](../../../difftest/net/url/userinfo/check.go)

Seven independent native and Lean cases passed and rejected seven altered observations in the shared evidence. The registered scope uses feasible raw strings, with constructor-reachable, nil and synthetic receiver states as applicable. [Reviewer control](../../../difftest/net/url/userinfo/review/main.go)
