# Independent review of `path/filepath.walkSymlinks`

**Verdict: FAITHFUL**

Evidence SHA-256: `fb0773942b7d4106742a7284bfaad955da4a16e7b447374f66cfbf233541ddec`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `4bc107ebace185773ba149a789516b0f1a2de796b39048fa64e0e6e7125e5eba`.

Reviewed 5b51289 against Go 1.27.0. The walker expands links before continuing through following components. Relative targets replace the current name within its containing directory. Absolute targets replace the destination root. Parent components are processed after the preceding link expansion. The walker returns lookup and readlink errors, detects non-directory intermediate components, and refuses the 256th traversal before reading that link. The final path passes through Clean.

Six independent model rows and four native private-walker scenarios passed. A physical chain of 255 links reaches the target, and starting one link earlier returns the exact too-many-links error. The nested target in `alias/../file` resolves to `real/file`. Cleaning the input first would instead resolve `file` in the starting directory. A link to a nondirectory before a parent component returns ENOTDIR.

The complete replay passed 271 native boundaries and 60,000 sampled paths using seeds 17, 2718 and 65537, plus eight fixed model observations. The claim assumes supplied filesystem outcomes and a sufficient model step bound. Windows behavior was inspected in pinned source and exercised through supplied model observations. Native Windows remains outside the local matrix.

Run the reviewer control from the repository root after building the Lean executable.

```sh
GOTOOLCHAIN=go1.27.0 go run difftest/path/filepath/walk-symlinks/review/main.go
```
