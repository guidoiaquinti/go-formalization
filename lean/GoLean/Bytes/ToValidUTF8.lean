import GoLean.Unicode.Utf8.DecodeRune

namespace GoLean.Bytes

private def scanValidUTF8 (replacement : List Nat) (invalid : Bool) (input : List Nat) :
    Nat → List Nat
  | 0 => []
  | fuel + 1 =>
    match input with
    | [] => []
    | lead :: rest =>
      if lead < 128 then lead :: scanValidUTF8 replacement false rest fuel
      else
        let width := (GoLean.Unicode.Utf8.decodeRune input).2
        if width == 1 then
          (if invalid then [] else replacement) ++ scanValidUTF8 replacement true rest fuel
        else
          input.take width ++ scanValidUTF8 replacement false (input.drop width) fuel

/-- Go 1.27.0 src/bytes/bytes.go:761–786. Return the bytes of the new,
    nonnil slice. A known capacity denotes completion without allocator
    growth. None denotes a capacity determined by runtime growth. -/
def toValidUTF8 (input replacement : List Nat) : List Nat × Option Nat :=
  let output := scanValidUTF8 replacement false input input.length
  let initialCap := input.length + replacement.length
  (output, if output.length ≤ initialCap then some initialCap else none)

theorem toValidUTF8_empty (replacement : List Nat) :
    toValidUTF8 [] replacement = ([], some replacement.length) := by
  simp [toValidUTF8, scanValidUTF8]

theorem toValidUTF8_ascii (lead : Nat) (rest replacement : List Nat) (ascii : lead < 128) :
    (toValidUTF8 (lead :: rest) replacement).1 = lead :: (toValidUTF8 rest replacement).1 := by
  simp [toValidUTF8, scanValidUTF8, ascii]

theorem toValidUTF8_capacity (input replacement : List Nat) :
    (toValidUTF8 input replacement).2 =
      if (toValidUTF8 input replacement).1.length ≤ input.length + replacement.length
      then some (input.length + replacement.length) else none := rfl

end GoLean.Bytes
