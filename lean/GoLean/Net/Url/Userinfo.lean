import GoLean.Net.Url.Values

namespace GoLean.Net.Url

structure Userinfo where
  username : URLBytes
  password : URLBytes
  passwordSet : Bool

/-- Go 1.27.0 src/net/url/url.go:311–313. -/
def user (username : URLBytes) : Userinfo :=
  ⟨username, [], false⟩

/-- Go 1.27.0 src/net/url/url.go:323–325. An empty supplied password is set. -/
def userPassword (username password : URLBytes) : Userinfo :=
  ⟨username, password, true⟩

/-- Go 1.27.0 src/net/url/url.go:338–343. -/
def userinfoUsername (receiver : Option Userinfo) : URLBytes :=
  receiver.map Userinfo.username |>.getD []

/-- Go 1.27.0 src/net/url/url.go:346–351. A live receiver returns both
    fields independently, even on a synthetic state with the flag false. -/
def userinfoPassword (receiver : Option Userinfo) : URLBytes × Bool :=
  receiver.map (fun u => (u.password, u.passwordSet)) |>.getD ([], false)

/-- Go 1.27.0 src/net/url/url.go:355–364. Escape each component in
    user-password mode and add the colon only when the password is set. -/
def userinfoString (receiver : Option Userinfo) : URLBytes :=
  receiver.map (fun u =>
    let name := escape u.username encodeUserPassword
    if u.passwordSet then name ++ [58] ++ escape u.password encodeUserPassword else name) |>.getD []

private def userinfoByte (c : BitVec 8) : Bool :=
  let n := c.toNat
  (65 ≤ n && n ≤ 90) || (97 ≤ n && n ≤ 122) || (48 ≤ n && n ≤ 57) ||
    [45, 46, 95, 58, 126, 33, 36, 38, 39, 40, 41, 42, 43, 44, 59, 61, 37, 64].contains c

/-- Go 1.27.0 src/net/url/url.go:1283–1309. Every allowed rune is ASCII.
    A non-ASCII byte either belongs to a non-ASCII rune or decodes as
    RuneError, so the rune loop is equivalent to a raw ASCII filter. -/
def validUserinfo (value : URLBytes) : Bool :=
  value.all userinfoByte

theorem user_noPassword (name : URLBytes) : user name = ⟨name, [], false⟩ := rfl

theorem userPassword_empty (name : URLBytes) : userPassword name [] = ⟨name, [], true⟩ := rfl

theorem userinfoUsername_nil : userinfoUsername none = [] := rfl

theorem userinfoPassword_nil : userinfoPassword none = ([], false) := rfl

theorem userinfoPassword_preservesStored (u : Userinfo) :
    userinfoPassword (some u) = (u.password, u.passwordSet) := rfl

theorem userinfoString_presentEmpty :
    userinfoString (some ⟨[], [], true⟩) = [58] := rfl

theorem userinfoString_nil : userinfoString none = [] := rfl

theorem validUserinfo_delimiters : validUserinfo [58, 64, 37, 43] = true := rfl

theorem validUserinfo_slash : validUserinfo [47] = false := rfl

end GoLean.Net.Url
