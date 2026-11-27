import GoLean.Compress.Flate.HCode

namespace GoLean.Compress.Flate

inductive CostResult where
  | value (cost : Nat)
  | nilReceiver
  | bounds
  deriving DecidableEq

private def lookup (codes : Option (List HCode)) (index : Nat) :
    Except CostResult HCode :=
  match codes with
  | none => .error .nilReceiver
  | some table => match table[index]? with
    | none => .error .bounds
    | some word => .ok word

private def frequencyCost (codes : Option (List HCode)) (requireCode : Bool) :
    Nat → Nat → List (BitVec 16) → CostResult
  | _, total, [] => .value total
  | index, total, freq :: rest =>
    if freq == 0 then frequencyCost codes requireCode (index + 1) total rest
    else match lookup codes index with
      | .error failure => failure
      | .ok word =>
        if requireCode && hCodeZero word then .value 2147483647
        else frequencyCost codes requireCode (index + 1)
          (total + freq.toNat * (hCodeLen word).toNat) rest

private def wrap (width : Nat) : CostResult → CostResult
  | .value n => .value (n % 2 ^ width)
  | outcome => outcome

def bitLength (width : Nat) (codes : Option (List HCode)) (freq : List (BitVec 16)) :
    CostResult := wrap width (frequencyCost codes false 0 0 freq)

def canEncodeLen (width : Nat) (codes : Option (List HCode)) (freq : List (BitVec 16)) :
    CostResult := wrap width (frequencyCost codes true 0 0 freq)

private def rawCost (codes : Option (List HCode)) : Nat → List (BitVec 8) → CostResult
  | total, [] => .value total
  | total, byte :: rest =>
    match lookup codes byte.toNat with
    | .error failure => failure
    | .ok word => rawCost codes (total + max 1 (hCodeLen word).toNat) rest

def bitLengthRaw (width : Nat) (codes : Option (List HCode)) (data : List (BitVec 8)) :
    CostResult := wrap width (rawCost codes 0 data)

theorem bitLength_skipsZero : bitLength 64 none [0, 0] = .value 0 := rfl

theorem canEncodeLen_fullWord :
    canEncodeLen 64 (some [256, 1]) [1, 2] = .value 2 ∧
    canEncodeLen 64 (some [0]) [1, 2] = .value 2147483647 := by decide

theorem bitLengthRaw_minimum :
    bitLengthRaw 64 (some [0, 256, 2]) [0, 1, 2] = .value 4 := rfl

end GoLean.Compress.Flate
