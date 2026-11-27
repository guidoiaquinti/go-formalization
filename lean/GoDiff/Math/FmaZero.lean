import GoLean.Math.FmaZero
import GoDiff.Math.Bits.Unary

def main : IO Unit :=
  runUnary 64 fun word => (GoLean.Math.fmaZero word).toNat
