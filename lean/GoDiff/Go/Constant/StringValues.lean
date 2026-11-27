import GoLean.Go.Constant.StringValues
import GoDiff.Math.Bits.Tuple

open GoLean.Go.Constant

private def wire (bs : StringBytes) : String := String.intercalate "," (bs.map toString)
private def parseBytes (input : List Int) : IO (StringBytes × List Int) := do
  let n :: rest := input | throw (IO.userError "missing length")
  let bs ← (rest.take n.toNat).mapM fun x => return (← wordInput 8 x).toNat
  return (bs, rest.drop n.toNat)

private def observe (render : α → String) : Except StringFault α → String
  | .ok result => render result
  | .error .nilReceiver => "nilpanic"
  | .error (.wrongType reason) => "panic:" ++ reason

private def flattenWire (value : StringBytes) (tree : Option StringTree) : String :=
  match tree with
  | none => wire value ++ ":unknown"
  | some (.leaf stored) => wire value ++ ":leaf:" ++ wire stored
  | some (.concat _ _) => "uncached"

def main : IO Unit := runTuple fun row => do
  let shape :: variant :: rest := row | throw (IO.userError "missing shape and variant")
  let (a, rest) ← parseBytes rest
  let (b, rest) ← parseBytes rest
  let (c, rest) ← parseBytes rest
  let (d, _) ← parseBytes rest
  let created := makeString a
  let l := StringTree.concat created.tree (.leaf b)
  let tree := if shape == 0 then created.tree else if shape == 1 then
    .concat (.concat l (.leaf c)) (.leaf d) else if shape == 2 then
    .concat created.tree (.concat (.leaf b) (.concat (.leaf c) (.leaf d))) else .concat l l
  let receiver := if variant == 4 then none else some tree
  let operand := if variant == 1 then StringOperand.unknown else if variant == 2 then
    .other "true" else if variant == 3 then .other "<nil>" else .string receiver
  return String.intercalate "|" [
    toString created.usesEmptySingleton, toString (stringTreeKind receiver),
    observe toString (stringTreeLen receiver),
    observe (fun (value, state) => flattenWire value (some state)) (stringTreeString receiver),
    observe toString (stringLength operand),
    observe (fun (value, state) => flattenWire value state) (stringValue operand)]
