import GoLean.Math.Bits.TrailingZeros
import GoDiff.Math.Bits.Unary

def main : IO Unit :=
  runUnary 32 fun x => (GoLean.Math.Bits.trailingZeros32 x)
