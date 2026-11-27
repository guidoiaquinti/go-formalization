# Independent review of `path/filepath.walk`

**Verdict: FAITHFUL**

Evidence SHA-256: `8dbc12de1ab091ede24c9cb8b35e24be29009acc2414a98afcefc8a57cb739f5`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `7de2d99e970afa930f347889aaa29f77ac14d0ee7c56f3735ac952a17e436f11`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `cd7d189eacae2b8d929e40dc296330795f099160f887fa1b4ba6624067bfd0eb`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `15aa0cf9de43ecef1a1060e03673a8d657ceb00769526ba6303d6e0834b8d2fc`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `d56062a34f78594b0691369de4d8e944238acde8a8c3143f2a3c639a1e6a4a7e`.

Reviewed e0fc147 against Go 1.27.0. The private walker lists directory names before the first directory callback. A listing error is sent to that callback along with the existing directory info. The callback's result controls the return and traversal does not descend into a failed listing. A child Lstat failure is reported with nil info. The parent consumes SkipDir on a failed child lookup or a directory child, while a non-directory child's SkipDir is returned through the surrounding traversal. SkipAll and other errors propagate without top-level normalization in this private helper.

The independent control tests six model observations and eight live native error scenarios. Removing a directory after capturing its FileInfo forces a listing failure in the private walker. Removing a child in the successful parent callback confirms that the child name was collected before the callback but looked up afterward. Nil, SkipDir, SkipAll and a distinct error are tested at each failure point, including exact callback arguments, result identity, and whether the remaining sibling is visited. Reading the listing after the parent callback would miss the removed child entirely.

The complete replay passed all 4,096 callback maps on the fixed tree, 60,000 seeded repeats under 17, 2718 and 65537, and eight supplied filesystem-error rows. The model takes sorted listing and Lstat observations and a sufficient traversal bound. Native Windows and Plan 9 remain outside the local matrix.

Run the reviewer control from the repository root after building the Lean executable.

```sh
GOTOOLCHAIN=go1.27.0 go run difftest/path/filepath/walk-private/review/main.go
```
