import GoLean.Math.Bits.Population
import GoDiff.Math.Bits.Unary

def main : IO Unit :=
  runUnary 64 fun x => (GoLean.Math.Bits.onesCount64 x)
