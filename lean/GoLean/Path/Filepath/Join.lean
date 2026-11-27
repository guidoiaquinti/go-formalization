import GoLean.Path.Filepath.Clean

namespace GoLean.Path.Filepath

private def isWindowsSeparator (byte : BitVec 8) : Bool := byte == slash || byte == 92

private def joinWindowsStep (acc element : List (BitVec 8)) : List (BitVec 8) :=
  match acc.getLast? with
  | none => element
  | some last =>
    if isWindowsSeparator last then
      let rest := element.dropWhile isWindowsSeparator
      let guard := acc.length = 1 && rest.take 2 == [63, 63] &&
        (rest.length == 2 || (rest[2]?).any isWindowsSeparator)
      acc ++ (if guard then [46, 92] else []) ++ rest
    else if last = 58 then acc ++ element
    else acc ++ [92] ++ element

def joinPath (windows : Bool) (elements : List (List (BitVec 8))) : List (BitVec 8) :=
  let raw :=
    if windows then elements.foldl joinWindowsStep []
    else List.intercalate [slash] (elements.dropWhile List.isEmpty)
  if raw.isEmpty then [] else cleanPath windows raw

theorem joinPath_empty (windows : Bool) : joinPath windows [] = [] := by
  cases windows <;> rfl

theorem joinPath_emptyElements (windows : Bool) : joinPath windows [[], []] = [] := by
  cases windows <;> decide

theorem joinPath_windowsDrive :
    joinPath true [[67, 58], [102]] = [67, 58, 102] ∧
    joinPath true [[67, 58], [92, 102]] = [67, 58, 92, 102] := by decide

theorem joinPath_windowsDeviceGuard :
    joinPath true [[92], [63, 63], [97]] = [92, 46, 92, 63, 63, 92, 97] := by decide

end GoLean.Path.Filepath
