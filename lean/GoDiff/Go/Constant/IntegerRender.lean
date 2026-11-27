import GoLean.Go.Constant.IntegerRender
import GoDiff.Math.Bits.Tuple

open GoLean.Go.Constant

def main : IO Unit := runTuple fun fields => do
  let [s, n, present] := fields | throw (IO.userError "expected signed word, large integer and presence")
  let word := BitVec.ofInt 64 s
  let value := if present == 0 then none else some n
  return String.intercalate "|" [int64ValString word, int64ValExactString word,
    intValString value, intValExactString value, toString (is32bit word), toString (is63bit word),
    toString (smallInt n.natAbs)]
