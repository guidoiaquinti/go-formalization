import Init.Data.BitVec.Basic

/-- Read one unsigned decimal word and emit one decimal result per line. -/
def runUnary (width : Nat) (model : BitVec width → Nat) : IO Unit := do
  let stdin ← IO.getStdin
  let stdout ← IO.getStdout
  repeat
    let line ← stdin.getLine
    if line.isEmpty then break
    let some x := line.trimAscii.toString.toNat?
      | throw (IO.userError s!"expected an unsigned decimal word: {line}")
    if x ≥ 2 ^ width then throw (IO.userError s!"word out of range: {line}")
    stdout.putStrLn (toString (model (BitVec.ofNat width x)))
