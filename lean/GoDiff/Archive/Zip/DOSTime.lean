import GoLean.Archive.Zip.DOSTime
import GoDiff.Math.Bits.Tuple

open GoLean.Archive.Zip

def main : IO Unit := runTuple fun row => do
  let [d, t, seconds, nanos, ly, lm, ld, lh, lmin, ls, uy, um, ud, uh, umin, us] := row |
    throw (IO.userError "expected DOS words and local/UTC civil fields")
  let date ← wordInput 16 d
  let clock ← wordInput 16 t
  let (localDate, localTime) := timeToMsDosTime ⟨ly, lm, ld, lh, lmin, ls⟩
  let set := fileHeaderSetModTime seconds nanos.toNat ⟨uy, um, ud, uh, umin, us⟩
  return s!"{msDosTimeToTime date clock}|{localDate.toNat}:{localTime.toNat}|{fileHeaderModTime date clock}|{set.seconds}:{set.nanoseconds}:{set.dateWord.toNat}:{set.timeWord.toNat}"
