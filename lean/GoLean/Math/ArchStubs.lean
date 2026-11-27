import Init.Data.BitVec.Basic

namespace GoLean.Math

inductive ArchStubProfile where
  | nonS390x
  | s390x

structure ArchStubOutcome where
  available : Bool
  panicMessage : String

private def archStubOutcome : ArchStubProfile → Option ArchStubOutcome
  | .nonS390x => some ⟨false, "not implemented"⟩
  | .s390x => none

/-- Go 1.27.0 src/math/stubs.go:14–16, selected by the !s390x build constraint. -/
def archAcosStub (profile : ArchStubProfile) (_ : BitVec 64) : Option ArchStubOutcome :=
  archStubOutcome profile

theorem archAcosStub_profiles (x : BitVec 64) :
    archAcosStub .nonS390x x = some ⟨false, "not implemented"⟩ ∧
    archAcosStub .s390x x = none := ⟨rfl, rfl⟩

/-- Go 1.27.0 src/math/stubs.go:20–22, selected by the !s390x build constraint. -/
def archAcoshStub (profile : ArchStubProfile) (_ : BitVec 64) : Option ArchStubOutcome :=
  archStubOutcome profile

theorem archAcoshStub_profiles (x : BitVec 64) :
    archAcoshStub .nonS390x x = some ⟨false, "not implemented"⟩ ∧
    archAcoshStub .s390x x = none := ⟨rfl, rfl⟩

/-- Go 1.27.0 src/math/stubs.go:26–28, selected by the !s390x build constraint. -/
def archAsinStub (profile : ArchStubProfile) (_ : BitVec 64) : Option ArchStubOutcome :=
  archStubOutcome profile

theorem archAsinStub_profiles (x : BitVec 64) :
    archAsinStub .nonS390x x = some ⟨false, "not implemented"⟩ ∧
    archAsinStub .s390x x = none := ⟨rfl, rfl⟩

/-- Go 1.27.0 src/math/stubs.go:32–34, selected by the !s390x build constraint. -/
def archAsinhStub (profile : ArchStubProfile) (_ : BitVec 64) : Option ArchStubOutcome :=
  archStubOutcome profile

theorem archAsinhStub_profiles (x : BitVec 64) :
    archAsinhStub .nonS390x x = some ⟨false, "not implemented"⟩ ∧
    archAsinhStub .s390x x = none := ⟨rfl, rfl⟩

/-- Go 1.27.0 src/math/stubs.go:38–40, selected by the !s390x build constraint. -/
def archAtanStub (profile : ArchStubProfile) (_ : BitVec 64) : Option ArchStubOutcome :=
  archStubOutcome profile

theorem archAtanStub_profiles (x : BitVec 64) :
    archAtanStub .nonS390x x = some ⟨false, "not implemented"⟩ ∧
    archAtanStub .s390x x = none := ⟨rfl, rfl⟩

/-- Go 1.27.0 src/math/stubs.go:44–46, selected by the !s390x build constraint. -/
def archAtan2Stub (profile : ArchStubProfile) (_ : BitVec 64) (_ : BitVec 64) : Option ArchStubOutcome :=
  archStubOutcome profile

theorem archAtan2Stub_profiles (x : BitVec 64) (y : BitVec 64) :
    archAtan2Stub .nonS390x x y = some ⟨false, "not implemented"⟩ ∧
    archAtan2Stub .s390x x y = none := ⟨rfl, rfl⟩

/-- Go 1.27.0 src/math/stubs.go:50–52, selected by the !s390x build constraint. -/
def archAtanhStub (profile : ArchStubProfile) (_ : BitVec 64) : Option ArchStubOutcome :=
  archStubOutcome profile

theorem archAtanhStub_profiles (x : BitVec 64) :
    archAtanhStub .nonS390x x = some ⟨false, "not implemented"⟩ ∧
    archAtanhStub .s390x x = none := ⟨rfl, rfl⟩

/-- Go 1.27.0 src/math/stubs.go:56–58, selected by the !s390x build constraint. -/
def archCbrtStub (profile : ArchStubProfile) (_ : BitVec 64) : Option ArchStubOutcome :=
  archStubOutcome profile

