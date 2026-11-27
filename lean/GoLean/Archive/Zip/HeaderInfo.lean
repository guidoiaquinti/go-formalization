import GoLean.Archive.Zip.FileMode
import Init.Data.BitVec.Lemmas

namespace GoLean.Archive.Zip

/-- Stable abstract identity of the header that the adapter reads. -/
structure HeaderInfo where
  header : Option Nat

def fileHeaderFileInfo (header : Option Nat) : HeaderInfo := ⟨header⟩

def headerInfoName (name : List Nat) : List Nat :=
  if name.isEmpty then [46] else
    let tail := name.reverse.dropWhile (· == 47)
    if tail.isEmpty then [47] else (tail.takeWhile (· != 47)).reverse

def headerInfoSize (legacy : BitVec 32) (wide : BitVec 64) : Int :=
  if wide == 0 then legacy.toNat else wide.toInt

def headerInfoMode (creator : BitVec 16) (attrs : BitVec 32) (name : List Nat) : BitVec 32 :=
  fileHeaderMode creator attrs (name.getLast? == some 47)

def headerInfoType (creator : BitVec 16) (attrs : BitVec 32) (name : List Nat) : BitVec 32 :=
  headerInfoMode creator attrs name &&& 0x8f280000

def headerInfoIsDir (creator : BitVec 16) (attrs : BitVec 32) (name : List Nat) : Bool :=
  (headerInfoMode creator attrs name).getLsbD 31

def headerInfoSys (info : HeaderInfo) : Option Nat := info.header

def headerInfoInfo (info : HeaderInfo) : HeaderInfo × Option String := (info, none)

theorem fileHeaderFileInfo_retains (header : Option Nat) :
    (fileHeaderFileInfo header).header = header := rfl

theorem headerInfoName_edges :
    headerInfoName [] = [46] ∧ headerInfoName [47, 47] = [47] ∧
    headerInfoName [97, 47, 255, 47] = [255] := ⟨rfl, rfl, rfl⟩

theorem headerInfoSize_legacy (legacy : BitVec 32) :
    headerInfoSize legacy 0 = Int.ofNat legacy.toNat := rfl

theorem headerInfoMode_delegates (creator : BitVec 16) (attrs : BitVec 32) (name : List Nat) :
    headerInfoMode creator attrs name = fileHeaderMode creator attrs (name.getLast? == some 47) := rfl

theorem headerInfoType_masks (creator : BitVec 16) (attrs : BitVec 32) (name : List Nat) :
    headerInfoType creator attrs name &&& 0x8f280000 = headerInfoType creator attrs name := by
  simp [headerInfoType, BitVec.and_assoc]

theorem headerInfoIsDir_directory (creator : BitVec 16) (attrs : BitVec 32) :
    headerInfoIsDir creator attrs [47] = true := fileHeaderMode_slash creator attrs

theorem headerInfoSys_retains (info : HeaderInfo) : headerInfoSys info = info.header := rfl

theorem headerInfoInfo_retains (info : HeaderInfo) : headerInfoInfo info = (info, none) := rfl

end GoLean.Archive.Zip
