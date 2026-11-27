import GoLean.Image.Color.PaletteIndex
import GoDiff.Math.Bits.Tuple

open GoLean.Image.Color

private def parseCall : List Int → IO (PaletteCall × List Int)
  | kind :: r :: g :: b :: a :: tail => do
    let call ← match kind with
      | 0 => pure (some (.ok (← wordInput 32 r, ← wordInput 32 g,
          ← wordInput 32 b, ← wordInput 32 a)))
      | 1 => pure (some (.error r.toNat))
      | 2 => pure none
      | _ => throw (IO.userError "unexpected callback kind")
    return (call, tail)
  | _ => throw (IO.userError "incomplete callback")

private def parsePalette (n : Nat) (fields : List Int) : IO (List PaletteCall) := do
  match n with
  | 0 =>
    unless fields.isEmpty do throw (IO.userError "unexpected trailing input")
    return []
  | n + 1 =>
    let (call, tail) ← parseCall fields
    return call :: (← parsePalette n tail)

def main : IO Unit := runTuple fun fields => do
  let n :: fields := fields | throw (IO.userError "expected palette size")
  let (query, tail) ← parseCall fields
  let (calls, result) := paletteIndex query (← parsePalette n.toNat tail)
  let result := match result with
    | .ok index => s!"index {index}"
    | .error payload => s!"panic {payload}"
  return String.intercalate " " (result :: calls.map toString)
