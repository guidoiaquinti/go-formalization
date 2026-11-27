# Independent review of `sort.heapSort_func`

**Verdict: FAITHFUL**

Evidence SHA-256: `8c4aed76fa30cac71c69b95f0c7ee8022ec2c0ffe57fa66120e4c2588c5c42f9`

I checked the exact pinned Go 1.27.0 declaration, the Lean state model and the private native overlay. The helper builds a max heap, then swaps the root with each decreasing right endpoint and sifts. The last self-swap is observable even though skipping it would preserve the final array. The overlay records every Less and Swap callback and checks the final array against its independent oracle.

I ran `go run ./difftest/sort/local-sorts/review` on darwin/arm64. Twelve independent literal traces matched native Go and compiled Lean and rejected six variant-specific mutant traces. The author comparison passed 6,653 boundary cases and 60,000 samples with seeds 17, 2718 and 65537 for this declaration. [Reviewer control](../../../difftest/sort/local-sorts/review/main.go)

The scope covers valid half-open ranges, representable index arithmetic and indexed callbacks that return normally.
