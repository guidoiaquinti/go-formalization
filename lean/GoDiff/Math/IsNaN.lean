import GoLean.Math.IsNaN
import GoDiff.Math.Bits.Unary

def main : IO Unit :=
  runUnary 64 fun x => (GoLean.Math.isNaN x).toNat
