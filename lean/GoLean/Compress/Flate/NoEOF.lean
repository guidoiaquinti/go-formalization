namespace GoLean.Compress.Flate

inductive ReadError where
  | none
  | eof
  | unexpectedEOF
  | other (identity : Nat)
  deriving DecidableEq

/-- Only the exact EOF sentinel is replaced. -/
def noEOF : ReadError → ReadError
  | .eof => .unexpectedEOF
  | error => error

theorem noEOF_eof : noEOF .eof = .unexpectedEOF := rfl

theorem noEOF_preserves (error : ReadError) (notEOF : error ≠ .eof) :
    noEOF error = error := by
  cases error <;> simp_all [noEOF]

end GoLean.Compress.Flate
