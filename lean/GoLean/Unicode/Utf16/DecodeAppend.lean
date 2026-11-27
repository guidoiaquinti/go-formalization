import GoLean.Unicode.Utf16.Decode

namespace GoLean.Unicode.Utf16

/-- The Go loop appends each decoded rune to the supplied buffer. A reversed
    list makes each individual append constant time in the model.
    https://github.com/golang/go/blob/go1.27.0/src/unicode/utf16/utf16.go#L123-L144 -/
def decodeAppend (words : List Nat) (dst : List Int) : List Int :=
  ((decode words).foldl (fun reversed (r : Nat) => (r : Int) :: reversed) dst.reverse).reverse

private theorem foldDecoded (runes : List Nat) (reversed : List Int) :
    runes.foldl (fun acc (r : Nat) => (r : Int) :: acc) reversed =
      (runes.map Int.ofNat).reverse ++ reversed := by
  induction runes generalizing reversed with
  | nil => rfl
  | cons r tail ih =>
    simp [List.foldl_cons, ih, List.reverse_cons, List.append_assoc]

theorem decodeAppend_spec (words : List Nat) (dst : List Int) :
    decodeAppend words dst = dst ++ (decode words).map Int.ofNat := by
  rw [decodeAppend, foldDecoded, List.reverse_append, List.reverse_reverse, List.reverse_reverse]

end GoLean.Unicode.Utf16
