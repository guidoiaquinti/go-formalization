import GoLean.Net.Netip.IPv6Unspecified

def main : IO Unit := do
  let ip := GoLean.Net.Netip.ipv6Unspecified
  IO.println s!"{ip.hi} {ip.lo} {ip.isV6} [{ip.zone}]"
