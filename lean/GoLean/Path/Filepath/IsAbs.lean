import GoLean.Path.Filepath.VolumeName

namespace GoLean.Path.Filepath

inductive AbsProfile where
  | unix
  | windows
  | plan9

private def windowsSeparator (byte : BitVec 8) : Bool :=
  byte == slash || byte == 92

def isAbs (profile : AbsProfile) (path : List (BitVec 8)) : Bool :=
  match profile with
  | .unix => path.head? == some slash
  | .plan9 => path.head?.any fun byte => byte == slash || byte == 35
  | .windows =>
    let volumeLength := (volumeName true path).length
    volumeLength != 0 &&
      ((path.head?.any windowsSeparator && (path[1]?).any windowsSeparator) ||
        (path[volumeLength]?).any windowsSeparator)

theorem isAbs_unix (path : List (BitVec 8)) :
    isAbs .unix path = (path.head? == some slash) := rfl

theorem isAbs_plan9Device (rest : List (BitVec 8)) :
    isAbs .plan9 (35 :: rest) = true := rfl

theorem isAbs_windowsRootWithoutVolume : isAbs .windows [47, 120] = false := by
  decide

theorem isAbs_windowsDriveExample :
    isAbs .windows [67, 58, 47, 120] = true ∧
    isAbs .windows [67, 58, 120] = false := by decide

end GoLean.Path.Filepath
