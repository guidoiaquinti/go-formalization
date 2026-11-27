# Independent review of `bytes.Reader.UnreadByte`

**Verdict: FAITHFUL**

Evidence SHA-256: `53f14838ac98841a8dd14d41ca49c46df0df9d944bdb535302d53701e2aa3c1f`

Go 1.27.0 decrements any positive cursor by one and clears the previous-rune marker. Positions past end are not clamped and no successful prior read is required. At zero the function returns a fresh error and leaves state intact. A nil receiver panics. The direct overlay matched all 63,859 fixtures, including 60,000 samples at seeds 17, 2718 and 65537. [Model](../../../lean/GoLean/Bytes/ReaderByte.lean) · [Comparison](../../../difftest/bytes/reader-byte/check.go)

Nine independent paired Go and Lean fixtures passed. Clearing the marker on a failed unread and refusing a cursor past end each failed three fixtures. The public control demonstrated two successive unreads after a two-byte rune and checked that two at-start failures had equal messages and distinct identities. Synthetic negative cursors are outside scope. [Reviewer control](../../../difftest/bytes/reader-byte/review/main.go)
