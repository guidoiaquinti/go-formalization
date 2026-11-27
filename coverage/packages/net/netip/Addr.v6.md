# Independent review of `net/netip.Addr.v6`

**Verdict: FAITHFUL**

Evidence SHA-256: `b4c49eca4e9336cb8efe93d680783ee3b16ef7b5d53cf00ac3feae47156b76db`

I compared the Lean model and native overlay with the pinned Go 1.27.0 declaration. The byte selector alternates the 64-bit limb every eight indices and selects the byte in network order.

I ran `go run ./difftest/net/netip/addr-fields/review` on darwin/arm64. Five independent literal rows matched pinned native Go and compiled Lean. The checks distinguish premature index rejection, unmapping a sixteen-byte address, retaining an IPv6 zone, and treating a valid all-zero IPv6 address as invalid. The author comparison passed 2,308 fixed cases and 60,000 sampled cases at seeds 17, 2718 and 65537. [Reviewer control](../../../../difftest/net/netip/addr-fields/review/main.go)

The domain contains reachable address states and sequential calls. The tests use raw byte zones and check the value returned by each operation.
