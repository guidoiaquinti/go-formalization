import GoLean.Math.RoundToEven
import GoDiff.Math.Bits.Unary

def main : IO Unit := runUnary 64 fun word =>
  let bits := GoLean.Math.roundToEvenWord word
  if bits / 2 ^ 52 % 2048 == 2047 && bits % 2 ^ 52 != 0 then 2 ^ 64 else bits
