import GoLean.Math.IsNegInt
import GoDiff.Math.Bits.Unary

def main : IO Unit :=
  runUnary 64 fun word => if GoLean.Math.isNegInt word then 1 else 0
