import GoLean.Math.Exp2

namespace GoLean.Math

/-- The Bessel amplitude helpers use four coefficient intervals. -/
def besselCoefficients (x : Nat) (atEight atMiddle atLow atTwo : List Nat × List Nat) :
    List Nat × List Nat :=
  if 0x4020000000000000 ≤ x then atEight
  else if 0x40122e7d566cf41f ≤ x then atMiddle
  else if 0x4006db573eab367a ≤ x then atLow
  else atTwo

/-- Go 1.27.0 src/math/j0.go and j1.go. Coefficients are in descending
    order. Arm64 contracts the Horner stages, and amd64 v3/v4 contracts
    the final addition of the denominator's leading one. -/
def besselRational (arm64 compileFused : Bool) (x : Nat) (p q : List Nat) : Nat :=
  let square := multiplyWords x x
  let z := if square = 0x7ff0000000000000 then 0 else divideWords 0x3ff0000000000000 square
  let evaluate := fun (coefficients : List Nat) => coefficients.tail.foldl
    (fun acc coefficient => multiplyAddWords arm64 z acc coefficient) (coefficients.headD 0)
  let denominator := multiplyAddWords (arm64 || compileFused) z (evaluate q) 0x3ff0000000000000
  divideWords (evaluate p) denominator

end GoLean.Math
