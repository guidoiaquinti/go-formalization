import GoLean.Math.Abs
import GoDiff.Math.Bits.Unary

def main : IO Unit :=
  runUnary 64 fun bits => (GoLean.Math.abs bits).toNat
