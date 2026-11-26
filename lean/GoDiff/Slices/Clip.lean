import GoLean.Slices.Clip
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let nilFlag :: spare :: elements := fields | throw (IO.userError "expected nil flag, spare capacity and elements")
  let view := GoLean.Slices.sliceClip ⟨elements, elements.length + spare.toNat, nilFlag == 1⟩
  let header := s!"{view.isNil} {view.capacity}"
  return String.intercalate " " (header :: view.elements.map toString)
