import GoLean.Net.Netip.AddrPortPort

open GoLean.Net.Netip

def main : IO Unit := do
  let stdout ← IO.getStdout
  for n in [:65536] do
    let port := addrPortPort ⟨⟨.invalid, 0, 0, ""⟩, n.toUInt16⟩
    stdout.putStrLn (toString port)
