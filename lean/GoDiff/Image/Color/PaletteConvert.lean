import GoLean.Image.Color.PaletteConvert
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

private def parseEntries : List Int → IO (List (PaletteCall × Option Nat))
  | [] => pure []
  | token :: kind :: r :: g :: b :: a :: tail => do
    let (call, _) ← parseCall [kind, r, g, b, a]
    let entry := (call, if token = 0 then none else some token.toNat)
    return entry :: (← parseEntries tail)
  | _ => throw (IO.userError "incomplete palette entry")

def main : IO Unit := runTuple fun fields => do
  let (query, tail) ← parseCall fields
  let (calls, result) := paletteConvert query (← parseEntries tail)
  let result := match result with
    | .ok none => "nil"
    | .ok (some token) => s!"color {token}"
    | .error payload => s!"panic {payload}"
  return String.intercalate " " (result :: calls.map toString)
