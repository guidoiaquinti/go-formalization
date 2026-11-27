import GoLean.Path.Filepath.Clean
import GoLean.Path.Filepath.SameWord

namespace GoLean.Path.Filepath

private def samePathWord (windows : Bool) (left right : List (BitVec 8)) : Bool :=
  sameWord windows (left.map BitVec.toNat) (right.map BitVec.toNat)

private def unmatchedPaths (windows : Bool) (separator : BitVec 8) :
    Nat → List (BitVec 8) → List (BitVec 8) → Option (List (BitVec 8) × List (BitVec 8))
  | _, [], [] => none
  | 0, base, target => some (base, target)
  | _ + 1, [], target => some ([], target)
  | fuel + 1, base, target =>
    let basePart := base.takeWhile (· != separator)
    let targetPart := target.takeWhile (· != separator)
    if samePathWord windows basePart targetPart then
      unmatchedPaths windows separator fuel (base.drop (basePart.length + 1))
        (target.drop (targetPart.length + 1))
    else some (base, target)

def relativePath (windows : Bool) (baseInput targetInput : List (BitVec 8)) :
    Option (Except (List (BitVec 8)) (List (BitVec 8))) :=
  let error := "Rel: can't make ".toUTF8.data.toList.map UInt8.toBitVec ++ targetInput ++
    " relative to ".toUTF8.data.toList.map UInt8.toBitVec ++ baseInput
  let separator := if windows then 92 else slash
  let baseVolume := volumeName windows baseInput
  let targetVolume := volumeName windows targetInput
  let baseClean := cleanPath windows baseInput
  let targetClean := cleanPath windows targetInput
  if samePathWord windows targetClean baseClean then some (.ok [46])
  else
    let base := baseClean.drop baseVolume.length
    let target := targetClean.drop targetVolume.length
    let base :=
      if base == [46] then []
      else if base.isEmpty && (volumeName windows baseVolume).length > 2 then [separator]
      else base
    if (base.head? == some separator) != (target.head? == some separator) ||
        !samePathWord windows baseVolume targetVolume then some (.error error)
    else do
      let (baseRest, targetRest) ← unmatchedPaths windows separator base.length base target
      if baseRest.takeWhile (· != separator) == [46, 46] then some (.error error)
      else if baseRest.isEmpty then some (.ok targetRest)
      else
        let levels := 1 + (baseRest.filter (· == separator)).length
        let up := List.intercalate [separator] (List.replicate levels [46, 46])
        some (.ok (cleanPath windows (up ++ (if targetRest.isEmpty then [] else separator :: targetRest))))

private def succeedsWith (wanted : List (BitVec 8)) :
    Option (Except (List (BitVec 8)) (List (BitVec 8))) → Bool
  | some (.ok actual) => actual == wanted
  | _ => false

theorem relativePath_sameUnix : succeedsWith [46] (relativePath false [97] [97]) = true := by decide

theorem relativePath_windowsCase : succeedsWith [46] (relativePath true [67, 58, 92, 97] [99, 58, 92, 65]) = true := by decide

theorem relativePath_upThenDown : succeedsWith [46, 46, 47, 99] (relativePath false [97, 47, 98] [97, 47, 99]) = true := by decide

theorem relativePath_windowsUNCStalls :
    (relativePath true [92, 92, 104, 92, 115] [92, 92, 104, 92, 115, 92]).isNone = true :=
  by decide

end GoLean.Path.Filepath
