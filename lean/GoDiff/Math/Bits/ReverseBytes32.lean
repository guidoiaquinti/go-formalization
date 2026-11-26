import GoLean.Math.Bits.Reversal
import GoDiff.Math.Bits.Unary

def main : IO Unit :=
  runUnary 32 fun x => (GoLean.Math.Bits.reverseBytes32 x).toNat
