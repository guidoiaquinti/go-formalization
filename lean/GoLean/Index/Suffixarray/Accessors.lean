import Init.Data.BitVec.Basic

namespace GoLean.Index.Suffixarray

structure View (α : Type) where
  backing : Option (List α)
  length : Nat

structure Ints where
  words32 : View (BitVec 32)
  words64 : View (BitVec 64)

structure IndexState where
  data : View (BitVec 8)
  sa : Ints

inductive AccessFailure where
  | nilReceiver
  | bounds

private def nilView : View α := ⟨none, 0⟩

private def selectIndex (v : View α) (i : Int) : Except AccessFailure α :=
  if i < 0 || v.length ≤ i.toNat then .error .bounds
  else match (v.backing.getD [])[i.toNat]? with
    | none => .error .bounds
    | some value => .ok value

private def replaceIndex (v : View α) (i : Int) (value : α) : Except AccessFailure (View α) :=
  (selectIndex v i).map fun _ => { v with backing := v.backing.map (·.set i.toNat value) }

private def subview (v : View α) (i j : Int) : Except AccessFailure (View α) :=
  if i < 0 || j < i || (v.backing.getD []).length < j.toNat then .error .bounds
  else .ok ⟨v.backing.map (·.drop i.toNat), (j - i).toNat⟩

/-- Go 1.27.0 src/index/suffixarray/suffixarray.go:48–50. Both lengths contribute. -/
def intsLen (s : Option Ints) : Except AccessFailure Nat :=
  match s with
  | none => .error .nilReceiver
  | some s => .ok (s.words32.length + s.words64.length)

/-- Go 1.27.0 suffixarray.go:52–57. Selection is by nilness, with signed
    widening for the 32-bit representation. -/
def intsGet (s : Option Ints) (i : Int) : Except AccessFailure Int :=
  match s with
  | none => .error .nilReceiver
  | some s =>
    if s.words32.backing.isSome then (selectIndex s.words32 i).map BitVec.toInt
    else (selectIndex s.words64 i).map BitVec.toInt

/-- Go 1.27.0 suffixarray.go:59–65. A successful write changes only one
    cell; conversion to int32 preserves the low word. -/
def intsSet (s : Option Ints) (i value : Int) : Except AccessFailure Ints :=
  match s with
  | none => .error .nilReceiver
  | some s =>
    if s.words32.backing.isSome then
      (replaceIndex s.words32 i (BitVec.ofInt 32 value)).map fun v => {s with words32 := v}
    else (replaceIndex s.words64 i (BitVec.ofInt 64 value)).map fun v => {s with words64 := v}

/-- Go 1.27.0 suffixarray.go:67–72. The selected view can extend through
    capacity, and the returned capacity shrinks by the starting offset. -/
def intsSlice (s : Option Ints) (i j : Int) : Except AccessFailure Ints :=
  match s with
  | none => .error .nilReceiver
  | some s =>
    if s.words32.backing.isSome then
      (subview s.words32 i j).map fun v => ⟨v, nilView⟩
    else (subview s.words64 i j).map fun v => ⟨nilView, v⟩

/-- Go 1.27.0 suffixarray.go:241–243. Return the original data view. -/
def indexBytes (s : Option IndexState) : Except AccessFailure (View (BitVec 8)) :=
  match s with
  | none => .error .nilReceiver
  | some s => .ok s.data

/-- Go 1.27.0 suffixarray.go:245–247. Read the suffix offset and slice
    through the visible end, retaining spare capacity. -/
def indexAt (s : Option IndexState) (i : Int) : Except AccessFailure (View (BitVec 8)) := do
  let .some s := s | throw .nilReceiver
  let start ← intsGet (some s.sa) i
  subview s.data start (Int.ofNat s.data.length)

theorem intsLen_sum (a : Ints) :
    intsLen (some a) = .ok (a.words32.length + a.words64.length) := rfl

theorem intsGet_negative :
    intsGet (some ⟨⟨some [0xffffffff], 1⟩, ⟨some [2], 1⟩⟩) 0 = .ok (-1) := rfl

theorem intsSet_truncates :
    intsSet (some ⟨⟨some [8, 9], 1⟩, nilView⟩) 0 4294967297 =
      .ok ⟨⟨some [1, 9], 1⟩, nilView⟩ := rfl

theorem intsSlice_spare :
    intsSlice (some ⟨⟨some [8, 9], 1⟩, nilView⟩) 1 2 =
      .ok ⟨⟨some [9], 1⟩, nilView⟩ := rfl

theorem indexBytes_original (s : IndexState) : indexBytes (some s) = .ok s.data := rfl

theorem indexAt_spare :
    indexAt (some ⟨⟨some [6, 7, 8], 2⟩, ⟨⟨some [1], 1⟩, nilView⟩⟩) 0 =
      .ok ⟨some [7, 8], 1⟩ := rfl

end GoLean.Index.Suffixarray
