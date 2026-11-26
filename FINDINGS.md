## POTENTIAL — net/netip: Addr.Next changes the stored bits of the zero Addr

Go 1.27.0 on darwin/arm64. The documentation says that Next returns the
zero Addr when there is no following address. If that includes a zero
receiver, the expected result is Addr{} with all-zero As16 bytes. The
observed result has BitLen 0 and As16 bytes ending in 1, so it differs
from Addr{}. The pinned implementation increments the address words
before checking overflow.

Reproducer from the repository root after building the Lean executables:

```sh
GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/net/netip/addr-next
```

The boundary case with the zero receiver checks the observed behavior.
Confirm whether the documented no-successor rule is intended to apply
to the zero receiver before treating this as a Go implementation bug.
