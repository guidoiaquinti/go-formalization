import GoLean.Path.Filepath.ScanChunk
import GoLean.Path.Filepath.MatchChunk

namespace GoLean.Path.Filepath

private inductive Attempt where
  | mismatch
  | badPattern
  | accepted (rest : List Nat)

private def attempt (windows final : Bool) (chunk name : List Nat) : Attempt :=
  let (rest, matched, bad) := matchChunk windows chunk name
  if bad then .badPattern
  else if matched && (!final || rest.isEmpty) then .accepted rest
  else .mismatch

private def seekAfterStar (windows final : Bool) (separator : Nat) (chunk : List Nat) :
    List Nat → Attempt
  | [] => .mismatch
  | byte :: tail =>
    if byte = separator then .mismatch
    else
      match attempt windows final chunk tail with
      | .mismatch => seekAfterStar windows final separator chunk tail
      | result => result

private def matchLoop (windows : Bool) (separator : Nat) :
    Nat → List (BitVec 8) → List Nat → Bool × Bool
  | 0, _, _ => (false, false)
  | _ + 1, [], name => (name.isEmpty, false)
  | fuel + 1, first :: tail, name =>
    let (star, rawChunk, pattern) := scanChunk windows (first :: tail)
    if star && rawChunk.isEmpty then (!(name.contains separator), false)
    else
      let chunk := rawChunk.map BitVec.toNat
      let final := pattern.isEmpty
      let current := attempt windows final chunk name
      let result :=
        match current with
        | .mismatch => if star then seekAfterStar windows final separator chunk name else .mismatch
        | other => other
      match result with
      | .accepted rest => matchLoop windows separator fuel pattern rest
      | .badPattern => (false, true)
      | .mismatch => (false, false)

def matchPath (windows : Bool) (pattern name : List (BitVec 8)) : Bool × Bool :=
  matchLoop windows (if windows then 92 else 47) (pattern.length + 1) pattern (name.map BitVec.toNat)

theorem matchPath_empty (windows : Bool) (name : List (BitVec 8)) :
    matchPath windows [] name = (name.isEmpty, false) := by
  simp [matchPath, matchLoop]

theorem matchPath_trailingStar :
    matchPath false [42] [97, 47, 98] = (false, false) ∧
    matchPath true [42] [97, 47, 98] = (true, false) := by decide

theorem matchPath_terminalConsumption :
    matchPath false [42, 97, 98] [97, 98, 97, 98] = (true, false) ∧
    matchPath false [97, 98] [97, 98, 97, 98] = (false, false) := by decide

theorem matchPath_errorOrder :
    matchPath false [97, 42, 91] [98] = (false, false) ∧
    matchPath false [97, 91] [98] = (false, true) := by decide

end GoLean.Path.Filepath
