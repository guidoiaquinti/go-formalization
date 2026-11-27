# Independent review of `index/suffixarray.map_32`

**Verdict: FAITHFUL**

Evidence SHA-256: `32abc2a8d4431701f30fa959660e6c2054677a37f3f2c8e64345dcac542c186d`

Go 1.27.0 scans the inclusive midpoint and earlier cells, keeping positive signed 32-bit IDs and writing each minus one at the visible end. The original order is preserved in the compacted tail. The count argument does not affect the trace, and an empty visible slice panics before writing even with spare capacity. The full replay agreed on 87,308 inputs per declaration, including 60,000 samples at seeds 17, 2718 and 65537. [Model](../../../lean/GoLean/Index/Suffixarray/MapLms.lean) · [Comparison](../../../difftest/index/suffixarray/map-lms/check.go)

The independent command `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/index/suffixarray/map-lms/review` passed six paired native and Lean fixtures and rejected six changed observations. The controls distinguish signed 32-bit truncation from 64-bit positivity, inclusion of the midpoint, singleton mutation, stable tail placement, an unused count, and the empty bounds fault. The scope covers valid backing geometry and sequential access. [Reviewer control](../../../difftest/index/suffixarray/map-lms/review/main.go)
