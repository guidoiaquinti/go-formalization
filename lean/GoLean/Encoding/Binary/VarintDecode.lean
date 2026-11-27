namespace GoLean.Encoding.Binary

private def decodeGroups (consumed : Nat) (value : BitVec 64) :
    List (BitVec 8) → BitVec 64 × Int
  | [] => (0, 0)
  | byte :: tail =>
      let next := consumed + 1
      let group := BitVec.ofNat 64 (byte.toNat % 128) <<< (7 * consumed)
      if consumed == 10 then (0, -(Int.ofNat next))
      else if byte.toNat < 128 then
        if consumed == 9 && byte.toNat > 1 then (0, -(Int.ofNat next))
        else (value ||| group, Int.ofNat next)
      else decodeGroups next (value ||| group) tail

def uvarint (input : List (BitVec 8)) : BitVec 64 × Int :=
  decodeGroups 0 0 input

def varint (input : List (BitVec 8)) : Int × Int :=
  let (value, count) := uvarint input
  let magnitude := Int.ofNat (value.toNat / 2)
  (if value.toNat % 2 == 1 then -magnitude - 1 else magnitude, count)

theorem uvarint_empty : uvarint [] = (0, 0) := rfl

theorem uvarint_one (byte : BitVec 8) (terminal : byte.toNat < 128) :
    uvarint [byte] = (BitVec.ofNat 64 byte.toNat, 1) := by
  simp [uvarint, decodeGroups, terminal, Nat.mod_eq_of_lt terminal]

theorem uvarint_ten_continued :
    uvarint (List.replicate 10 128) = (0, 0) ∧
    uvarint (List.replicate 10 128 ++ [0]) = (0, -11) :=
  by decide

theorem varint_report (input : List (BitVec 8)) : (varint input).2 = (uvarint input).2 := rfl

theorem varint_value (input : List (BitVec 8)) :
    (varint input).1 =
      let value := (uvarint input).1.toNat
      let magnitude := Int.ofNat (value / 2)
      if value % 2 == 1 then -magnitude - 1 else magnitude :=
  rfl

end GoLean.Encoding.Binary
