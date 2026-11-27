import GoLean.Encoding.Binary.EndianNames

open GoLean.Encoding.Binary

def main : IO Unit := do
  IO.println bigEndianString
  IO.println bigEndianGoString
  IO.println littleEndianString
  IO.println littleEndianGoString
  IO.println nativeEndianString
  IO.println nativeEndianGoString
