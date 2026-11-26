import GoLean.Math.Bits.Length

/-- Row x is the model's result for input x, in ascending unsigned order. -/
def main : IO Unit := do
  let stdout ← IO.getStdout
  for x in [:65536] do
    stdout.putStrLn (toString (GoLean.Math.Bits.len16 (BitVec.ofNat 16 x)))
