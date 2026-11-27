import GoLean.Bytes.BufferGrow
import GoDiff.Math.Bits.Tuple

open GoLean.Bytes

private def wire (bytes : List (BitVec 8)) : String :=
  String.intercalate "," (bytes.map fun byte => toString byte.toNat)

private def describe (status : String) (state : BufferReadState) (reused : Bool) : String :=
  let tail := state.tail
  let marker := if tail.stored.isNone then "nil" else "slice"
  let stored := tail.stored.getD []
  s!"{status}:{state.lastRead}:{tail.offset}:{stored.length + tail.spare.length}:{marker}:" ++
    s!"{reused}:{wire stored}:{wire tail.spare}"

def main : IO Unit := runTuple fun row => do
  let (receiver, count, maxInt) ← match row with
    | [0, count, maxInt] => pure (none, count, maxInt.toNat)
    | tag :: count :: maxInt :: length :: offset :: lastRead :: raw => do
      let bytes ← raw.mapM (wordInput 8)
      let stored := if tag == 1 then none else some (bytes.take length.toNat)
      let tail : BufferTailState := ⟨stored, bytes.drop length.toNat, offset.toNat⟩
      pure (some (⟨tail, lastRead⟩ : BufferReadState), count, maxInt.toNat)
    | _ => throw (IO.userError "expected receiver geometry and grow count")
  match bufferGrow receiver count maxInt with
  | .negative =>
    match receiver with
    | none => return "negative"
    | some state =>
      let reused := !(state.tail.stored.getD []).isEmpty || !state.tail.spare.isEmpty
      return describe "negative" state reused
  | .nilReceiver => return "nil-receiver"
  | .tooLarge state => return describe "too-large" state true
  | .success result => return describe "ok" result.state result.reused
