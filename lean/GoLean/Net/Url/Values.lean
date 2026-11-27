import GoLean.Net.Url.Escape
import Init.Data.List.Sort

namespace GoLean.Net.Url

abbrev URLBytes := List (BitVec 8)
abbrev ValuesEntry := URLBytes × Option (List URLBytes)
abbrev Values := Option (List ValuesEntry)

private def valuesFind (values : Values) (key : URLBytes) : Option ValuesEntry :=
  (values.getD []).find? (fun e => e.1 == key)

private def removeKey (entries : List ValuesEntry) (key : URLBytes) : List ValuesEntry :=
  entries.filter (fun e => e.1 != key)

/-- Go 1.27.0 src/net/url/url.go:898–904. Missing and empty lists read empty. -/
def valuesGet (values : Values) (key : URLBytes) : URLBytes :=
  ((valuesFind values key).bind Prod.snd |>.getD []).headD []

/-- Go 1.27.0 src/net/url/url.go:924–927. Membership is independent of list length. -/
def valuesHas (values : Values) (key : URLBytes) : Bool :=
  (valuesFind values key).isSome

/-- Go 1.27.0 src/net/url/url.go:908–910. None is a nil-map assignment panic. -/
def valuesSet (values : Values) (key value : URLBytes) : Values :=
  values.map fun entries => (key, some [value]) :: removeKey entries key

/-- Go 1.27.0 src/net/url/url.go:914–916. Append to the sequence for the current key. -/
def valuesAdd (values : Values) (key value : URLBytes) : Values :=
  let old := ((valuesFind values key).bind Prod.snd).getD []
  values.map fun entries => (key, some (old ++ [value])) :: removeKey entries key

/-- Go 1.27.0 src/net/url/url.go:919–921. A nil map remains nil. -/
def valuesDel (values : Values) (key : URLBytes) : Values :=
  values.map (fun entries => removeKey entries key)

/-- Go 1.27.0 src/net/url/url.go:930–939. Copy the map and each list of
    strings. Allocation identity is observed on native Go. -/
def valuesClone (values : Values) : Values :=
  values.map fun entries =>
    entries.map fun (key, items) => (key, items.map (List.map id))

def sortValuesEntries (entries : List ValuesEntry) : List ValuesEntry :=
  entries.mergeSort fun x y => !y.1.lex x.1 (fun a b => decide (a.toNat < b.toNat))

/-- Go 1.27.0 src/net/url/url.go:1015–1043. Sort unescaped keys, preserve
    each list's order and insert ampersands between emitted key-value pairs. -/
def valuesEncode (values : Values) : URLBytes :=
  let entries := sortValuesEntries (values.getD [])
  let pairs := entries.flatMap fun (key, words) =>
    let escaped := queryEscape key
    (words.getD []).map fun word => escaped ++ [61] ++ queryEscape word
  (List.intersperse [38] pairs).flatten

theorem valuesGet_nil (key : URLBytes) : valuesGet none key = [] := rfl
theorem valuesHas_nil (key : URLBytes) : valuesHas none key = false := rfl
theorem valuesHas_nilEntry (key : URLBytes) :
    valuesHas (some [(key, none)]) key = true := by simp [valuesHas, valuesFind]

theorem valuesSet_nil (key value : URLBytes) : valuesSet none key value = none := rfl
theorem valuesSet_overwrite (v : URLBytes) :
    valuesSet (some [([65], some [[66], [67]])]) [65] v = some [([65], some [v])] := rfl

theorem valuesAdd_nil (key value : URLBytes) : valuesAdd none key value = none := rfl
theorem valuesAdd_order :
    valuesAdd (some [([65], some [[66], [67]])]) [65] [68] =
      some [([65], some [[66], [67], [68]])] := rfl

theorem valuesDel_nil (key : URLBytes) : valuesDel none key = none := rfl
theorem valuesDel_single (key : URLBytes) (words : Option (List URLBytes)) :
    valuesDel (some [(key, words)]) key = some [] := by simp [valuesDel, removeKey]

theorem valuesClone_sameData (v : Values) : valuesClone v = v := by
  cases v with
  | none => rfl
  | some entries =>
    induction entries with
    | nil => rfl
    | cons e es ih =>
      obtain ⟨key, words⟩ := e
      cases words <;> simp [valuesClone]

theorem valuesEncode_nil : valuesEncode none = [] := by simp [valuesEncode, sortValuesEntries]

theorem valuesEncode_keyOrder :
    valuesEncode (some [([255], some [[32]]), ([65], some [[], [43]]), ([0], none)]) =
      [65, 61, 38, 65, 61, 37, 50, 66, 38, 37, 70, 70, 61, 43] := by
  simp [valuesEncode, sortValuesEntries, List.mergeSort] <;> rfl

end GoLean.Net.Url
