import GoLean.Net.Url.Userinfo

namespace GoLean.Net.Url

structure URLFields where
  scheme : URLBytes
  opaqueValue : URLBytes
  host : URLBytes
  path : URLBytes
  fragment : URLBytes
  rawQuery : URLBytes
  rawPath : URLBytes
  rawFragment : URLBytes
  forceQuery : Bool
  omitHost : Bool

structure URLObject where
  fields : URLFields
  identity : Nat
  user : Option (Nat × Userinfo)

/-- Go 1.27.0 net/url.(*URL).Clone. Fresh identities represent the two possible allocations. -/
def urlClone (receiver : Option URLObject) (freshURL freshUser : Nat) : Option URLObject :=
  receiver.map fun u => ⟨u.fields, freshURL, u.user.map fun old => (freshUser, old.2)⟩

theorem urlClone_nil (a b : Nat) : urlClone none a b = none := rfl

theorem urlClone_values (u : URLObject) (a b : Nat) :
    (urlClone (some u) a b).map (fun c => (c.fields, c.user.map Prod.snd)) =
      some (u.fields, u.user.map Prod.snd) := by
  cases h : u.user with
  | none => simp [urlClone, h]
  | some entry => cases entry <;> simp [urlClone, h]

theorem urlClone_fresh (u : URLObject) (a b : Nat) (ha : a ≠ u.identity)
    (hb : ∀ entry ∈ u.user, b ≠ entry.1) :
    (urlClone (some u) a b).all (fun c =>
      c.identity != u.identity && c.user.all (fun entry =>
        u.user.all (fun old => entry.1 != old.1))) = true := by
  cases h : u.user with
  | none => simp [urlClone, h, ha]
  | some old =>
    have different : b ≠ old.1 := hb old (by simp [h])
    simp [urlClone, h, ha, different]

end GoLean.Net.Url
