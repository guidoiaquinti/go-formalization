# Independent review of `path/filepath.Base`

**Verdict: FAITHFUL**

Evidence SHA-256: `576806126f27d276a447c8cb6b7e35613fadb5274481b43e7434b826dad5970e`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `999d224d549a98ec9dd6bcf3b79f0a18060bef5a00516f4bc15b3107b8f534c0`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `d0baa14c3e88dd6544d559a1f5d639506fa39c0fc91ac0b117efdf530ac43fd0`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `30609da3443f0cbdf615e0e0ea81bbc0a2088486f15dc1ec6d7f8e8e313c139d`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `18df92bd9f5cba3e7315c345ab5c144d950ae12eb1d4a2f01ba563c540ba2e0a`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `0548a16f24ec9fa39e6fe8943024131c6de6c2e3589bd67de1e8d21a6444c9d4`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `4353d05beca1fb170d3f499c97dc0116e7514ab48c8664b2869b1005222f964a`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `40d842366801fc7a1335d885c110245a49f01ac2bd6d03ea1ad7aca65c8342fa`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `e8d717f73a1193f048cccb75f61d921f16ecc175ef52ea1297ec3518bb1156eb`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `41085b7a85302207fe344c9a8b07aa9e0d5a6aeda9c8f0504d9f3288cecedf23`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `69dabefc776e0b9f6846861c51a1d8c05d568ca5e577c9369c0ad8b9c2fa399c`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `d511afea84f96e90e60b3e460a2c8eea67fe3d80a8c55ce560c81147647129fb`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `3218afca8837d312eb4c2482d15dad21332aa7a329c660f421d8207b8f022206`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `7a21c62791e1a1eb5931212fe86aeb62b97cfc13983a22a127599906f0179a17`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `36dd3ef03a57e1bd0265534ccd56293c2761c8db086e5661140c587775bccb38`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `9ee0665e7de509957fd907f6e48efbf1e3cef63f8329c7f0fa294ed8cb095f68`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `9090b6f0088293f0ea8983fee7cde63bc386103e3361037a18d8f6cd47385671`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `445dbb3e35d13db9ec2e133be358f4458ed585324b808de998b1d2e89e91d11d`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `ced7dc3b87f5f6473734553d5d494f43a7d2812b0f1aba7795a049b9c33cb749`.

Pinned Go 1.27.0 first returns a dot for empty input. Otherwise it
strips trailing native separators before computing the volume length
on the trimmed path, discards that prefix, then scans backwards for
the last native separator. The model follows the same order. It uses
the byte length of the modeled volume, so slash-to-backslash rewriting
in the Windows volume model does not alter the discarded prefix
length. The original leaf bytes are preserved, including dot and
parent components. A path reduced to no leaf returns the platform's
separator.

The author comparison checks 528 boundaries and 60,000 cases with
seeds 17, 2718 and 65537, plus 16 explicit Windows fixtures. The
independent control passed 16 host Go and 32 dual-profile Lean
fixtures, including trailing separators on device UNC roots and
leaves, drive-relative and nonletter drive cases, differing
backslash handling, raw dot and parent components, and malformed
and NUL bytes. Treating backslash as a separator on Unix and
cleaning before extracting the base each failed four cases.
Requiring drive letters on Windows failed two. Windows was checked
against pinned source and explicit Lean expectations, with native
Windows execution open.

Run the installed control from the repository root with
`GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/path/filepath/base/review`
after building the model executable.
