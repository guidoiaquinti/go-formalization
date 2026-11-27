import GoLean.Math.IsOddInt
import GoDiff.Math.Bits.Unary

def main : IO Unit :=
  runUnary 64 fun word => if GoLean.Math.isOddInt word then 1 else 0
