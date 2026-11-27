import Init.Data.BitVec.Basic

namespace GoLean.Internal.Bytealg

private def prime : BitVec 32 := 16777619

private def extendHash (hash : BitVec 32) (byte : BitVec 8) : BitVec 32 :=
  hash * prime + BitVec.ofNat 32 byte.toNat

private def powerLoop : Nat → Nat → BitVec 32 → BitVec 32 → BitVec 32
  | 0, _, result, _ => result
  | fuel + 1, exponent, result, square =>
    if exponent = 0 then result
    else
      let next := if exponent % 2 = 1 then result * square else result
      powerLoop fuel (exponent / 2) next (square * square)

/-- Go 1.27.0 src/internal/bytealg/bytealg.go:34–50. All arithmetic
    wraps at 32 bits, including the square-and-multiply power calculation. -/
def hashStr (bytes : List (BitVec 8)) : BitVec 32 × BitVec 32 :=
  (bytes.foldl extendHash 0, powerLoop bytes.length bytes.length 1 prime)

/-- The same recurrence over the bytes in reverse order. -/
def hashStrRev (bytes : List (BitVec 8)) : BitVec 32 × BitVec 32 :=
  hashStr bytes.reverse

private def rollingSearch (separator : List (BitVec 8)) (want power : BitVec 32)
    (offset : Nat) (hash : BitVec 32) (window : List (BitVec 8)) :
    List (BitVec 8) → Int
  | [] => if hash == want && window.take separator.length == separator then Int.ofNat offset else -1
  | byte :: rest =>
    if hash == want && window.take separator.length == separator then Int.ofNat offset
    else match window with
    | [] => -1
    | outgoing :: next =>
      let updated := extendHash hash byte - power * BitVec.ofNat 32 outgoing.toNat
      rollingSearch separator want power (offset + 1) updated next rest

private def rollingIndex (input separator : List (BitVec 8)) : Int :=
  let (want, power) := hashStr separator
  let hash := (input.take separator.length).foldl extendHash 0
  rollingSearch separator want power 0 hash input (input.drop separator.length)

/-- Go 1.27.0 src/internal/bytealg/bytealg.go:71–92. None is the
    bounds panic from initializing a window longer than the input. -/
def indexRabinKarp (input separator : List (BitVec 8)) : Option Int :=
  if input.length < separator.length then none else some (rollingIndex input separator)

/-- The reversed scan follows Go's decreasing window position. -/
def lastIndexRabinKarp (input separator : List (BitVec 8)) : Option Int :=
  if input.length < separator.length then none
  else
    let reversed := rollingIndex input.reverse separator.reverse
    some (if reversed < 0 then -1 else Int.ofNat (input.length - separator.length) - reversed)

theorem hashStr_empty : hashStr [] = (0, 1) := rfl

theorem hashStr_collision :
    hashStr [0x28, 0x9e, 0x1a, 0x30, 0xb8, 0xe4, 0x2a] =
      hashStr [0x20, 0x0c, 0xf2, 0x59, 0x3d, 0x9d, 0x05] := rfl

theorem hashStrRev_reverse (bytes : List (BitVec 8)) :
    hashStrRev bytes = hashStr bytes.reverse := rfl

theorem indexRabinKarp_long (input separator : List (BitVec 8))
    (tooLong : input.length < separator.length) : indexRabinKarp input separator = none := by
  simp [indexRabinKarp, tooLong]

theorem lastIndexRabinKarp_long (input separator : List (BitVec 8))
    (tooLong : input.length < separator.length) : lastIndexRabinKarp input separator = none := by
  simp [lastIndexRabinKarp, tooLong]

theorem indexRabinKarp_collision :
    indexRabinKarp [0x28, 0x9e, 0x1a, 0x30, 0xb8, 0xe4, 0x2a]
      [0x20, 0x0c, 0xf2, 0x59, 0x3d, 0x9d, 0x05] = some (-1) := rfl

theorem lastIndexRabinKarp_collision :
    lastIndexRabinKarp [0x2a, 0xe4, 0xb8, 0x30, 0x1a, 0x9e, 0x28]
      [0x05, 0x9d, 0x3d, 0x59, 0xf2, 0x0c, 0x20] = some (-1) := rfl

end GoLean.Internal.Bytealg
