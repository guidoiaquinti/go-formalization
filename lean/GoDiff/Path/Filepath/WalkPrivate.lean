import GoLean.Path.Filepath.Walk
import GoDiff.Math.Bits.Tuple

open GoLean.Path.Filepath

private def raw (name : String) : List (BitVec 8) :=
  name.toUTF8.data.toList.map UInt8.toBitVec

private def node (failure : Int) : WalkNode :=
  .directory (raw "root") (if failure == 1 then .error 7 else .ok
    [if failure == 2 then .statFailed (raw "a") 8 else .file (raw "a"),
      .directory (raw "d") (.ok [.file (raw "b"), .file (raw "c")]), .file (raw "z")])

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
  let callback := fun event : WalkDirPublicVisit =>
    if event.readError.isSome then response retry
    else response (((table.find? fun item => item.1 == event.path).map Prod.snd).getD 0)
  let (visits, result) := walkTree windows 2 (raw "root") (node failure) callback
  let events := visits.map fun event =>
    let path := String.fromUTF8! (ByteArray.mk ((event.path.map fun byte => UInt8.ofNat byte.toNat).toArray))
    let kind := match event.entryIsDirectory with
      | none => ":nil"
      | some true => ":d"
      | some false => ":f"
    path ++ kind ++ ((event.readError.map fun code => ":" ++ toString code).getD "")
  return observe result ++ "|" ++ String.intercalate "," events
