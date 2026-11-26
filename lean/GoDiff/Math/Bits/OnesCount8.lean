import GoLean.Math.Bits.OnesCount8

/-- Row x contains the model's result for byte x, in ascending unsigned order. -/
def main : IO Unit := do
  let stdout ← IO.getStdout
  for x in [:256] do
    stdout.putStrLn (toString (GoLean.Math.Bits.onesCount8 (BitVec.ofNat 8 x)))
