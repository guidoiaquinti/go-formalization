import GoLean.Unicode.Utf8.DecodeRuneInString

namespace GoLean.Path.Filepath

private def decodeEsc : List Nat → Nat × List Nat × Bool
  | [] => (0, [], true)
  | byte :: tail =>
    let (rune, width) := GoLean.Unicode.Utf8.decodeNonempty byte tail
    let rest := (byte :: tail).drop width
    (rune, rest, (rune == 0xfffd && width == 1) || rest.isEmpty)

def getEsc (windows : Bool) : List Nat → Nat × List Nat × Bool
  | [] => (0, [], true)
  | first :: tail =>
    if first = 45 || first = 93 then (0, [], true)
    else if !windows && first = 92 then decodeEsc tail
    else decodeEsc (first :: tail)

theorem getEsc_empty (windows : Bool) : getEsc windows [] = (0, [], true) := rfl

theorem getEsc_hyphen (windows : Bool) (rest : List Nat) :
    getEsc windows (45 :: rest) = (0, [], true) := rfl

theorem getEsc_escapes :
    getEsc false [92, 93, 120] = (93, [120], false) ∧
    getEsc true [92, 93, 120] = (92, [93, 120], false) := by decide

theorem getEsc_errorRemainder :
    getEsc false [255, 120] = (0xfffd, [120], true) ∧
    getEsc false [0xef, 0xbf, 0xbd, 120] = (0xfffd, [120], false) := by decide

end GoLean.Path.Filepath
