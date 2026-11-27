namespace GoLean.Math.RandV2

inductive WordSource
  | word (value : Nat)
  | raises
  deriving DecidableEq

structure RandState where
  source : Option WordSource

def newRand (source : Option WordSource) : RandState := ⟨source⟩

structure WordResult where
  status : String
  value : Nat
  pulls : Nat

private def pull (receiver : Option RandState) (f : Nat → Nat) : WordResult :=
  match receiver with
  | none => ⟨"nil", 0, 0⟩
  | some r => match r.source with
    | none => ⟨"nil", 0, 0⟩
    | some .raises => ⟨"callback", 0, 1⟩
    | some (.word w) => ⟨"ok", f w, 1⟩

def randInt64 (r : Option RandState) := pull r (· % 2^63)
def randUint32 (r : Option RandState) := pull r (· / 2^32)
def randUint64 (r : Option RandState) := pull r id
def randInt32 (r : Option RandState) := pull r (· / 2^33)
def randInt (width : Nat) (r : Option RandState) := pull r (· % 2^(width - 1))
def randUint (width : Nat) (r : Option RandState) := pull r (· % 2^width)

/-- The stored number represents the numerator over the fixed denominator
    2^53 or 2^24. Every numerator in range is exactly representable. -/
def randFloat64 (r : Option RandState) := pull r (· % 2^53)
def randFloat32 (r : Option RandState) := pull r (fun w => w / 2^32 % 2^24)

theorem newRand_source (s : Option WordSource) : (newRand s).source = s := rfl

theorem randInt64_word (w : Nat) : (randInt64 (some (newRand (some (.word w))))).value = w % 2^63 := rfl
theorem randUint32_word (w : Nat) : (randUint32 (some (newRand (some (.word w))))).value = w / 2^32 := rfl
theorem randUint64_word (w : Nat) : (randUint64 (some (newRand (some (.word w))))).value = w := rfl
theorem randInt32_word (w : Nat) : (randInt32 (some (newRand (some (.word w))))).value = w / 2^33 := rfl
theorem randInt_word (width w : Nat) :
    (randInt width (some (newRand (some (.word w))))).value = w % 2^(width - 1) := rfl
theorem randUint_word (width w : Nat) :
    (randUint width (some (newRand (some (.word w))))).value = w % 2^width := rfl
theorem randFloat64_word (w : Nat) :
    (randFloat64 (some (newRand (some (.word w))))).value = w % 2^53 := rfl
theorem randFloat32_word (w : Nat) :
    (randFloat32 (some (newRand (some (.word w))))).value = w / 2^32 % 2^24 := rfl

end GoLean.Math.RandV2
