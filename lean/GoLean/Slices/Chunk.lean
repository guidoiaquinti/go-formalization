namespace GoLean.Slices

structure ChunkFragment (α : Type) where
  values : List α
  capacity : Nat

structure ChunkRun (α : Type) where
  yields : List (ChunkFragment α)
  backing : List α

private def chunkSteps (size length : Nat) (write : Option α) :
    Nat → Nat → List α → ChunkRun α
  | 0, _, backing => ⟨[], backing⟩
  | fuel + 1, begin, backing =>
    if length ≤ begin then ⟨[], backing⟩
    else
      let width := min size (length - begin)
      let fragment := ⟨(backing.drop begin).take width, width⟩
      let backing := (write.map fun value => backing.set begin value).getD backing
      let rest := chunkSteps size length write fuel (begin + size) backing
      ⟨fragment :: rest.yields, rest.backing⟩

/-- Go 1.27.0 src/slices/iter.go:97–113. The size guard happens before
    iterator construction. Each invocation starts at zero and produces
    capacity-clipped aliases. The callback budget includes its final false
    result. The optional write models overwriting the first cell of each alias. -/
def sliceChunk (backing : List α) (length : Nat) (size : Int) (budget : Nat)
    (write : Option α) : Except String (ChunkRun α) :=
  if size < 1 then .error "cannot be less than 1"
  else .ok (chunkSteps size.toNat length write budget 0 backing)

theorem sliceChunk_badSize (backing : List α) (length budget : Nat) (size : Int)
    (write : Option α) (bad : size < 1) :
    sliceChunk backing length size budget write = .error "cannot be less than 1" := by
  simp [sliceChunk, bad]

theorem sliceChunk_shortLast :
    sliceChunk [1, 2, 3, 99] 3 2 5 none =
      .ok ⟨[⟨[1, 2], 2⟩, ⟨[3], 1⟩], [1, 2, 3, 99]⟩ := rfl

theorem sliceChunk_stoppedWrite :
    sliceChunk [1, 2, 3, 99] 3 2 1 (some 8) =
      .ok ⟨[⟨[1, 2], 2⟩], [8, 2, 3, 99]⟩ := rfl

end GoLean.Slices
