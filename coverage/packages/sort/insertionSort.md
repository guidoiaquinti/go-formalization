# Independent review of `sort.insertionSort`

**Verdict: FAITHFUL**

Evidence SHA-256: `ae7ffd9160e6c42cab07328dafd9841399545d3e5d2ed24a78d936cb43af13d6`

I checked the exact pinned Go 1.27.0 declaration, the Lean state model and the private native overlay. The loop compares adjacent values strictly and swaps backwards only while Less is true. Equal values leave the order and callback trace undisturbed. The sentinel fixture confines all swaps to the requested subrange. The overlay records every Less and Swap callback and checks the final array against its independent oracle.

I ran `go run ./difftest/sort/local-sorts/review` on darwin/arm64. Twelve independent literal traces matched native Go and compiled Lean and rejected six variant-specific mutant traces. The author comparison passed 6,653 boundary cases and 60,000 samples with seeds 17, 2718 and 65537 for this declaration. [Reviewer control](../../../difftest/sort/local-sorts/review/main.go)

The scope covers valid half-open ranges, representable index arithmetic and indexed callbacks that return normally.
