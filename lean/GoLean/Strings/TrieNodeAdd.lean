import GoLean.Strings.MakeGenericReplacer

namespace GoLean.Strings

private def commonBytes : List (BitVec 8) → List (BitVec 8) → Nat
  | a :: as, b :: bs => if a == b then 1 + commonBytes as bs else 0
  | _, _ => 0

private def installTerminal (node : GenericTrie) (value : Nat × List (BitVec 8)) : GenericTrie :=
  match node with
  | .leaf terminal => .leaf (terminal <|> some value)
  | .edge terminal prefixBytes child => .edge (terminal <|> some value) prefixBytes child
  | .table terminal children => .table (terminal <|> some value) children

private def updateChild (byte : BitVec 8) (child : GenericTrie) :
    List (BitVec 8 × GenericTrie) → List (BitVec 8 × GenericTrie)
  | [] => [(byte, child)]
  | (key, old) :: rest =>
    if byte == key then (byte, child) :: rest
    else if byte.toNat < key.toNat then (byte, child) :: (key, old) :: rest
    else (key, old) :: updateChild byte child rest

private def insertTrieNode (value : Nat × List (BitVec 8)) :
    Nat → GenericTrie → List (BitVec 8) → GenericTrie
  | _, node, [] => installTerminal node value
  | 0, node, _ :: _ => node
  | fuel + 1, node, key@(_ :: _) =>
    let byte := key.head!
    let rest := key.tail
    match node with
    | .leaf terminal => .edge terminal key (.leaf (some value))
    | .table terminal children =>
      let child := ((children.find? fun (b, _) => b == byte).map Prod.snd).getD (.leaf none)
      .table terminal (updateChild byte (insertTrieNode value fuel child rest) children)
    | .edge terminal prefixBytes child =>
      let common := commonBytes prefixBytes key
      if common == prefixBytes.length then
        .edge terminal prefixBytes (insertTrieNode value fuel child (key.drop common))
      else if common == 0 then
        let oldByte := prefixBytes.head!
        let oldRest := prefixBytes.tail
        let oldChild := if oldRest.isEmpty then child else .edge none oldRest child
        let newChild := insertTrieNode value fuel (.leaf none) rest
        .table terminal (updateChild byte newChild [(oldByte, oldChild)])
      else
        let next : GenericTrie := .edge none (prefixBytes.drop common) child
        .edge terminal (prefixBytes.take common)
          (insertTrieNode value fuel next (key.drop common))

/-- Go 1.27.0 src/strings/replace.go:158–218. Insert into a well-formed
    compressed trie with the key bytes present in the replacer's byte map.
    Branch entries use original bytes in ascending order. -/
def trieNodeAdd (node : GenericTrie) (key replacement : List (BitVec 8)) (priority : Nat) :
    GenericTrie :=
  insertTrieNode (priority, replacement) key.length node key

theorem trieNodeAdd_emptyExisting (priority : Nat) (value replacement : List (BitVec 8)) :
    trieNodeAdd (.leaf (some (priority, value))) [] replacement (priority + 1) =
      .leaf (some (priority, value)) := rfl

theorem trieNodeAdd_leaf (head : BitVec 8) (tail replacement : List (BitVec 8)) (priority : Nat) :
    trieNodeAdd (.leaf none) (head :: tail) replacement priority =
      .edge none (head :: tail) (.leaf (some (priority, replacement))) := rfl

theorem trieNodeAdd_split (value replacement : List (BitVec 8)) :
    trieNodeAdd (.edge none [1, 2] (.leaf (some (2, value)))) [1, 3] replacement 1 =
      .edge none [1] (.table none [(2, .leaf (some (2, value))), (3, .leaf (some (1, replacement)))]) :=
  rfl

theorem trieNodeAdd_splitTerminal (value replacement : List (BitVec 8)) :
    trieNodeAdd (.edge none [1, 2] (.leaf (some (2, value)))) [1] replacement 1 =
      .edge none [1] (.edge (some (1, replacement)) [2] (.leaf (some (2, value)))) := rfl

end GoLean.Strings
