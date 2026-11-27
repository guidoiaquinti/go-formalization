# Independent review of `math.archHypot`

**Verdict: FAITHFUL**

Evidence SHA-256: `5417ec340f178a58aa8ae39a4586c8fde9bc2b15c69f96b4d8767840808272ff`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `29226e54c2b0889fe51afaddee6dd49af5a29312dd1f6cba1bd12398e0400a0f`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `3a2085b46f64875e76fb5ae7635175193f282d9b3e34fff8946a521f2adc7da6`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `c4ddf93ff8e12c766911e1618f212331fdbc0aca2a80ad2f114ef27eb6d1d866`.

I checked the Go 1.27.0 amd64 SSE2 source and the fallback build
constraint. The assembly takes absolute magnitudes, gives infinity
precedence over NaN and uses the canonical quiet NaN word. Finite inputs
use magnitude ordering and split binary64 arithmetic. The fallback
reports false availability and panics with the exact string.

The control passed 12 native darwin/arm64 calls and 24 Lean rows.
Its fixed oracles cover exact triples, both infinities, canonical
NaNs, signed zero, subnormals, and overflow. The author's 1,936
boundaries and 60,000 seeded rows passed using 17, 2718 and 65537.
The independent midpoint oracle found a maximum difference of two
encoded steps for finite values and four classification differences
near overflow. Native amd64 and 386 execution remain open.
