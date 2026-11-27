import GoLean.Strings.Compare

namespace GoLean.Archive.Zip

abbrev EntryName := List (BitVec 8)

structure EntrySplit where
  directory : EntryName
  element : EntryName
  trailingSlash : Bool

/-- Go 1.27.0 reader.go:893. Strip one final slash, then split at the last slash. -/
def splitEntryName (name : EntryName) : EntrySplit :=
  let (reversed, slash) :=
    match name.reverse with
    | 47 :: tail => (tail, true)
    | rest => (rest, false)
  let elem := reversed.takeWhile (· != 47)
  match reversed.drop elem.length with
  | [] => ⟨[46], elem.reverse, slash⟩
  | _ :: dir => ⟨dir.reverse, elem.reverse, slash⟩

/-- Go 1.27.0 reader.go:860. The trailing-slash flag is not a sort key. -/
def fileEntryCompare (x y : EntryName) : Int :=
  let x := splitEntryName x
  let y := splitEntryName y
  if x.directory == y.directory then GoLean.Strings.compare x.element y.element
  else GoLean.Strings.compare x.directory y.directory

structure ListEntryTime where
  seconds : Int
  nanoseconds : Nat

/-- Identity and the two fields observed by the directory accessors. -/
structure ListEntry where
  identity : Nat
  name : EntryName
  modified : Option ListEntryTime

/-- Directory accessors in Go 1.27.0 reader.go:748–766. -/
def listEntryName (e : ListEntry) : EntryName := (splitEntryName e.name).element
def listEntrySize (_ : ListEntry) : Int := 0
def listEntryMode (_ : ListEntry) : BitVec 32 := 0x8000016d
def listEntryType (_ : ListEntry) : BitVec 32 := 0x80000000
def listEntryIsDir (_ : ListEntry) : Bool := true
def listEntrySys (_ : ListEntry) : Option Nat := none

/-- An absent file contributes Go's zero instant, January 1 of year 1 UTC. -/
def listEntryModTime (e : ListEntry) : ListEntryTime :=
  e.modified.getD ⟨-62135596800, 0⟩

def listEntryInfo (e : ListEntry) : Nat × Option String := (e.identity, none)

def listEntryString (e : ListEntry) : EntryName :=
  [100, 32] ++ listEntryName e ++ [47]

theorem splitEntryName_doubleSlash :
    (splitEntryName [97, 47, 47]).directory = [97] ∧
    (splitEntryName [97, 47, 47]).element = [] ∧
    (splitEntryName [97, 47, 47]).trailingSlash = true := ⟨rfl, rfl, rfl⟩

theorem fileEntryCompare_directoryFirst :
    fileEntryCompare [97, 47, 122] [128] = 1 := rfl

theorem listEntryName_slash (identity : Nat) (modified : Option ListEntryTime) :
    listEntryName ⟨identity, [47], modified⟩ = [] := rfl

theorem listEntrySize_zero (e : ListEntry) : listEntrySize e = 0 := rfl
theorem listEntryMode_permissions (e : ListEntry) : listEntryMode e = 0x8000016d := rfl
theorem listEntryType_directory (e : ListEntry) : listEntryType e = 0x80000000 := rfl
theorem listEntryIsDir_true (e : ListEntry) : listEntryIsDir e = true := rfl
theorem listEntrySys_none (e : ListEntry) : listEntrySys e = none := rfl

theorem listEntryModTime_absent (identity : Nat) (name : EntryName) :
    listEntryModTime ⟨identity, name, none⟩ = ⟨-62135596800, 0⟩ := rfl

theorem listEntryInfo_identity (e : ListEntry) :
    listEntryInfo e = (e.identity, none) := rfl

theorem listEntryString_slash (identity : Nat) (modified : Option ListEntryTime) :
    listEntryString ⟨identity, [47], modified⟩ = [100, 32, 47] := rfl

end GoLean.Archive.Zip
