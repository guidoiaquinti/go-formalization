import GoLean.Bytes.ReaderSetup
import GoDiff.Math.Bits.Tuple

open GoLean.Bytes

private def wire (xs : List (BitVec 8)) : String :=
  String.intercalate "," (xs.map (toString ∘ BitVec.toNat))

private def describe (state : ReaderState) : String :=
  let marker := if state.bytes.isNone then "nil" else "slice"
  s!"{marker}:{state.cursor}:{state.previousRune}:{wire (state.bytes.getD [])}:{wire state.spare}"

def main : IO Unit := runTuple fun row => do
  let tag :: length :: cursor :: previous :: raw := row | throw (IO.userError "expected slice and old reader state")
  let data ← raw.mapM (wordInput 8)
  let input := if tag == 0 then none else some (data.take length.toNat)
  let spare := data.drop length.toNat
  let state := newReader input spare
  let prior := {state with cursor := cursor.toNat, previousRune := previous}
  let reset := (readerReset (some prior) input spare).map describe |>.getD "panic"
  let nilReset := (readerReset none input spare).map describe |>.getD "panic"
  let size := (readerSize (some prior)).map toString |>.getD "panic"
  let nilSize := (readerSize none).map toString |>.getD "panic"
  return s!"{describe state}|{reset}|{nilReset}|{size}|{nilSize}"
