# Independent review of `strings.byteStringReplacer.Replace`

**Verdict: FAITHFUL**

Evidence SHA-256: `1f7668c9c61c9c684942a8efd269e02f3a6b43a6bea5924c37e368ac2cd2218d`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `6523cf99c3a5b947c3d2734ce0bbee5ec03bb6ade5fe48984f7f86aad0039914`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `8682c4f077c07dd47a1abbdafad4c1dfaf7181d230f31edfcb59a59ce180a57b`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `c3ae1e717d36b30c0e09cbff49ee1bdce416fc9d5afd51c21c3ddb36d1b911de`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `5a4e0792c5ad2de0e987d10b0599fa5c5d122ed35254d8dcf9202dc1c14f0cba`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `2b6506cb0f651e9dfa5cd53b9bb4c4aa40a5fea0062e259103e7a25991cfbe4d`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `9e4d9055c44826872f087e19dc7c9541e49e4a6b98bcc5ad7ae020a549928dc6`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `c16943aa622e3614f8296f21741b28e82114f081b434aa731dfb1fa7f44d8682`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `80b73d61ac77d7fb2c736971550dad8e06048702adad4f80428edcf6b4997768`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `acda09d935008475039cb0d19958d9f6954839b3b0c724ad5172ecfd455b1f68`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `7d5ffcab5eddf695828de663ab581d89f89a047039a528d7e2a0f0dec38c2e3e`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `19dd216f65e047b22a9c1a0c9a7318652fc904ecc50b99e96e3035cbe8fa5705`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `68b55c82a10624382f1e548631056166860def1d6a4c85cfcb6fcde8ea12199f`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `2ba4eea62233e3723a4904e713be91269bfb9e6ad69ba40b241836fda815c5e1`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `ecf53393656381929d687edfbbe43008008173d493fd9c83818a12bbc399d8b4`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `033dacc6ef2f82238c9c66a76d445815cd0ac29103dea9fb042855bb27cf8539`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `081689144b64f8a9fdfdb712ab8f1e27698decc72b5dd8e0461d205c0e36ee48`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `a024fa8cdcd39952ab5be8076c4f4c7d72877d5dd458992165c612c92641c990`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `39e889e1fef9ac5228125abf51c3a3204b658cbb1920068bd461ec746f1c95cf`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `b8285004e834305c04ba18651f8f1f3c0705d2beb12e7aecd54c732000d56cab`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `3a8f1ebcfd9272b5fa7c8182e329facc02b29bc8cc886de8fe95603c40897c52`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `829b22693e7ea65b4517d58996acdc30e87955b251d807010fa3600439a1f7b5`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `d514f6ecb0bf18f2062467feeb2deb616341edf4f4cb77cfe727bf5ab4e313ba`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `fed532ee33fc7d0c9a9c93edd6d724d8c9ac6cf727ebb4fe1d978c46847c3cb4`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `7ac830d4ead3ec06b030aa5b2f5f69f41808dfd7cceadfac745bd60dc5941044`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `8f0ca635d8f5fe10719a82772139e2a7033026ec4c4e16c88971acb7ecc41cc3`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `80dd70882ec62a706fa1e939783c285f8a0ed92d99e2917bca769bb64c7a6f2e`.

Pinned Go 1.27.0 derives the output length by counting listed keys at
input lengths at least eight times the number of keys. Below that
threshold it walks the bytes and checks the table directly. A present
entry with zero length deletes that byte. A nil entry leaves its
original byte in place. Replacement bytes are copied without another
lookup. The Lean model preserves both size passes, the no-change
return and raw byte construction. The claim assumes a consistent
table with unique listed keys and a representable output size.
Allocation identity and failures are outside scope.

The author comparison checks 7,770 boundaries and 60,000 samples
with seeds 17, 2718 and 65537 using direct same-package calls.
The independent control passed 12 fixed private Go and Lean cases
on darwin/arm64. Fixtures include typed nil receiver on empty and
nonempty input, present-empty deletion, raw malformed and embedded
NUL input, no mappings and no matches, and both sides of the exact
15/16 threshold with two keys. Interpreting empty replacements as
absent failed four cases. Reprocessing generated bytes failed three.

Run the installed control from the repository root with
`GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/strings/byte-string-replacer-replace/review`
after building the model executable.
