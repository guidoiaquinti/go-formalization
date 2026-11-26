import GoLean.Math.Bits.TrailingZeros

/-- Row x is the model's result for input x, in ascending unsigned order. -/
def main : IO Unit := do
  let stdout ← IO.getStdout
  for x in [:65536] do
    stdout.putStrLn (toString (GoLean.Math.Bits.trailingZeros16 (BitVec.ofNat 16 x)))
