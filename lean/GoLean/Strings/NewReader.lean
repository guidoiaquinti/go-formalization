import GoLean.Strings.ReaderReset

namespace GoLean.Strings

/-- Construct the initial state of the returned Reader.
    https://github.com/golang/go/blob/go1.27.0/src/strings/reader.go#L154-L156 -/
def newReader (contents : List (BitVec 8)) : ReaderState :=
  ⟨contents, 0, -1⟩

theorem newReader_spec (contents : List (BitVec 8)) :
    (newReader contents).contents = contents ∧
      (newReader contents).cursor = 0 ∧ (newReader contents).previousRune = -1 :=
  ⟨rfl, rfl, rfl⟩

end GoLean.Strings
