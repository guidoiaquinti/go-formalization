import GoLean.Bytes.Map
import GoLean.Strings.ToLower
import GoLean.Strings.ToUpper
import GoLean.Strings.ToTitle

namespace GoLean.Bytes

structure CaseBytes where
  bytes : List Nat
  exactCapacity : Option Nat

def mapCaseBytes (mapping : Nat → Nat) (input : List (BitVec 8)) : CaseBytes :=
  let output := ((mapRunes (fun r => Int.ofNat (mapping r)) input).2).getD []
  ⟨output, if output.length ≤ input.length then some input.length else none⟩

private def asciiRange (lo hi : Nat) (byte : BitVec 8) : Bool :=
  decide (lo ≤ byte.toNat ∧ byte.toNat ≤ hi)

private def changeASCII (lo hi : Nat) (adjust mapping : Nat → Nat)
    (input : List (BitVec 8)) : CaseBytes :=
  if input.all (fun b => b.toNat < 128) then
    if input.any (asciiRange lo hi) then
      ⟨input.map (fun b => if asciiRange lo hi b then adjust b.toNat else b.toNat), some input.length⟩
    else
      ⟨input.map BitVec.toNat, if input.isEmpty then some 0 else none⟩
  else mapCaseBytes mapping input

/-- Go 1.27.0 src/bytes/bytes.go:679–706. ASCII transformations use
    a length-clipped result. Unchanged ASCII input uses append to copy. -/
def toUpperBytes (input : List (BitVec 8)) : CaseBytes :=
  changeASCII 97 122 (· - 32) GoLean.Strings.upperRune input

/-- Go 1.27.0 src/bytes/bytes.go:710–735. Non-ASCII input delegates
    to the byte-slice Map implementation. -/
def toLowerBytes (input : List (BitVec 8)) : CaseBytes :=
  changeASCII 65 90 (· + 32) GoLean.Strings.lowerRune input

/-- Go 1.27.0 src/bytes/bytes.go:739. Titlecase is applied to every rune. -/
def toTitleBytes (input : List (BitVec 8)) : CaseBytes :=
  mapCaseBytes GoLean.Strings.titleRune input

theorem toUpperBytes_ascii (input : List (BitVec 8))
    (ascii : input.all (fun b => b.toNat < 128) = true)
    (changes : input.any (asciiRange 97 122) = true) :
    toUpperBytes input =
      ⟨input.map (fun b => if asciiRange 97 122 b then b.toNat - 32 else b.toNat), some input.length⟩ := by
  simp [toUpperBytes, changeASCII, ascii, changes]

theorem toLowerBytes_ascii (input : List (BitVec 8))
    (ascii : input.all (fun b => b.toNat < 128) = true)
    (changes : input.any (asciiRange 65 90) = true) :
    toLowerBytes input =
      ⟨input.map (fun b => if asciiRange 65 90 b then b.toNat + 32 else b.toNat), some input.length⟩ := by
  simp [toLowerBytes, changeASCII, ascii, changes]

theorem toUpperBytes_empty : toUpperBytes [] = ⟨[], some 0⟩ := rfl
theorem toLowerBytes_empty : toLowerBytes [] = ⟨[], some 0⟩ := rfl
theorem toTitleBytes_empty : toTitleBytes [] = ⟨[], some 0⟩ := rfl

theorem toTitleBytes_map (input : List (BitVec 8)) :
    (toTitleBytes input).bytes =
      ((mapRunes (fun r => Int.ofNat (GoLean.Strings.titleRune r)) input).2).getD [] := rfl

end GoLean.Bytes
