# Independent review of `archive/zip.NewWriter`

**Verdict: FAITHFUL**

Evidence SHA-256: `f159eeab60eb17cbaabdef5a614428a041dcb392c8dcf9ecfeb9b9967e49a6bc`

I reviewed evidence commit `2f6cd5a` against Go 1.27.0 `src/archive/zip/writer.go` and the relevant buffer operations in `src/bufio/bufio.go`. Pinned Go initializes the outer ZIP writer and count wrapper and delegates internal buffer creation to bufio.NewWriter. The native tests check zero writer state, no callback and exact buffer identity. The fixed 4,095 and 4,096 capacities distinguish allocating a fresh buffer from retaining the existing one and its pending length.

Pinned Go and Lean agreed on 672 fixed rows and 20,000 sampled rows per seed for 17, 2718 and 65537 on darwin/arm64 with `go run ./difftest/archive/zip/writer-io`. Four independent fixtures passed with `go run ./difftest/archive/zip/writer-io/review` and rejected nine incorrect observations across the group. The scoped kernel audit passed. This review covers lawful in-memory writer replies, initialized states where needed and native error identity. Archive generation, invalid delegated replies, delegated panics and concurrent writes are outside scope. [Independent control](../../../../difftest/archive/zip/writer-io/review/main.go)
