import GoLean.Math.Bits.Length

/-- Row x is the model's result for input x, in ascending unsigned order. -/
def main : IO Unit := do
  let stdout ← IO.getStdout
  for x in [:256] do
    stdout.putStrLn (toString (GoLean.Math.Bits.len8 (BitVec.ofNat 8 x)))
