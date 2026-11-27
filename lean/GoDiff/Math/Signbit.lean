import GoLean.Math.Signbit
import GoDiff.Math.Bits.Unary

def main : IO Unit :=
  runUnary 64 fun x => (GoLean.Math.signbit x).toNat
