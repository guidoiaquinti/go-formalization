# Independent review of `internal/bytealg.abigen_runtime_memequal_varlen`

**Verdict: FAITHFUL**

Evidence SHA-256: `cd8aed5e3de78b4b55f5bc7a5ed42e165ddc850de75437ecd4d74c89033d1566`

I checked the valid closure and readable window scope. The comparison
uses the exact byte count stored in the generated closure. The native
fixtures verify the compiler descriptor's entry name before invoking it.

Seven independent native closure fixtures and eight Lean observations
passed. They cover displaced initialized views, aliasing, a difference
at the last included byte, ignored bytes beyond the window and a
synthetic zero-count model row. The author's 2,496 boundary cases and
60,000 seeded cases passed using 17, 2718 and 65537. Zero-count closure
generation and native execution on other architectures remain open.
