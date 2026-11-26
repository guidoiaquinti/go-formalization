import GoLean.Net.Netip.IPv6Loopback

def main : IO Unit := do
  let ip := GoLean.Net.Netip.ipv6Loopback
  IO.println s!"{ip.hi} {ip.lo} {ip.isV6} [{ip.zone}]"
