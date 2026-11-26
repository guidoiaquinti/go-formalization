import GoLean.Math.Bits.Population
import GoDiff.Math.Bits.Unary

def main : IO Unit :=
  runUnary 32 fun x => (GoLean.Math.Bits.onesCount32 x)
