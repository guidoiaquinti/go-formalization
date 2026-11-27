import GoLean.Math.RandV2.Bounded

namespace GoLean.Math.RandV2

def globalInt64N (bound : Nat) (words : List Nat) :=
  randInt64N bound words

theorem globalInt64N_delegates (bound : Nat) (words : List Nat) :
    globalInt64N bound words = randInt64N bound words := rfl

def globalUint64N (bound : Nat) (words : List Nat) :=
  randUint64N bound words

theorem globalUint64N_delegates (bound : Nat) (words : List Nat) :
    globalUint64N bound words = randUint64N bound words := rfl

def globalInt32N (bound : Nat) (words : List Nat) :=
  randInt32N bound words

theorem globalInt32N_delegates (bound : Nat) (words : List Nat) :
    globalInt32N bound words = randInt32N bound words := rfl

def globalUint32N (bound : Nat) (words : List Nat) :=
  randUint32N bound words

theorem globalUint32N_delegates (bound : Nat) (words : List Nat) :
    globalUint32N bound words = randUint32N bound words := rfl

def globalIntN (width : Nat) (bound : Nat) (words : List Nat) :=
  randIntN width bound words

theorem globalIntN_delegates (width : Nat) (bound : Nat) (words : List Nat) :
    globalIntN width bound words = randIntN width bound words := rfl

def globalUintN (width : Nat) (bound : Nat) (words : List Nat) :=
  randUintN width bound words

theorem globalUintN_delegates (width : Nat) (bound : Nat) (words : List Nat) :
    globalUintN width bound words = randUintN width bound words := rfl

def globalN (signed : Bool) (width : Nat) (bound : Nat) (words : List Nat) :=
  randN signed width bound words

theorem globalN_delegates (signed : Bool) (width : Nat) (bound : Nat) (words : List Nat) :
    globalN signed width bound words = randN signed width bound words := rfl

end GoLean.Math.RandV2
