import GoLean.Go.Constant.Integer
import GoDiff.Math.Bits.Tuple

open GoLean.Go.Constant

private def wire (x : IntegerValue) : String :=
  match x with
  | .small w => s!"small:{w.toInt}"
  | .large n => s!"large:{n}"
  | .unknown => "unknown"
  | .other => "other"

private def convert [ToString α] (q : Option (α × Bool)) : String :=
  (q.map (fun (x, exact) => s!"{x}:{exact}")).getD "panic"

def main : IO Unit := runTuple fun row => do
  let tag :: input :: raw := row | throw (IO.userError "expected tag, input and bytes")
  let bytes ← raw.mapM fun v => UInt8.ofBitVec <$> wordInput 8 v
  let signed := makeInt64 (BitVec.ofInt 64 input)
  let unsigned := makeUint64 (BitVec.ofInt 64 input)
  let canonical := makeInt input
  let fromBytes := makeFromBytes bytes
  let selected := match tag with
    | 0 => canonical
    | 1 => fromBytes
    | 2 => IntegerValue.unknown
    | _ => IntegerValue.other
  let written := (integerBytes selected).map (fun data =>
    String.intercalate "," (data.map (toString ∘ UInt8.toNat)))
  return String.intercalate "|" [wire signed, wire unsigned, wire canonical, wire fromBytes,
    convert (int64Val selected), convert (uint64Val selected),
    (bitLen selected).map toString |>.getD "panic", written.getD "panic"]
