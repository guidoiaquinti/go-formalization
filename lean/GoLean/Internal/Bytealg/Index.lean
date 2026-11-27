import GoLean.Strings.Substring

namespace GoLean.Internal.Bytealg

/-- Go 1.27.0 src/internal/bytealg/index_native.go:13. This models
    the result on the native small-needle domain. None represents -1. -/
def index (haystack needle : List Nat) : Option Nat :=
  GoLean.Strings.index haystack needle

/-- Go 1.27.0 src/internal/bytealg/index_native.go:19. Search raw bytes. -/
def indexString (haystack needle : List Nat) : Option Nat :=
  index haystack needle

theorem index_spec (haystack needle : List Nat) :
    GoLean.Strings.FirstMatch haystack needle (index haystack needle) :=
  GoLean.Strings.index_spec haystack needle

theorem indexString_spec (haystack needle : List Nat) :
    GoLean.Strings.FirstMatch haystack needle (indexString haystack needle) :=
  index_spec haystack needle

end GoLean.Internal.Bytealg
