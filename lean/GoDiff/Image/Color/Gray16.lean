import GoLean.Image.Color.Gray16

/-- Produce the four channel values for each possible uint16 input, in ascending order. -/
def main : IO Unit := do
  let stdout ← IO.getStdout
  for y in [:65536] do
    let (r, g, b, a) := GoLean.Image.Color.gray16RGBA (BitVec.ofNat 16 y)
    stdout.putStrLn s!"{r.toNat} {g.toNat} {b.toNat} {a.toNat}"
