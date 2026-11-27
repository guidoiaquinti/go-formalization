import GoLean.Strings.IndexFunc

namespace GoLean.Bytes

/-- Go 1.27.0 src/bytes/bytes.go:885–906. With fixed input bytes,
    DecodeRune has the same width and scalar as the string decoder.
    None denotes the nil-callback panic. -/
def indexFunc (input : List (BitVec 8)) (predicate : Option (Nat → Bool)) :
    Option (Int × List Nat) :=
  GoLean.Strings.indexFunc input predicate

theorem indexFunc_empty (predicate : Option (Nat → Bool)) :
    indexFunc [] predicate = some (-1, []) := rfl

theorem indexFunc_nil (lead : BitVec 8) (rest : List (BitVec 8)) :
    indexFunc (lead :: rest) none = none := rfl

theorem indexFunc_first (lead : BitVec 8) (rest : List (BitVec 8)) (f : Nat → Bool)
    (matched : f (GoLean.Unicode.Utf8.decodeNonempty lead.toNat (rest.map BitVec.toNat)).1 = true) :
    indexFunc (lead :: rest) (some f) =
      some (0, [(GoLean.Unicode.Utf8.decodeNonempty lead.toNat (rest.map BitVec.toNat)).1]) :=
  GoLean.Strings.indexFunc_first lead rest f matched

theorem indexFunc_noMatch (input : List (BitVec 8)) :
    (indexFunc input (some (fun _ => false))).map Prod.fst = some (-1) :=
  GoLean.Strings.indexFunc_noMatch input

end GoLean.Bytes
