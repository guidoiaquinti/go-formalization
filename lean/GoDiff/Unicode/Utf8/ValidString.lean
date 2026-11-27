import GoLean.Unicode.Utf8.ValidString
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let width :: words := fields | throw (IO.userError "expected word width and bytes")
  if width != 4 && width != 8 then throw (IO.userError "expected 4- or 8-byte machine word")
  let bytes ← words.mapM fun n => return (← wordInput 8 n).toNat
  return toString (GoLean.Unicode.Utf8.validString width.toNat bytes)
