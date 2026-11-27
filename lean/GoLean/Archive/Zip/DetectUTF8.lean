import GoLean.Unicode.Utf8.DecodeRuneInString

namespace GoLean.Archive.Zip

/-- Byte-preserving model of the scan in Go 1.27.0 archive/zip.
    https://github.com/golang/go/blob/go1.27.0/src/archive/zip/writer.go#L247-L268 -/
def detectUTF8 (bytes : List Nat) : Bool × Bool :=
  scan bytes.length bytes false
where
  scan (fuel : Nat) (bytes : List Nat) (require : Bool) : Bool × Bool :=
    match fuel, bytes with
    | remaining + 1, lead :: rest =>
      let (rune, size) := Unicode.Utf8.decodeRuneInString (lead :: rest)
      if rune < 0x20 || rune > 0x7d || rune == 0x5c then
        if rune > 0x10ffff || (0xd800 ≤ rune && rune ≤ 0xdfff) ||
            (rune == 0xfffd && size == 1) then
          (false, false)
        else
          scan remaining ((lead :: rest).drop size) true
      else
        scan remaining ((lead :: rest).drop size) require
    | _, _ => (true, require)

/-- A malformed continuation byte clears the UTF-8 requirement, even after
    the backslash has already set it. -/
theorem detectUTF8_malformed_after_backslash (tail : List Nat) :
    detectUTF8 (0x5c :: 0x80 :: tail) = (false, false) := by
  rw [detectUTF8, List.length_cons, detectUTF8.scan,
    Unicode.Utf8.decodeRuneInString_ascii _ _ (by decide)]
  simp only [show (0x5c < 0x20 || 0x5c > 0x7d || (0x5c : Nat) == 0x5c) = true by decide,
    show ((0x5c : Nat) > 0x10ffff || (0xd800 ≤ (0x5c : Nat) && (0x5c : Nat) ≤ 0xdfff) ||
      ((0x5c : Nat) == 0xfffd && (1 : Nat) == 1)) = false by decide,
    ite_true, List.drop_succ_cons, List.drop_zero, List.length_cons]
  rw [detectUTF8.scan, Unicode.Utf8.decodeRuneInString_error _ _ (by decide)]
  · rfl
  · cases tail with
    | nil => rfl
    | cons a tail =>
      cases tail with
      | nil => rfl
      | cons b tail =>
        cases tail <;> rfl

end GoLean.Archive.Zip
