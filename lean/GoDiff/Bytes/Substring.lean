import GoLean.Bytes.Substring
import GoDiff.Math.Bits.Tuple

private def describe : GoLean.Bytes.Slice → String
  | none => "nil"
  | some (data, cap) => String.intercalate " " (toString cap :: data.map (fun b => toString b.toNat))

private def split (value : GoLean.Bytes.Slice × GoLean.Bytes.Slice × Bool) : String :=
  s!"{value.2.2}|{describe value.1}|{describe value.2.1}"

def main : IO Unit := runTuple fun fields => do
  let cap :: length :: words := fields | throw (IO.userError "expected capacity and needle length")
  if cap < -1 || length < -1 || words.length < length.toNat then
    throw (IO.userError "invalid capacity or needle length")
  let data ← words.mapM (wordInput 8)
  let body := data.drop length.toNat
  if (cap == -1 && !body.isEmpty) || (cap ≥ 0 && cap.toNat < body.length) then
    throw (IO.userError "invalid input capacity")
  let buffer := if cap == -1 then none else some (body, cap.toNat)
  let needle := if length == -1 then none else some (data.take length.toNat, length.toNat)
  let first := (GoLean.Bytes.index buffer needle).map toString |>.getD "-1"
  let last := (GoLean.Bytes.lastIndex buffer needle).map toString |>.getD "-1"
  return s!"{first}|{last}|{GoLean.Bytes.contains buffer needle}|{split (GoLean.Bytes.cut buffer needle)}|{split (GoLean.Bytes.cutLast buffer needle)}"
