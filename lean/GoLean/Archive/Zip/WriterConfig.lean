import Init.Data.BitVec.Lemmas

namespace GoLean.Archive.Zip

abbrev Comment := List (BitVec 8)

structure ZipWriterConfig where
  comment : Comment
  count : Option Int

inductive ZipConfigFault where
  | nilAccess
  | tooLong
  | afterWrite
  deriving DecidableEq

def writerSetComment (writer : Option ZipWriterConfig) (comment : Comment) :
    Except ZipConfigFault ZipWriterConfig :=
  if comment.length > 65535 then .error .tooLong
  else
    match writer with
    | none => .error .nilAccess
    | some w => .ok {w with comment := comment}

def writerSetOffset (writer : Option ZipWriterConfig) (offset : Int) :
    Except ZipConfigFault ZipWriterConfig :=
  match writer with
  | some w =>
    match w.count with
    | some 0 => .ok {w with count := some offset}
    | some _ => .error .afterWrite
    | none => .error .nilAccess
  | none => .error .nilAccess

theorem writerSetComment_tooLong (writer : Option ZipWriterConfig) (comment : Comment)
    (long : comment.length > 65535) :
    writerSetComment writer comment = .error .tooLong := by simp [writerSetComment, long]

theorem writerSetOffset_zero (comment : Comment) (n : Int) :
    writerSetOffset (some ⟨comment, some 0⟩) n = .ok ⟨comment, some n⟩ := rfl

end GoLean.Archive.Zip
