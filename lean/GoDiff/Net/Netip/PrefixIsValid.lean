import GoLean.Net.Netip.PrefixIsValid

open GoLean.Net.Netip

def main : IO Unit := do
  for tag in [:130] do
    IO.println (prefixIsValid ⟨⟨.invalid, 0, 0, ""⟩, tag.toUInt8⟩)
