import GoLean.Bytes.Join
import GoDiff.Math.Bits.Tuple

private def parts : List Nat → List (List Nat) → List Nat → List (List Nat)
  | [], done, _ => done.reverse
  | size :: sizes, done, raw => parts sizes (raw.take size :: done) (raw.drop size)

def main : IO Unit := runTuple fun row => do
  let maxInt :: count :: sepSize :: fields := row | throw (IO.userError "expected int max, counts and bytes")
  let sizes := (fields.take count.toNat).map Int.toNat
  let bytes ← (fields.drop count.toNat).mapM (wordInput 8)
  let bytes := bytes.map BitVec.toNat
  let sep := bytes.take sepSize.toNat
  let elems := parts sizes [] (bytes.drop sepSize.toNat)
  match GoLean.Bytes.join maxInt.toNat elems sep with
  | none => return "overflow"
  | some none => return "nil"
  | some (some result) => return "bytes" ++ String.intercalate "," (result.map toString)
