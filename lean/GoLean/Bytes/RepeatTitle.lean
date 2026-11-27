import GoLean.Bytes.SimpleCase
import GoLean.Strings.IsSeparator

namespace GoLean.Bytes

inductive RepeatError where
  | negative
  | overflow

private def fillRepeat (target chunkMax : Nat) (filled : List Nat) : Nat → List Nat
  | 0 => filled
  | fuel + 1 =>
    if filled.length < target then
      let amount := min (target - filled.length) (min filled.length chunkMax)
      fillRepeat target chunkMax (filled ++ filled.take amount) fuel
    else filled

/-- Go 1.27.0 src/bytes/bytes.go:629–675. The full product guard
    implements both the high-word and signed-low-word checks. The output
    has exact clipped capacity and is nonnil even when empty. -/
def repeatBytes (input : List Nat) (count : Int) (maxInt : Nat) : Except RepeatError (List Nat) :=
  if count = 0 then .ok []
  else if count < 0 then .error .negative
  else
    let target := input.length * count.toNat
    if maxInt < target then .error .overflow
    else if input.isEmpty then .ok []
    else
      let aligned := 8192 / input.length * input.length
      let chunkMax := if target ≤ 8192 then target else if aligned = 0 then input.length else aligned
      .ok (fillRepeat target chunkMax input count.toNat)

theorem repeatBytes_zero (input : List Nat) (maxInt : Nat) :
    repeatBytes input 0 maxInt = .ok [] := by simp [repeatBytes]

theorem repeatBytes_negative (input : List Nat) (count : Int) (maxInt : Nat) (negative : count < 0) :
    repeatBytes input count maxInt = .error .negative := by
  simp [repeatBytes, negative, show count ≠ 0 by omega]

theorem repeatBytes_overflow (input : List Nat) (count : Int) (maxInt : Nat)
    (positive : 0 < count) (large : maxInt < input.length * count.toNat) :
    repeatBytes input count maxInt = .error .overflow := by
  simp [repeatBytes, large, show count ≠ 0 by omega, show ¬ count < 0 by omega]

theorem repeatBytes_chunkExample :
    repeatBytes [255, 0, 13] 3 9 = .ok [255, 0, 13, 255, 0, 13, 255, 0, 13] := rfl

private def scanTitle (previous : Nat) (input : List Nat) : Nat → List Nat
  | 0 => []
  | fuel + 1 =>
    if input.isEmpty then []
    else
      let (rune, width) := GoLean.Unicode.Utf8.decodeRune input
      let mapped := if GoLean.Strings.isSeparator (Int.ofNat previous)
        then GoLean.Strings.titleRune rune else rune
      GoLean.Unicode.Utf8.appendRune [] (Int.ofNat mapped) ++
        scanTitle rune (input.drop width) fuel

/-- Go 1.27.0 src/bytes/bytes.go:818–832. The stateful Map closure
    retains the previous original rune and starts after a space. -/
def titleBytes (input : List (BitVec 8)) : CaseBytes :=
  let output := scanTitle 32 (input.map BitVec.toNat) input.length
  ⟨output, if output.length ≤ input.length then some input.length else none⟩

theorem titleBytes_empty : titleBytes [] = ⟨[], some 0⟩ := rfl

theorem titleBytes_boundary :
    titleBytes [97, 98, 45, 99, 100, 95, 101, 102] =
      ⟨[65, 98, 45, 67, 100, 95, 101, 102], some 8⟩ := rfl

end GoLean.Bytes
