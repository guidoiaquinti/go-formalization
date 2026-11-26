import GoLean.Net.Netip.IPv6LinkLocalAllNodes

def main : IO Unit := do
  let ip := GoLean.Net.Netip.ipv6LinkLocalAllNodes
  IO.println s!"{ip.hi} {ip.lo} {ip.isV6} [{ip.zone}]"
