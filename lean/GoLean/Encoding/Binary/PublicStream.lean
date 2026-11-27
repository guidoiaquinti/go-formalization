import GoLean.Encoding.Binary.PublicBuffer

namespace GoLean.Encoding.Binary

structure StreamReadResult where
  status : String
  state : List Nat
  requests : List Nat

private def encodedWidth (kind : FastKind) (form : FastForm) (words : List Nat) : Nat :=
  if form = .slice then words.length * fastWidth kind else fastWidth kind

private def readFully (chunk : Nat) (terminal : String) :
    Nat → Nat → Nat → List Nat × String
  | 0, _, _ => ([], "other")
  | fuel + 1, needed, available =>
    if needed = 0 then ([], "ok")
    else
      let used := min (min needed chunk) available
      if used = needed then ([needed], "ok")
      else if used = available then ([needed], terminal)
      else
        let (rest, result) := readFully chunk terminal fuel (needed - used) (available - used)
        (needed :: rest, result)

def binaryRead (big namedByte : Bool) (kind : FastKind) (form : FastForm)
    (words : List Nat) (source : List (BitVec 8)) (chunk : Nat) (terminal : String) :
    StreamReadResult :=
  if form = .unsupported ∨ (namedByte ∧ form = .scalar) then ⟨"type", words, []⟩
  else
    let n := encodedWidth kind form words
    let (requests, err) := readFully (max 1 chunk) terminal (n + 1) n source.length
    if err != "ok" then
      ⟨if err == "eof" && source.length != 0 then "unexpected" else err, words, requests⟩
    else
      let decoded := binaryDecode big namedByte kind form (source.take n) words
      ⟨decoded.status, decoded.state, requests⟩

structure StreamWriteResult where
  status : String
  output : List (BitVec 8)
  calls : Nat
  aliases : Bool

def binaryWrite (big namedByte : Bool) (kind : FastKind) (form : FastForm)
    (words : List Nat) (response : String) : StreamWriteResult :=
  if form = .unsupported then ⟨"type", [], 0, false⟩
  else
    let n := encodedWidth kind form words
    let encoded := binaryEncode big kind form words (List.replicate n 0)
    ⟨response, encoded.state, 1, !namedByte && kind == .uint8 && form == .slice && n != 0⟩

theorem binaryRead_bad_form (big namedByte : Bool) (kind : FastKind) (words : List Nat)
    (source : List (BitVec 8)) (chunk : Nat) (terminal : String) :
    (binaryRead big namedByte kind .unsupported words source chunk terminal).requests = [] := by
  simp [binaryRead]

theorem binaryWrite_bad_form (big namedByte : Bool) (kind : FastKind) (words : List Nat)
    (response : String) :
    (binaryWrite big namedByte kind .unsupported words response).calls = 0 :=
  rfl

end GoLean.Encoding.Binary
