# Independent review of `strings.genericReplacer.Replace`

**Verdict: FAITHFUL**

Evidence SHA-256: `ec0d4a3375274e6dfc4be501dbe210a9f7294778cf3fdd2ddd12f0102c83243e`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `c964f1fe4a59734c80d1e3308f534377551cf55ea6c22d61c770d2be1752ea6d`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `4afb745509208961fba14eabfe2170ab4b0c50851249735544fc1479009e767d`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `c6720fb323fc3cf64845a7d557dc52306c0bedfe9412b57d2971b0fea8f2454a`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `8d5808f8deb6ba12612b4f7058b79101ff26eb649acbd6ab57e412efe7568778`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `6300db1f4ae98afb6f3572d7a55c1fadeddd0df2bac9cfa56266d28a4e12b909`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `20bbf7119aea0a46b91f05585ed03b89b4f4932b2ae72c91c9bc1516373b43bc`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `3141fb409bb3f9bfbdb0605ce9c7ff4c2e3412bc50b80a6761cd5986130f5416`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `4aa0c1eb85a1da8e9828757cfb8984af180eaa70dd8785d0f6e553601ff2217b`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `14627e6463d8264c3ff53664b6bb9c79a3b65622a788e3f197b05b219b48bede`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `083863acfe7280ba7d75e2bad6d35e3e8c8b6cb5a3639be23a1078db358e9d60`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `ebe53124909b8b57ac0e80c9129ea85d8146d477e91e1d8b9582206baae95897`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `636fa7c6ae696e2ef66b10af0ca8648b8f371d8a44c2d6b95525a0fa20a6b883`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `487d361fba5ab5f24c157fe2120a7b193ff735462474d85cfddb17de4d676b31`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `369fdce1725783e68c209c4d4ab3e29d02adc59613703f7205e28817883f4aa5`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `d1e23f91b8a01a5e306ef471ae1dc08b00ef11e49c4cf6f0639deec38b72c333`.

Reviewed 8d8c975 against Go 1.27.0. The private Replace writes into an appendSliceWriter and returns the accumulated string. The Lean model concatenates the exact segment sequence from the reviewed generic writer model. The constructor's first matching priority applies at each raw byte position and replacements are not rescanned. A selected empty key can precede a nonempty match at the same offset and can match at the terminal position. The scan's byte alignment includes the interior of valid multibyte encodings.

The independent control compares 13 fixed native and Lean outputs. It covers empty input, malformed raw bytes, the position inside UTF-8 for an empty key, opposite priorities between empty and nonempty keys, overlapping key lengths, duplicate keys with an empty first replacement, and noncascading replacements. Rune-boundary empty matching, longest-match priority, and cascading replacement each differ from a fixed pinned result.

The registered replay checks 656 boundaries and 60,000 samples across seeds 17, 2718 and 65537. The claim covers exact output bytes for constructor-built tries from even old/new lists with representable allocation success. Mutated tries, nil receivers, allocation, and storage identity remain outside the modeled scope.

Run the reviewer control from the repository root after building the Lean executable.

```sh
GOTOOLCHAIN=go1.27.0 go run difftest/strings/generic-replacer-replace/review/main.go
```
