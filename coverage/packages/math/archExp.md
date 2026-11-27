# Independent review of `math.archExp`

**Verdict: FAITHFUL**

Evidence SHA-256: `e86b4c8618eefc25ee8605ae7e5dfe2b8da637bc40904eefcf8430b4b8a552e5`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `b5a49815516a5a889fa5005e62b11b87ad3e6667350b8041a420a0003fd835c8`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `ad22f979d3c6da1f9bd8affbee3ef17cf1fefc8a43db9bbabcace1ec77fcdc18`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `717b1b6ad9bdb7253ebabf7874072b69f3e092afa89027e353c8ca076df479db`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `74e66875102bf1cbda1f052afcfcacf742b3e3a6d7649fc5a87031b4e363bde2`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `4cc5ba9b09964863dfcf5b0b0108c1d9da70644135a504f9acac64e435b7fcf7`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `c82fc542b4104a879c74f7e2f36268a97bcd88ad9aaa9d715b07b5c84528807d`.

Reviewed c477f03 against Go 1.27.0's arm64 and amd64 assembly and the existing Lean Exp model. The arm64 branch preserves the original NaN bits, applies near-zero addition, uses fused reduction and rational correction, and stages the subnormal scale. The amd64 branch checks nonfinite inputs before reduction, uses nearest-even integer exponent conversion, selects separate or fused polynomial operations via the pinned useFMA flag, and returns infinity when the reduced exponent reaches 1024. Both profiles preserve the specified infinity outcomes.

The control checks 33 fixed model rows across arm64 and both amd64 branches. It verifies signed zeros, NaN sign/payload, infinities, and near-zero and subnormal boundaries. At -745 the pinned arm64 path returns zero and the amd64 model returns the minimum positive subnormal. At 709.5 the arm64 path stays finite while both amd64 models return infinity. These cases detect conflating the architecture-specific scale paths. All eleven arm64 fixed rows also went through the actual native private assembly entry.

The independent full replay passed 69 boundaries and 60,000 samples under seeds 17, 2718 and 65537. Darwin/arm64 matched every model bit and differed from the independent high-precision finite reference by at most one encoded step. Native amd64 was not executed locally. The amd64 models were checked against the pinned assembly and profile controls. Loong64, RISC-V, s390x and configurations without the modeled assembly remain outside scope.

Run the reviewer control from the repository root after building the Lean executable.

```sh
GOTOOLCHAIN=go1.27.0 go run difftest/math/arch-exp/review/main.go
```
