import GoLean.Math.RandV2.Bounded

namespace GoLean.Math.RandV2

inductive SwapScript where
  | normal
  | failAt (index : Nat)
  | missing

structure ShuffleResult where
  status : String
  pulls : Nat
  calls : List (Nat × Nat)

private def swapSteps (width : Nat) (script : SwapScript) :
    Nat → List Nat → Nat → List (Nat × Nat) → ShuffleResult
  | 0, _, pulls, rev => ⟨"ok", pulls, rev.reverse⟩
  | i + 1, words, pulls, rev =>
    let index := i + 1
    let choice := randUint64n width (index + 1) words
    let pulls := pulls + choice.pulls
    if choice.status != "ok" then ⟨"source", pulls, rev.reverse⟩
    else
      let args := (index, choice.value)
      match script with
      | .missing => ⟨"nil", pulls, rev.reverse⟩
      | .failAt n =>
        if rev.length = n then ⟨"swap", pulls, (args :: rev).reverse⟩
        else swapSteps width script i (words.drop choice.pulls) pulls (args :: rev)
      | .normal => swapSteps width script i (words.drop choice.pulls) pulls (args :: rev)

def randShuffle (width : Nat) (n : Int) (words : List Nat) (script : SwapScript) : ShuffleResult :=
  if n < 0 then ⟨"negative", 0, []⟩
  else swapSteps width script (n.toNat - 1) words 0 []

def randPerm (width : Nat) (n : Int) (words : List Nat) : String × Nat × Option (List Nat) :=
  let result := randShuffle width n words .normal
  if result.status != "ok" then (result.status, result.pulls, none)
  else
    let values := result.calls.foldl (fun xs (i, j) =>
        (xs.set! i xs[j]!).set! j xs[i]!) (List.range n.toNat).toArray
    ("ok", result.pulls, some values.toList)

theorem randShuffle_empty (width : Nat) (words : List Nat) (script : SwapScript) :
    randShuffle width 0 words script = ⟨"ok", 0, []⟩ := rfl

theorem randPerm_singleton (width : Nat) (words : List Nat) :
    randPerm width 1 words = ("ok", 0, some [0]) := rfl

end GoLean.Math.RandV2
