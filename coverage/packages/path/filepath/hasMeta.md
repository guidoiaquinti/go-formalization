# Independent review of `path/filepath.hasMeta`

**Verdict: FAITHFUL**

Evidence SHA-256: `45426e1ba344873b37380d230a43102110b8f31b3e1a6e204fe212c5094ec5bd`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `232ec3433ce4b24291cd05a029546c6dd29775c0401a95f6553128c3a3744ce7`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `11e97b2430c8aa0631479d4bece1d19f3be84e15ed9e1836e97018717200ea2a`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `e1c11c2a6a566e9f120c20f360f334e2aaf5ef6bd32525f4ad7107a52b8ab311`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `b827b8d8f03794dc6680b9cc7ee67d0453d92000ea45399a5b6be9f821fdfd35`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `01aa4a2734baf530075eb244df7ce977afd8c3aa83cf0a06cd7904ea51622070`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `162d1bddb91613c5a96b8d5d9ef59fc35cd6d86c498c272b054b4d818e1c7b03`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `8c03770560b1909bcbecc9cfdb503b959683691c7f7873ad4898010b8c071a99`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `0b87e2613e60e2bff4accf8621f1abb4ad4dde50e1e90a1aaffad405413695ac`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `257df90dc7d1331bafe71d0e091f3270c35e7e2351e28925f0ccc966563f4652`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `dcf82aad7d9b5e4bed8ca06c80d8f69601a38df37034c65d7ff9211367515c86`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `a6dc42103fe7551da2cf97b0b85aaee3f6b4640898ef2905664630547052a8e4`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `dad940c493ff9506912b9a848f52a587995be944f56bd00034c004f299a5687c`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `0d81a171080ce28ce45705a41f8541cbef677aff10642925308a455a6df9ec08`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `c77111fde8ada8a8a0e63d0fe22d3f0b34f534206a4ee2485c2e5a9f2ad7766c`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `68a77dc8b9f081979e3df8f728f15d51de57106cafac8071e00f64def5350f55`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `07a9cd347b3faadf01bdf0118f017d85e682041926fb079d45fd5651b25578af`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `6195f9c8784c3d1ceba5c4dd2381367e2ef56a688d19bd9657f2eaa003414326`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `2fb78fcc6e8a91f0d2646a6dad6e6cb3b16da213666b9c99790c4a67363dd360`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `454ca9bebc431727f4d1ffd3f17c71a1dce9848f3b0c30a72ca916ea1e858699`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `5d656ce577d1220feb1a9c1d09ddeae949cd50db390aa971010f7e4a77746bc4`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `d8e4852b114f23d70cf964dd9ebb94f3d35df6a08bff6af82066447409da6bb5`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `7c0296f43f85d9b5c47f2b2c359b4dcc1905c57fc72497137f94f7a8f8f1ec35`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `5ff6fbdea4d8590b42d861f5eda6e7c94e15622af5cd776315877d63b23daf88`.

Pinned Go 1.27.0 checks whether the path contains any of star,
question mark and left bracket, adding backslash on non-Windows
systems. Right bracket is absent from the set. The Go Unicode
search yields the same answer as the Lean raw-byte check because
all four accepted values are ASCII and cannot be embedded in the
bytes of a non-ASCII decoded rune. A malformed byte does not mask
an adjacent ASCII metacharacter. Backslash is a set member only
on the non-Windows profile and does not hide adjacent members.

The author comparison checks 520 boundaries and 60,000 samples with
seeds 17, 2718 and 65537. The independent control passed 13 fixed
direct private Go and dual-profile Lean fixtures. They cover right
brackets, lone backslashes, backslashes next to actual meta bytes,
ASCII meta appearing among malformed UTF-8 and NUL, and Unicode
lookalikes. Including right bracket, dropping the non-Windows
backslash and interpreting escapes each failed two cases. Native
Go ran on darwin/arm64. The Windows result was exercised in Lean
and checked against the pinned source. Native Windows remains open.

Run the installed control from the repository root with
`GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/path/filepath/has-meta/review`
after building the model executable.
