namespace GoLean.Time

/-- The entire backing from the slice start through capacity, including cells
    that are not currently visible. A bounds failure retains prior writes. -/
structure FormatResult where
  backing : List UInt8
  result : Option (Nat × Nat)

private def digit (n : Nat) : UInt8 := UInt8.ofNat (48 + n % 10)

private def intDigits : Nat → Nat → Nat → List UInt8 → FormatResult
  | 0, _, pos, buf => ⟨buf, some (pos, 0)⟩
  | _ + 1, _, 0, buf => ⟨buf, none⟩
  | fuel + 1, n, pos + 1, buf =>
    let buf := buf.set pos (digit n)
    let rest := n / 10
    if rest == 0 then ⟨buf, some (pos, 0)⟩
    else intDigits fuel rest pos buf

/-- Go 1.27.0 time.go:1056–1069. The fuel of 20 covers every uint64. -/
def fmtInt (value length : Nat) (backing : List UInt8) : FormatResult :=
  intDigits 20 value length backing

private def fracDigits : Nat → Nat → Nat → Bool → List UInt8 → FormatResult
  | 0, n, pos, false, buf => ⟨buf, some (pos, n)⟩
  | 0, _, 0, true, buf => ⟨buf, none⟩
  | 0, n, pos + 1, true, buf => ⟨buf.set pos 46, some (pos, n)⟩
  | fuel + 1, n, pos, printed, buf =>
    if printed || n % 10 != 0 then
      match pos with
      | 0 => ⟨buf, none⟩
      | p + 1 => fracDigits fuel (n / 10) p true (buf.set p (digit n))
    else fracDigits fuel (n / 10) pos false buf

/-- Go 1.27.0 time.go:1034–1052. Nonpositive precision takes no iterations. -/
def fmtFrac (value length : Nat) (precision : Int) (backing : List UInt8) : FormatResult :=
  fracDigits precision.toNat value length false backing

theorem fmtInt_empty (n : Nat) (backing : List UInt8) :
    (fmtInt n 0 backing).result = none := by rfl

theorem fmtInt_partial :
    (fmtInt 567 2 [7, 8, 9]).backing = [54, 55, 9] ∧
    (fmtInt 567 2 [7, 8, 9]).result = none := by decide

theorem fmtFrac_zeroPrecision (n length : Nat) (backing : List UInt8) :
    (fmtFrac n length 0 backing).backing = backing ∧
    (fmtFrac n length 0 backing).result = some (length, n) := by
  exact ⟨rfl, rfl⟩

theorem fmtFrac_partial :
    (fmtFrac 120 2 3 [7, 8, 9]).backing = [49, 50, 9] ∧
    (fmtFrac 120 2 3 [7, 8, 9]).result = none := by decide

end GoLean.Time
