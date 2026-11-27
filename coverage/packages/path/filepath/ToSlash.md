# Independent review of `path/filepath.ToSlash`

**Verdict: FAITHFUL**

Evidence SHA-256: `a9b6685f9c98d374144751df8fb35e8fdadb41a73f455f78c2b96443e3c7a1b1`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `154496ab3a9620cfbe76f82b23ce02b1127507ac6d9b5cc01b2e341e35684639`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `412eeca98d56c8bb3fa0c6ab15d9a0a68a45e54dd1a8fa8586c83ed33a2aa672`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `416e34ec05e45500ce71a194889c1caba7f80be36fe853b08c032aa5797791aa`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `aa6066eaa5cc9a8c2e4d2d45bf125ff637e95df60cbae5afdfb83d71dcf01f01`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `01b91dadebc47cc864de93bf55566a665cc74f6c941f565233173782fb3ca9cb`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `9353e66e83e4982ad619a911564cf6325be80829224a478c41a7cc2e3a651d33`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `90250812c311cbd7843591836acd74d83efbe123f2305052b5a5f31510d54394`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `993f3cd923ba8287e154635e90a9f8937eaf3f34dcf23217509e2a82342a9e99`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `92f2590be871b5341d824957d225e5b23f38db1dbd2d5229c5f6ddc55acd88ae`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `13e79f0b125fb35ac4f1a23da0d6a3cabe4a8deb7cc5cf6f01b336b127160ad5`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `98dd157ecc8a79cc4523f11dc5747080b8725badc3d02a8a1cbf1151fc151bc7`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `399a9aa13de623e6462df901359b166b5e3c39d1cd5a8b07a847f96a98e84373`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `1d12b89145199c4f30b2eab71785dacc4e8bf86f69e271248ab43c5c8b13b014`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `25b8cb0fab5cf13262cf96b0b37a85172a275c2951d2e1c13cf3e6ee9d930adc`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `6ee8511b45fde9787c130f52657a03c735f56345e0edd174e20d7ae50304079c`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `2f3e6051008e22a203c29ce70faa394cdfc47ac547f64150120f3529fe50331b`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `6163e2373df0915a8b007404154c320ae059dec0b1a073729f725d53b322cdc1`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `e7de4f4a4d0505d34674ba245230e1a04a64b58aa300fdc8b4037464b8d012a7`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `a94b37ab440225466020d064411a033db67cc180a1e012b5f75170eb292257be`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `af8a5882ebf7419a485a17045cad5084715a4c1fabfaabe8a23c5b6661a16c9d`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `30ac8f7b13b7549eeaa21149feb76c15cc8a65b43d1855464a76c8857c174dc4`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `93ec05c13668b12da818c553108c1366dae94915632298ac64b1077596a2f452`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `4261cda5220b66aa408876b9d32a6ac90d6b5261cadd9b375031e99d619f4b20`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `8226e62be076fee25102b07e52e5172aa4a6ca6e5c171a318a8fc732624da544`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `09c1ce6e8116bd917c9e969ec8e88d5fb10becd4747f133caac808c393666ae5`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `f4f15d12c258c83bbd116249946c01f2348bb4acd9836dfd2de5dcbae0cdc253`.

Pinned Go 1.27.0 filepath.ToSlash delegates to filepathlite.ToSlash.
With slash as the native separator it returns the original string.
With backslash as the native separator it replaces each backslash
with slash, preserving length, repeated separators, dot segments and
every other byte. The Lean model proves slash-profile identity and
length preservation, and its byte-step theorem specifies the
replacement. Backing identity is outside the returned-byte claim.

The author comparison passed 520 boundaries and 60,000 samples with
seeds 17, 2718 and 65537. The independent control passed ten fixed
darwin/arm64 fixtures against native Go and both Lean profiles.
It checks drive-like and UNC-like paths, dot segments, repeated
separators, embedded NUL, malformed UTF-8 and slash lookalikes.
Cleaning after replacement failed six Windows-profile fixtures,
and sanitizing malformed bytes failed one. The backslash profile
was checked through the explicit byte oracle and Lean executable.
Native Windows execution remains open.

Run the installed control from the repository root with
`GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/path/filepath/to-slash/review`
after building the model executable.
