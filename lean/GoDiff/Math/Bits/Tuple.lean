import Init.Data.BitVec.Basic

def wordInput (width : Nat) (x : Int) : IO (BitVec width) := do
  if x < 0 || x.toNat ≥ 2 ^ width then
    throw (IO.userError s!"unsigned input out of range: {x}")
  return BitVec.ofNat width x.toNat

def runTuple (evaluate : List Int → IO String) : IO Unit := do
  let stdin ← IO.getStdin
  let stdout ← IO.getStdout
  repeat
    let line ← stdin.getLine
    if line.isEmpty then break
    let some fields := (line.trimAscii.toString.splitOn " ").mapM String.toInt?
      | throw (IO.userError s!"expected decimal integers: {line}")
    stdout.putStrLn (← evaluate fields)
