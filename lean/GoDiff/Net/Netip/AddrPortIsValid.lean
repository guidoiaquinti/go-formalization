import GoLean.Net.Netip.AddrPortIsValid

open GoLean.Net.Netip

def main : IO Unit := do
  let stdout ← IO.getStdout
  for detail in [AddrDetailClass.invalid, .ipv4, .ipv6] do
    for n in [:65536] do
      stdout.putStrLn (toString (addrPortIsValid ⟨⟨detail, 0, 0, ""⟩, n.toUInt16⟩))
