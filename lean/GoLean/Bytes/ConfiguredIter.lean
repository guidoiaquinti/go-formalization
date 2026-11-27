import GoLean.Strings.Substring

namespace GoLean.Bytes

structure IterState where
  bytes : Option (List Nat)
  spare : List Nat

structure IterFragment where
  bytes : Option (List Nat)
  capacity : Nat

structure IterRun where
  fragments : List IterFragment
  remaining : IterState
  boundsPanic : Bool

private def clipped (bytes : Option (List Nat)) : IterFragment :=
  ⟨bytes, (bytes.getD []).length⟩

/-- Go 1.27.0 src/bytes/iter.go:18–33. Lines consumes its captured
    suffix before calling the callback and clips the yielded capacity. -/
def lines (input : IterState) : Nat → IterRun
  | 0 => ⟨[], input, false⟩
  | steps + 1 =>
    let bytes := input.bytes.getD []
    if bytes.isEmpty then ⟨[], input, false⟩
    else
      let (part, next) := match GoLean.Strings.index bytes [10] with
        | none => (clipped input.bytes, ⟨none, []⟩)
        | some i => (clipped (some (bytes.take (i + 1))),
            ⟨some (bytes.drop (i + 1)), input.spare⟩)
      let result := lines next steps
      {result with fragments := part :: result.fragments}

private def configuredStep (input : IterState) (separator : List Nat) (save : Int) :
    Except Unit (Option (IterFragment × Option IterState)) :=
  let bytes := input.bytes.getD []
  if separator.isEmpty then
    match bytes with
    | [] => .ok none
    | lead :: rest =>
      let width := GoLean.Strings.runeWidth lead rest
      .ok (some (clipped (some (bytes.take width)),
        some ⟨some (bytes.drop width), input.spare⟩))
  else
    match GoLean.Strings.index bytes separator with
    | none => .ok (some (clipped input.bytes, none))
    | some i =>
      let endpoint := Int.ofNat i + save
      let storage := bytes ++ input.spare
      if endpoint < 0 || Int.ofNat storage.length < endpoint then .error ()
      else .ok (some (clipped (some (storage.take endpoint.toNat)),
        some ⟨some (bytes.drop (i + separator.length)), input.spare⟩))

/-- Go 1.27.0 src/bytes/iter.go:36–60. The two-index fragment expression
    can reach spare backing up to capacity. Stopping leaves the current
    suffix intact. Positive signed overflow is already out of bounds. -/
def splitSeqConfigured (input : IterState) (separator : List Nat) (save : Int) : Nat → IterRun
  | 0 => ⟨[], input, false⟩
  | steps + 1 =>
    match configuredStep input separator save with
    | .error _ => ⟨[], input, true⟩
    | .ok none => ⟨[], input, false⟩
    | .ok (some (part, none)) => ⟨[part], input, false⟩
    | .ok (some (part, some next)) =>
      if steps = 0 then ⟨[part], input, false⟩
      else
        let result := splitSeqConfigured next separator save steps
        {result with fragments := part :: result.fragments}

theorem lines_stop_after_newline :
    (lines ⟨some [10, 65], [92]⟩ 1).remaining = ⟨some [65], [92]⟩ := rfl

theorem lines_no_phantom_final :
    (lines ⟨some [10], [92]⟩ 3).fragments = [⟨some [10], 1⟩] := rfl

theorem lines_empty (spare : List Nat) (budget : Nat) :
    (lines ⟨some [], spare⟩ budget).fragments = [] := by cases budget <;> rfl

theorem splitSeqConfigured_stopped (input : IterState) (separator : List Nat) (save : Int) :
    (splitSeqConfigured input separator save 1).remaining = input := by
  unfold splitSeqConfigured
  split <;> simp

theorem splitSeqConfigured_spare :
    (splitSeqConfigured ⟨some [97, 45], [88, 99]⟩ [45] 3 1).fragments =
      [⟨some [97, 45, 88, 99], 4⟩] := rfl

theorem splitSeqConfigured_partialPanic :
    let result := splitSeqConfigured ⟨some [97, 98, 45, 45], [99]⟩ [45] (-1) 8
    result.fragments = [⟨some [97], 1⟩] ∧ result.remaining = ⟨some [45], [99]⟩ ∧
      result.boundsPanic = true := ⟨rfl, rfl, rfl⟩

end GoLean.Bytes
