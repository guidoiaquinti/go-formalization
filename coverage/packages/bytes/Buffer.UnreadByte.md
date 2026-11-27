# Independent review of `bytes.Buffer.UnreadByte`

**Verdict: FAITHFUL**

Evidence SHA-256: `4cef4c35ab83e896296761479ba4cee4adcf3f2144128c5e5d546ebb71611019`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `edcc00658ede6ab1e180ead9fd404d75bc0af0be97e91f56de5858f189e6f7e8`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `8639714924bb2727f2d42a43d6cc35f8c90a32216182dc0f08ff9bd4948ded7c`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `9feaf6bedc8e243f8224f1f820b0cc2ef539754b94ae58ce4fc771bf39d118c5`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `7900511daf25efec5c5b04be742f7e8ee0b65d589888c2a3f4b6fd543951fd08`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `51ce239840692e29cb72fdc0fe78b1c9c34da600c044c49c61306ced6b5f6e50`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `143eafd064d2e517b54c6b54dc9423e9c8dd9f4a440d3e3fc2b0c2ed3480b630`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `3b434edfa90dae876bdc51b44f48e6bdea0e9754f89a04be0b4a485aa671fdef`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `17185649d303ce540bac68e0da1935ea05edc7ef75ac0e9d0f9125200d930d24`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `ccbc0c3d055be1ea728b97b1d9bcb334eb16d2a5be818e15677c6a36bcafe943`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `1a2632860750bfda8fea77e987116a235efd8e3b5c7b3e45846643319c7875fc`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `b049c0a536b0663ccac32e6231cb5ffc193e58ba68726a56453df303be341ddf`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `d9b7fca19e385e938bb8f8b883ff39d462c7ef908814f97d396f66c185bdd140`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `7fbcb10ead1bce69fdf4cab38fb15bb00ab8c4092154da1aa60c4257c42089cb`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `b632805c574f3598d1d39d31d332f8f3ea05982d518ac5648afc190d3e71379a`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `3bb079cb80f69ed88fd54b430348061930a705d1f240f827d82f76071277ddad`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `e395e3d490505e02007a72ae21c9359ea66bcc9299153c9a04b5c391be245ce8`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `364588dc95c3c7deea45d0dcb6a2d8a26f2e7562ce589eed09e25fbd259c219d`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `74bafe40917c1313370725225dde7f3820331bdd09fcaccd0463d5ff66e459a8`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `3b64c5fae2720af5aaf2a823d7d97a0909537d56c379d2849378711be846ab0b`.

The model follows the pinned Go 1.27.0 implementation. Only the zero read flag rejects the operation. Every other signed flag permits it and clears the flag. A positive offset decreases by one and zero stays zero. The nil receiver panics. The Lean theorems cover these branches, while the Go overlay checks the exact sentinel identity and observes retained storage, effective capacity and slice nilness. The 33,281 boundary inputs cover all signed byte flags against every valid offset in a 128-byte buffer, a nil backing slice and a nil receiver. The 60,000 sampled states use seeds 17, 2718 and 65537. All matched. [Model](../../../lean/GoLean/Bytes/BufferUnreadByte.lean) · [Comparison](../../../difftest/bytes/buffer-unread-byte/check.go)

Eleven independent fixtures matched both the overlay and Lean. A native read of a two-byte rune followed by UnreadByte left exactly its last byte unread. Repeated rejected operations preserved the cursor and returned the same live error as the initial rejection. Three incorrect variants each failed two fixtures by rejecting negative flags, rewinding the full rune width or rejecting a valid flag at a zero offset. The review covers valid slice and cursor states with all signed flag values. It does not exercise memory corruption or concurrent mutation. [Reviewer control](../../../difftest/bytes/buffer-unread-byte/review/main.go)
