import GoLean.Bytes.BufferTryGrowByReslice
import GoDiff.Math.Bits.Tuple

open GoLean.Bytes

private def wire (bytes : List (BitVec 8)) : String :=
  String.intercalate "," (bytes.map fun byte => toString byte.toNat)

private def describe (status : String) (state : BufferReadState) : String :=
  let tail := state.tail
  let marker := if tail.stored.isNone then "nil" else "slice"
  let stored := tail.stored.getD []
  s!"{status}:{state.lastRead}:{tail.offset}:{stored.length + tail.spare.length}:{marker}:" ++
    s!"{wire stored}:{wire tail.spare}"

def main : IO Unit := runTuple fun row => do
  let (receiver, count) ← match row with
    | [0, count] => pure (none, count)
    | tag :: count :: length :: offset :: lastRead :: raw => do
      let bytes ← raw.mapM (wordInput 8)
      let stored := if tag == 1 then none else some (bytes.take length.toNat)
      let tail : BufferTailState := ⟨stored, bytes.drop length.toNat, offset.toNat⟩
      pure (some (⟨tail, lastRead⟩ : BufferReadState), count)
    | _ => throw (IO.userError "expected receiver geometry and grow count")
  match bufferTryGrowByReslice receiver count with
  | .nilReceiver => return "nil-receiver"
  | .bounds state => return describe "bounds" state
  | .returned index ok state => return describe s!"{index}:{ok}" state
