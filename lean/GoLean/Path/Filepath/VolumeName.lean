import GoLean.Path.Filepath.FromSlash

namespace GoLean.Path.Filepath

private def isSeparator (byte : BitVec 8) : Bool :=
  byte == slash || byte == 92

private def asciiUpper (byte : BitVec 8) : BitVec 8 :=
  if 97 ≤ byte.toNat && byte.toNat ≤ 122 then byte - 32 else byte

private def prefixFold : List (BitVec 8) → List (BitVec 8) → Bool
  | [], [] => true
  | [], byte :: _ => isSeparator byte
  | _ :: _, [] => false
  | first :: rest, byte :: path =>
    (if isSeparator first then isSeparator byte else asciiUpper first == asciiUpper byte) &&
      prefixFold rest path

private def uncEnd (index separators : Nat) : List (BitVec 8) → Nat
  | [] => index
  | byte :: rest =>
    if isSeparator byte then
      if separators = 1 then index
      else uncEnd (index + 1) (separators + 1) rest
    else uncEnd (index + 1) separators rest

private def componentEnd (index : Nat) : List (BitVec 8) → Nat
  | [] => index
  | byte :: rest =>
    if isSeparator byte then index else componentEnd (index + 1) rest

private def noParent (part : List (BitVec 8)) : List (BitVec 8) → Bool
  | [] => part != [46, 46]
  | byte :: rest =>
    if isSeparator byte then part != [46, 46] && noParent [] rest
    else noParent (byte :: part) rest

private def validate (path : List (BitVec 8)) (length : Nat) : Nat :=
  if noParent [] (path.take length) then length else 0

private def windowsVolumeLength (path : List (BitVec 8)) : Nat :=
  if path[1]? == some 58 then 2
  else if !(path.head?.any isSeparator) then 0
  else if prefixFold [92, 92, 46] path || prefixFold [92, 92, 63] path ||
      prefixFold [92, 63, 63] path then
    if path.length = 3 then 3
    else if prefixFold [85, 78, 67] (path.drop 4) then
      validate path (uncEnd (min 8 path.length) 0 (path.drop 8))
    else validate path (componentEnd 4 (path.drop 4))
  else if (path[1]?).any isSeparator then
    validate path (uncEnd 2 0 (path.drop 2))
  else 0

def volumeName (windows : Bool) (path : List (BitVec 8)) : List (BitVec 8) :=
  if windows then fromSlash 92 (path.take (windowsVolumeLength path)) else []

theorem volumeName_unix (path : List (BitVec 8)) : volumeName false path = [] := rfl

theorem volumeName_drive (first : BitVec 8) (rest : List (BitVec 8)) :
    volumeName true (first :: 58 :: rest) = fromSlash 92 [first, 58] := by
  rfl

theorem volumeName_unc : volumeName true [47, 47, 97, 47, 98, 47, 99] = [92, 92, 97, 92, 98] := by
  decide

theorem volumeName_uncParent : volumeName true [47, 47, 97, 47, 46, 46] = [] := by
  decide

end GoLean.Path.Filepath
