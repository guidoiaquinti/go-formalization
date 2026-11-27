import GoLean.Unicode.Utf8.Word
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let width :: values := fields | throw (IO.userError "expected word width and bytes")
  if width != 4 && width != 8 then throw (IO.userError "expected 4- or 8-byte word")
  let bytes ← values.mapM (wordInput 8)
  return ((GoLean.Unicode.Utf8.word width.toNat bytes).map toString).getD "bounds"
