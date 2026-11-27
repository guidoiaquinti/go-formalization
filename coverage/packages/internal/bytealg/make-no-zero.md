# Independent review of `internal/bytealg.MakeNoZero`

**Verdict: FAITHFUL**

Evidence SHA-256: `c53220baf6b63763f6bfd4676eee429c5de10711217cf4fa16052952949761d9`

Go 1.27.0 links the declaration to `runtime.bytealg_MakeNoZero`. The implementation checks the length after unsigned conversion, rounds small valid requests to the pinned allocator classes and reslices to the requested length. The zero allocation produces a nonnil empty slice. The Lean model represents that geometry and the explicit length failure. The runtime threshold is parameterized in the model. Native tests use the ordinary local 64-bit allocation profile and only request successful small allocations or values known to exceed the length bound.

The author tested 5,125 boundary cases and 60,000 samples using seeds 17, 2718 and 65537. Nine independent native and Lean fixtures check nilness at zero, length and capacity at allocator boundaries, and negative and maximum machine-int panic cases. Replacing rounded capacity with exact length fails four controls. The native bridge writes each capacity byte before inspecting its value.

Initial contents, resource-dependent failures and other successful allocation sizes remain outside scope.
