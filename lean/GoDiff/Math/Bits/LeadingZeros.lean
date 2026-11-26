import GoLean.Math.Bits.Length
import GoDiff.Math.Bits.Unary

def main : IO Unit :=
  runUnary 64 fun x => (GoLean.Math.Bits.leadingZeros x)
