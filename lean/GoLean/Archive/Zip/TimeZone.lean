import GoLean.Time.DurationNumeric

namespace GoLean.Archive.Zip

structure ZipZone where
  seconds : Int
  cached : Bool
  deriving DecidableEq

/-- Zip rounds to quarter hours before checking the permitted offset band. -/
def timeZone (nanoseconds : Int) : ZipZone :=
  let rounded := GoLean.Time.durationRound nanoseconds 900000000000
  let seconds := if rounded < -43200000000000 || 50400000000000 < rounded then 0
    else rounded.tdiv 1000000000
  ⟨seconds, seconds % 3600 == 0⟩

theorem timeZone_roundBeforeLimit :
    timeZone 50849999999999 = ⟨50400, true⟩ ∧
      timeZone 50850000000000 = ⟨0, true⟩ := by decide

theorem timeZone_negativeTie :
    timeZone (-450000000000) = ⟨-900, false⟩ := by decide

end GoLean.Archive.Zip
