import GoLean.Image.Color.Alpha

def main : IO Unit := do
  let stdout ← IO.getStdout
  for input in [:256] do
    let (r, g, b, a) := GoLean.Image.Color.alphaRGBA (BitVec.ofNat 8 input)
    stdout.putStrLn s!"{r.toNat} {g.toNat} {b.toNat} {a.toNat}"
