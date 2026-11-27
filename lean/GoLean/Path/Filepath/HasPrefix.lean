import GoLean.Strings.ToLower

namespace GoLean.Path.Filepath

private def beginsWith [BEq α] (path leading : List α) : Bool :=
  path.take leading.length == leading

def hasPathPrefix (windows : Bool) (path leading : List (BitVec 8)) : Bool :=
  beginsWith path leading ||
    (windows && beginsWith (GoLean.Strings.toLower path) (GoLean.Strings.toLower leading))

theorem hasPathPrefix_empty (windows : Bool) (path : List (BitVec 8)) :
    hasPathPrefix windows path [] = true := rfl

theorem hasPathPrefix_windowsRawLead :
    hasPathPrefix true [195, 137] [195] = true := by decide

theorem hasPathPrefix_windowsLowercaseFallback (path leading : List (BitVec 8))
    (raw : (path.take leading.length == leading) = false) :
    hasPathPrefix true path leading =
      ((GoLean.Strings.toLower path).take (GoLean.Strings.toLower leading).length ==
        GoLean.Strings.toLower leading) := by
  simp [hasPathPrefix, beginsWith, raw]

theorem hasPathPrefix_ignoresBoundary :
    hasPathPrefix false [97, 98] [97] = true := by decide

end GoLean.Path.Filepath
