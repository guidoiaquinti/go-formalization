import GoLean.Math.Float64Bits
import GoDiff.Math.Bits.Unary

def main : IO Unit :=
  runUnary 64 fun bits => (GoLean.Math.float64Bits ⟨bits⟩).toNat
