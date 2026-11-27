import GoLean.Image.Color.Alpha16

def main : IO Unit := do
  let stdout ← IO.getStdout
  for input in [:65536] do
    let (r, g, b, a) := GoLean.Image.Color.alpha16RGBA (BitVec.ofNat 16 input)
    stdout.putStrLn s!"{r.toNat} {g.toNat} {b.toNat} {a.toNat}"
