# Independent review of `strings.splitSeq`

**Verdict: FAITHFUL**

Evidence SHA-256: `6a40215afe0776a0aa3e9d7b8ca6daddf30973d1d5270b3bee2a566ad942e2f5`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `288c9d8ac51b956235c438cf6061315447d80a7260a14ba9ecfa0c2248a01f83`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `c5c1509e80d820f209ae1cc32530607ece4ff1bf37ec8d5def7bcfea15dad662`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `cf78f4103705d5fe459fb02aef3e24d4bea3ec0b7d03f317ba0b0803ad61b956`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `f83a6d8ad03cfa8703960473de1ceb2d1b5cff8f5a879b105fb3094acc450622`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `110fc74106426d7136c75917bd338a26705422ed4c80506391197c553df36249`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `72a74a3072c20b16e2688851682facbe1fd56bc377e15b1893a73a4670dd93b9`.

In the pinned Go 1.27.0 source, a nonempty separator selects the first byte match. The signed save value controls the yielded endpoint, and the full separator width controls the advance after a successful callback. A false callback returns before changing the captured string. A bounds panic similarly leaves that fragment unconsumed. An unmatched final fragment is yielded without changing the capture even if the callback returns true. For an empty separator, Go uses UTF-8 decoding widths and never reads the save value. Under legal native string lengths, a positive mathematical endpoint that overflows signed int wraps negative and still panics.

The author checked 16,000 boundary inputs and 60,000 sampled traces based on seeds 17, 2718 and 65537. Seven independent literal controls passed pinned darwin/arm64 Go via the actual private helper and Lean. These cases distinguish partial inclusion of a two-byte separator, a negative save and panic after an earlier yield, maximum signed overflow, save beyond the separator width, malformed UTF-8 and callback resumption, final-fragment persistence, and an empty input and separator. Scratch Lean mutations that consumed a paused fragment, ignored save, or consumed the final fragment failed three, four, and two controls respectively.

The comparison covers sequential use of a single closure and a slice bounds panic class. Callback panics, reentrancy, and concurrent invocation remain outside scope.
