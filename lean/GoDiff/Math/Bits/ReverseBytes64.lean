import GoLean.Math.Bits.Reversal
import GoDiff.Math.Bits.Unary

def main : IO Unit :=
  runUnary 64 fun x => (GoLean.Math.Bits.reverseBytes64 x).toNat
