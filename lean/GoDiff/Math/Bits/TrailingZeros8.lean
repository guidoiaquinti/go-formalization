import GoLean.Math.Bits.TrailingZeros8

/-- Row x contains the result for byte x, in ascending unsigned order. -/
def main : IO Unit := do
  let stdout ← IO.getStdout
  for x in [:256] do
    stdout.putStrLn (toString (GoLean.Math.Bits.trailingZeros8 (BitVec.ofNat 8 x)))
