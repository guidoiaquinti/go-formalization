import GoLean.Bytes.BufferReadFrom
import GoDiff.Math.Bits.Tuple

open GoLean.Bytes

private def byteRuns : List (BitVec 8) → List (BitVec 8 × Nat)
  | [] => []
  | b :: bs =>
    match byteRuns bs with
    | [] => [(b, 1)]
    | (c, n) :: rest => if b = c then (c, n + 1) :: rest else (b, 1) :: (c, n) :: rest

private def wire (xs : List (BitVec 8)) : String :=
  String.intercalate "," ((byteRuns xs).map fun (b, n) => s!"{b.toNat}.{n}")

private def parseSteps : List Int → IO (List ReadFromStep)
  | [] => pure []
  | mode :: count :: writes :: value :: rest => do
    let b ← wordInput 8 value
    let reply := if mode = 3 then ReadFromReply.panics
      else ReadFromReply.returns count (mode = 1) (mode = 2)
    return ⟨List.replicate writes.toNat b, reply⟩ :: (← parseSteps rest)
  | _ => throw (IO.userError "expected complete read reply quadruple")

def main : IO Unit := runTuple fun row => do
  let tag :: len :: off :: cap :: flag :: nilReader :: maxInt :: raw := row |
    throw (IO.userError "expected buffer geometry and read script")
  let steps ← parseSteps raw
  let stored := List.replicate off.toNat (11 : BitVec 8) ++
    List.replicate (len.toNat - off.toNat) 22
  let spare := List.replicate (cap.toNat - len.toNat) (33 : BitVec 8)
  let state : Option BufferReadState := if tag = 0 then none else
    some ⟨⟨(if tag = 1 then none else some stored), spare, off.toNat⟩, flag⟩
  let r := bufferReadFrom state (if nilReader = 1 then none else some steps) maxInt.toNat
  match r with
  | none => return "nilReceiver"
  | some r =>
    let status := match r.status with
      | .ok => "ok"
      | .readError => "readError"
      | .readerPanic => "readerPanic"
      | .nilReader => "nilReader"
      | .negativeRead => "negativeRead"
      | .bounds => "bounds"
      | .tooLarge => "tooLarge"
      | .incomplete => "incomplete"
    let tail := r.state.tail
    let marker := if tail.stored.isNone then "nil" else "slice"
    let trace := String.intercalate ";" (r.calls.map fun offer => s!"{offer.length}:{wire offer}")
    -- Go cannot return a partial count from a panicking invocation.
    let count := if status = "ok" || status = "readError" then r.count else 0
    return s!"{status}|{count}|{marker}|{tail.offset}|{r.state.lastRead}|{wire (tail.stored.getD [])}|{wire tail.spare}|{trace}"
