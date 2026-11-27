# Independent review of `net/url.validUserinfo`

**Verdict: FAITHFUL**

Evidence SHA-256: `87c25700a6083bf538cb18ff3966a2a58cb54e00dd9272481fb8d6eb50d03733`

Accepts the pinned ASCII allowlist including at signs and percent without checking the following digits. It rejects slash, NUL, valid non-ASCII and malformed UTF-8. Pinned native Go and Lean agreed on all 126,043 author rows, including 60,000 samples at seeds 17, 2718 and 65537. [Model](../../../lean/GoLean/Net/Url/Userinfo.lean) · [Comparison](../../../difftest/net/url/userinfo/check.go)

Seven independent native and Lean cases passed and rejected seven altered observations in the shared evidence. The registered scope uses feasible raw strings, with constructor-reachable, nil and synthetic receiver states as applicable. [Reviewer control](../../../difftest/net/url/userinfo/review/main.go)
