import Init.Data.BitVec.Basic

namespace GoLean.Strings

private structure RankedKey where
  suffix : List (BitVec 8)
  value : List (BitVec 8)
  priority : Nat

abbrev GenericTerminal := Option (Nat × List (BitVec 8))

inductive GenericTrie where
  | leaf (terminal : GenericTerminal)
  | edge (terminal : GenericTerminal) (prefixBytes : List (BitVec 8)) (child : GenericTrie)
  | table (terminal : GenericTerminal) (children : List (BitVec 8 × GenericTrie))

structure GenericReplacer where
  mapping : List (BitVec 8)
  tableSize : Nat
  root : GenericTrie

private def insertByte (byte : BitVec 8) : List (BitVec 8) → List (BitVec 8)
  | [] => [byte]
  | first :: rest =>
    if byte == first then first :: rest
    else if byte.toNat < first.toNat then byte :: first :: rest
    else first :: insertByte byte rest

private def sortPresent (bytes : List (BitVec 8)) : List (BitVec 8) :=
  bytes.foldl (fun alphabet byte => insertByte byte alphabet) []


private def terminalAt (keys : List RankedKey) : GenericTerminal :=
  (keys.find? fun key => key.suffix.isEmpty).map fun key => (key.priority, key.value)

private def descend (byte : BitVec 8) (keys : List RankedKey) : List RankedKey :=
  keys.filterMap fun key =>
    match key.suffix with
    | head :: tail => if head == byte then some { key with suffix := tail } else none
    | [] => none

private def buildGenericTrie (root : Bool) (keys : List RankedKey) : Nat → GenericTrie
  | 0 =>
    if root then .table (terminalAt keys) [] else .leaf (terminalAt keys)
  | fuel + 1 =>
    let terminal := terminalAt keys
    let firstBytes := sortPresent (keys.filterMap fun key => key.suffix.head?)
    let children := firstBytes.map fun byte => (byte, buildGenericTrie false (descend byte keys) fuel)
    if root then .table terminal children
    else match children with
      | [] => .leaf terminal
      | [(byte, .edge none prefixBytes child)] => .edge terminal (byte :: prefixBytes) child
      | [(byte, child)] => .edge terminal [byte] child
      | _ => .table terminal children

/-- Go 1.27.0 src/strings/replace.go:268–298. Ascending key-byte
    order gives the dense mapping. Form the trie by grouping keys
    by next byte and compressing stretches with no terminal or branch. -/
def makeGenericReplacer (pairs : List (List (BitVec 8) × List (BitVec 8))) :
    GenericReplacer :=
  let data := pairs.zipIdx.map fun ((key, value), index) =>
    ⟨key, value, 2 * (pairs.length - index)⟩
  let alphabet := sortPresent (pairs.flatMap Prod.fst)
  let mapping := (List.range 256).map fun byte =>
    BitVec.ofNat 8 (alphabet.idxOf (BitVec.ofNat 8 byte))
  let depth := pairs.foldl (fun n pair => max n pair.1.length) 0
  ⟨mapping, alphabet.length, buildGenericTrie true data depth⟩

theorem makeGenericReplacer_emptyTable :
    (makeGenericReplacer []).tableSize = 0 := rfl

theorem makeGenericReplacer_emptyDuplicate (first second : List (BitVec 8)) :
    (makeGenericReplacer [([], first), ([], second)]).root = .table (some (4, first)) [] := rfl

theorem makeGenericReplacer_compressed (value : List (BitVec 8)) :
    (makeGenericReplacer [([97, 98, 99], value)]).root =
      .table none [(97, .edge none [98, 99] (.leaf (some (2, value))))] := rfl

set_option maxRecDepth 4096 in
theorem makeGenericReplacer_mapping :
    ((makeGenericReplacer [([255, 0], [])]).mapping[0]?,
      (makeGenericReplacer [([255, 0], [])]).mapping[1]?,
      (makeGenericReplacer [([255, 0], [])]).mapping[255]?) =
      (some 0, some 2, some 1) := by decide

end GoLean.Strings
