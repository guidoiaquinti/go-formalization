import GoLean.Compress.Flate.HCode
import GoLean.Compress.Flate.ReverseBits

namespace GoLean.Compress.Flate

structure FixedHuffmanTable where
  codes : List HCode
  capacity : Nat
  zeroWorkspace : Bool

private def packed (code size : Nat) : HCode :=
  let length := BitVec.ofNat 8 size
  newHCode (reverseBits (BitVec.ofNat 16 code) length) length

private def literal (ch : Nat) : HCode :=
  if ch < 144 then packed (ch + 48) 8
  else if ch < 256 then packed (ch + 256) 9
  else if ch < 280 then packed (ch - 256) 7
  else packed (ch - 88) 8

/-- The implementation stores symbols 0 through 285 and reserves 512 entries. -/
def generateFixedLiteralEncoding : FixedHuffmanTable :=
  ⟨(List.range 286).map literal, 512, true⟩

/-- Only the 30 usable distance symbols appear in the fixed Go encoder. -/
def generateFixedOffsetEncoding : FixedHuffmanTable :=
  ⟨(List.range 30).map (fun ch => packed ch 5), 32, true⟩

theorem generateFixedLiteralEncoding_geometry :
    generateFixedLiteralEncoding.codes.length = 286 ∧
    generateFixedLiteralEncoding.capacity = 512 ∧
    generateFixedLiteralEncoding.zeroWorkspace = true := by
  simp [generateFixedLiteralEncoding]

theorem generateFixedLiteralEncoding_transitions :
    [143, 144, 255, 256, 279, 280, 285].map
      (fun i => (hCodeLen generateFixedLiteralEncoding.codes[i]!).toNat) =
      [8, 9, 9, 7, 7, 8, 8] := by
  simp +decide

theorem generateFixedOffsetEncoding_geometry :
    generateFixedOffsetEncoding.codes.length = 30 ∧
    generateFixedOffsetEncoding.capacity = 32 ∧
    generateFixedOffsetEncoding.zeroWorkspace = true := by decide

theorem generateFixedOffsetEncoding_lengths :
    generateFixedOffsetEncoding.codes.all (fun word => hCodeLen word == 5) = true := by decide

end GoLean.Compress.Flate
