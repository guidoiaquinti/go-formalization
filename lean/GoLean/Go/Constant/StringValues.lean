namespace GoLean.Go.Constant

abbrev StringBytes := List Nat

inductive StringTree where
  | leaf (value : StringBytes)
  | concat (left right : StringTree)

inductive StringOperand where
  | string (tree : Option StringTree)
  | unknown
  | other (printed : String)

inductive StringFault where
  | nilReceiver
  | wrongType (message : String)

structure StringCreation where
  tree : StringTree
  usesEmptySingleton : Bool

/-- Go 1.27.0 src/go/constant/value.go:387–392. -/
def makeString (bytes : StringBytes) : StringCreation :=
  ⟨.leaf bytes, bytes.isEmpty⟩

def stringTreeKind (_ : Option StringTree) : Nat := 2

private def length : StringTree → Nat
  | .leaf value => value.length
  | .concat left right => length left + length right

private def pieces : StringTree → List StringBytes → List StringBytes
  | .leaf value, suffix => value :: suffix
  | .concat left right, suffix => pieces right (pieces left suffix)

def stringTreeLen (receiver : Option StringTree) : Except StringFault Nat :=
  match receiver with
  | none => .error .nilReceiver
  | some tree => .ok (length tree)

/-- The root is replaced by its flattened leaf. The model abstracts locks and
    uses immutable subtrees, which retain their prior values. -/
def stringTreeString (receiver : Option StringTree) : Except StringFault (StringBytes × StringTree) :=
  match receiver with
  | none => .error .nilReceiver
  | some (.leaf value) => .ok (value, .leaf value)
  | some tree@(.concat _ _) =>
    let value := (pieces tree []).reverse.flatten
    .ok (value, .leaf value)

def stringValue : StringOperand → Except StringFault (StringBytes × Option StringTree)
  | .unknown => .ok ([], none)
  | .other text => .error (.wrongType (text ++ " not a String"))
  | .string tree => (stringTreeString tree).map (fun (value, state) => (value, some state))

def stringLength : StringOperand → Except StringFault Nat
  | .unknown => .ok 0
  | .other text => .error (.wrongType (text ++ " not a String"))
  | .string tree => stringTreeLen tree

theorem makeString_leaf (value : StringBytes) :
    (makeString value).tree = .leaf value ∧
      (makeString value).usesEmptySingleton = value.isEmpty := ⟨rfl, rfl⟩

theorem stringTreeKind_string (receiver : Option StringTree) :
    stringTreeKind receiver = 2 := rfl

theorem stringTreeLen_children (left right : StringTree) :
    stringTreeLen (some (.concat left right)) = .ok (length left + length right) := rfl

theorem stringTreeString_twoLeaves (left right : StringBytes) :
    stringTreeString (some (.concat (.leaf left) (.leaf right))) =
      .ok (left ++ right, .leaf (left ++ right)) := by
  simp [stringTreeString, pieces]

theorem stringValue_unknown : stringValue .unknown = .ok ([], none) := rfl

theorem stringValue_wrongType (printed : String) :
    stringValue (.other printed) = .error (.wrongType (printed ++ " not a String")) := rfl

theorem stringLength_unknown : stringLength .unknown = .ok 0 := rfl

theorem stringLength_wrongType (printed : String) :
    stringLength (.other printed) = .error (.wrongType (printed ++ " not a String")) := rfl

end GoLean.Go.Constant
