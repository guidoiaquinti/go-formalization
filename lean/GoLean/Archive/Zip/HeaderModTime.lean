import GoLean.Archive.Zip.DOSTime
import GoLean.Time.TimeCompare

namespace GoLean.Archive.Zip

structure ModificationTime where
  seconds : Int
  nanoseconds : Nat
  deriving DecidableEq

/-- The selected instant is returned in UTC, with no monotonic reading. -/
def headerFileInfoModTime (modified : GoLean.Time.ClockWords) (date time : BitVec 16) :
    ModificationTime :=
  if GoLean.Time.clockIsZero modified then ⟨fileHeaderModTime date time, 0⟩
  else ⟨GoLean.Time.clockUnix modified, GoLean.Time.clockNanosecond modified⟩

theorem headerFileInfoModTime_legacy (date time : BitVec 16) :
    headerFileInfoModTime ⟨0, 0⟩ date time = ⟨fileHeaderModTime date time, 0⟩ := rfl

theorem headerFileInfoModTime_modern (modified : GoLean.Time.ClockWords) (date time : BitVec 16)
    (nonzero : GoLean.Time.clockIsZero modified = false) :
    headerFileInfoModTime modified date time =
      ⟨GoLean.Time.clockUnix modified, GoLean.Time.clockNanosecond modified⟩ := by
  simp [headerFileInfoModTime, nonzero]

end GoLean.Archive.Zip
