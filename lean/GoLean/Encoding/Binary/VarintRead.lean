namespace GoLean.Encoding.Binary

inductive ReaderError where
  | eof
  | unexpectedEOF
  | overflow
  | other
  deriving DecidableEq, Repr

structure VarintReadResult (α : Type) where
  value : α
  error : Option ReaderError
  calls : Nat
  deriving DecidableEq, Repr

private def readGroups (onEnd : ReaderError) (consumed : Nat) (value : BitVec 64) :
    List (BitVec 8) → VarintReadResult (BitVec 64)
  | [] =>
    let err := if consumed > 0 && onEnd == .eof then ReaderError.unexpectedEOF else onEnd
    ⟨value, some err, consumed + 1⟩
  | byte :: tail =>
    let calls := consumed + 1
    let group := BitVec.ofNat 64 (byte.toNat % 128) <<< (7 * consumed)
    let combined := value ||| group
    if byte.toNat < 128 then
      if consumed == 9 && byte.toNat > 1 then ⟨value, some .overflow, calls⟩
      else ⟨combined, none, calls⟩
    else if calls == 10 then ⟨combined, some .overflow, calls⟩
    else readGroups onEnd calls combined tail

def readUvarint (bytes : List (BitVec 8)) (onEnd : ReaderError) : VarintReadResult (BitVec 64) :=
  readGroups onEnd 0 0 bytes

def readVarint (bytes : List (BitVec 8)) (onEnd : ReaderError) : VarintReadResult Int :=
  let result := readUvarint bytes onEnd
  let magnitude := Int.ofNat (result.value.toNat / 2)
  let value := if result.value.toNat % 2 == 0 then magnitude else -magnitude - 1
  ⟨value, result.error, result.calls⟩

theorem readUvarint_empty (err : ReaderError) :
    readUvarint [] err = ⟨0, some err, 1⟩ := rfl

theorem readUvarint_eof_after_continuation :
    readUvarint [129] .eof = ⟨1, some .unexpectedEOF, 2⟩ := by decide

theorem readUvarint_tenth_continuation :
    readUvarint (List.replicate 9 128 ++ [129, 0]) .eof =
      ⟨BitVec.ofNat 64 (2^63), some .overflow, 10⟩ := by decide

theorem readVarint_status (bytes : List (BitVec 8)) (err : ReaderError) :
    ((readVarint bytes err).error, (readVarint bytes err).calls) =
      ((readUvarint bytes err).error, (readUvarint bytes err).calls) := rfl

theorem readVarint_partial :
    readVarint [129] .unexpectedEOF = ⟨-1, some .unexpectedEOF, 2⟩ := by decide

end GoLean.Encoding.Binary
