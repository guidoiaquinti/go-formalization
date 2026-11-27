import GoLean.Bytes.ByteSearch
import GoDiff.Math.Bits.Tuple

private def index (value : Option Nat) : String :=
  value.map toString |>.getD "-1"

def main : IO Unit := runTuple fun fields => do
  let needle :: words := fields | throw (IO.userError "expected a needle and bytes")
  let key := (← wordInput 8 needle).toNat.toUInt8
  let data ← words.mapM fun n => do
    return (← wordInput 8 n).toNat.toUInt8
  return s!"{index (GoLean.Bytes.indexByte data key)} {index (GoLean.Bytes.lastIndexByte data key)}"
