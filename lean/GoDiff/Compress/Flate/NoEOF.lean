import GoLean.Compress.Flate.NoEOF
import GoDiff.Math.Bits.Tuple

open GoLean.Compress.Flate

def main : IO Unit := runTuple fun row => do
  let [tag] := row | throw (IO.userError "expected error identity")
  let error : ReadError := match tag.toNat with
    | 0 => .none
    | 1 => .eof
    | 2 => .unexpectedEOF
    | other => .other other
  return match noEOF error with
    | .none => "0"
    | .eof => "1"
    | .unexpectedEOF => "2"
    | .other identity => toString identity
