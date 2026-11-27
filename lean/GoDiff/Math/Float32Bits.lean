import GoLean.Math.Float32Bits
import GoDiff.Math.Bits.Unary

def main : IO Unit :=
  runUnary 32 fun bits => (GoLean.Math.float32Bits ⟨bits⟩).toNat
