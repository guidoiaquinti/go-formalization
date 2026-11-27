import GoLean.Math.ArchStubs
import GoDiff.Math.Bits.Tuple

open GoLean.Math

def main : IO Unit := runTuple fun fields => do
  let [platform, selector, a, b] := fields |
    throw (IO.userError "expected architecture profile, selector and operands")
  let profile := if platform == 1 then ArchStubProfile.s390x else ArchStubProfile.nonS390x
  let x ← wordInput 64 a
  let result ← match selector with
  | 0 => pure (archAcosStub profile x)
  | 1 => pure (archAcoshStub profile x)
  | 2 => pure (archAsinStub profile x)
  | 3 => pure (archAsinhStub profile x)
  | 4 => pure (archAtanStub profile x)
  | 5 => pure (archAtan2Stub profile x (← wordInput 64 b))
  | 6 => pure (archAtanhStub profile x)
  | 7 => pure (archCbrtStub profile x)
  | 8 => pure (archCosStub profile x)
  | 9 => pure (archCoshStub profile x)
  | 10 => pure (archErfStub profile x)
  | 11 => pure (archErfcStub profile x)
  | 12 => pure (archExpm1Stub profile x)
  | 13 => pure (archFrexpStub profile x)
  | 14 => pure (archLdexpStub profile x b)
  | 15 => pure (archLog10Stub profile x)
  | 16 => pure (archLog2Stub profile x)
  | 17 => pure (archLog1pStub profile x)
  | 18 => pure (archModStub profile x (← wordInput 64 b))
  | 19 => pure (archPowStub profile x (← wordInput 64 b))
  | 20 => pure (archRemainderStub profile x (← wordInput 64 b))
  | 21 => pure (archSinStub profile x)
  | 22 => pure (archSinhStub profile x)
  | 23 => pure (archTanStub profile x)
  | 24 => pure (archTanhStub profile x)
  | _ => throw (IO.userError "unknown architecture stub selector")
  match result with
  | none => return "outside-scope"
  | some outcome => return s!"{outcome.available}:panic:{outcome.panicMessage}"