theorem archCbrtStub_profiles (x : BitVec 64) :
    archCbrtStub .nonS390x x = some ⟨false, "not implemented"⟩ ∧
    archCbrtStub .s390x x = none := ⟨rfl, rfl⟩

/-- Go 1.27.0 src/math/stubs.go:62–64, selected by the !s390x build constraint. -/
def archCosStub (profile : ArchStubProfile) (_ : BitVec 64) : Option ArchStubOutcome :=
  archStubOutcome profile

theorem archCosStub_profiles (x : BitVec 64) :
    archCosStub .nonS390x x = some ⟨false, "not implemented"⟩ ∧
    archCosStub .s390x x = none := ⟨rfl, rfl⟩

/-- Go 1.27.0 src/math/stubs.go:68–70, selected by the !s390x build constraint. -/
def archCoshStub (profile : ArchStubProfile) (_ : BitVec 64) : Option ArchStubOutcome :=
  archStubOutcome profile

theorem archCoshStub_profiles (x : BitVec 64) :
    archCoshStub .nonS390x x = some ⟨false, "not implemented"⟩ ∧
    archCoshStub .s390x x = none := ⟨rfl, rfl⟩

/-- Go 1.27.0 src/math/stubs.go:74–76, selected by the !s390x build constraint. -/
def archErfStub (profile : ArchStubProfile) (_ : BitVec 64) : Option ArchStubOutcome :=
  archStubOutcome profile

theorem archErfStub_profiles (x : BitVec 64) :
    archErfStub .nonS390x x = some ⟨false, "not implemented"⟩ ∧
    archErfStub .s390x x = none := ⟨rfl, rfl⟩

/-- Go 1.27.0 src/math/stubs.go:80–82, selected by the !s390x build constraint. -/
def archErfcStub (profile : ArchStubProfile) (_ : BitVec 64) : Option ArchStubOutcome :=
  archStubOutcome profile

theorem archErfcStub_profiles (x : BitVec 64) :
    archErfcStub .nonS390x x = some ⟨false, "not implemented"⟩ ∧
    archErfcStub .s390x x = none := ⟨rfl, rfl⟩

/-- Go 1.27.0 src/math/stubs.go:86–88, selected by the !s390x build constraint. -/
def archExpm1Stub (profile : ArchStubProfile) (_ : BitVec 64) : Option ArchStubOutcome :=
  archStubOutcome profile

theorem archExpm1Stub_profiles (x : BitVec 64) :
    archExpm1Stub .nonS390x x = some ⟨false, "not implemented"⟩ ∧
    archExpm1Stub .s390x x = none := ⟨rfl, rfl⟩

/-- Go 1.27.0 src/math/stubs.go:92–94, selected by the !s390x build constraint. -/
def archFrexpStub (profile : ArchStubProfile) (_ : BitVec 64) : Option ArchStubOutcome :=
  archStubOutcome profile

theorem archFrexpStub_profiles (x : BitVec 64) :
    archFrexpStub .nonS390x x = some ⟨false, "not implemented"⟩ ∧
    archFrexpStub .s390x x = none := ⟨rfl, rfl⟩

/-- Go 1.27.0 src/math/stubs.go:98–100, selected by the !s390x build constraint. -/
def archLdexpStub (profile : ArchStubProfile) (_ : BitVec 64) (_ : Int) : Option ArchStubOutcome :=
  archStubOutcome profile

theorem archLdexpStub_profiles (x : BitVec 64) (exponent : Int) :
    archLdexpStub .nonS390x x exponent = some ⟨false, "not implemented"⟩ ∧
    archLdexpStub .s390x x exponent = none := ⟨rfl, rfl⟩

/-- Go 1.27.0 src/math/stubs.go:104–106, selected by the !s390x build constraint. -/
def archLog10Stub (profile : ArchStubProfile) (_ : BitVec 64) : Option ArchStubOutcome :=
  archStubOutcome profile

theorem archLog10Stub_profiles (x : BitVec 64) :
    archLog10Stub .nonS390x x = some ⟨false, "not implemented"⟩ ∧
    archLog10Stub .s390x x = none := ⟨rfl, rfl⟩

/-- Go 1.27.0 src/math/stubs.go:110–112, selected by the !s390x build constraint. -/
def archLog2Stub (profile : ArchStubProfile) (_ : BitVec 64) : Option ArchStubOutcome :=
  archStubOutcome profile

