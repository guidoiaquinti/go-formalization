import GoLean.Go.Types.ContainerTypes
import GoDiff.Math.Bits.Tuple

open GoLean.Go.Types

private def inputRef (tag : Int) : TypeRef :=
  if tag == 0 then .nilInterface
  else if tag < 6 then .typedNil tag.toNat
  else .value tag.toNat

private def render : TypeRef → String
  | .nilInterface => "nil"
  | .typedNil tag => s!"typed:{tag}"
  | .value id => toString id

private def renderField : Except ReadFault TypeRef → String
  | .error .nilReceiver => "panic"
  | .ok ref => render ref

private def receiver (mode : Int) (value zero : α) : Option α :=
  if mode == 0 then none else if mode == 2 then some zero else some value

def main : IO Unit := runTuple fun fields => do
  let [keyTag, elemTag, mode] := fields | throw (IO.userError "expected key, element and receiver tags")
  let key := inputRef keyTag
  let elem := inputRef elemTag
  let s := newSlice 31 elem
  let p := newPointer 32 elem
  let m := newMap 33 key elem
  let sr := receiver mode s (newSlice 31 .nilInterface)
  let pr := receiver mode p (newPointer 32 .nilInterface)
  let mr := receiver mode m (newMap 33 .nilInterface .nilInterface)
  return String.intercalate "|" [
    toString (s.identity != (newSlice 41 elem).identity), render s.elem,
    toString (p.identity != (newPointer 42 elem).identity), render p.base,
    toString (m.identity != (newMap 43 key elem).identity), render m.key, render m.elem,
    renderField (sliceElem sr), render (sliceUnderlying sr),
    renderField (pointerElem pr), render (pointerUnderlying pr),
    renderField (mapKey mr), renderField (mapElem mr), render (mapUnderlying mr)]
