import GoLean.Math.Bits.Rem

open GoLean.Math.Bits

private def render {w : Nat} : RemOutcome w → String
  | .divideByZero => "divide-by-zero"
  | .value r => toString r.toNat

private def run (width : UintWidth) : IO Unit := do
  let stdin ← IO.getStdin
  let stdout ← IO.getStdout
  repeat
    let line ← stdin.getLine
    if line.isEmpty then break
    let [some hi, some lo, some y] := (line.trimAscii.toString.splitOn " ").map String.toNat?
      | throw (IO.userError s!"expected three unsigned words: {line}")
    if max hi (max lo y) ≥ 2 ^ width.bits then
      throw (IO.userError s!"word out of range: {line}")
    let word := BitVec.ofNat width.bits
    let expected := render (remUint width (word hi) (word lo) (word y))
    let source := (remUintSource width (word hi) (word lo) (word y)).map render
    match source with
    | some result => stdout.putStrLn s!"{expected} {result}"
    | none => stdout.putStrLn s!"{expected} correction-exhausted"

def main (args : List String) : IO Unit :=
  match args with
  | ["32"] => run .bits32
  | ["64"] => run .bits64
  | _ => throw (IO.userError "expected word width 32 or 64")
