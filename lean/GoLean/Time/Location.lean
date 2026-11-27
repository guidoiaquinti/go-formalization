import GoLean.Time.TimeCompare

namespace GoLean.Time

structure LocatedTime where
  words : ClockWords
  loc : Nat

/-- Nil and UTC have distinct pointer tags. Other tags retain their identity. -/
def clockStripMono (s : LocatedTime) : LocatedTime :=
  if s.words.wall.getLsbD 63 then
    {s with words := ⟨s.words.wall &&& (2^30 - 1), BitVec.ofInt 64 (clockSec s.words)⟩}
  else s

def clockSetLoc (s : LocatedTime) (loc : Nat) : LocatedTime :=
  {clockStripMono s with loc := if loc = 1 then 0 else loc}

def clockUTC (s : LocatedTime) : LocatedTime := clockSetLoc s 1

def clockLocal (s : LocatedTime) : LocatedTime := clockSetLoc s 2

def clockIn (s : LocatedTime) (loc : Nat) : Option LocatedTime :=
  if loc = 0 then none else some (clockSetLoc s loc)

def clockLocation (s : LocatedTime) : Nat := if s.loc = 0 then 1 else s.loc

theorem clockStripMono_unflagged (s : LocatedTime) (clear : s.words.wall.getLsbD 63 = false) :
    clockStripMono s = s := by
  unfold clockStripMono
  rw [clear]
  rfl

theorem clockSetLoc_utc (s : LocatedTime) : (clockSetLoc s 1).loc = 0 := rfl

theorem clockUTC_canonical (s : LocatedTime) : (clockUTC s).loc = 0 := rfl

theorem clockLocal_identity (s : LocatedTime) : (clockLocal s).loc = 2 := rfl

theorem clockIn_nil (s : LocatedTime) : clockIn s 0 = none := rfl

theorem clockLocation_nil (s : LocatedTime) (missing : s.loc = 0) : clockLocation s = 1 := by
  simp [clockLocation, missing]

end GoLean.Time
