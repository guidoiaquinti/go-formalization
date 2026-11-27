# Independent review of `strings.Replacer.WriteString`

**Verdict: FAITHFUL**

Evidence SHA-256: `11e0e651521be24c0d3738772ea2a76ddcd87005e834bd678e73c01952d78c4f`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `d61296bfe941d63e941e4da3b3fccfbc6d79f229673debb5779db6070b207a94`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `bb26b8bfd98d4c11d97c78d18224e0006d092a7b9ecafbb85beb2b5934f325e9`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `76080734d531c73c033c3e7f86336eee55c95c0567474e33c9354e53b84a19af`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `e0230cbf9b536027a05c76e180a5d4ba9d88c08b317fb30abbff5c467eccd7bf`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `5825ce09fac28ddd7b815eab0756ca0b681b50f8f0ab07ac40345782e2c555e6`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `9562833df6db4403a8bd7b6b48e3a156aeb021af116f0f653554a52c66024943`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `3fee0953ec17a9c706820237749c2b0d3812a0819164ecfa6f161e47e6cc256d`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `154292f800286115781d2c7dd54cdc48451aed57b43aac59b5586650ff41d6c7`.

Reviewed dc28ede against Go 1.27.0. The wrapper completes sync.Once initialization before invoking the selected writer strategy, so an ordinary returned write error leaves the implementation available for the next call. The model delegates exact writer routing, count and error semantics to the four reviewed implementation models. A byte or byte-string replacement calls Write for the replacement even if the target implements StringWriter. The unchanged spans and single/generic implementation spans follow their own dispatch rules.

The independent control checks eight fixed native and Lean call sequences. Cases cover nil and zero receivers, all four routes, first duplicate priority, a present empty byte-slice write, empty spans from single and generic writers, a short successful write, and a partial failed write followed by reuse with a different writer capability. The count includes accepted bytes from the failing call and trailing dispatch stops on error. Routing a replacement through WriteString or writing an extra suffix after the first error disagrees with fixed traces.

The independent differential checks 2,136 boundaries and 60,000 paired calls using seeds 17, 2718 and 65537. The claim covers representable sequential outcomes and writers reporting counts within payload length. Native observations exercise first and subsequent calls, while the pending/ready distinction follows from pinned wrapper source and the Lean transition. Concurrent initialization and writer panics remain outside scope.

Run the reviewer control from the repository root after building the Lean executable.

```sh
GOTOOLCHAIN=go1.27.0 go run difftest/strings/replacer-write-string/review/main.go
```
