import Init.Data.BitVec.Basic
import Init.Data.List.Lemmas

namespace GoLean.Path.Filepath

private def splitChunk (windows inRange : Bool) (reversed : List (BitVec 8)) :
    List (BitVec 8) → List (BitVec 8) × List (BitVec 8)
  | [] => (reversed.reverse, [])
  | byte :: rest =>
    if byte = 42 && !inRange then (reversed.reverse, byte :: rest)
    else if byte = 92 && !windows then
      match rest with
      | [] => ((byte :: reversed).reverse, [])
      | escaped :: tail => splitChunk windows inRange (escaped :: byte :: reversed) tail
    else
      let next := if byte = 91 then true else if byte = 93 then false else inRange
      splitChunk windows next (byte :: reversed) rest

def scanChunk (windows : Bool) (pattern : List (BitVec 8)) :
    Bool × List (BitVec 8) × List (BitVec 8) :=
  let body := pattern.dropWhile (· == 42)
  (pattern.head? == some 42, splitChunk windows false [] body)

theorem scanChunk_empty (windows : Bool) : scanChunk windows [] = (false, [], []) := rfl

theorem scanChunk_stars (windows : Bool) :
    scanChunk windows [42, 42] = (true, [], []) := by
  cases windows <;> decide

theorem scanChunk_escape :
    scanChunk false [92, 42, 120] = (false, [92, 42, 120], []) ∧
    scanChunk true [92, 42, 120] = (false, [92], [42, 120]) := by decide

theorem scanChunk_class (windows : Bool) :
    scanChunk windows [91, 42, 93, 42] = (false, [91, 42, 93], [42]) := by
  cases windows <;> decide

end GoLean.Path.Filepath
