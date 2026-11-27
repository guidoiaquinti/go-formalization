# Independent review of `index/suffixarray.map_64`

**Verdict: FAITHFUL**

Evidence SHA-256: `8246cf2c17aa8fe73d97e49d3e013d35987acc1de8eca91676874e073f93ecc9`

Go 1.27.0 scans the inclusive midpoint and earlier cells, keeping positive signed 64-bit IDs and writing each minus one at the visible end. The original order is preserved in the compacted tail. The count argument does not affect the trace, and an empty visible slice panics before writing even with spare capacity. The full replay agreed on 87,308 inputs per declaration, including 60,000 samples at seeds 17, 2718 and 65537. [Model](../../../lean/GoLean/Index/Suffixarray/MapLms.lean) · [Comparison](../../../difftest/index/suffixarray/map-lms/check.go)

The independent command `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/index/suffixarray/map-lms/review` passed six paired native and Lean fixtures and rejected six changed observations. The controls distinguish signed 32-bit truncation from 64-bit positivity, inclusion of the midpoint, singleton mutation, stable tail placement, an unused count, and the empty bounds fault. The scope covers valid backing geometry and sequential access. [Reviewer control](../../../difftest/index/suffixarray/map-lms/review/main.go)
