import GoLean.Math.Bits.ReverseBytes16

/-- Row x is the model's result for input x, in ascending unsigned order. -/
def main : IO Unit := do
  let stdout ← IO.getStdout
  for x in [:65536] do
    let result := GoLean.Math.Bits.reverseBytes16 (BitVec.ofNat 16 x)
    stdout.putStrLn (toString result.toNat)
