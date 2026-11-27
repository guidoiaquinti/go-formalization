import GoLean.Index.Suffixarray.IntegerIO
import GoDiff.Math.Bits.Tuple

open GoLean.Index.Suffixarray

private def wire (bytes : List (BitVec 8)) := String.intercalate "," (bytes.map (toString ∘ BitVec.toNat))

private def writer (r : IntegerWrite) : String :=
  let offer := r.offered.map (fun (n, cap, b) => s!"{n},{cap}:{wire b}") |>.getD "none"
  s!"{r.status}/{r.count}/{wire r.backing}/{offer}"

def main : IO Unit := runTuple fun fields => do
  let width :: length :: cap :: value :: chunk :: terminal :: reply :: srcSize :: profile :: raw := fields |
    throw (IO.userError "expected slice and I/O profiles")
  let backing ← (raw.take cap.toNat).mapM (wordInput 8)
  let src ← ((raw.drop cap.toNat).take srcSize.toNat).mapM (wordInput 8)
  let integers := raw.drop (cap + srcSize).toNat
  let data : Ints := if profile == 0 then ⟨⟨some (integers.map (BitVec.ofInt 32)), integers.length⟩, ⟨none, 0⟩⟩
    else ⟨⟨none, 0⟩, ⟨some (integers.map (BitVec.ofInt 64)), integers.length⟩⟩
  let readStatus := if terminal == 0 then "eof" else if terminal == 1 then "other" else "nilIO"
  let reply := if reply == 0 then "ok" else if reply == 1 then "other" else "nilIO"
  let r := readInt backing length.toNat src chunk.toNat readStatus
  let x := BitVec.ofInt 64 ((BitVec.ofInt width.toNat value).toInt)
  let read := s!"{r.status}/{r.value}/{wire r.backing}/{String.intercalate "," (r.requests.map toString)}"
  return writer (writeInt backing length.toNat x reply) ++ "|" ++ read ++ "|" ++
    writer (writeSlice backing length.toNat data reply)
