import GoLean.Strconv.FormatBool

def main : IO Unit := do
  for b in [false, true] do
    IO.println (GoLean.Strconv.formatBool b)
