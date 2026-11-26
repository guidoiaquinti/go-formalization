import GoLean.Net.Netip.PrefixBits

open GoLean.Net.Netip

def main : IO Unit := do
  for tag in [:130] do
    IO.println (prefixBits ⟨⟨.invalid, 0, 0, ""⟩, tag.toUInt8⟩)
