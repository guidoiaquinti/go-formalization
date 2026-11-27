import GoLean.Go.Constant.Scalars
import GoDiff.Math.Bits.Tuple

open GoLean.Go.Constant

def main : IO Unit := runTuple fun fields => do
  let [code, flag, variant] := fields |
    throw (IO.userError "expected Kind code, boolean and accessor input")
  let value := flag != 0
  let unknown := match makeUnknown with | .unknown => "unknownVal" | _ => "wrong"
  let boolean := match makeBool value with | .boolean b => s!"boolVal:{b}" | _ => "wrong"
  let input := if variant == 0 then makeUnknown else if variant == 1 then makeBool value
    else .other (if variant == 2 then "<nil>" else if variant == 3 then "\"x\"" else toString code)
  let observed := match boolVal input with | .ok b => toString b | .error msg => "panic:" ++ msg
  return String.intercalate "|" [unknown, boolean, toString unknownKind, unknownString,
    unknownExactString, toString (boolKind value), boolString value, boolExactString value,
    observed, kindString code]
