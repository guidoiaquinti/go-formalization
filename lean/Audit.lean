import Lean
import GoLean

open Lean in
run_cmd do
  for (model, name) in GoLean.registeredTheorems do
    let theoremInfo ← getConstInfo name
    unless theoremInfo.isTheorem do
      throwError "registered property {name} must be a theorem"
    unless theoremInfo.type.getUsedConstants.contains model do
      throwError "registered property {name} must mention its model {model}"
  let allowed := [``propext, ``Classical.choice, ``Quot.sound]
  for (name, _) in (← getEnv).constants.toList do
    if `GoLean |>.isPrefixOf name then
      for axiomName in ← collectAxioms name do
        unless allowed.contains axiomName do
          throwError "{name} depends on forbidden axiom {axiomName}"

