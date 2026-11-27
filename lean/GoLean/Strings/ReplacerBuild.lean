import GoLean.Strings.MakeGenericReplacer
import GoLean.Strings.MakeSingleStringReplacer
import GoLean.Strings.ByteStringReplacerReplace

namespace GoLean.Strings

abbrev ReplacerPairs := List (List (BitVec 8) × List (BitVec 8))

inductive BuiltReplacer where
  | single (state : SingleStringReplacer)
  | generic (pairs : ReplacerPairs) (state : GenericReplacer)
  | bytes (lookup : BitVec 8 → BitVec 8)
  | byteStrings (state : ByteStringReplacement)

private def firstByteReplacement (pairs : ReplacerPairs) (byte : BitVec 8) :
    Option (List (BitVec 8)) :=
  (pairs.find? fun pair => pair.1 == [byte]).map Prod.snd

/-- Go 1.27.0 src/strings/replace.go:44–91. Select the implementation
    from raw key/value byte lengths. Reverse insertion into the byte
    tables gives the first pair priority, and orders toReplace by the
    last appearances of its distinct keys. -/
def replacerBuild (pairs : ReplacerPairs) : BuiltReplacer :=
  match pairs with
  | [(key, value)] =>
    if key.length > 1 then .single (makeSingleStringReplacer key value)
    else buildByteTables pairs
  | _ => buildByteTables pairs
where
  buildByteTables (pairs : ReplacerPairs) : BuiltReplacer :=
    if pairs.any (fun pair => pair.1.length != 1) then
      .generic pairs (makeGenericReplacer pairs)
    else if pairs.all (fun pair => pair.2.length == 1) then
      .bytes fun byte => ((firstByteReplacement pairs byte).getD []).headD byte
    else
      let keys := (pairs.reverse.filterMap fun pair => pair.1.head?).eraseDups
      .byteStrings ⟨keys, firstByteReplacement pairs⟩

theorem replacerBuild_empty : replacerBuild [] = .bytes id := rfl

theorem replacerBuild_single (key value : List (BitVec 8)) (long : 1 < key.length) :
    replacerBuild [(key, value)] = .single (makeSingleStringReplacer key value) := by
  simp [replacerBuild, long]

theorem replacerBuild_firstByte :
    (match replacerBuild [([0xff], [1]), ([0xff], [2])] with
      | .bytes lookup => some (lookup 0xff)
      | _ => none) = some 1 := rfl

theorem replacerBuild_emptyKey (value : List (BitVec 8)) :
    replacerBuild [([], value)] = .generic [([], value)] (makeGenericReplacer [([], value)]) :=
  rfl

end GoLean.Strings
