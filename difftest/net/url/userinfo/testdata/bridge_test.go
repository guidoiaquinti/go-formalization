package url

func GoLeanSyntheticUserinfo(username, password string, set bool) *Userinfo {
	return &Userinfo{username: username, password: password, passwordSet: set}
}
func GoLeanValidUserinfo(s string) bool { return validUserinfo(s) }

// GoLeanUserinfoFields reads the constructor state independently of the accessors.
func GoLeanUserinfoFields(u *Userinfo) (string, string, bool) {
	return u.username, u.password, u.passwordSet
}
