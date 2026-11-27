import Init.Data.BitVec.Lemmas
import Init.Omega

namespace GoLean.Archive.Zip

private def modeDir : BitVec 32 := 0x80000000
private def modeSymlink : BitVec 32 := 0x08000000
private def modeDevice : BitVec 32 := 0x04000000
private def modePipe : BitVec 32 := 0x02000000
private def modeSocket : BitVec 32 := 0x01000000
private def modeSetuid : BitVec 32 := 0x00800000
private def modeSetgid : BitVec 32 := 0x00400000
private def modeChar : BitVec 32 := 0x00200000
private def modeSticky : BitVec 32 := 0x00100000
private def modeIrregular : BitVec 32 := 0x00080000

def msdosModeToFileMode (attrs : BitVec 32) : BitVec 32 :=
  let mode := if attrs &&& 16 != 0 then modeDir ||| 0o777 else 0o666
  if attrs &&& 1 != 0 then mode &&& ~~~0o222 else mode

def unixModeToFileMode (attrs : BitVec 32) : BitVec 32 :=
  let typ := attrs &&& 0xf000
  let mode :=
    if typ == 0x6000 then modeDevice
    else if typ == 0x2000 then modeDevice ||| modeChar
    else if typ == 0x4000 then modeDir
    else if typ == 0x1000 then modePipe
    else if typ == 0xa000 then modeSymlink
    else if typ == 0xc000 then modeSocket
    else 0
  mode ||| (attrs &&& 0o777) |||
    (if attrs &&& 0x800 != 0 then modeSetuid else 0) |||
    (if attrs &&& 0x400 != 0 then modeSetgid else 0) |||
    (if attrs &&& 0x200 != 0 then modeSticky else 0)

def fileModeToUnixMode (mode : BitVec 32) : BitVec 32 :=
  let typ := mode &&&
    (modeDir ||| modeSymlink ||| modeDevice ||| modePipe ||| modeSocket ||| modeChar ||| modeIrregular)
  let attrs : BitVec 32 :=
    if typ == modeDir then 0x4000
    else if typ == modeSymlink then 0xa000
    else if typ == modePipe then 0x1000
    else if typ == modeSocket then 0xc000
    else if typ == modeDevice then 0x6000
    else if typ == (modeDevice ||| modeChar) then 0x2000
    else 0x8000
  (attrs &&& 0xf000) ||| (mode &&& 0o777) |||
    (if mode &&& modeSetuid != 0 then 0x800 else 0) |||
    (if mode &&& modeSetgid != 0 then 0x400 else 0) |||
    (if mode &&& modeSticky != 0 then 0x200 else 0)

def fileHeaderMode (creatorVersion : BitVec 16) (attrs : BitVec 32)
    (nameEndsInSlash : Bool) : BitVec 32 :=
  let creator := creatorVersion >>> 8
  let mode :=
    if creator == 3 || creator == 19 then unixModeToFileMode (attrs >>> 16)
    else if creator == 11 || creator == 14 || creator == 0 then msdosModeToFileMode attrs
    else 0
  if nameEndsInSlash then mode ||| modeDir else mode

def fileHeaderSetMode (creatorVersion : BitVec 16) (mode : BitVec 32) :
    BitVec 16 × BitVec 32 :=
  ((creatorVersion &&& 255) ||| 0x300,
    (fileModeToUnixMode mode <<< 16) |||
    (if mode &&& modeDir != 0 then 16 else 0) |||
    (if mode &&& 0o200 == 0 then 1 else 0))

def fileHeaderHasDataDescriptor (flags : BitVec 16) : Bool :=
  flags &&& 8 != 0

private theorem bits16 (i : Nat) (h : i < 16) :
    i = 0 ∨ i = 1 ∨ i = 2 ∨ i = 3 ∨ i = 4 ∨ i = 5 ∨ i = 6 ∨ i = 7 ∨
    i = 8 ∨ i = 9 ∨ i = 10 ∨ i = 11 ∨ i = 12 ∨ i = 13 ∨ i = 14 ∨ i = 15 := by
  omega

private theorem maskedIf (p : Prop) [Decidable p] (x y m : BitVec w) :
    (if p then x else y) &&& m = (if p then x &&& m else y &&& m) := by
  split <;> rfl

theorem fileHeaderSetMode_creator (creator : BitVec 16) (mode : BitVec 32) :
    (fileHeaderSetMode creator mode).1 >>> 8 = 3 ∧
    ((fileHeaderSetMode creator mode).1 &&& 255) = (creator &&& 255) := by
  constructor <;> apply BitVec.eq_of_getLsbD_eq_iff.mpr <;> intro i hi <;>
    rcases bits16 i hi with rfl | rfl | rfl | rfl | rfl | rfl | rfl | rfl |
      rfl | rfl | rfl | rfl | rfl | rfl | rfl | rfl <;>
    simp [fileHeaderSetMode]

theorem fileHeaderHasDataDescriptor_bit (flags : BitVec 16) :
    fileHeaderHasDataDescriptor flags = flags.getLsbD 3 := by
  have maskEq : (flags &&& 8 = 0) ↔ flags.getLsbD 3 = false := by
    constructor
    · intro h
      have hbit := congrArg (fun a : BitVec 16 => a.getLsbD 3) h
      simpa using hbit
    · intro h
      apply BitVec.eq_of_getLsbD_eq_iff.mpr
      intro i hi
      simp at h
      rcases bits16 i hi with rfl | rfl | rfl | rfl | rfl | rfl | rfl | rfl |
        rfl | rfl | rfl | rfl | rfl | rfl | rfl | rfl <;> simp [h]
  cases h : flags.getLsbD 3
  · have hz := maskEq.mpr h
    simpa [fileHeaderHasDataDescriptor] using hz
  · have hn : flags &&& (8 : BitVec 16) ≠ 0 := fun hz => by
      have hb := maskEq.mp hz
      rw [h] at hb
      contradiction
    simpa [fileHeaderHasDataDescriptor] using hn

theorem fileHeaderMode_slash (creator : BitVec 16) (attrs : BitVec 32) :
    (fileHeaderMode creator attrs true).getLsbD 31 = true := by
  simp [fileHeaderMode, modeDir]

theorem msdosModeToFileMode_mask (attrs : BitVec 32) :
    msdosModeToFileMode attrs = msdosModeToFileMode (attrs &&& 17) := by
  simp [msdosModeToFileMode, BitVec.and_assoc]

theorem unixModeToFileMode_permissions (attrs : BitVec 32) :
    unixModeToFileMode attrs &&& 0o777 = attrs &&& 0o777 := by
  simp [unixModeToFileMode, BitVec.and_or_distrib_right, BitVec.and_assoc, maskedIf,
    modeDevice, modeChar, modeDir, modePipe, modeSymlink, modeSocket,
    modeSetuid, modeSetgid, modeSticky]

theorem fileModeToUnixMode_permissions (mode : BitVec 32) :
    fileModeToUnixMode mode &&& 0o777 = mode &&& 0o777 := by
  simp [fileModeToUnixMode, BitVec.and_or_distrib_right, BitVec.and_assoc, maskedIf,
    modeDir, modeSymlink, modeDevice, modePipe, modeSocket, modeChar,
    modeIrregular, modeSetuid, modeSetgid, modeSticky]

end GoLean.Archive.Zip
