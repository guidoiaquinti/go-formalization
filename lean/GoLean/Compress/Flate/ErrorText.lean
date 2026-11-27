import Init.Data.BitVec.Lemmas

namespace GoLean.Compress.Flate

abbrev ErrorBytes := List (BitVec 8)

private def text (s : String) : ErrorBytes := s.toUTF8.data.toList.map UInt8.toBitVec

def corruptInputError (offset : BitVec 64) : ErrorBytes :=
  text ("flate: corrupt input before offset " ++ toString offset.toInt)

def internalError (message : ErrorBytes) : ErrorBytes :=
  text "flate: internal error: " ++ message

inductive ErrorCall where
  | absent
  | returns (message : ErrorBytes)
  | panics
  deriving DecidableEq

inductive ErrorResult where
  | nilAccess
  | panicked (calls : Nat)
  | returned (calls : Nat) (message : ErrorBytes)
  deriving DecidableEq

private def ioError (kind : String) (offset : Option (BitVec 64)) (callback : ErrorCall) :
    ErrorResult :=
  match offset, callback with
  | none, _ | _, .absent => .nilAccess
  | some _, .panics => .panicked 1
  | some word, .returns content =>
    .returned 1 (text ("flate: " ++ kind ++ " error at offset " ++ toString word.toInt ++ ": ") ++ content)

def readError := ioError "read"

def writeError := ioError "write"

theorem corruptInputError_sign :
    corruptInputError (BitVec.ofInt 64 (-1)) = text "flate: corrupt input before offset -1" := rfl

theorem internalError_raw (message : ErrorBytes) :
    internalError message = text "flate: internal error: " ++ message := rfl

theorem readError_callback (offset : BitVec 64) :
    readError (some offset) .panics = .panicked 1 := rfl

theorem writeError_callback (offset : BitVec 64) (bytes : ErrorBytes) :
    writeError (some offset) (.returns bytes) =
      .returned 1 (text ("flate: write error at offset " ++ toString offset.toInt ++ ": ") ++ bytes) := rfl

end GoLean.Compress.Flate
