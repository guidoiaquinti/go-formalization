import GoLean.Strings.SplitSeqPrivate
import GoDiff.Math.Bits.Tuple

open GoLean.Strings

private def describe (run : ConfiguredSplitRun) : String :=
  (if run.boundsPanic then "p:" else "o:") ++
    String.intercalate ";" (run.parts.map fun part => "f" ++ String.intercalate "," (part.map toString))

def main : IO Unit := runTuple fun row => do
  let size :: save :: first :: second :: raw := row |
    throw (IO.userError "expected separator size, save, budgets and bytes")
  let bytes ← raw.mapM (wordInput 8)
  let values := bytes.map BitVec.toNat
  let separator := values.take size.toNat
  let input := values.drop size.toNat
  let one := splitSeqConfigured input separator save first.toNat
  let two := splitSeqConfigured one.remaining separator save second.toNat
  let full := splitSeqConfigured two.remaining separator save (input.length + 2)
  let after := splitSeqConfigured full.remaining separator save 1
  return String.intercalate "|" ([one, two, full, after].map describe)