theorem archLog2Stub_profiles (x : BitVec 64) :
    archLog2Stub .nonS390x x = some ⟨false, "not implemented"⟩ ∧
    archLog2Stub .s390x x = none := ⟨rfl, rfl⟩

/-- Go 1.27.0 src/math/stubs.go:116–118, selected by the !s390x build constraint. -/
def archLog1pStub (profile : ArchStubProfile) (_ : BitVec 64) : Option ArchStubOutcome :=
  archStubOutcome profile

theorem archLog1pStub_profiles (x : BitVec 64) :
    archLog1pStub .nonS390x x = some ⟨false, "not implemented"⟩ ∧
    archLog1pStub .s390x x = none := ⟨rfl, rfl⟩

/-- Go 1.27.0 src/math/stubs.go:122–124, selected by the !s390x build constraint. -/
def archModStub (profile : ArchStubProfile) (_ : BitVec 64) (_ : BitVec 64) : Option ArchStubOutcome :=
  archStubOutcome profile

theorem archModStub_profiles (x : BitVec 64) (y : BitVec 64) :
    archModStub .nonS390x x y = some ⟨false, "not implemented"⟩ ∧
    archModStub .s390x x y = none := ⟨rfl, rfl⟩

/-- Go 1.27.0 src/math/stubs.go:128–130, selected by the !s390x build constraint. -/
def archPowStub (profile : ArchStubProfile) (_ : BitVec 64) (_ : BitVec 64) : Option ArchStubOutcome :=
  archStubOutcome profile

theorem archPowStub_profiles (x : BitVec 64) (y : BitVec 64) :
    archPowStub .nonS390x x y = some ⟨false, "not implemented"⟩ ∧
    archPowStub .s390x x y = none := ⟨rfl, rfl⟩

/-- Go 1.27.0 src/math/stubs.go:134–136, selected by the !s390x build constraint. -/
def archRemainderStub (profile : ArchStubProfile) (_ : BitVec 64) (_ : BitVec 64) : Option ArchStubOutcome :=
  archStubOutcome profile

theorem archRemainderStub_profiles (x : BitVec 64) (y : BitVec 64) :
    archRemainderStub .nonS390x x y = some ⟨false, "not implemented"⟩ ∧
    archRemainderStub .s390x x y = none := ⟨rfl, rfl⟩

/-- Go 1.27.0 src/math/stubs.go:140–142, selected by the !s390x build constraint. -/
def archSinStub (profile : ArchStubProfile) (_ : BitVec 64) : Option ArchStubOutcome :=
  archStubOutcome profile

theorem archSinStub_profiles (x : BitVec 64) :
    archSinStub .nonS390x x = some ⟨false, "not implemented"⟩ ∧
    archSinStub .s390x x = none := ⟨rfl, rfl⟩

/-- Go 1.27.0 src/math/stubs.go:146–148, selected by the !s390x build constraint. -/
def archSinhStub (profile : ArchStubProfile) (_ : BitVec 64) : Option ArchStubOutcome :=
  archStubOutcome profile

theorem archSinhStub_profiles (x : BitVec 64) :
    archSinhStub .nonS390x x = some ⟨false, "not implemented"⟩ ∧
    archSinhStub .s390x x = none := ⟨rfl, rfl⟩

/-- Go 1.27.0 src/math/stubs.go:152–154, selected by the !s390x build constraint. -/
def archTanStub (profile : ArchStubProfile) (_ : BitVec 64) : Option ArchStubOutcome :=
  archStubOutcome profile

theorem archTanStub_profiles (x : BitVec 64) :
    archTanStub .nonS390x x = some ⟨false, "not implemented"⟩ ∧
    archTanStub .s390x x = none := ⟨rfl, rfl⟩

/-- Go 1.27.0 src/math/stubs.go:158–160, selected by the !s390x build constraint. -/
def archTanhStub (profile : ArchStubProfile) (_ : BitVec 64) : Option ArchStubOutcome :=
  archStubOutcome profile

theorem archTanhStub_profiles (x : BitVec 64) :
    archTanhStub .nonS390x x = some ⟨false, "not implemented"⟩ ∧
    archTanhStub .s390x x = none := ⟨rfl, rfl⟩

end GoLean.Math
