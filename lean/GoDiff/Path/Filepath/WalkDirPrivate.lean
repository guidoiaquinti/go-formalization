import GoLean.Path.Filepath.WalkDir
import GoDiff.Math.Bits.Tuple

open GoLean.Path.Filepath

private def raw (name : String) : List (BitVec 8) :=
  name.toUTF8.data.toList.map UInt8.toBitVec

private def node (readFailure : Bool) : WalkDirNode :=
  .directory (raw "root") (if readFailure then some 7 else none)
    [.entry (raw "a"), .directory (raw "d") none [.entry (raw "b"), .entry (raw "c")],
     .entry (raw "z")]

private def response : Int → Option WalkDirError
  | 0 => none
  | 1 => some .skipDir
  | 2 => some .skipAll
  | _ => some (.other 9)

private def observe : Option WalkDirError → String
  | none => "ok"
  | some .skipDir => "skipDir"
  | some .skipAll => "skipAll"
  | some (.other code) => toString code

def main : IO Unit := runTuple fun fields => do
  let profile :: failure :: root :: first :: dir :: innerOne :: innerTwo :: last :: retry :: [] := fields |
    throw (IO.userError "expected profile, read failure and callback outcomes")
  let windows := profile == 1
  let separator := if windows then "\\" else "/"
  let table :=
    [("root", root), ("root" ++ separator ++ "a", first), ("root" ++ separator ++ "d", dir),
      (String.intercalate separator ["root", "d", "b"], innerOne),
      (String.intercalate separator ["root", "d", "c"], innerTwo),
      ("root" ++ separator ++ "z", last)].map fun entry => (raw entry.1, entry.2)
  let callback := fun event : WalkDirVisit =>
    if event.readError.isSome then response retry
    else response (((table.find? fun item => item.1 == event.path).map Prod.snd).getD 0)
  let (visits, result) := walkDirTree windows 2 (raw "root") (node (failure == 1)) callback
  let events := visits.map fun event =>
    let path := String.fromUTF8! (ByteArray.mk ((event.path.map fun byte => UInt8.ofNat byte.toNat).toArray))
    path ++ (if event.directory then ":d" else ":f") ++ (if event.readError.isSome then ":7" else "")
  return observe result ++ "|" ++ String.intercalate "," events
