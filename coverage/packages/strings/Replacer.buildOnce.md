# Independent review of `strings.Replacer.buildOnce`

**Verdict: FAITHFUL**

Evidence SHA-256: `04f00615009785657a8f5c0e99ac873ac6e857ed964bd52ce6830d1242a70130`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `817f2d2f7f69cf678e7802ad04771da1116337447bc5f4b8f6df0aff303dcf9e`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `3fffcae17d515edba79e873563038659bf9ba4cfa2b2bbc49e61eb28fd5aa427`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `2b9b247f6dff9f8e1bacae4b1e3f958e6e6e4e421f16871afe8e7a7eece6f278`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `b523b10d73b80276bc123edb938a0aff15cbdf09151ba5196da084d3eb29dd1a`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `a06d4020b5c66a40749eb83e9275fd93730b893206021ccbb494c135c29b52a6`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `4b4fabbe81a1c9c4bdc44d9efdd1fdcf78d5985e205f032bce8b136c3e7f0521`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `76d462fca8f701573d8a288b8bbeb030ecdd2c74c015952954af1be5c63d4738`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `caf731024a09cc4103927a3de326e901a5cda1d3805cab0a8057859b0e59b34e`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `13325ef5b862a8fd4944b1c82c90cc1030bc3574b7e07b511f5fde3836f9a43f`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `29831bb2f65ba641483b1d9e12a0b41964ab8888a609fa5d2c9249d23df3ca87`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `6cce886d4e0533abbca529551492491513b15521ab745ee64a711b251b955b18`.

Pinned Go 1.27.0 assigns the implementation from `build` using the pending arguments, then clears those arguments to nil. The callback has no `sync.Once` operation of its own. The model delegates construction to the reviewed route model and retains the supplied zero or completed marker. Calling the callback again builds the byte identity table from the cleared pending state. A direct call on a nil receiver panics.

The author comparison exercised 1,098 boundaries and 60,000 samples with seeds 17, 2718 and 65537, comparing route-specific state and the native Once observation. Six reviewer Lean fixtures check nil receiver, single replacement, unchanged marker for zero and completed observations, both repeated-call states, and duplicate malformed-byte precedence. A separate pinned private overlay in darwin/arm64 observed the callback's installed implementation and confirmed that public `Replace` subsequently rebuilds identity if the Once was zero, but retains the installed replacement if the Once was already completed. Two direct calls installed identity from either starting state.

The review covers successful allocation and legal constructor pairs. Once internals, concurrent calls, and malformed manually constructed private arguments remain outside scope.
