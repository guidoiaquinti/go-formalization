import GoLean.Image.Color.Gray

def main : IO Unit := do
  let stdout ← IO.getStdout
  for y in [:256] do
    let (r, g, b, a) := GoLean.Image.Color.grayRGBA (BitVec.ofNat 8 y)
    stdout.putStrLn s!"{r.toNat} {g.toNat} {b.toNat} {a.toNat}"
