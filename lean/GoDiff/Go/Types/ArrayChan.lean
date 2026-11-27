import GoLean.Go.Types.ArrayChan
import GoDiff.Math.Bits.Tuple

open GoLean.Go.Types

private def refWire : TypeRef → String
  | .nilInterface => "nil"
  | .typedNil 1 => "nilArray"
  | .typedNil 2 => "nilChan"
  | .typedNil _ => "otherNil"
  | .value i => toString i

private def readWire (render : α → String) : Except ReadFault α → String
  | .error .nilReceiver => "panic"
  | .ok value => render value

def main : IO Unit := runTuple fun fields => do
  let [width, len, dir, element, mode] := fields |
    throw (IO.userError "expected width, length, direction, element and receiver mode")
  let elem := if element == 0 then TypeRef.nilInterface
    else if element == 1 then TypeRef.typedNil 1
    else if element == 2 then TypeRef.typedNil 2 else TypeRef.value element.toNat
  let a := newArray 31 elem (BitVec.ofInt 64 len)
  let c := newChan 32 (BitVec.ofInt width.toNat dir) elem
  let a' := newArray 33 elem (BitVec.ofInt 64 len)
  let c' := newChan 34 (BitVec.ofInt width.toNat dir) elem
  let ar := if mode == 0 then none else if mode == 2 then some (newArray 31 .nilInterface 0) else some a
  let cr := if mode == 0 then none else if mode == 2 then some (newChan 32 0 .nilInterface) else some c
  return String.intercalate "|" [
    toString (a.identity != a'.identity), toString a.len.toInt, refWire a.elem,
    toString (c.identity != c'.identity), toString c.dir.toInt, refWire c.elem,
    readWire (toString ∘ BitVec.toInt) (arrayLen ar), readWire refWire (arrayElem ar), refWire (arrayUnderlying ar),
    readWire (toString ∘ BitVec.toInt) (chanDir cr), readWire refWire (chanElem cr), refWire (chanUnderlying cr)]
