import GoLean.Compress.Flate.Histogram
import GoDiff.Math.Bits.Tuple

open GoLean.Compress.Flate

def main : IO Unit := runTuple fun fields => do
  let cap :: words := fields | throw (IO.userError "expected capacity and backing")
  if cap < 0 || cap.toNat > words.length then
    throw (IO.userError "invalid backing capacity")
  let cap := cap.toNat
  let backing ← (words.take cap).toArray.mapM fun n => do
    let w ← wordInput 16 n
    return UInt16.ofNat w.toNat
  let data ← (words.drop cap).toArray.mapM fun n => do
    let w ← wordInput 8 n
    return UInt8.ofNat w.toNat
  let encode (out : HistogramResult) :=
    String.intercalate " " ((if out.panicked then "1" else "0") ::
      out.backing.toList.map (toString ·.toNat))
  return s!"{encode (histogram data backing)} {encode (histogramSplit data backing)}"
