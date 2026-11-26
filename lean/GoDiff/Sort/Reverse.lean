import GoLean.Sort.Reverse
import GoDiff.Math.Bits.Tuple

open GoLean.Sort

def main : IO Unit := runTuple fun fields => do
  let [i, j] := fields | throw (IO.userError "expected two indices")
  let data : SortInterface String String String :=
    ⟨fun _ => "len 7", fun a b => s!"less {a} {b} {decide (a < b)}",
      fun a b => s!"swap {a} {b}"⟩
  let rev := reverseInterface data
  return s!"{rev.length ()}; {rev.less i j}; {rev.swap i j}"
