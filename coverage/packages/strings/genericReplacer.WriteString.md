# Independent review of `strings.genericReplacer.WriteString`

**Verdict: FAITHFUL**

Evidence SHA-256: `94f0a348c1689b8660ed09b7480f8d6c089d6c58a8b39e0383c485aadfe017d0`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `19c8fabe65c7286eec931d6b1d8d5e188b3fb3fa9fe575e7c996e75f7fabf670`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `84fa18e9d88bbac3cc3e793e0f525246c0099c2948585202b4f0dd21a2129033`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `1b17728f976b4ef50773d2f97f28a0edde900b1196fc6999bfc7105bfca7bb15`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `7edf00b2a4b3dda85ee9f821e92bcc3ca219253210c985d9ac4f79c092c1328a`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `0a135e58bbfc79bcfa17cd65649ade0846e090ad589844b69873bf05f8097422`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `8a88de64122c0c945c314654c713de980073a0a09ae09f0fb4359712d895fd61`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `2b33518b4d7eb78b6d25fa7fafb4eb934b3313430969daf81790b286307b4d8d`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `0892d7144e4937376592701ef6f2be985d40f2571363947f52cab07148fbaa7b`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `c42735f82a8101b351a35c756358f6db75838c0d8547d38febc6cc8e906caf06`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `02234ff33ebb06ca2d775d5a364db31b5ca1f6dcb547cea77defd143d613d195`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `8416bb5d89516d459031aac4fdffcb74dc3df4aac8deca74c8dad7252686072c`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `2111acac943082b894fd142f3e2f6b88d1b98270063b86a0ca07450236dc4d22`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `ecbd73a72cb7dcb54b10080cac3dd63ca7924fb722396e294653c1a97c1f7421`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `2862aaf657125094439c7f029c28659e789a273f23061a24e0808706f41a639f`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `cd877a8472642e0b9868d033a998ef0679cfde164759d95c8b6b7369053b01cf`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `a2438359ef9789a1fc10742a5e4a16e8b730a0877a4cbb386d23fe15a3c1cf9c`.

Reviewed e9d494d against Go 1.27.0. The model preserves the two calls for every match, including zero-length unchanged spans and zero-length replacements. A zero-length root match enables one root-suppressed lookup at the same byte offset before scanning can advance. There are at most two such iterations at each byte position, including the terminal position. The fast root-table skip has the same observable result as a miss in the scoped constructor-built trie. The writer route remains fixed, count includes the first failing write's accepted bytes, and a short successful write does not stop subsequent dispatch.

The independent control checks ten native and Lean traces. It covers opposite empty-key priorities, final-position empty matches, malformed input scanned by bytes, an empty replacement returning an error, exact first and second sentinel identities, short success, and failures in prefix, replacement, and tail writes. Removing empty writes changes four traces; forcing the byte adapter changes seven. A trace with a zero-count success on a nonempty prefix still writes the replacement and tail.

The author comparison checks 1,592 boundaries and 60,000 deterministic scripts from seeds 17, 2718 and 65537. The claim is scoped to constructor-built tries, writer counts within payload length, representable total count, and the observable calls and error identity for nonpanicking deterministic writers.

Run the reviewer control from the repository root after building the Lean executable.

```sh
GOTOOLCHAIN=go1.27.0 go run difftest/strings/generic-replacer-write-string/review/main.go
```
