# Scoped verification progress

The inventory uses Go 1.27.0 on a Linux/amd64 (v1) reference profile with cgo disabled and default experiments. Differential checks run on the native host. Each named function or method declared in a selected production Go file counts once, including unexported functions and declarations with assembly implementations. Tests, closures and interface method signatures are excluded.

Started counts functions with a recorded model. Reviewed progress counts declarations with a current FAITHFUL verdict for the recorded claim. The percentage measures review progress across the selected declarations. Each claim below states its scope and execution configuration.

| Target | Verification focus | Declarations | Started | FAITHFUL | Reviewed progress |
| --- | --- | ---: | ---: | ---: | ---: |
| [math/bits](https://github.com/golang/go/tree/go1.27.0/src/math/bits) | Bitwise arithmetic and overflow | 49 | 49 | 0 | 0.000% |
| [slices](https://github.com/golang/go/tree/go1.27.0/src/slices) | Sorting, searching and boundary conditions | 81 | 24 | 0 | 0.000% |
| [sort](https://github.com/golang/go/tree/go1.27.0/src/sort) | Ordering and permutation guarantees | 71 | 22 | 0 | 0.000% |
| [strconv](https://github.com/golang/go/tree/go1.27.0/src/strconv) | Numeric parsing and round-trip correctness | 136 | 5 | 0 | 0.000% |
| [encoding/json](https://github.com/golang/go/tree/go1.27.0/src/encoding/json) | Parsing and serialization properties | 469 | 0 | 0 | 0.000% |
| [net/netip](https://github.com/golang/go/tree/go1.27.0/src/net/netip) | IP address representation and arithmetic | 108 | 40 | 0 | 0.000% |
| [sync](https://github.com/golang/go/tree/go1.27.0/src/sync) | Synchronization invariants | 114 | 0 | 0 | 0.000% |
| [runtime](https://github.com/golang/go/tree/go1.27.0/src/runtime) | Maps, channels and scheduler invariants | 2,819 | 0 | 0 | 0.000% |
| [go/types](https://github.com/golang/go/tree/go1.27.0/src/go/types) | Type checking and generic constraints | 811 | 0 | 0 | 0.000% |
| [math/big](https://github.com/golang/go/tree/go1.27.0/src/math/big) | Arbitrary-precision arithmetic | 311 | 0 | 0 | 0.000% |
| [time](https://github.com/golang/go/tree/go1.27.0/src/time) | Calendar and duration arithmetic | 200 | 0 | 0 | 0.000% |
| [math](https://github.com/golang/go/tree/go1.27.0/src/math) | Floating-point functions and special values | 157 | 0 | 0 | 0.000% |
| [regexp/syntax](https://github.com/golang/go/tree/go1.27.0/src/regexp/syntax) | Regular expression parsing and simplification | 112 | 0 | 0 | 0.000% |
| [net/url](https://github.com/golang/go/tree/go1.27.0/src/net/url) | URL parsing and escaping | 63 | 0 | 0 | 0.000% |
| [path/filepath](https://github.com/golang/go/tree/go1.27.0/src/path/filepath) | Path cleaning and matching | 39 | 0 | 0 | 0.000% |
| [encoding/binary](https://github.com/golang/go/tree/go1.27.0/src/encoding/binary) | Byte order and variable-length integers | 66 | 0 | 0 | 0.000% |
| [unicode/utf8](https://github.com/golang/go/tree/go1.27.0/src/unicode/utf8) | UTF-8 encoding and decoding | 20 | 0 | 0 | 0.000% |
| [unicode/utf16](https://github.com/golang/go/tree/go1.27.0/src/unicode/utf16) | Surrogate pairs and decoding | 8 | 0 | 0 | 0.000% |
| [bytes](https://github.com/golang/go/tree/go1.27.0/src/bytes) | Byte-slice search and manipulation | 117 | 0 | 0 | 0.000% |
| [strings](https://github.com/golang/go/tree/go1.27.0/src/strings) | String search and manipulation | 118 | 0 | 0 | 0.000% |
| [container/heap](https://github.com/golang/go/tree/go1.27.0/src/container/heap) | Heap ordering invariants | 7 | 0 | 0 | 0.000% |
| [go/constant](https://github.com/golang/go/tree/go1.27.0/src/go/constant) | Exact constant arithmetic | 97 | 0 | 0 | 0.000% |
| [image/color](https://github.com/golang/go/tree/go1.27.0/src/image/color) | Color model conversions | 31 | 0 | 0 | 0.000% |
| [math/rand/v2](https://github.com/golang/go/tree/go1.27.0/src/math/rand/v2) | Random number generation and range reduction | 64 | 0 | 0 | 0.000% |
| [compress/flate](https://github.com/golang/go/tree/go1.27.0/src/compress/flate) | DEFLATE compression and decompression | 138 | 0 | 0 | 0.000% |
| [archive/zip](https://github.com/golang/go/tree/go1.27.0/src/archive/zip) | ZIP archive reading and writing | 105 | 0 | 0 | 0.000% |
| [index/suffixarray](https://github.com/golang/go/tree/go1.27.0/src/index/suffixarray) | Suffix ordering and substring lookup | 74 | 0 | 0 | 0.000% |
| [internal/bytealg](https://github.com/golang/go/tree/go1.27.0/src/internal/bytealg) | Byte search and comparison | 23 | 0 | 0 | 0.000% |
| [internal/dag](https://github.com/golang/go/tree/go1.27.0/src/internal/dag) | Graph ordering and cycle detection | 15 | 0 | 0 | 0.000% |
| [cmd/compile](https://github.com/golang/go/tree/go1.27.0/src/cmd/compile) | Compiler optimization equivalence | 10,458 | 0 | 0 | 0.000% |
| **Selected standard library total** | | 6,423 | 140 | 0 | 0.000% |
| **Selected compiler toolchain total** | | 10,458 | 0 | 0 | 0.000% |

## Package scope

- `strconv` includes `strconv`, `internal/strconv`.
- `encoding/json` includes `encoding/json`, `encoding/json/v2`, `encoding/json/jsontext`, `encoding/json/internal`, `encoding/json/internal/jsonflags`, `encoding/json/internal/jsonopts`, `encoding/json/internal/jsonwire`.
- `sync` includes `sync`, `internal/sync`.
- `runtime` includes `runtime`, `internal/runtime/maps`.
- `cmd/compile` includes `cmd/compile/...`.

Other targets cover their named package. Use `make inventory TARGET=slices` to list function names, verification states and pinned source lines for one target.

## Verification records

### math/bits.Add (AWAITING_REVIEW)

The sum word and carry reconstruct the unsigned sum including the input carry.

64-bit uint only. Go requires the input flag to be zero or one. The model uses widened arithmetic. The source bitwise flag calculation is tested on Go. Timing is outside the claim. Boundary cases and 60,000 seeded inputs (17, 2718, 65537).

Native Go with default lowering and with intrinsics disabled in the driver and math/bits.

### math/bits.Add32 (AWAITING_REVIEW)

The sum word and carry reconstruct the unsigned sum including the input carry.

Go requires the input flag to be zero or one. Timing is outside the claim. Boundary cases and 60,000 seeded inputs (17, 2718, 65537).

Native Go with default lowering and with intrinsics disabled in the driver and math/bits.

### math/bits.Add64 (AWAITING_REVIEW)

The sum word and carry reconstruct the unsigned sum including the input carry.

Go requires the input flag to be zero or one. The model uses widened arithmetic. The source bitwise flag calculation is tested on Go. Timing is outside the claim. Boundary cases and 60,000 seeded inputs (17, 2718, 65537).

Native Go with default lowering and with intrinsics disabled in the driver and math/bits.

### math/bits.Div (AWAITING_REVIEW)

Return the Euclidean quotient and remainder when the quotient fits. Panic for a zero divisor or quotient overflow, with zero divisor taking precedence.

64-bit uint only. Observable outcomes only. The model abstracts the normalized long-division steps. The compiled Go paths are compared by sampling. Boundary cases and 60,000 seeded inputs (17, 2718, 65537).

Native Go with default lowering and with intrinsics disabled in the driver and math/bits.

### math/bits.Div32 (AWAITING_REVIEW)

Return the Euclidean quotient and remainder when the quotient fits. Panic for a zero divisor or quotient overflow, with zero divisor taking precedence.

Observable outcomes only. Boundary cases and 60,000 seeded inputs (17, 2718, 65537).

Native Go with default lowering and with intrinsics disabled in the driver and math/bits.

### math/bits.Div64 (AWAITING_REVIEW)

Return the Euclidean quotient and remainder when the quotient fits. Panic for a zero divisor or quotient overflow, with zero divisor taking precedence.

Observable outcomes only. The model abstracts the normalized long-division steps. The compiled Go paths are compared by sampling. Boundary cases and 60,000 seeded inputs (17, 2718, 65537).

Native Go with default lowering and with intrinsics disabled in the driver and math/bits.

### math/bits.LeadingZeros (AWAITING_REVIEW)

Count consecutive zero bits from the most significant end, including the word width at zero.

Return value only. 64-bit uint only. The length model factors the threshold and shift steps into recursive half-word selection. Boundary cases and 60,000 seeded inputs (17, 2718, 65537).

Native Go with default lowering and with intrinsics disabled in the driver and math/bits.

### math/bits.LeadingZeros16 (AWAITING_REVIEW)

Count consecutive zero bits from the most significant end, including sixteen at zero.

Return value only. Exhaustive comparison.

Native Go with default lowering and with intrinsics disabled in the driver and math/bits.

### math/bits.LeadingZeros32 (AWAITING_REVIEW)

Count consecutive zero bits from the most significant end, including the word width at zero.

Return value only. The length model factors the threshold and shift steps into recursive half-word selection. Boundary cases and 60,000 seeded inputs (17, 2718, 65537).

Native Go with default lowering and with intrinsics disabled in the driver and math/bits.

### math/bits.LeadingZeros64 (AWAITING_REVIEW)

Count consecutive zero bits from the most significant end, including the word width at zero.

Return value only. The length model factors the threshold and shift steps into recursive half-word selection. Boundary cases and 60,000 seeded inputs (17, 2718, 65537).

Native Go with default lowering and with intrinsics disabled in the driver and math/bits.

### math/bits.LeadingZeros8 (AWAITING_REVIEW)

Return the number of consecutive zero bits from the most significant end of every uint8 input, including eight at zero.

Return value only. Exhaustive comparison.

Native Go, with default lowering and with intrinsics disabled in the driver and math/bits. CI matrix covers Linux/amd64, Linux/arm64 and macOS/arm64.

### math/bits.Len (AWAITING_REVIEW)

Return the minimum unsigned bit width, including zero bits at zero.

Return value only. 64-bit uint only. The length model factors the threshold and shift steps into recursive half-word selection. Boundary cases and 60,000 seeded inputs (17, 2718, 65537).

Native Go with default lowering and with intrinsics disabled in the driver and math/bits.

### math/bits.Len16 (AWAITING_REVIEW)

Return the minimum unsigned bit width, including zero bits at zero.

Return value only. Exhaustive comparison.

Native Go with default lowering and with intrinsics disabled in the driver and math/bits.

### math/bits.Len32 (AWAITING_REVIEW)

Return the minimum unsigned bit width, including zero bits at zero.

Return value only. The length model factors the threshold and shift steps into recursive half-word selection. Boundary cases and 60,000 seeded inputs (17, 2718, 65537).

Native Go with default lowering and with intrinsics disabled in the driver and math/bits.

### math/bits.Len64 (AWAITING_REVIEW)

Return the minimum unsigned bit width, including zero bits at zero.

Return value only. The length model factors the threshold and shift steps into recursive half-word selection. Boundary cases and 60,000 seeded inputs (17, 2718, 65537).

Native Go with default lowering and with intrinsics disabled in the driver and math/bits.

### math/bits.Len8 (AWAITING_REVIEW)

Return the smallest number of bits needed to represent every uint8 input, including zero bits at zero.

Return value only. Exhaustive comparison.

Native Go, with default lowering and with intrinsics disabled in the driver and math/bits. CI matrix covers Linux/amd64, Linux/arm64 and macOS/arm64.

### math/bits.Mul (AWAITING_REVIEW)

The high and low words reconstruct the full unsigned product.

64-bit uint only. Timing is outside the claim. Boundary cases and 60,000 seeded inputs (17, 2718, 65537).

Native Go with default lowering and with intrinsics disabled in the driver and math/bits.

### math/bits.Mul32 (AWAITING_REVIEW)

The high and low words reconstruct the full unsigned product.

Timing is outside the claim. Boundary cases and 60,000 seeded inputs (17, 2718, 65537).

Native Go with default lowering and with intrinsics disabled in the driver and math/bits.

### math/bits.Mul64 (AWAITING_REVIEW)

The source-shaped uint64 model is equivalent to the arithmetic model and its outputs reconstruct the unsigned product for every input pair.

The source guard checks the pinned Go declaration against the reviewed transcription. Boundary and seeded samples cover compiled Go and compiler intrinsics. Timing is outside the claim.

Native Go, with default lowering and with intrinsics disabled in the driver and math/bits. CI matrix covers Linux/amd64, Linux/arm64 and macOS/arm64.

[Review history](../coverage/packages/math/bits/Mul64.md)

### math/bits.OnesCount (AWAITING_REVIEW)

Return the number of set bits.

64-bit uint only. The model sums disjoint halves, and the individual mask and shift intermediates in Go are outside the claim. Boundary cases and 60,000 seeded inputs (17, 2718, 65537).

Native Go with default lowering and with intrinsics disabled in the driver and math/bits.

### math/bits.OnesCount16 (AWAITING_REVIEW)

Return the number of set bits.

Return value only. Exhaustive comparison.

Native Go with default lowering and with intrinsics disabled in the driver and math/bits.

### math/bits.OnesCount32 (AWAITING_REVIEW)

Return the number of set bits.

Return value only. The model groups the four byte-table counts into two half-word sums. Boundary cases and 60,000 seeded inputs (17, 2718, 65537).

Native Go with default lowering and with intrinsics disabled in the driver and math/bits.

### math/bits.OnesCount64 (AWAITING_REVIEW)

Return the number of set bits.

Return value only. The model sums disjoint halves. The individual mask and shift intermediates in Go are outside the claim. Boundary cases and 60,000 seeded inputs (17, 2718, 65537).

Native Go with default lowering and with intrinsics disabled in the driver and math/bits.

### math/bits.OnesCount8 (AWAITING_REVIEW)

Return the number of set bits for every byte.

Return value only. Exhaustive comparison.

Native Go, with default lowering and with intrinsics disabled in the driver and math/bits. CI matrix covers Linux/amd64, Linux/arm64 and macOS/arm64.

[Review history](../coverage/packages/math/bits/OnesCount8.md)

### math/bits.Rem (AWAITING_REVIEW)

The source-shaped model agrees with the arithmetic model at both uint widths. It returns the Euclidean remainder or reports a divide-by-zero panic.

Lean proves equivalence of the source-shaped and arithmetic models on all inputs, including the normalized 64-bit path. Source correspondence is checked by inspection. Comparisons with compiled Go are sampled.

Native Go, with default lowering and with intrinsics disabled in the driver and math/bits. CI matrix covers Linux/amd64, Linux/arm64 and macOS/arm64. CI also runs the native Linux/386 wrapper.

[Review history](../coverage/packages/math/bits/Rem.md)

### math/bits.Rem32 (AWAITING_REVIEW)

Return the Euclidean remainder even when the quotient overflows. Panic for a zero divisor.

Observable outcomes only. Boundary cases and 60,000 seeded inputs (17, 2718, 65537).

Native Go with default lowering and with intrinsics disabled in the driver and math/bits.

### math/bits.Rem64 (AWAITING_REVIEW)

Return the Euclidean remainder even when the quotient overflows. Panic for a zero divisor.

Observable outcomes only. The model abstracts the normalized long-division steps. The compiled Go paths are compared by sampling. Boundary cases and 60,000 seeded inputs (17, 2718, 65537).

Native Go with default lowering and with intrinsics disabled in the driver and math/bits.

### math/bits.Reverse (AWAITING_REVIEW)

Map every input bit to the opposite position in the word.

Return value only. 64-bit uint only. Boundary cases and 60,000 seeded inputs (17, 2718, 65537).

Native Go with default lowering and with intrinsics disabled in the driver and math/bits.

### math/bits.Reverse16 (AWAITING_REVIEW)

Map every input bit to the opposite position in the word.

Return value only. Exhaustive comparison.

Native Go with default lowering and with intrinsics disabled in the driver and math/bits.

### math/bits.Reverse32 (AWAITING_REVIEW)

Map every input bit to the opposite position in the word.

Return value only. Boundary cases and 60,000 seeded inputs (17, 2718, 65537).

Native Go with default lowering and with intrinsics disabled in the driver and math/bits.

### math/bits.Reverse64 (AWAITING_REVIEW)

Map every input bit to the opposite position in the word.

Return value only. Boundary cases and 60,000 seeded inputs (17, 2718, 65537).

Native Go with default lowering and with intrinsics disabled in the driver and math/bits.

### math/bits.Reverse8 (AWAITING_REVIEW)

Map every input bit to the opposite position in the word.

Return value only. Exhaustive comparison.

Native Go with default lowering and with intrinsics disabled in the driver and math/bits.

### math/bits.ReverseBytes (AWAITING_REVIEW)

Reverse the byte order while preserving each bit’s position within its byte.

Return value only. 64-bit uint only. Boundary cases and 60,000 seeded inputs (17, 2718, 65537).

Native Go with default lowering and with intrinsics disabled in the driver and math/bits.

### math/bits.ReverseBytes16 (AWAITING_REVIEW)

Swap both bytes of every 16-bit word and recover the input after two swaps.

Return values only. Exhaustive comparison. Timing is outside the claim.

Native Go, with default lowering and with intrinsics disabled in the driver and math/bits. CI matrix covers Linux/amd64, Linux/arm64 and macOS/arm64.

[Review history](../coverage/packages/math/bits/ReverseBytes16.md)

### math/bits.ReverseBytes32 (AWAITING_REVIEW)

Reverse the byte order while preserving each bit’s position within its byte.

Return value only. Boundary cases and 60,000 seeded inputs (17, 2718, 65537).

Native Go with default lowering and with intrinsics disabled in the driver and math/bits.

### math/bits.ReverseBytes64 (AWAITING_REVIEW)

Reverse the byte order while preserving each bit’s position within its byte.

Return value only. Boundary cases and 60,000 seeded inputs (17, 2718, 65537).

Native Go with default lowering and with intrinsics disabled in the driver and math/bits.

### math/bits.RotateLeft (AWAITING_REVIEW)

Move each bit left by the signed count modulo the word width, wrapping at the word boundary.

64-bit uint only. Return value only. Timing is outside the claim. Boundary cases and 60,000 seeded inputs (17, 2718, 65537).

Native Go with default lowering and with intrinsics disabled in the driver and math/bits.

### math/bits.RotateLeft16 (AWAITING_REVIEW)

Move each bit left by the signed count modulo the word width, wrapping at the word boundary.

Return value only. Timing is outside the claim. Boundary cases and 60,000 seeded inputs (17, 2718, 65537).

Native Go with default lowering and with intrinsics disabled in the driver and math/bits.

### math/bits.RotateLeft32 (AWAITING_REVIEW)

Move each bit left by the signed count modulo the word width, wrapping at the word boundary.

Return value only. Timing is outside the claim. Boundary cases and 60,000 seeded inputs (17, 2718, 65537).

Native Go with default lowering and with intrinsics disabled in the driver and math/bits.

### math/bits.RotateLeft64 (AWAITING_REVIEW)

Move each bit left by the signed count modulo the word width, wrapping at the word boundary.

Return value only. Timing is outside the claim. Boundary cases and 60,000 seeded inputs (17, 2718, 65537).

Native Go with default lowering and with intrinsics disabled in the driver and math/bits.

### math/bits.RotateLeft8 (AWAITING_REVIEW)

Move each bit left by the signed count modulo the word width, wrapping at the word boundary.

Return value only. Timing is outside the claim. Exhausts all bytes and rotation residues. Boundary cases and 60,000 seeded inputs (17, 2718, 65537).

Native Go with default lowering and with intrinsics disabled in the driver and math/bits.

### math/bits.Sub (AWAITING_REVIEW)

The difference word and borrow reconstruct the unsigned difference including the input borrow.

64-bit uint only. Go requires the input flag to be zero or one. The model uses widened arithmetic. The source bitwise flag calculation is tested on Go. Timing is outside the claim. Boundary cases and 60,000 seeded inputs (17, 2718, 65537).

Native Go with default lowering and with intrinsics disabled in the driver and math/bits.

### math/bits.Sub32 (AWAITING_REVIEW)

The difference word and borrow reconstruct the unsigned difference including the input borrow.

Go requires the input flag to be zero or one. The model uses widened arithmetic. The source bitwise flag calculation is tested on Go. Timing is outside the claim. Boundary cases and 60,000 seeded inputs (17, 2718, 65537).

Native Go with default lowering and with intrinsics disabled in the driver and math/bits.

### math/bits.Sub64 (AWAITING_REVIEW)

The difference word and borrow reconstruct the unsigned difference including the input borrow.

Go requires the input flag to be zero or one. The model uses widened arithmetic. The source bitwise flag calculation is tested on Go. Timing is outside the claim. Boundary cases and 60,000 seeded inputs (17, 2718, 65537).

Native Go with default lowering and with intrinsics disabled in the driver and math/bits.

### math/bits.TrailingZeros (AWAITING_REVIEW)

Count consecutive zero bits from the least significant end, including the word width at zero.

Return value only. 64-bit uint only. Boundary cases and 60,000 seeded inputs (17, 2718, 65537).

Native Go with default lowering and with intrinsics disabled in the driver and math/bits.

### math/bits.TrailingZeros16 (AWAITING_REVIEW)

Count consecutive zero bits from the least significant end, including sixteen at zero.

Return value only. Exhaustive comparison.

Native Go with default lowering and with intrinsics disabled in the driver and math/bits.

### math/bits.TrailingZeros32 (AWAITING_REVIEW)

Count consecutive zero bits from the least significant end, including the word width at zero.

Return value only. Boundary cases and 60,000 seeded inputs (17, 2718, 65537).

Native Go with default lowering and with intrinsics disabled in the driver and math/bits.

### math/bits.TrailingZeros64 (AWAITING_REVIEW)

Count consecutive zero bits from the least significant end, including the word width at zero.

Return value only. Boundary cases and 60,000 seeded inputs (17, 2718, 65537).

Native Go with default lowering and with intrinsics disabled in the driver and math/bits.

### math/bits.TrailingZeros8 (AWAITING_REVIEW)

Return eight at zero and the number of low-order zero bits for every other byte.

Return value only. Exhaustive comparison.

Native Go, with default lowering and with intrinsics disabled in the driver and math/bits. CI matrix covers Linux/amd64, Linux/arm64 and macOS/arm64.

[Review history](../coverage/packages/math/bits/TrailingZeros8.md)

### net/netip.Addr.As16 (AWAITING_REVIEW)

Returns each stored address word as eight big-endian bytes, high word first. The zero receiver gives zero bytes and IPv4 uses its mapped representation.

The model takes the stored words directly because the source does not read address detail. The proof covers every pair of words and every output position. Go comparison uses 133 boundary cases and 60,000 samples across seeds 17, 2718 and 65537.

Native Go calling netip.Addr.As16 on zero, IPv4, IPv6 and scoped IPv6 receivers.

### net/netip.Addr.As4 (AWAITING_REVIEW)

Returns the four IPv4 bytes for IPv4 and IPv4-mapped IPv6 receivers. The zero receiver and other IPv6 receivers panic with distinct messages.

Models panic payloads and returned bytes. The four-byte truncated store is modeled by the last four bytes of the big-endian low word and proved against base-256 digits. The Go check uses 100 boundary cases and 60,000 samples across seeds 17, 2718 and 65537.

Native Go calling netip.Addr.As4 and comparing its bytes or recovered string panic payload.

### net/netip.Addr.AsSlice (AWAITING_REVIEW)

Returns nil for the zero address, the four IPv4 bytes for IPv4, and all sixteen IPv6 bytes for IPv6, including mapped and scoped IPv6.

The proof covers nilness, length and contents for every address class and word pair. Slice capacity, allocation and storage identity are outside the claim. The Go comparison uses 133 boundary cases and 60,000 samples with seeds 17, 2718 and 65537.

Native Go comparing netip.Addr.AsSlice nilness and byte sequences.

### net/netip.Addr.BitLen (AWAITING_REVIEW)

Returns 0 for the zero Addr, 32 for IPv4 and 128 for IPv6, including scoped and IPv4-mapped IPv6 addresses.

Models which address-detail handle class the source switch observes. The theorem exhausts those classes. Address bytes and zone contents are irrelevant to the switch. The Go check uses five boundary cases and 60,000 samples across seeds 17, 2718 and 65537.

Native Go directly invoking netip.Addr.BitLen on zero, IPv4 and IPv6 receivers.

### net/netip.Addr.Compare (AWAITING_REVIEW)

Returns -1, 0 or 1 according to lexicographic order by bit length, unsigned address words, and zone for tied IPv6 addresses.

The proof covers modeled words and Lean strings. The Go check uses empty and decimal zones in 64 boundary pairs and 60,000 samples with seeds 17, 2718 and 65537. Arbitrary Go zone byte sequences are outside the comparison.

Native Go calling Addr.Compare on zero, IPv4, mapped IPv6 and scoped and unscoped IPv6 receivers.

### net/netip.Addr.Is4 (AWAITING_REVIEW)

Returns true exactly for IPv4 receivers. IPv4-mapped IPv6 receivers return false.

The model groups handles into zero, IPv4 and IPv6 detail classes. The proof covers all classes. The Go comparison uses five boundary cases and 60,000 samples across seeds 17, 2718 and 65537.

Native Go calling netip.Addr.Is4 on zero, IPv4, IPv6, and scoped IPv6 receivers.

### net/netip.Addr.Is4In6 (AWAITING_REVIEW)

Returns true exactly for IPv6 receivers whose first 96 bits match ::ffff:0:0/96.

The model uses the address-detail class and both address words. Zone contents do not affect the predicate. The Go check uses 100 boundary cases and 60,000 samples across seeds 17, 2718 and 65537.

Native Go calling netip.Addr.Is4In6 on zero, IPv4 and IPv6 receivers.

### net/netip.Addr.Is6 (AWAITING_REVIEW)

Returns true exactly for IPv6 receivers, including scoped and IPv4-mapped IPv6 addresses.

The model groups handles into zero, IPv4 and IPv6 detail classes. The proof covers all classes. The Go comparison uses five boundary cases and 60,000 samples across seeds 17, 2718 and 65537.

Native Go calling netip.Addr.Is6 on zero, IPv4, IPv6, and scoped IPv6 receivers.

### net/netip.Addr.IsGlobalUnicast (AWAITING_REVIEW)

After unmapping, excludes invalid and unspecified addresses, the IPv4 limited broadcast, loopback, multicast and link-local unicast. Scoped IPv6 zero and private unicast addresses can qualify.

The model compares canonical address details by contents. The check uses net.IP.IsGlobalUnicast as the value oracle with an explicit adjustment for scoped IPv6 zero. It exercises 31 boundary cases and 60,000 samples with seeds 17, 2718 and 65537.

Native Go calling Addr.IsGlobalUnicast on zero, IPv4 and IPv6 receivers, including scoped and mapped IPv6.

### net/netip.Addr.IsInterfaceLocalMulticast (AWAITING_REVIEW)

Recognizes IPv6 multicast with scope 1 regardless of flags.

Proof covers all modeled words, families and zones. Go comparison uses unzoned receivers in 2,049 boundary cases and 60,000 samples with seeds 17, 2718 and 65537.

Native Go calling Addr.IsInterfaceLocalMulticast on zero, IPv4, mapped IPv6 and ordinary IPv6 receivers.

### net/netip.Addr.IsLinkLocalMulticast (AWAITING_REVIEW)

Recognizes 224.0.0.0/24 for IPv4 and mapped IPv6 payloads, and IPv6 multicast with scope 2 regardless of flags.

Proof covers all modeled words, families and zones. Go comparison uses unzoned receivers in 2,049 boundary cases and 60,000 samples with seeds 17, 2718 and 65537.

Native Go calling Addr.IsLinkLocalMulticast on zero, IPv4, mapped IPv6 and ordinary IPv6 receivers.

### net/netip.Addr.IsLinkLocalUnicast (AWAITING_REVIEW)

Recognizes 169.254.0.0/16 for IPv4 and mapped IPv6 payloads, and fe80::/10 for IPv6.

Proof covers all modeled words, families and zones. Go comparison uses unzoned receivers in 1,537 boundary cases and 60,000 samples with seeds 17, 2718 and 65537.

Native Go calling Addr.IsLinkLocalUnicast on zero, IPv4, mapped IPv6 and ordinary IPv6 receivers.

### net/netip.Addr.IsLoopback (AWAITING_REVIEW)

Returns true exactly for IPv4 and mapped IPv6 addresses with IPv4 payload in 127/8, and for IPv6 ::1.

The proof covers every modeled address and arbitrary zone contents. The Go check uses unzoned receivers in 233 boundary cases and 60,000 samples with seeds 17, 2718 and 65537.

Native Go calling Addr.IsLoopback on zero, IPv4, mapped IPv6 and ordinary IPv6 receivers.

### net/netip.Addr.IsMulticast (AWAITING_REVIEW)

Recognizes 224.0.0.0/4 for IPv4 and mapped IPv6 payloads, and ff00::/8 for IPv6.

The proof covers all modeled addresses and ignores zone contents. The Go check uses unzoned receivers and 769 boundary cases, including all first octets for each address form, plus 60,000 samples with seeds 17, 2718 and 65537.

Native Go calling Addr.IsMulticast on zero, IPv4, mapped IPv6 and ordinary IPv6 receivers.

### net/netip.Addr.IsPrivate (AWAITING_REVIEW)

Recognizes the three RFC 1918 IPv4 ranges, including mapped IPv6 payloads, and fc00::/7 for IPv6.

Proof covers all modeled words, families and zones. Go comparison uses unzoned receivers in 277 boundary cases and 60,000 samples with seeds 17, 2718 and 65537.

Native Go calling Addr.IsPrivate on zero, IPv4, mapped IPv6 and ordinary IPv6 receivers.

### net/netip.Addr.IsUnspecified (AWAITING_REVIEW)

Returns true exactly for 0.0.0.0 and unzoned ::.

Models Go address equality using words, the detail class and zone contents. The proof covers all modeled receivers. The Go comparison uses 133 boundary cases and 60,000 samples across seeds 17, 2718 and 65537, using empty and decimal zones.

Native Go calling Addr.IsUnspecified on zero, IPv4, IPv6 and scoped IPv6 receivers.

### net/netip.Addr.IsValid (AWAITING_REVIEW)

Returns true exactly when the receiver is initialized, including both unspecified addresses.

The model groups handles into zero, IPv4 and IPv6 detail classes. The proof covers all classes. The Go comparison uses five boundary cases and 60,000 samples across seeds 17, 2718 and 65537.

Native Go calling netip.Addr.IsValid on zero, IPv4, IPv6, and scoped IPv6 receivers.

### net/netip.Addr.Less (AWAITING_REVIEW)

Orders addresses lexicographically by bit length, unsigned address words and, for tied IPv6 addresses, zone.

The model uses Lean strings for zones. Go checks empty and decimal zones in 64 boundary cases and 60,000 samples using seeds 17, 2718 and 65537. Arbitrary Go string byte sequences are outside the comparison.

Native Go invoking Addr.Less directly on zero, IPv4, mapped IPv6 and scoped and unscoped IPv6 receivers.

### net/netip.Addr.Next (AWAITING_REVIEW)

Increments the stored address with carry, returns the zero Addr on family-width overflow, and otherwise retains the original address detail.

The model covers all stored words and observable details. The zero receiver retains invalid detail with incremented bits in the pinned implementation. Go checks 16 boundary cases and 60,000 samples using seeds 17, 2718 and 65537.

Native Go calling Addr.Next and observing BitLen, As16 and Zone.

### net/netip.Addr.Prev (AWAITING_REVIEW)

Returns the zero Addr at the start of the address range. Otherwise subtracts one with borrow and preserves address detail.

The model covers all stored words and observable address details. Go checks 16 boundary cases and 60,000 samples using seeds 17, 2718 and 65537.

Native Go calling Addr.Prev and observing BitLen, As16 and Zone.

### net/netip.Addr.Unmap (AWAITING_REVIEW)

Converts mapped IPv6 receivers to IPv4 with an empty zone and identical As16 bytes. All other receivers retain their family, address bytes and zone.

Models the address-detail handle by family class and zone contents. Handle identity is abstracted. Go comparison uses 101 boundary cases and 60,000 samples across seeds 17, 2718 and 65537, with empty and decimal zone strings.

Native Go calling netip.Addr.Unmap and observing BitLen, As16 and Zone.

### net/netip.Addr.WithZone (AWAITING_REVIEW)

Preserves address bytes and family. Replaces the zone for IPv6, including removing it for an empty request, and leaves IPv4 and the zero address unchanged.

Models observable family and zone contents. Handle identity and interning costs are outside the claim. The proof covers arbitrary Lean strings. The Go check uses empty and decimal zone strings in six boundary cases and 60,000 samples with seeds 17, 2718 and 65537.

Native Go calling netip.Addr.WithZone and observing BitLen, As16 and Zone.

### net/netip.Addr.Zone (AWAITING_REVIEW)

Returns the stored zone for IPv6 and an empty string for IPv4 and zero receivers.

Models the observable contents of the detail handle. IPv4 uses z4, whose zone is empty. The proof covers arbitrary Lean strings. The Go check uses empty and decimal zones across five boundary cases and 60,000 samples with seeds 17, 2718 and 65537.

Native Go calling netip.Addr.Zone on zero, IPv4 and IPv6 receivers.

### net/netip.AddrFrom16 (AWAITING_REVIEW)

Returns an unzoned IPv6 address whose two words concatenate the sixteen input bytes in network order, including IPv4-mapped addresses.

The proof covers all 128-bit inputs. The Go comparison samples 60,000 inputs with seeds 17, 2718 and 65537, plus 132 boundary cases. The address detail is modeled by its family and zone contents; handle identity and allocation are outside the claim.

Native Go calling netip.AddrFrom16 and observing As16, Is6 and Zone.

### net/netip.AddrFrom4 (AWAITING_REVIEW)

Returns an unzoned IPv4 address whose As16 bytes have the IPv4-mapped prefix followed by the four input bytes in order.

The proof covers all 32-bit inputs. The Go comparison samples 60,000 inputs with seeds 17, 2718 and 65537, plus 36 boundary cases. The model represents address detail by family and zone contents. Handle identity and allocation are outside the claim.

Native Go calling netip.AddrFrom4 and observing As16, Is6 and Zone.

### net/netip.AddrPort.Addr (AWAITING_REVIEW)

Returns the stored address, preserving bytes, family and zone for every port.

The proof covers all modeled addresses and ports. The Go comparison uses six boundary cases and 60,000 samples with seeds 17, 2718 and 65537, and empty or decimal zones.

Native Go calling AddrPort.Addr and observing BitLen, As16 and Zone.

### net/netip.AddrPort.Compare (AWAITING_REVIEW)

Orders address-port pairs by address and breaks ties by unsigned port, returning -1, 0 or 1.

The proof covers all modeled addresses and ports. Go comparisons use empty and decimal zones in 256 boundary pairs and 60,000 samples with seeds 17, 2718 and 65537.

Native Go calling AddrPort.Compare on zero-address, IPv4, mapped IPv6, and scoped and unscoped IPv6 pairs.

### net/netip.AddrPort.IsValid (AWAITING_REVIEW)

Returns true exactly when the stored address is initialized, for every port including zero.

The proof covers every modeled address and port. The comparison exhausts all 65,536 ports in each of the three modeled address classes and uses six representative Go addresses.

Native Go calling netip.AddrPort.IsValid on receivers built with AddrPortFrom.

### net/netip.AddrPort.Port (AWAITING_REVIEW)

Returns the stored uint16 port independently of the address, including zero for a zero AddrPort.

The proof quantifies over all modeled addresses and ports. The comparison exhausts all 65,536 port values for one Lean address and five Go addresses, with a separate Go zero-value check.

Native Go invoking netip.AddrPort.Port on the zero value and values built with AddrPortFrom.

### net/netip.AddrPortFrom (AWAITING_REVIEW)

Preserves the supplied address bytes, family, zone and uint16 port.

The proof covers every modeled address and port. The Go check uses six boundary cases and 60,000 samples across seeds 17, 2718 and 65537, with empty and decimal zone strings. Allocation is outside the claim.

Native Go calling AddrPortFrom and observing Addr and Port, with BitLen, As16 and Zone on the returned address.

### net/netip.IPv4Unspecified (AWAITING_REVIEW)

Returns the unzoned IPv4 address 0.0.0.0 with IPv4-mapped bytes in As16.

Exhaustive for this zero-argument function. The model represents the interned address detail by its IPv6 flag and zone contents. Handle identity and allocation are outside the claim.

Native Go calling netip.IPv4Unspecified and observing As16, Is6 and Zone.

### net/netip.IPv6LinkLocalAllNodes (AWAITING_REVIEW)

Returns the unzoned IPv6 multicast address ff02::1.

Exhaustive for this zero-argument function. The model represents the interned address detail by its IPv6 flag and zone contents. Handle identity and allocation are outside the claim.

Native Go calling netip.IPv6LinkLocalAllNodes and observing As16, Is6 and Zone.

### net/netip.IPv6LinkLocalAllRouters (AWAITING_REVIEW)

Returns the unzoned IPv6 multicast address ff02::2.

Exhaustive for this zero-argument function. The model represents the interned address detail by its IPv6 flag and zone contents. Handle identity and allocation are outside the claim.

Native Go calling netip.IPv6LinkLocalAllRouters and observing As16, Is6 and Zone.

### net/netip.IPv6Loopback (AWAITING_REVIEW)

Returns the unzoned IPv6 address ::1, with high word zero and low word one.

Exhaustive for this zero-argument function. The model represents the interned address detail by its IPv6 flag and zone contents. Handle identity and allocation are outside the claim.

Native Go calling netip.IPv6Loopback and observing As16, Is6 and Zone.

### net/netip.IPv6Unspecified (AWAITING_REVIEW)

Returns the unzoned IPv6 address ::, with both address words zero.

Exhaustive for this zero-argument function. The model represents the interned address detail by its IPv6 flag and zone contents. Handle identity and allocation are outside the claim.

Native Go calling netip.IPv6Unspecified and observing As16, Is6 and Zone.

### net/netip.Prefix.Addr (AWAITING_REVIEW)

Returns the stored address regardless of prefix length or validity, preserving stored host bits.

The proof covers all modeled stored addresses and length tags. Fixtures use PrefixFrom, which strips the input zone before Addr reads the stored address. The Go check uses six boundary cases and 60,000 samples with seeds 17, 2718 and 65537.

Native Go calling Prefix.Addr on valid and invalid PrefixFrom results and observing BitLen, As16 and Zone.

### net/netip.Prefix.Bits (AWAITING_REVIEW)

Decodes the stored prefix length tag by subtracting one after widening, yielding -1 for an invalid prefix.

The proof covers all byte tags. Go and Lean checks exhaust the 130 reachable tags, with Go checking five representative addresses at requested lengths -1 through 129 and both int extrema, plus a zero Prefix.

Native Go calling Prefix.Bits on the zero value and PrefixFrom results.

### net/netip.Prefix.IsSingleIP (AWAITING_REVIEW)

For PrefixFrom results, returns true exactly for initialized IPv4 addresses at /32 and IPv6 addresses at /128.

Proof covers every modeled address and requested integer length. The comparison exhausts reachable length tags by family using five representative addresses at requested lengths -1 through 129 and both int extrema. A separate Go check covers the zero Prefix.

Native Go invoking Prefix.IsSingleIP on the zero value and PrefixFrom results.

### net/netip.Prefix.IsValid (AWAITING_REVIEW)

Returns true exactly when the stored prefix length decodes to a nonnegative length, including /0.

The proof covers all byte tags. Go and Lean checks exercise all 130 reachable tags. Go uses five representative addresses at requested lengths -1 through 129 and both int extrema, plus a zero Prefix.

Native Go calling Prefix.IsValid on the zero value and PrefixFrom results.

### net/netip.PrefixFrom (AWAITING_REVIEW)

Preserves family and address bits including host bits, strips IPv6 zones, and reports the requested prefix length if valid or -1 otherwise.

The model represents observable address details and uses mathematical integers, covering every Go int input. Allocation is outside the claim. Go checks 11 boundary cases and 60,000 samples with seeds 17, 2718 and 65537.

Native Go calling PrefixFrom and observing Bits and the stored address through Addr, BitLen, As16 and Zone.

### slices.BinarySearch (AWAITING_REVIEW)

Returns the lower bound and reports a match only at an in-range equal element, with NaN keys matching NaN elements for float64.

The proof assumes a partitioned less-than predicate and abstracts typed comparisons. Go checks int and float64 directly against linear scans, with 18 boundary cases and 120,000 samples across both types using seeds 17, 2718 and 65537. Other element types and cross-runtime primitive equivalence are outside the proof.

Native Go calling slices.BinarySearch on ascending integer and binary64 slices.

### slices.BinarySearchFunc (AWAITING_REVIEW)

For a pure element-to-target comparator with a negative prefix and nonnegative suffix, returns the first nonnegative index or the slice length. The found flag is true exactly when the index is in bounds and the comparator returns zero there.

Models a pure comparator and return values. Comparator effects, panics, nil callbacks and inputs outside the comparison partition are outside the claim. The direct Go check covers nine boundary cases and 60,000 samples with seeds 17, 2718 and 65537.

Native Go calling slices.BinarySearchFunc with int elements and a struct target in ascending and descending order, using results -7, 0 and 13.

### slices.Clip (AWAITING_REVIEW)

Preserves visible elements and nilness and sets capacity to length.

The model abstracts backing storage. Go separately checks backing address preservation for nonempty slices. Five boundary cases and 60,000 samples use seeds 17, 2718 and 65537.

Native Go calling slices.Clip with a named slice of int, including subslices with nonzero offsets.

### slices.Clone (AWAITING_REVIEW)

Preserves the original element sequence and nilness.

Models shallow copying by visible contents. Capacity and allocation are outside the proof. The Go check also writes to the first element of each nonempty clone and checks that the original is unaffected. Five boundary cases and 60,000 samples use seeds 17, 2718 and 65537.

Native Go calling slices.Clone on a named slice of int, including offset subslices.

### slices.Compact (AWAITING_REVIEW)

Keeps the first element and each element that differs from its original predecessor. The original-length backing region contains the retained prefix followed by zeros.

Lean models values and the original backing region. Slice identity, nilness, capacity, and cells outside the view are checked directly in Go. The comparison uses named int slices and covers nine boundary cases and 60,000 samples with seeds 17, 2718 and 65537. Other Go type instantiations and equality panics are outside the evidence.

Native Go calling slices.Compact with named int slices, including nil and nonnil empty slices, offset views, and spare capacity.

### slices.Compare (AWAITING_REVIEW)

For int slices, returns -1, 0 or 1 for lexicographic less-than, equality or greater-than, with a shorter strict prefix ordered first.

The scan is parametric in element comparison. The proof and Go check instantiate it for int. Floating-point ordering and other element types are outside this evidence. Go uses eight boundary pairs and 60,000 samples with seeds 17, 2718 and 65537.

Native Go calling slices.Compare instantiated with int and comparing with a suffix-to-prefix oracle.

### slices.CompareFunc (AWAITING_REVIEW)

Returns the exact first nonzero comparator result, propagates a comparator panic, or compares slice lengths if all aligned pairs tie.

Models fixed elements and a panic category. Go checks call order, call count and panic payload identity on int and string slices. Nil comparators and mutation of slice elements are outside the comparison. Go checks twelve boundary cases and 60,000 samples with seeds 17, 2718 and 65537.

Native Go calling slices.CompareFunc with int and string slices and a recording comparator.

### slices.Concat (AWAITING_REVIEW)

Appends parts in order, returns nil for an empty result, and panics when the total length exceeds MaxInt.

Models visible elements and a length bound. Capacity rounding and allocator failures are outside the claim. Go checks seven boundary cases and 60,000 samples with seeds 17, 2718 and 65537, plus native overflow using zero-sized elements. Lean separately checks overflow at a reduced bound. Go also probes storage independence at each nonempty part.

Native Go calling slices.Concat on a named int slice type and on zero-sized elements for overflow.

### slices.Contains (AWAITING_REVIEW)

Returns true exactly when at least one slice element compares equal to the key.

The proof abstracts typed equality as a total Boolean operation. Go checks int slices in eight boundary cases and 60,000 samples with seeds 17, 2718 and 65537. Interface comparisons that panic are outside the model.

Native Go directly calling slices.Contains instantiated with int.

### slices.ContainsFunc (AWAITING_REVIEW)

Reports whether an element satisfies the predicate and stops calling it at the first true result or panic.

The model uses fixed elements and a panic category. Go checks panic payload identity and callback order. Nil predicates and mutation of slice elements are outside the comparison. Go checks eight boundary cases and 60,000 samples with seeds 17, 2718 and 65537.

Native Go calling slices.ContainsFunc on int slices with scripted callback outcomes.

### slices.Delete (AWAITING_REVIEW)

Valid signed bounds remove the indicated range and zero the vacated cells in the original backing region. Invalid bounds cause a bounds panic before any mutation.

Lean models values, the original backing region and the panic category. The Go check also observes nilness, capacity, backing identity and cells outside offset views. Exact panic text and other Go element instantiations are outside the evidence. Go checks nine boundary cases and 60,000 samples with seeds 17, 2718 and 65537.

Native Go calling slices.Delete with named int slices and valid and invalid signed bounds, including bounds within capacity but past length.

### slices.DeleteFunc (AWAITING_REVIEW)

Visits original elements in order. On success retains the elements rejected by the deletion predicate and clears the tail. On predicate panic it stops, preserves completed prefix writes, and leaves the unwritten suffix intact.

The Lean model assumes predicate outcomes depend on fixed element values and represents panic as a category. The Go check observes callback order, partial writes, panic payload, nilness, capacity, backing identity and cells outside offset views. Nil callbacks and predicates mutating the backing array are outside the evidence. Go checks nine boundary cases and 60,000 samples with seeds 17, 2718 and 65537.

Native Go calling slices.DeleteFunc on named int slices with an even-value predicate that optionally panics on 13.

### slices.Equal (AWAITING_REVIEW)

For lawful element equality, returns true exactly when the slices have identical element sequences. Nil and empty slices compare equal.

The direct comparison uses int slices. NaN equality and panicking interface comparisons are outside this evidence. Go checks eight boundary cases and 60,000 samples with seeds 17, 2718 and 65537.

Native Go calling slices.Equal instantiated with int and checking sequence equality against a reflection-based oracle.

### slices.EqualFunc (AWAITING_REVIEW)

Unequal lengths return false without comparator calls. Equal-length inputs are compared in order until the first false result or panic.

Models fixed elements and scripted comparator outcomes. Go checks argument order, call count and panic payload identity on int and string slices. Nil comparators and mutation of slice elements are outside the comparison. Go checks nine boundary cases and 60,000 samples with seeds 17, 2718 and 65537.

Native Go calling slices.EqualFunc with int and string slices and a recording comparator.

### slices.Index (AWAITING_REVIEW)

Returns the first matching index or -1 exactly when no element matches.

The proof abstracts typed equality as a total Boolean operation. Go checks int slices in eight boundary cases and 60,000 samples with seeds 17, 2718 and 65537. Interface comparisons that panic are outside the model.

Native Go directly calling slices.Index instantiated with int.

### slices.IndexFunc (AWAITING_REVIEW)

Visits elements in order, returns the first true index or -1, and stops calling the predicate after a true result or panic.

The model uses fixed elements and a panic category. Go checks panic payload identity and callback order. Nil predicates and mutation of slice elements are outside the comparison. Go checks eight boundary cases and 60,000 samples with seeds 17, 2718 and 65537.

Native Go calling slices.IndexFunc on int slices with scripted callback outcomes.

### slices.IsSorted (AWAITING_REVIEW)

Returns true exactly when the slice has no adjacent inversions under the element ordering.

The proof abstracts the typed cmp.Less specialization as a pure comparison. The Go comparison covers int slices on nine boundary cases and 60,000 samples with seeds 17, 2718 and 65537. Other Go type instantiations are outside this comparison.

Native Go directly invoking slices.IsSorted instantiated with int.

### slices.IsSortedFunc (AWAITING_REVIEW)

Returns true exactly when every descending adjacent comparison returns a nonnegative result. Stops at the first negative result or comparator panic.

Models fixed elements and comparator outcomes, including panic as a category. The Go trace also checks the panic payload and earlier callback effects. Comparators that mutate slice elements and nil comparators are outside the comparison. Go checks six boundary cases and 60,000 samples with seeds 17, 2718 and 65537.

Native Go calling slices.IsSortedFunc on int slices with scripted comparator outcomes.

### slices.Max (AWAITING_REVIEW)

For int slices, returns a member no smaller than any other element, or panics with the empty-list message for empty input.

The fold is parametric in the typed max operation. The theorem and direct Go check specialize to int. Floating-point NaN and signed-zero rules are outside this evidence. Go checks eight boundary cases and 60,000 samples with seeds 17, 2718 and 65537.

Native Go calling slices.Max instantiated with int.

### slices.MaxFunc (AWAITING_REVIEW)

For a pure comparator whose positive results agree with a strict integer-key order, returns the first element with a maximum key. Empty input panics with the pinned message.

The theorem assumes an integer-key ordering and a pure comparator. Comparator effects, nil callbacks, and comparator panics are outside this evidence. The direct Go comparison covers nine boundary cases and 60,000 samples with seeds 17, 2718 and 65537.

Native Go calling slices.MaxFunc on records with int keys and observable original indexes, comparing only by key.

### slices.Min (AWAITING_REVIEW)

For int slices, returns a member no greater than any other element, or panics with the empty-list message for empty input.

The fold is parametric in the typed min operation. The theorem and direct Go check specialize to int. Floating-point NaN and signed-zero rules are outside this evidence. Go checks eight boundary cases and 60,000 samples with seeds 17, 2718 and 65537.

Native Go calling slices.Min instantiated with int.

### slices.MinFunc (AWAITING_REVIEW)

For a pure comparator whose negative results agree with a strict integer-key order, returns the first element with a minimum key. Empty input panics with the pinned message.

The theorem assumes an integer-key ordering and a pure comparator. Comparator effects, nil callbacks, and comparator panics are outside this evidence. The direct Go comparison covers nine boundary cases and 60,000 samples with seeds 17, 2718 and 65537.

Native Go calling slices.MinFunc on records with int keys and observable original indexes, comparing only by key.

### slices.Repeat (AWAITING_REVIEW)

Returns the requested repetitions with length equal to the product, or panics for a negative count or a product exceeding MaxInt.

The model summarizes two-word multiplication as an exact product and the doubling copies as repeated concatenation. The Go check observes contents, nilness, capacity and panic payloads for int slices. Allocator failures are outside the claim. Ten boundary cases and 60,000 samples use seeds 17, 2718 and 65537.

Native Go directly calling slices.Repeat on int slices, including native int extrema.

### slices.Reverse (AWAITING_REVIEW)

Reverses the element sequence in place.

The proof relates endpoint swaps on the modeled view to list reversal. Go separately checks nilness, length, capacity and writes outside the view. Six boundary cases and 60,000 samples use seeds 17, 2718 and 65537.

Native Go calling slices.Reverse with a named slice of int, including offset subslices.

### sort.Find (AWAITING_REVIEW)

Returns the first nonpositive comparison index or the search length, with found true exactly for an in-range zero comparison. Nonpositive lengths return zero and false.

The proof covers normally returning comparators with a partitioned sign predicate. Callback panics and state effects are outside the claim. Go checks callback bounds, ten boundary cases and 60,000 samples with seeds 17, 2718 and 65537, including native int extrema.

Native Go calling sort.Find on comparator sign partitions.

### sort.Float64Slice.Len (AWAITING_REVIEW)

For a valid subslice x[lo:hi], Len returns hi - lo, including zero for empty and nil slices.

Models visible elements as a list of binary64 bit patterns and Go's builtin len as its length. Capacity, backing storage identity and nilness are abstracted. The proof covers arbitrary elements and valid bounds. Go comparisons use backing lengths 0 through 256, positive infinity and NaN payloads, seven boundary cases and 60,000 samples across seeds 17, 2718 and 65537.

Native Go calling sort.Float64Slice.Len on two-index subslices.

### sort.Float64Slice.Less (AWAITING_REVIEW)

Compares using binary64 less-than with NaNs ordered before non-NaNs and panics at the first invalid index.

The proof covers bounds and the comparison formula using Lean Float primitives. Go and Lean floating-point behavior is compared over sampled bit patterns, not proved equivalent. Go checks 109 boundary cases and 60,000 samples with seeds 17, 2718 and 65537.

Native Go calling sort.Float64Slice.Less and comparing Boolean results or recovered runtime bounds panic text.

### sort.Float64Slice.Search (AWAITING_REVIEW)

Returns the first index satisfying the binary64 >= predicate, or the slice length if none qualifies, when the predicate partitions the slice.

The proof assumes a partitioned predicate and uses Lean Float comparisons. The Go check uses ascending slices with NaNs first, includes NaN keys whose predicate is always false, and checks the cross-runtime comparison on ten boundary cases and 60,000 samples with seeds 17, 2718 and 65537.

Native Go directly calling sort.Float64Slice.Search on sorted binary64 slices and comparing against a linear scan.

### sort.Float64Slice.Swap (AWAITING_REVIEW)

Exchanges the selected binary64 bit patterns and preserves slice length and all other elements. An invalid index panics before either write.

Models visible slice elements by their exact binary64 encodings. Go checks nine boundary cases and 60,000 samples with seeds 17, 2718 and 65537.

Native Go calling sort.Float64Slice.Swap and comparing post-call bit patterns and recovered bounds panic text.

### sort.Float64sAreSorted (AWAITING_REVIEW)

Returns true exactly when no adjacent pair is inverted in ascending binary64 order with NaNs first.

The proof uses Lean Float primitives. Their correspondence with Go is tested over binary64 bit patterns. Go checks nine boundary cases and 60,000 samples with seeds 17, 2718 and 65537, using an all-pairs order oracle.

Native Go directly calling sort.Float64sAreSorted on nil, empty, sorted and unsorted slices.

### sort.IntSlice.Len (AWAITING_REVIEW)

For a valid subslice x[lo:hi], Len returns hi - lo, including zero for empty and nil slices.

Models visible elements as a list and Go's builtin len as its length. Capacity, backing storage identity and nilness are abstracted. The proof covers arbitrary elements and valid bounds. Go comparisons use backing lengths 0 through 256, seven boundary cases and 60,000 samples across seeds 17, 2718 and 65537.

Native Go calling sort.IntSlice.Len on two-index subslices.

### sort.IntSlice.Less (AWAITING_REVIEW)

Compares signed elements at valid indices and panics at the first invalid index in evaluation order.

Lists model visible elements, with nil and empty slices represented by the empty list. Lean models the failed index and the driver formats its panic payload. Go checks nine boundary cases and 60,000 samples with seeds 17, 2718 and 65537.

Native Go calling sort.IntSlice.Less and comparing its Boolean or recovered runtime bounds panic text.

### sort.IntSlice.Search (AWAITING_REVIEW)

Returns the first index whose element is at least the key, or the slice length if there is none, for ascending sorted slices.

The proof assumes the Boolean search predicate partitions at the lower bound, as it does for an ascending slice. The unsigned midpoint is represented by exact natural arithmetic on valid nonnegative int indices. Go checks eight boundary cases and 60,000 samples with seeds 17, 2718 and 65537.

Native Go directly calling sort.IntSlice.Search on ascending integer slices.

### sort.IntSlice.Swap (AWAITING_REVIEW)

Exchanges the two indexed elements, preserving length and all other elements. An invalid index panics before any write.

Models visible slice contents and the failed index. The driver observes the original list on panic. Go uses nine boundary cases and 60,000 samples with seeds 17, 2718 and 65537.

Native Go calling sort.IntSlice.Swap and observing contents after success or a recovered runtime bounds panic.

### sort.IntsAreSorted (AWAITING_REVIEW)

Returns true exactly when every adjacent pair is in nondecreasing signed integer order.

The proof covers arbitrary integer lists. Nil and empty slices have the same visible elements. Go checks nine boundary cases and 60,000 samples with seeds 17, 2718 and 65537, using an all-pairs order oracle.

Native Go directly calling sort.IntsAreSorted on nil, empty, sorted and unsorted slices.

### sort.IsSorted (AWAITING_REVIEW)

After Len returns a nonnegative length, checks adjacent pairs from right to left, stopping on an inversion or callback panic. Returns true exactly when every adjacent comparison returns false.

Models a completed nonnegative Len call and scripted Less outcomes including panic as a category. Go also checks Len call count and the Less panic marker. Negative lengths and Len panics are outside the claim. Comparison uses six boundary cases and 60,000 samples with seeds 17, 2718 and 65537.

Native Go calling sort.IsSorted with a recording sort.Interface.

### sort.Reverse (AWAITING_REVIEW)

Delegates Len and Swap unchanged and delegates Less with the two indices exchanged.

The model abstracts interface identity and allocation. Callback outcomes are parametric. The Go comparison uses a recorder on four boundary pairs and 60,000 samples with seeds 17, 2718 and 65537, plus a Go Less panic and side-effect check.

Native Go calling sort.Reverse and then invoking all three methods on the returned interface.

### sort.Search (AWAITING_REVIEW)

Returns the first true index for a partitioned predicate or the search length if all indices are false. Nonpositive lengths return zero.

The theorem models the returned index for normally returning partitioned callbacks. Panic propagation and callback state effects are outside the claim. Go asserts callback indices are in bounds and checks eight boundary cases and 60,000 samples with seeds 17, 2718 and 65537, including native int extremes.

Native Go directly calling sort.Search with threshold predicates.

### sort.SearchFloat64s (AWAITING_REVIEW)

Returns the first index satisfying the binary64 >= predicate, or the slice length if none qualifies, when the predicate partitions the slice.

The proof assumes a partitioned predicate and uses Lean Float comparisons. The Go check uses ascending slices with NaNs first, includes NaN keys whose predicate is always false, and checks the cross-runtime comparison on ten boundary cases and 60,000 samples with seeds 17, 2718 and 65537.

Native Go directly calling sort.SearchFloat64s on sorted binary64 slices and comparing against a linear scan.

### sort.SearchInts (AWAITING_REVIEW)

Returns the first index whose element is at least the key, or the slice length if there is none, for ascending sorted slices.

The proof assumes the Boolean search predicate partitions at the lower bound, as it does for an ascending slice. The unsigned midpoint is represented by exact natural arithmetic on valid nonnegative int indices. Go checks eight boundary cases and 60,000 samples with seeds 17, 2718 and 65537.

Native Go directly calling sort.SearchInts on ascending integer slices.

### sort.SearchStrings (AWAITING_REVIEW)

On a bytewise ascending slice, returns the first index at or above the key, or the slice length if no element qualifies.

The proof assumes the search predicate partitions at the lower bound, as it does for ascending byte sequences. The unsigned midpoint is represented by exact natural arithmetic on valid nonnegative int indices. Go checks eight boundary cases and 60,000 samples with seeds 17, 2718 and 65537.

Native Go directly calling sort.SearchStrings on sorted slices of arbitrary byte strings.

### sort.StringSlice.Len (AWAITING_REVIEW)

For a valid subslice x[lo:hi], Len returns hi - lo, including zero for empty and nil slices.

Models visible elements as a list and Go's builtin len as its length. Capacity, backing storage identity and nilness are abstracted. The proof covers arbitrary elements and valid bounds. Go comparisons use backing lengths 0 through 256, seven boundary cases and 60,000 samples across seeds 17, 2718 and 65537.

Native Go calling sort.StringSlice.Len on two-index subslices.

### sort.StringSlice.Less (AWAITING_REVIEW)

Compares valid operands lexicographically by unsigned bytes and panics at the first invalid index.

Models Go strings as arbitrary byte sequences and nil and empty slices as empty lists. Go checks 11 boundary cases and 60,000 samples with seeds 17, 2718 and 65537.

Native Go calling sort.StringSlice.Less and comparing its Boolean or recovered runtime bounds panic text.

### sort.StringSlice.Search (AWAITING_REVIEW)

On a bytewise ascending slice, returns the first index at or above the key, or the slice length if no element qualifies.

The proof assumes the search predicate partitions at the lower bound, as it does for ascending byte sequences. The unsigned midpoint is represented by exact natural arithmetic on valid nonnegative int indices. Go checks eight boundary cases and 60,000 samples with seeds 17, 2718 and 65537.

Native Go directly calling sort.StringSlice.Search on sorted slices of arbitrary byte strings.

### sort.StringSlice.Swap (AWAITING_REVIEW)

Exchanges the two indexed elements, preserving length and all other elements. An invalid index panics before any write.

Models Go strings as arbitrary byte sequences and compares visible slice contents, including on panic. Go checks 11 boundary cases and 60,000 samples with seeds 17, 2718 and 65537.

Native Go calling sort.StringSlice.Swap and observing contents after success or a recovered runtime bounds panic.

### sort.StringsAreSorted (AWAITING_REVIEW)

Returns true exactly when no adjacent pair is inverted under unsigned bytewise lexicographic ordering.

Models strings as arbitrary byte sequences. Go checks nine boundary cases and 60,000 samples with seeds 17, 2718 and 65537, using an all-pairs byte order oracle.

Native Go directly calling sort.StringsAreSorted on nil, empty, sorted and unsorted slices.

### strconv.AppendBool (AWAITING_REVIEW)

Returns the original destination bytes followed by the lowercase Boolean spelling, extending the length by four or five bytes.

Proof covers arbitrary destination bytes and both Boolean values. The claim covers returned contents and length. Backing array identity, capacity and allocation are outside the model. Go uses 18 boundary cases and 60,000 samples with seeds 17, 2718 and 65537.

Native Go calling strconv.AppendBool on nil, empty and populated slices with varied spare capacity.

### strconv.CanBackquote (AWAITING_REVIEW)

Accepts exactly valid UTF-8 without backquotes, DEL, BOM or ASCII control characters other than tab.

Lean validates the whole byte array before checking characters. Equivalence of the Go and Lean UTF-8 decoders is tested rather than proved. Go checks 272 boundary cases and 60,000 samples with seeds 17, 2718 and 65537.

Native Go calling strconv.CanBackquote on valid and malformed UTF-8 byte strings.

### strconv.FormatBool (AWAITING_REVIEW)

Returns exactly the lowercase string "false" for false and "true" for true.

Exhaustive over both boolean values. The claim covers returned string contents.

Native Go calling the public strconv.FormatBool entry point.

### strconv.NumError.Unwrap (AWAITING_REVIEW)

Returns the stored error interface unchanged for non-nil receivers and panics on a nil receiver.

The proof treats error values as opaque and models nil dereference as a panic outcome. Go compares comparable error values including nil, typed nil and distinct equal-message errors, and checks the runtime panic. It uses 14 boundary cases and 60,000 samples with seeds 17, 2718 and 65537.

Native Go calling (*strconv.NumError).Unwrap directly.

### strconv.ParseBool (AWAITING_REVIEW)

Accepts exactly the twelve documented spellings and returns the corresponding Boolean. Other inputs return false and NumError with Func ParseBool, the original input bytes and ErrSyntax.

Proof covers arbitrary byte strings. Error message rendering and allocation are outside the claim. Go checks 313 boundary cases and 60,000 samples with seeds 17, 2718 and 65537.

Native Go calling the public strconv.ParseBool wrapper and comparing returned values and error payloads.
