# Independent review of `bytes.Buffer.grow`

**Verdict: FAITHFUL**

Evidence SHA-256: `03a458d92d27bf862f737bbfb02f0c031e46cc2bc7bb9aaa54d6b902f6f173c7`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `55c1e8794345b79379ae26b1e226cac42e665e471d22d9a92a25c5e91a4df8e1`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `02987ce9d4ce8c4dc57a14a0d408d650226eeb7dd21ea1a1e5a653663c133ba5`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `0b5a7cc512f1d8d04117edf017d0d20330df9290c77712f941ce8333ae436034`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `fe9fc99bbcadcf3c0ab153c72a7934055ee0ef9259df39df2f58a44098320675`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `e56d68ee811202417b7f64866bcdc1e607da0b9a8ac10027f6f42c90c6039c9b`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `0656005b54177f29bd86efc1957759406a118c8fa4e016cfa31ed2958c47a2e4`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `c6ebfd7a2b7182dabc93d220f1e3c7988e0b6248c98b48d9b04fe11985ed3a58`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `4fdbbe9d507ac2351a35d6bb7e96a84b54226701231c63bca852b199a3b2d743`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `11b7713dcf579ab015c659a4be1466c325c9464ed3035baa405da2fabfdc3de4`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `d4b5add10d9b6ff383d7f7cdecd094110417b798a8cb86d6a3a80f1fe6d1a966`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `7e72f6e2ef3c2e63b5f50e5261d535a8a14f8026866bbf9491ae1a64428ceae3`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `90e7b45ca58b00df9d790ad88f08bed7139fd68dba80deee99ddcc1d4d141389`.

Go 1.27.0 resets an exhausted buffer with a positive cursor before testing spare capacity. That reset remains visible if a negative length subsequently panics or if the explicit growth guard panics with ErrTooLarge. For a nonexhausted buffer the signed reslice can shorten storage below the cursor without changing the previous read flag. Successful growth returns the insertion index and exposes the reserved bytes. The model's use of the public Grow model on the slow path reconstructs those bytes from spare storage after the shared growth branches. [Model](../../../lean/GoLean/Bytes/BufferGrowPrivate.lean)

The native overlay and Lean matched all 43,165 boundaries and 60,000 samples from seeds 17, 2718 and 65537. Ten independent fixtures matched both paths and distinguished reslicing before reset, rolling back on panic, and hiding the reservation as public Grow does. The cases exercise retained aliases, in-place compaction, new zero-filled allocation, and native minimum and maximum counts. Requests beyond 512 bytes except for the explicit overflow guard, allocator failure, and invalid starting geometry remain outside scope. [Comparison](../../../difftest/bytes/buffer-grow-private/check.go) · [Reviewer control](../../../difftest/bytes/buffer-grow-private/review/main.go)
