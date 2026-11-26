import GoLean.Net.Netip.IPv6LinkLocalAllRouters

def main : IO Unit := do
  let ip := GoLean.Net.Netip.ipv6LinkLocalAllRouters
  IO.println s!"{ip.hi} {ip.lo} {ip.isV6} [{ip.zone}]"
