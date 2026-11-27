import GoLean.Unicode.Utf16.DecodeAppend
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let _ :: length :: count :: rest := fields | throw (IO.userError "expected capacity and input lengths")
  if length < -1 || count < -1 || rest.length != length.toNat + count.toNat then
    throw (IO.userError "invalid destination length")
  let dst := rest.take length.toNat
  let words ← (rest.drop length.toNat).mapM fun word => do
    return (← wordInput 16 word).toNat
  return String.intercalate " " ("runes" :: (GoLean.Unicode.Utf16.decodeAppend words dst).map toString)
