import GoLean.Compress.Flate.TokenState
import GoDiff.Math.Bits.Tuple

open GoLean.Compress.Flate

def main : IO Unit := runTuple fun fields => do
  let n :: filled :: hist :: literal :: bytes := fields |
    throw (IO.userError "expected count, filled count, histogram seed, literal and byte list")
  let n := UInt16.ofNat (← wordInput 16 n).toNat
  let hist := UInt16.ofNat (← wordInput 16 hist).toNat
  let literal := UInt8.ofNat (← wordInput 8 literal).toNat
  let bytes ← bytes.mapM fun b => return UInt8.ofNat (← wordInput 8 b).toNat
  let initial : TokenState := ⟨n, filled, (fun _ => hist), (fun _ => hist), (fun _ => hist),
    fun i => UInt32.ofNat (65536 + i)⟩
  let indices := [0, n.toNat, 65535] ++
    (List.range bytes.length).map (fun i => (n.toNat + i) % 65536)
  let bins := literal :: bytes
  let encode (t : TokenState) :=
    let fields := [toString t.n.toNat, toString t.nFilled, toString (t.literalHist literal.toNat).toNat,
      toString (t.extraHist 0).toNat, toString (t.offsetHist 0).toNat]
    String.intercalate " " (fields ++ bins.map (fun b => toString (t.literalHist b.toNat).toNat) ++
      indices.map (fun i => toString (t.backing i).toNat))
  let view := tokensSlice initial
  let changed := tokenViewWrite view initial 65535 73
  let readAfterChange := (tokenViewRead view changed 65535).getD 0
  let rewritten := { changed with backing := fun i => if i == 65535 then 29 else changed.backing i }
  let readReceiverWrite := (tokenViewRead view rewritten 65535).getD 0
  let readOriginal := (tokenViewRead view initial 65535).getD 0
  return s!"{encode (tokensReset initial)} {encode (tokensAddLiteral initial literal)} {encode (tokensAddEOB initial)} {encode (emitLiterals initial bytes)} {view.length} {view.capacity} {readOriginal.toNat} {changed.backing 65535 |>.toNat} {readAfterChange.toNat} {readReceiverWrite.toNat}"
