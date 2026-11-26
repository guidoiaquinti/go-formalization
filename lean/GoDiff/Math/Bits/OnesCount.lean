import GoLean.Math.Bits.Population
import GoDiff.Math.Bits.Unary

def main : IO Unit :=
  runUnary 64 GoLean.Math.Bits.onesCount
