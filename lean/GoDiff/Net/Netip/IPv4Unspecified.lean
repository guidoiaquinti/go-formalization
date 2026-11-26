import GoLean.Net.Netip.IPv4Unspecified

def main : IO Unit := do
  let ip := GoLean.Net.Netip.ipv4Unspecified
  IO.println s!"{ip.hi} {ip.lo} {ip.isV6} [{ip.zone}]"
