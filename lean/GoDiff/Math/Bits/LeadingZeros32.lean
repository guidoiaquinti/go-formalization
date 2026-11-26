import GoLean.Math.Bits.Length
import GoDiff.Math.Bits.Unary

def main : IO Unit :=
  runUnary 32 fun x => (GoLean.Math.Bits.leadingZeros32 x)
