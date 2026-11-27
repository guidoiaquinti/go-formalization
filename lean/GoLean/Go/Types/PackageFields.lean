import GoLean.Go.Types.ArrayChan

namespace GoLean.Go.Types

abbrev PackageBytes := List Nat

structure PackageState where
  path : PackageBytes
  name : PackageBytes
  scope : Option Nat
  complete : Bool
  version : PackageBytes

/-- Go 1.27.0 src/go/types/package.go:28–31. The scope token abstracts
    a fresh, empty child of Universe. -/
def newPackage (path name : PackageBytes) (scope : Nat) : PackageState :=
  ⟨path, name, some scope, false, []⟩

private def packageRead (field : PackageState → α) (pkg : Option PackageState) : Except ReadFault α :=
  match pkg with
  | none => .error .nilReceiver
  | some p => .ok (field p)

def packagePath := packageRead PackageState.path
def packageName := packageRead PackageState.name
def packageGoVersion := packageRead PackageState.version
def packageComplete := packageRead PackageState.complete

def packageScope (pkg : Option PackageState) (universeID : Nat) : Option Nat :=
  match pkg with
  | none => some universeID
  | some p => p.scope

def packageSetName (pkg : Option PackageState) (name : PackageBytes) : Except ReadFault PackageState :=
  packageRead (fun p => { p with name }) pkg

def packageMarkComplete (pkg : Option PackageState) : Except ReadFault PackageState :=
  packageRead (fun p => { p with complete := true }) pkg

theorem newPackage_initial (path name : PackageBytes) (scope : Nat) :
    (newPackage path name scope).complete = false ∧
      (newPackage path name scope).version = [] ∧
      (newPackage path name scope).scope = some scope := ⟨rfl, rfl, rfl⟩

theorem packageScope_identity (p : PackageState) (universeID : Nat) :
    packageScope (some p) universeID = p.scope ∧ packageScope none universeID = some universeID :=
  ⟨rfl, rfl⟩

theorem packageSetName_only_name (p : PackageState) (name : PackageBytes) :
    packageSetName (some p) name = .ok { p with name } ∧
      packageSetName none name = .error .nilReceiver := ⟨rfl, rfl⟩

theorem packageMarkComplete_only_complete (p : PackageState) :
    packageMarkComplete (some p) = .ok { p with complete := true } ∧
      packageMarkComplete none = .error .nilReceiver := ⟨rfl, rfl⟩

theorem packagePath_reads (p : PackageState) :
    packagePath (some p) = .ok p.path ∧ packagePath none = .error .nilReceiver :=
  ⟨rfl, rfl⟩

theorem packageName_reads (p : PackageState) :
    packageName (some p) = .ok p.name ∧ packageName none = .error .nilReceiver :=
  ⟨rfl, rfl⟩

theorem packageGoVersion_reads (p : PackageState) :
    packageGoVersion (some p) = .ok p.version ∧ packageGoVersion none = .error .nilReceiver :=
  ⟨rfl, rfl⟩

theorem packageComplete_reads (p : PackageState) :
    packageComplete (some p) = .ok p.complete ∧ packageComplete none = .error .nilReceiver :=
  ⟨rfl, rfl⟩

end GoLean.Go.Types
