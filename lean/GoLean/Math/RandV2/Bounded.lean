namespace GoLean.Math.RandV2

structure BoundedResult where
  status : String
  value : Nat
  pulls : Nat

private def scale (n : Nat) (pulls : Nat) : List Nat → BoundedResult
  | [] => ⟨"callback", 0, pulls + 1⟩
  | word :: rest =>
    let product := word * n
    if product % 2^64 < 2^64 % n then scale n (pulls + 1) rest
    else ⟨"ok", product / 2^64, pulls + 1⟩

private def reduce (width bound : Nat) (words : List Nat) : BoundedResult :=
  if bound &&& (bound - 1) = 0 then
    match words with
    | [] => ⟨"callback", 0, 1⟩
    | word :: _ =>
      let mask := if bound = 0 then 2^width - 1 else bound - 1
      ⟨"ok", word &&& mask, 1⟩
  else scale bound 0 words

def randUint64n (machineWidth bound : Nat) (words : List Nat) :=
  reduce (if machineWidth = 32 ∧ bound < 2^32 then 32 else 64) bound words
def randUint32n (bound : Nat) (words : List Nat) := reduce 32 bound words

private def checked (label : String) (signed : Bool) (width bound : Nat)
    (words : List Nat) : BoundedResult :=
  let bound := bound % 2^width
  if bound = 0 ∨ (signed ∧ 2^(width-1) ≤ bound) then
    ⟨"invalid argument to " ++ label, 0, 0⟩
  else randUint64n 64 bound words

def randInt64N := checked "Int64N" true 64
def randUint64N := checked "Uint64N" false 64
def randInt32N := checked "Int32N" true 32
def randUint32N := checked "Uint32N" false 32
def randIntN (width : Nat) := checked "IntN" true width
def randUintN (width : Nat) := checked "UintN" false width
def randN (signed : Bool) (width : Nat) := checked "N" signed width

theorem randUint64n_zero (word : Nat) (tail : List Nat) :
    randUint64n 64 0 (word :: tail) = ⟨"ok", word &&& (2^64-1), 1⟩ := rfl

theorem randUint32n_zero (word : Nat) (tail : List Nat) :
    randUint32n 0 (word :: tail) = ⟨"ok", word &&& (2^32-1), 1⟩ := rfl

theorem randInt64N_zero  (words : List Nat) :
    randInt64N 0 words = ⟨"invalid argument to Int64N", 0, 0⟩ := by
  simp [randInt64N, checked]

theorem randUint64N_zero  (words : List Nat) :
    randUint64N 0 words = ⟨"invalid argument to Uint64N", 0, 0⟩ := by
  simp [randUint64N, checked]

theorem randInt32N_zero  (words : List Nat) :
    randInt32N 0 words = ⟨"invalid argument to Int32N", 0, 0⟩ := by
  simp [randInt32N, checked]

theorem randUint32N_zero  (words : List Nat) :
    randUint32N 0 words = ⟨"invalid argument to Uint32N", 0, 0⟩ := by
  simp [randUint32N, checked]

theorem randIntN_zero (width : Nat) (words : List Nat) :
    randIntN width 0 words = ⟨"invalid argument to IntN", 0, 0⟩ := by
  simp [randIntN, checked]

theorem randUintN_zero (width : Nat) (words : List Nat) :
    randUintN width 0 words = ⟨"invalid argument to UintN", 0, 0⟩ := by
  simp [randUintN, checked]

theorem randN_zero (signed : Bool) (width : Nat) (words : List Nat) :
    randN signed width 0 words = ⟨"invalid argument to N", 0, 0⟩ := by
  simp [randN, checked]

end GoLean.Math.RandV2
