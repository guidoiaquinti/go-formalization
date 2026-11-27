import GoLean.Math.FmaNonzero
import GoDiff.Math.Bits.Unary

def main : IO Unit :=
  runUnary 64 fun word => (GoLean.Math.fmaNonzero word).toNat
