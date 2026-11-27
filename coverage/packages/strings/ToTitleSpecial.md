# Independent review of `strings.ToTitleSpecial`

**Verdict: FAITHFUL**

Evidence SHA-256: `3d1f3e8b3782154381adf82c5d134818fcc119c4f6da2855046311eb33ce7f93`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `0d0848a3c8b36fead28cd3eae1a944c5616ff39e6013b9814538cb4243e3c27e`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `fd851367d0d5585527e8c90ae30efd5463f4940da800f0db5c879a2897e4e0d7`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `df4c83dcf59b8bb9bf8a83dc4af79cdfc127bb745b4e98bf88088bae4b700efc`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `c10b0993d1006b38caf956a94e774b60fa71f6c7a25c2d03d5133371d68db08c`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `2b5f6558d4d754530782cffabaf882805aa8f28b1491886a0c9903b961e48279`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `7595689114565a78964490b9d745a9a85a486c2e2ca25de7f6f9fa33c637e727`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `127ef8b43fb4dec668569a9a08562e6bb172872d204ce7b19944f23afbd0affb`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `b1fc9b395d437266750b8db6e90b38d04ac9ad90889300d3d1f4a0a0496e85f1`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `5a3d89a21c11db4a5b98441b42b58cb5ab2811ea48a8d32000230348612a19bf`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `587bdf77a88a30abf619b93c213d4e1d9d55eccf189a3ad2c53f5890b1609373`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `b922c080769b8b7930fa9c95c70e94d609d1da10b15afaec1539cb67510972b1`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `a0545d5bc82a993d1695c604c7d1ce120e9d052d4435b03042f298f270266add`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `bc93e5037a0c8b0b800643e1ea52b4b82fd8940772312435f9887a0953afda93`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `6c3b2a5e4cba9c229237ab2adad540162c327334ecd08132d91db416f73bceff`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `72a0012b30a4c272b58b8f7ff24d92966526a5df540ca973e698dfb41dbc7825`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `0356ec44465397994f6a4de961d1e672598fbd87a16933d72fef492cf9e512e4`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `dbc9bd433d32bb18d7e7e90ed5d5c707edab9b3180ef52e6a1b761fb26ce1a39`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `d263da2ba69918f279ec7b9e376e89e0c7e639d17136d5e302261aac835c4be2`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `448a86acf9d2bec6cff4d029de485ddbd03a2971595d0450e5fd18ab9e54889b`.

Reviewed 268a7d8 against pinned Go 1.27.0. The wrapper maps each decoded rune through SpecialCase.ToTitle. The selected title delta is the third delta in the first entry reached by the pinned half-open midpoint search. Signed 32-bit range bounds, wrapping deltas and even offsets in alternating ranges agree with unicode.ToTitle. A hit suppresses the ordinary title mapping even when the selected delta is zero. A miss uses the ordinary pinned title table. Negative mapped values are removed, and invalid nonnegative values encode as RuneError. Consecutive malformed bytes are decoded one at a time.

The control checks ten fixed pinned Go and Lean results. It distinguishes ordinary title fallback from uppercase on the digraph `ǳ`, exercises the title delta field separately from the upper and lower fields, and separates midpoint traversal from first matching range. An unsorted miss, zero-delta identity, negative drop, surrogate output, signed endpoint wrapping, Turkish special case and adjacent malformed bytes are included. The control detects first-match traversal, fallback on a zero-delta hit, the lower delta, and the uppercase fallback.

The differential checks 1,262 boundaries and 60,000 samples across seeds 17, 2718 and 65537. The verdict concerns exact output bytes over representable strings and supplied SpecialCase tables. It excludes allocation behavior and unsafe backing-storage observations.

Run the reviewer control from the repository root after building the Lean executable.

```sh
GOTOOLCHAIN=go1.27.0 go run difftest/strings/to-title-special/review/main.go
```
