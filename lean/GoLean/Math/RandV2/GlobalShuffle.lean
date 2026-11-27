import GoLean.Math.RandV2.Shuffle

namespace GoLean.Math.RandV2

def globalShuffle (width : Nat) (n : Int) (words : List Nat) (script : SwapScript) :=
  randShuffle width n words script

def globalPerm (width : Nat) (n : Int) (words : List Nat) :=
  randPerm width n words

theorem globalShuffle_delegates (width : Nat) (n : Int) (words : List Nat) (script : SwapScript) :
    globalShuffle width n words script = randShuffle width n words script := rfl

theorem globalPerm_delegates (width : Nat) (n : Int) (words : List Nat) :
    globalPerm width n words = randPerm width n words := rfl

end GoLean.Math.RandV2
