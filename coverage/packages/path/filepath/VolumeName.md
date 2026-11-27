# Independent review of `path/filepath.VolumeName`

**Verdict: FAITHFUL**

Evidence SHA-256: `2e4028e7d10861e4671b9dd1bd6d0854f51379a5213757c37baf386688044211`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `1ed075e817e99415ab85df1986872579bb045b494783f569c10c95864a4a0ea1`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `5edb1a3f2bed523d7936580f64c38caa882444af7bcc3224640a3176b66dd7df`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `b1736ec3ff00be34e6d7fb4b3d1a4e5c719c33e2d8a44390ef6ba256d0d4f12e`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `392611c30eaa26043e378136b98d4e229744a33ac69a9821a7df491bb890dbdf`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `2ff160d1f30fc33d9db8d32270976036d3c65205b766e2210a1b35cbf111796e`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `43be190952f1fa1bffbd0a0e56bd3a76f483a3e79bb8ac6c134fce5f01f28e8e`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `81328794f3bf5ec6f291a1dafebde9e42126a37fc621aa642178fee075f74a0c`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `d7fe19de8bb1d3f8a3334fa7482326dd895340d393738746104366c6042eee85`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `365bbf65938b8cb2ea42aea2ef9033b5f4f9ebaed4776556a05a688d90f39a57`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `24159daf41b96eb4a9f00f9205c662901aeacca497c3c69821238da315fc323e`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `18c312b51ee8693eecb189f5df490fb2b4ad4623358266d8913b50cc369fe65a`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `5e4656387dbea65b0f8e5d1ba1ff63cdbca43a11301b4529a23f54e7b7dc63bf`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `4297e03ac178b2c64a7ffa7281974d6af0a52d2adaf7078be6444bb1d6b64467`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `a3d20243e79e5b814c59f65a514f46b3ca87a3dccac78095a1bdd4d70dd04f34`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `e5982e121ba1806d510b8b6fc08914c1dbaaa47c58c0c27d534b866606059664`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `a82c815c8b0e511c3e0b5031e36d9cf25af6960f2afc5f35b7d6993620c0123d`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `a51aa90e68ccfd27e0ac81216d006966c39fc32aa7fae029245f026d5a76f931`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `56f8bdb8043c50a62696a121fb2cb45fc752c804ee5351cee9ed3bb573e08fd4`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `c95a7753ecfeaad4746005b56b3730e23283f3ceddbcf1b28e16a265dc8839f5`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `6759d29dbd74c89df4039a2e517c58de677be1a5f756e52b5a04c05b6c3f2dfb`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `bc23cb189108d153b771e9200bc932d57f464238676ee28a0fa8b7510074e0b3`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `ab89bafce6003cd22cc862640e669f37460d6c15e23dfe74d4e1b64d399f2ddd`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `0e3110858cd24b3b86e398e8805d6f771a83372c7ac0a87c20b1dec78e66616d`.

Pinned Go 1.27.0 returns the empty string on non-Windows systems.
The Windows source first accepts any byte followed by a colon.
It then checks device prefixes using ASCII case and equivalent
separator bytes, or a leading UNC separator pair, and selects the
specified volume components. Device UNC recognition requires a
component boundary. The Windows validity check excludes an exact
parent component inside the volume length. A parent outside that
span does not affect the result. The wrapper converts slash to
backslash while preserving other raw bytes.

The author comparison checks 534 boundaries and 60,000 sampled
paths with seeds 17, 2718 and 65537. The independent control passed
20 fixed native and two-profile Lean cases. They include nonletter
and malformed drive bytes, mixed and repeated separators, device
and UNC prefixes, case-folded UNC and a nonmatching continuation,
parent versus longer dot components, and embedded NUL. Enforcing
drive letters, allowing parent volumes and keeping raw slashes
failed 3, 3 and 11 fixtures respectively. Native Go ran on
darwin/arm64. Windows was checked by pinned source, explicit
expected volumes and Lean. Native Windows execution remains open.

Run the installed control from the repository root with
`GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/path/filepath/volume-name/review`
after building the model executable.
