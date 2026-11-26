import GoLean.Math.Bits.TrailingZeros
import GoDiff.Math.Bits.Unary

def main : IO Unit :=
  runUnary 64 fun x => (GoLean.Math.Bits.trailingZeros64 x)
