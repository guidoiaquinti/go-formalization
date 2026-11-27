package url_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"
)

func TestGoLeanUserinfo(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_URL_USERINFO_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	wire := func(raw string) string {
		fields := make([]string, len(raw))
		for i := range fields {
			fields[i] = strconv.Itoa(int(raw[i]))
		}
		return fmt.Sprintf("%d:%s", len(raw), strings.Join(fields, ","))
	}
	password := func(u *url.Userinfo) string {
		pw, set := u.Password()
		return wire(pw) + ":" + strconv.FormatBool(set)
	}
	describe := func(u *url.Userinfo) string {
		if u == nil {
			t.Fatal("constructor returned nil")
		}
		name, secret, set := url.GoLeanUserinfoFields(u)
		return wire(name) + "/" + wire(secret) + ":" + strconv.FormatBool(set)
	}
	for i := 0; ; i++ {
		var c struct {
			Tag                    int
			Name, Secret, Validate []byte
			Want                   string
		}
		if err := decoder.Decode(&c); errors.Is(err, io.EOF) {
			if i == 0 {
				t.Fatal("missing userinfo fixtures")
			}
			t.Logf("six pinned URL userinfo declarations matched %d rows", i)
			return
		} else if err != nil {
			t.Fatal(err)
		}
		name, secret := string(c.Name), string(c.Secret)
		var receiver *url.Userinfo
		if c.Tag != 0 {
			receiver = url.GoLeanSyntheticUserinfo(name, secret, c.Tag == 2)
		}
		got := strings.Join([]string{describe(url.User(name)), describe(url.UserPassword(name, secret)),
			wire(receiver.Username()), password(receiver), wire(receiver.String()),
			strconv.FormatBool(url.GoLeanValidUserinfo(string(c.Validate)))}, "|")
		if got != c.Want {
			t.Fatalf("row %d %+v: got %q want %q", i, c, got, c.Want)
		}
	}
}
