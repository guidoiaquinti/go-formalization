import GoLean.Index.Suffixarray.Buckets
import GoDiff.Math.Bits.Tuple

open GoLean.Index.Suffixarray

private def runs : List (BitVec width) → List (BitVec width × Nat)
  | [] => []
  | x :: xs =>
    match runs xs with
    | (y, n) :: rest => if x = y then (x, n + 1) :: rest else (x, 1) :: (y, n) :: rest
    | [] => [(x, 1)]

private def wire (v : View (BitVec width)) : String :=
  match v.backing with
  | none => "nil"
  | some data => s!"{v.length}:{String.intercalate "," ((runs data).map fun (x,n) => s!"{x.toInt}.{n}")}"

private def render (status : String) (s : BucketState width) :=
  status ++ "/" ++ wire s.frequency ++ "/" ++ wire s.bucket

private def suite
    (freq : List Int → BucketState width → Option (View (BitVec width)) × BucketState width)
    (lower upper : List Int → BucketState width → Bool × BucketState width)
    (text : List Int) (state : BucketState width) : String :=
  let (f, fs) := freq text state
  let (l, ls) := lower text state
  let (u, us) := upper text state
  String.intercalate "|" [
    render (f.map wire |>.getD "panic") fs,
    render (if l then "ok" else "panic") ls,
    render (if u then "ok" else "panic") us]

def main : IO Unit := runTuple fun fields => do
  let profile :: fsize :: fcap :: bsize :: bcap :: textSize :: cells := fields |
    throw (IO.userError "expected profile, slice geometries, and cells")
  let text := cells.take textSize.toNat
  let data := cells.drop textSize.toNat
  let state (width : Nat) : BucketState width :=
    let f := (data.take fcap.toNat).map (BitVec.ofInt width)
    let b := ((data.drop fcap.toNat).take bcap.toNat).map (BitVec.ofInt width)
    ⟨⟨if fsize < 0 then none else some f, fsize.toNat⟩,
      ⟨if bsize < 0 then none else some b, bsize.toNat⟩⟩
  return match profile with
  | 0 => suite freqByte32 bucketMinByte32 bucketMaxByte32 (text.map fun n => (Int.ofNat (BitVec.ofInt 8 n).toNat)) (state 32)
  | 1 => suite freqByte64 bucketMinByte64 bucketMaxByte64 (text.map fun n => (Int.ofNat (BitVec.ofInt 8 n).toNat)) (state 64)
  | 2 => suite freq32 bucketMin32 bucketMax32 (text.map fun n => (BitVec.ofInt 32 n).toInt) (state 32)
  | _ => suite freq64 bucketMin64 bucketMax64 (text.map fun n => (BitVec.ofInt 64 n).toInt) (state 64)
