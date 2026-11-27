import GoLean.Internal.Dag.AddEdge
import GoDiff.Math.Bits.Tuple

private def keys := ["", String.singleton '\x00']

private def row (code : Nat) : Option (List (String × Bool)) :=
  if code < 2 then none
  else some <|
    (keys.zip [(code - 2) % 3, (code - 2) / 3]).filterMap fun (key, val) =>
      if val == 0 then none else some (key, val == 2)

def main : IO Unit := runTuple fun fields => do
  let [receiver, source, to, first, second] := fields
    | throw (IO.userError "expected receiver, two keys and two row codes")
  let outer := (keys.zip [first.toNat, second.toNat]).filterMap fun (key, code) =>
    if code == 0 then none else some (key, row code)
  let graph := if receiver == 0 then none else if receiver == 1 then some none else some (some outer)
  let some fromKey := keys[source.toNat]? | throw (IO.userError "invalid source key")
  let some toKey := keys[to.toNat]? | throw (IO.userError "invalid destination key")
  return match GoLean.Internal.Dag.addEdge graph fromKey toKey with
    | none => "panic"
    | some rows =>
      let code (key : String) : Nat :=
        match rows.lookup key with
        | none => 0
        | some none => 1
        | some (some entries) =>
          let digit (key : String) :=
            match entries.lookup key with
            | none => 0
            | some false => 1
            | some true => 2
          2 + digit "" + 3 * digit (String.singleton '\x00')
      s!"{code ""},{code (String.singleton '\x00')}"
