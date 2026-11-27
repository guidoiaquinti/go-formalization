namespace GoLean.Compress.Flate

/-- The counter and array widths are those in Go 1.27.0 token.go:131–138.
    Index functions model fixed arrays on their valid index ranges. -/
structure TokenState where
  n : UInt16
  nFilled : Int
  literalHist : Nat → UInt16
  extraHist : Nat → UInt16
  offsetHist : Nat → UInt16
  backing : Nat → UInt32

private def put {α : Type} (array : Nat → α) (i : Nat) (value : α) : Nat → α :=
  fun j => if j == i then value else array j

/-- Go 1.27.0 token.go:141–150. -/
def tokensReset (t : TokenState) : TokenState :=
  if t.n == 0 then t
  else { t with
    n := 0
    nFilled := 0
    literalHist := fun _ => 0
    extraHist := fun _ => 0
    offsetHist := fun _ => 0 }

/-- Go 1.27.0 token.go:181–185. -/
def tokensAddLiteral (t : TokenState) (lit : UInt8) : TokenState :=
  { t with
    backing := put t.backing t.n.toNat (UInt32.ofNat lit.toNat)
    literalHist := put t.literalHist lit.toNat (t.literalHist lit.toNat + 1)
    n := t.n + 1 }

/-- Go 1.27.0 token.go:172–178. -/
def emitLiterals (dst : TokenState) (lit : List UInt8) : TokenState :=
  lit.foldl tokensAddLiteral dst

/-- Go 1.27.0 token.go:277–281, with endBlockMarker from huffman_bit_writer.go. -/
def tokensAddEOB (t : TokenState) : TokenState :=
  { t with
    backing := put t.backing t.n.toNat 256
    extraHist := put t.extraHist 0 (t.extraHist 0 + 1)
    n := t.n + 1 }

/-- A view rooted at the receiver's backing array. Reads and writes take the
    current receiver so an old view observes later writes to that storage. -/
structure TokenView where
  length : Nat
  capacity : Nat

/-- Go 1.27.0 token.go:284–286. -/
def tokensSlice (t : TokenState) : TokenView := ⟨t.n.toNat, 65536⟩

/-- Read after reslicing to capacity. -/
def tokenViewRead (view : TokenView) (t : TokenState) (i : Nat) : Option UInt32 :=
  if i < view.capacity then some (t.backing i) else none

/-- Write after reslicing to capacity. -/
def tokenViewWrite (view : TokenView) (t : TokenState) (i : Nat) (value : UInt32) :
    TokenState :=
  if i < view.capacity then { t with backing := put t.backing i value } else t

theorem tokensReset_zero_count (t : TokenState) (zero : t.n = 0) :
    tokensReset t = t := by simp [tokensReset, zero]

theorem tokensAddLiteral_written (t : TokenState) (lit : UInt8) :
    (tokensAddLiteral t lit).backing t.n.toNat = UInt32.ofNat lit.toNat := by
  simp [tokensAddLiteral, put]

theorem emitLiterals_empty (t : TokenState) : emitLiterals t [] = t := rfl

theorem tokensAddEOB_written (t : TokenState) :
    (tokensAddEOB t).backing t.n.toNat = 256 := by simp [tokensAddEOB, put]

theorem tokensSlice_last_cell (t : TokenState) :
    tokenViewRead (tokensSlice t) t 65535 = some (t.backing 65535) := by
  change (if 65535 < 65536 then some (t.backing 65535) else none) = _
  have bound : (65535 : Nat) < 65536 := by omega
  simp [bound]

end GoLean.Compress.Flate
