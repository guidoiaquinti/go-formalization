import Init.Data.BitVec.Lemmas

namespace GoLean.Time

private def fraction (value precision : Nat) : String :=
  let digits := (List.range precision).map fun i =>
    Char.ofNat (48 + value / 10 ^ (precision - i - 1) % 10)
  let trimmed := (digits.reverse.dropWhile (· == '0')).reverse
  if trimmed.isEmpty then "" else "." ++ String.ofList trimmed

/-- The representation uses the magnitude, including the minimum signed word. -/
def durationString (duration : Int) : String :=
  let magnitude := duration.natAbs
  if magnitude = 0 then "0s"
  else
    let sign := if duration < 0 then "-" else ""
    let formatSmall (scale precision : Nat) (unit : String) :=
      toString (magnitude / scale) ++ fraction magnitude precision ++ unit
    let text :=
      if magnitude < 1000 then formatSmall 1 0 "ns"
      else if magnitude < 1000000 then formatSmall 1000 3 "µs"
      else if magnitude < 1000000000 then formatSmall 1000000 6 "ms"
      else
        let whole := magnitude / 1000000000
        let hours := if whole ≥ 3600 then toString (whole / 3600) ++ "h" else ""
        let minutes := if whole ≥ 60 then toString (whole / 60 % 60) ++ "m" else ""
        hours ++ minutes ++ toString (whole % 60) ++ fraction magnitude 9 ++ "s"
    sign ++ text

def durationFormat (duration : Int) (backing : Option (List (BitVec 8))) :
    Option (Nat × List (BitVec 8)) :=
  backing.map fun previous =>
    let bytes := (durationString duration).toUTF8.data.toList.map UInt8.toBitVec
    let first := previous.length - bytes.length
    (first, previous.take first ++ bytes)

theorem durationString_zero : durationString 0 = "0s" := by decide

theorem durationString_micro : durationString (-1001) = "-1.001µs" := by decide

theorem durationFormat_zero (first : List (BitVec 8)) (old1 old2 : BitVec 8) :
    durationFormat 0 (some (first ++ [old1, old2])) =
      some (first.length, first ++ [48, 115]) := by
  have bytes : "0s".toUTF8.data.toList.map UInt8.toBitVec = [48, 115] := by decide
  rw [durationFormat, Option.map_some, durationString_zero, bytes]
  simp

end GoLean.Time
