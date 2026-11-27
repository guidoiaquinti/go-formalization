namespace GoLean.Archive.Zip

structure WriteReply where
  count : Nat
  failed : Bool

structure CountWriteResult where
  count : BitVec 64
  returned : WriteReply
  calls : Nat

/-- Go 1.27.0 writer.go:706. The reply represents a lawful delegated write. -/
def countWriterWrite (count : BitVec 64) (reply : WriteReply) : CountWriteResult :=
  ⟨count + BitVec.ofNat 64 reply.count, reply, 1⟩

/-- Go 1.27.0 writer.go:716. The observed callback count is unchanged. -/
def nopCloserClose (calls : Nat) : Bool × Nat := (false, calls)

structure FlushResult where
  buffered : Nat
  failed : Bool
  calls : Nat
  archiveCount : BitVec 64

/-- Go 1.27.0 writer.go:66 delegates to bufio.Writer.Flush. A lawful successful
    write consumes the entire visible buffer. The error flag is sticky. -/
def zipWriterFlush (count : BitVec 64) (buffered : Nat) (sticky : Bool)
    (reply : WriteReply) : FlushResult :=
  if sticky || buffered == 0 then ⟨buffered, sticky, 0, count⟩
  else if reply.failed then ⟨buffered - reply.count, true, 1, count⟩
  else ⟨0, false, 1, count⟩

structure NewWriterResult where
  capacity : Nat
  buffered : Nat
  reused : Bool
  archiveCount : Nat
  calls : Nat

/-- Go 1.27.0 writer.go:49, delegating construction to bufio.NewWriter. -/
def zipNewWriter (existing : Option (Nat × Nat)) : NewWriterResult :=
  match existing with
  | some (capacity, buffered) =>
    if 4096 ≤ capacity then ⟨capacity, buffered, true, 0, 0⟩ else ⟨4096, 0, false, 0, 0⟩
  | none => ⟨4096, 0, false, 0, 0⟩

theorem countWriterWrite_return (count : BitVec 64) (reply : WriteReply) :
    (countWriterWrite count reply).returned = reply ∧
    (countWriterWrite count reply).count = count + BitVec.ofNat 64 reply.count := ⟨rfl, rfl⟩

theorem nopCloserClose_noEffect (calls : Nat) :
    nopCloserClose calls = (false, calls) := rfl

theorem zipWriterFlush_sticky (count : BitVec 64) (buffered : Nat) (reply : WriteReply) :
    zipWriterFlush count buffered true reply = ⟨buffered, true, 0, count⟩ := rfl

theorem zipNewWriter_reuse : zipNewWriter (some (4096, 7)) = ⟨4096, 7, true, 0, 0⟩ := rfl

end GoLean.Archive.Zip
