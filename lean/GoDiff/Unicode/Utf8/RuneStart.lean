import GoLean.Unicode.Utf8.RuneStart

def main : IO Unit := do
  let stdout ← IO.getStdout
  for b in [:256] do
    let result := GoLean.Unicode.Utf8.runeStart (BitVec.ofNat 8 b)
    stdout.putStrLn (if result then "1" else "0")
