import GoLean.Strings.Substring

namespace GoLean.Strings

/-- U+FFFD is emitted by the range loop for either a malformed byte or
    the valid encoding EF BF BD. Width agrees with Go's UTF-8 acceptance rules. -/
def runeErrorHere (lead : Nat) (rest : List Nat) : Bool :=
  (lead ≥ 128 && runeWidth lead rest == 1) || ([lead] ++ rest.take 2 == [239, 191, 189])

def firstRuneError : List Nat → Option Nat
  | [] => none
  | lead :: rest =>
    if runeErrorHere lead rest then some 0
    else (firstRuneError (rest.drop (runeWidth lead rest - 1))).map
      (fun i => runeWidth lead rest + i)
termination_by xs => xs.length
decreasing_by simp [List.length_drop]; omega

def validRune (r : Int) : Bool :=
  0 ≤ r && r ≤ 0x10ffff && (r < 0xd800 || 0xdfff < r)

def encodeRune (n : Nat) : List Nat :=
  if n < 128 then [n]
  else if n < 2048 then [192 + n / 64, 128 + n % 64]
  else if n < 65536 then [224 + n / 4096, 128 + n / 64 % 64, 128 + n % 64]
  else [240 + n / 262144, 128 + n / 4096 % 64, 128 + n / 64 % 64, 128 + n % 64]

/-- Go searches for ASCII directly, range-scans for RuneError, and searches the
    encoded bytes of other valid runes. The candidate optimizations are abstracted.
    https://github.com/golang/go/blob/go1.27.0/src/strings/strings.go#L106-L179
    https://github.com/golang/go/blob/go1.27.0/src/unicode/utf8/utf8.go#L368-L408 -/
def indexRune (input : List Nat) (r : Int) : Option Nat :=
  if r = 0xfffd then firstRuneError input
  else if validRune r then index input (encodeRune r.toNat)
  else none

/-- https://github.com/golang/go/blob/go1.27.0/src/strings/strings.go#L71-L74 -/
def containsRune (input : List Nat) (r : Int) : Bool :=
  (indexRune input r).isSome

theorem indexRune_valid (input : List Nat) (r : Int) (h : validRune r = true)
    (ordinary : r ≠ 0xfffd) :
    FirstMatch input (encodeRune r.toNat) (indexRune input r) := by
  simpa [indexRune, ordinary, h] using index_spec input (encodeRune r.toNat)

theorem indexRune_invalid (input : List Nat) (r : Int) (h : validRune r = false)
    (ordinary : r ≠ 0xfffd) : indexRune input r = none := by
  simp [indexRune, ordinary, h]

theorem indexRune_error_cons (lead : Nat) (rest : List Nat) :
    indexRune (lead :: rest) 0xfffd =
      if runeErrorHere lead rest then some 0
      else (indexRune (rest.drop (runeWidth lead rest - 1)) 0xfffd).map
        (fun i => runeWidth lead rest + i) := by
  simp only [indexRune, ite_true]
  rw [firstRuneError.eq_2]

theorem containsRune_spec (input : List Nat) (r : Int) :
    containsRune input r = true ↔ ∃ i, indexRune input r = some i := by
  cases h : indexRune input r <;> simp [containsRune, h]

end GoLean.Strings
