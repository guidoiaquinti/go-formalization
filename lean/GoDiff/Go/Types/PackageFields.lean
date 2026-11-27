import GoLean.Go.Types.PackageFields
import GoDiff.Math.Bits.Tuple

open GoLean.Go.Types

private def wire (xs : List Nat) : String := String.intercalate "," (xs.map toString)
private def render (f : α → String) (value : Except ReadFault α) : String :=
  match value with
  | .error .nilReceiver => "panic"
  | .ok value => f value

private def scopeWire (s : Option Nat) : String := (s.map toString).getD "nil"

private def fields (p : Option PackageState) : String :=
  String.intercalate "|" [
    render wire (packagePath p), render wire (packageName p), scopeWire (packageScope p 0),
    render toString (packageComplete p), render wire (packageGoVersion p)]

def main : IO Unit := runTuple fun row => do
  let mode :: complete :: pathLen :: nameLen :: versionLen :: bytes := row |
    throw (IO.userError "expected package state and field lengths")
  let bs ← bytes.mapM fun b => return (← wordInput 8 b).toNat
  let path := bs.take pathLen.toNat
  let bs := bs.drop pathLen.toNat
  let name := bs.take nameLen.toNat
  let bs := bs.drop nameLen.toNat
  let version := bs.take versionLen.toNat
  let replacement := bs.drop versionLen.toNat
  let created := newPackage path name 31
  let scope := if mode == 1 then created.scope else none
  let pkg : Option PackageState :=
    if mode == 0 then none else some ⟨path, name, scope, complete == 1, version⟩
  let renamed := render (fields ∘ some) (packageSetName pkg replacement)
  let marked := render (fields ∘ some) (packageMarkComplete pkg)
  return String.intercalate "#" [fields (some created), fields pkg, renamed, marked]
