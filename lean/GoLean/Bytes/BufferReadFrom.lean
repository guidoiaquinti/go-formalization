import GoLean.Bytes.BufferReset

namespace GoLean.Bytes

inductive ReadFromReply where
  | returns (count : Int) (eof failed : Bool)
  | panics

structure ReadFromStep where
  written : List (BitVec 8)
  reply : ReadFromReply

inductive ReadFromStatus where
  | ok
  | readError
  | readerPanic
  | nilReader
  | negativeRead
  | bounds
  | tooLarge
  | incomplete

structure ReadFromResult where
  count : Nat
  status : ReadFromStatus
  state : BufferReadState
  calls : List (List (BitVec 8))

/-- Pinned Go 1.27.0 noscan byte-slice class sizes from
    src/internal/runtime/gc/sizeclasses.go, through 32 KiB. -/
private def readFromCapacity (requested : Nat) : Nat :=
  ([512, 576, 640, 704, 768, 896, 1024, 1152, 1280, 1408, 1536, 1792, 2048,
    2304, 2688, 3072, 3200, 3456, 4096, 4864, 5376, 6144, 6528, 6784, 6912,
    8192, 9472, 9728, 10240, 10880, 12288, 13568, 14336, 16384, 18432, 19072,
    20480, 21760, 24576, 27264, 28672, 32768].find? (requested ≤ ·)).getD requested

/-- The effect of grow(MinRead) immediately followed by shortening the slice
    to the insertion index. MinRead is larger than the small initial buffer. -/
private def prepareReadFrom (s : BufferReadState) (maxInt : Nat) :
    Except BufferReadState BufferReadState :=
  let prior := s.tail
  let s := if prior.offset > 0 && prior.offset = (prior.stored.getD []).length
    then (bufferReset (some s)).getD s else s
  let tail := s.tail
  let bytes := tail.stored.getD []
  let backing := bytes ++ tail.spare
  let unread := bytes.drop tail.offset
  if 512 ≤ tail.spare.length then .ok s
  else if unread.length + 512 ≤ backing.length / 2 then
    .ok { s with tail := ⟨some unread, backing.drop unread.length, 0⟩ }
  else if maxInt < 2 * backing.length + 512 then .error s
  else
    let requested := max (bytes.length + 512) (2 * (backing.length - tail.offset))
    let capacity := readFromCapacity requested
    .ok { s with tail := ⟨some unread, List.replicate (capacity - unread.length) 0, 0⟩ }

private def readFromLoop (maxInt count : Nat) (state : BufferReadState) :
    List ReadFromStep → ReadFromResult
  | [] => ⟨count, .incomplete, state, []⟩
  | step :: rest =>
    match prepareReadFrom state maxInt with
    | .error prior => ⟨count, .tooLarge, prior, []⟩
    | .ok prepared =>
      let tail := prepared.tail
      let offered := tail.spare
      let written := step.written.take offered.length
      let scratch := written ++ offered.drop written.length
      let scratchState := { prepared with tail := { tail with spare := scratch } }
      match step.reply with
      | .panics => ⟨count, .readerPanic, scratchState, [offered]⟩
      | .returns n eof failed =>
        if n < 0 then ⟨count, .negativeRead, scratchState, [offered]⟩
        else if scratch.length < n.toNat then ⟨count, .bounds, scratchState, [offered]⟩
        else
          let stored := tail.stored.getD [] ++ scratch.take n.toNat
          let next := { prepared with tail :=
            ⟨some stored, scratch.drop n.toNat, tail.offset⟩ }
          let count := count + n.toNat
          if eof then ⟨count, .ok, next, [offered]⟩
          else if failed then ⟨count, .readError, next, [offered]⟩
          else
            let result := readFromLoop maxInt count next rest
            { result with calls := offered :: result.calls }

/-- Go 1.27.0 src/bytes/buffer.go:224–244. Clear the read flag, reserve
    capacity, offer all spare storage, commit the reported bytes, and
    suppress only io.EOF. A missing script models a nil io.Reader. -/
def bufferReadFrom (receiver : Option BufferReadState) (script : Option (List ReadFromStep))
    (maxInt : Nat) : Option ReadFromResult :=
  receiver.map fun prior =>
    let state := { prior with lastRead := 0 }
    match script with
    | none =>
      match prepareReadFrom state maxInt with
      | .error saved => ⟨0, .tooLarge, saved, []⟩
      | .ok saved => ⟨0, .nilReader, saved, []⟩
    | some steps => readFromLoop maxInt 0 state steps

theorem bufferReadFrom_nil (script : Option (List ReadFromStep)) (maxInt : Nat) :
    bufferReadFrom none script maxInt = none := rfl

set_option maxRecDepth 4096 in
theorem bufferReadFrom_eof :
    (bufferReadFrom (some ⟨⟨none, [], 0⟩, 3⟩)
      (some [⟨[7, 8], .returns 1 true false⟩]) 2147483647).map
      (fun r => (r.count, r.status, r.state.tail.stored, r.state.tail.spare.take 1, r.state.lastRead)) =
    some (1, .ok, some [7], [8], 0) := rfl

set_option maxRecDepth 4096 in
theorem bufferReadFrom_negative_keepsScratch :
    (bufferReadFrom (some ⟨⟨none, [], 0⟩, 3⟩)
      (some [⟨[7, 8], .returns (-1) true false⟩]) 2147483647).map
      (fun r => (r.status, r.state.tail.stored, r.state.tail.spare.take 2)) =
    some (.negativeRead, some [], [7, 8]) := rfl

end GoLean.Bytes
