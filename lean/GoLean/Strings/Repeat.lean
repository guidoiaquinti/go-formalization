import Init.Data.List.TakeDrop
import Init.Data.Int.Lemmas

namespace GoLean.Strings

inductive RepeatError where
  | negative
  | overflow

private def staticRepeatSize (byte : Nat) : Nat :=
  if byte = 32 || byte = 45 then 128
  else if byte = 48 || byte = 61 || byte = 9 then 64
  else 0

private def repeatChunks (target chunkMax : Nat) (buffer : List Nat) : Nat → List Nat
  | 0 => buffer
  | fuel + 1 =>
    if buffer.length < target then
      let chunk := min (target - buffer.length) (min buffer.length chunkMax)
      repeatChunks target chunkMax (buffer ++ buffer.take chunk) fuel
    else buffer

/-- Go 1.27.0 strings/strings.go:597–669. The full product guard is
    equivalent to checking the high unsigned limb and the signed range
    of the low limb. The static strings and chunk loop are byte models. -/
def repeatBytes (input : List Nat) (count : Int) (maxInt : Nat) : Except RepeatError (List Nat) :=
  if count = 0 then .ok []
  else if count = 1 then .ok input
  else if count < 0 then .error .negative
  else
    let size := input.length * count.toNat
    if maxInt < size then .error .overflow
    else
      match input with
      | [] => .ok []
      | first :: rest =>
        if size ≤ staticRepeatSize first && rest.all (· == first) then
          .ok (List.replicate size first)
        else
          let perChunk := 8192 / input.length * input.length
          let chunkMax := if size ≤ 8192 then size else if perChunk = 0 then input.length else perChunk
          .ok (repeatChunks size chunkMax input count.toNat)

theorem repeat_zero (input : List Nat) (maxInt : Nat) :
    repeatBytes input 0 maxInt = .ok [] := by simp [repeatBytes]

theorem repeat_one (input : List Nat) (maxInt : Nat) :
    repeatBytes input 1 maxInt = .ok input := by simp [repeatBytes]

theorem repeat_negative (input : List Nat) (count : Int) (maxInt : Nat) (negative : count < 0) :
    repeatBytes input count maxInt = .error .negative := by
  simp [repeatBytes, negative, show count ≠ 0 by omega, show count ≠ 1 by omega]

theorem repeat_overflow (input : List Nat) (count : Int) (maxInt : Nat)
    (positive : 1 < count) (large : maxInt < input.length * count.toNat) :
    repeatBytes input count maxInt = .error .overflow := by
  simp [repeatBytes, large, show count ≠ 0 by omega, show count ≠ 1 by omega,
    show ¬count < 0 by omega]

end GoLean.Strings
