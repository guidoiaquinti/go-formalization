import GoLean.Math.Bits.Arithmetic

def main : IO Unit := do
  let stdin ← IO.getStdin
  let stdout ← IO.getStdout
  repeat
    let line ← stdin.getLine
    if line.isEmpty then break
    let [some x, some y] := (line.trimAscii.toString.splitOn " ").map String.toNat?
      | throw (IO.userError s!"expected two uint64 inputs: {line}")
    if max x y ≥ 2^64 then throw (IO.userError s!"word out of range: {line}")
    let (hi, lo) := GoLean.Math.Bits.mul64Source (BitVec.ofNat 64 x) (BitVec.ofNat 64 y)
    stdout.putStrLn s!"{hi.toNat} {lo.toNat}"
