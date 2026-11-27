namespace GoLean.Archive.Zip

structure StatEntry where
  identity : Nat
  name : List UInt8
  duplicate : Bool
  directory : Bool
  header : Option Nat

/-- Return observations classify interface type and identity without modeling
    Go's allocation of the error or its raw runtime panic payload. -/
inductive StatResult where
  | header (identity : Nat)
  | directory (identity : Nat)
  | duplicate (message : List UInt8)
  | nilAccess

private def duplicateSuffix : List UInt8 :=
  ": duplicate entries in zip file".toUTF8.data.toList

/-- Go 1.27.0 reader.go:737–746. -/
def listEntryStat (entry : StatEntry) : StatResult :=
  if entry.duplicate then .duplicate (entry.name ++ duplicateSuffix)
  else if entry.directory then .directory entry.identity
  else match entry.header with
    | some identity => .header identity
    | none => .nilAccess

/-- Go 1.27.0 reader.go:955. -/
def openDirStat (entry : Option StatEntry) : StatResult :=
  entry.map listEntryStat |>.getD .nilAccess

/-- Go 1.27.0 reader.go:290–292. -/
def checksumReaderStat (header : Option Nat) : StatResult :=
  match header with
  | some identity => .header identity
  | none => .nilAccess

theorem listEntryStat_duplicate (entry : StatEntry) (h : entry.duplicate = true) :
    listEntryStat entry = .duplicate (entry.name ++ duplicateSuffix) := by
  simp [listEntryStat, h]

theorem listEntryStat_missingFile (identity : Nat) (name : List UInt8) :
    listEntryStat ⟨identity, name, false, false, none⟩ = .nilAccess := rfl

theorem openDirStat_missingEntry : openDirStat none = .nilAccess := rfl

theorem openDirStat_delegates (entry : StatEntry) :
    openDirStat (some entry) = listEntryStat entry := rfl

theorem checksumReaderStat_header (identity : Nat) :
    checksumReaderStat (some identity) = .header identity := rfl

end GoLean.Archive.Zip
