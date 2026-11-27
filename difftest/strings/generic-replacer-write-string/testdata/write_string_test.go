package strings_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"testing"
)

var writeFailures = [...]error{nil, errors.New("first scripted error"), errors.New("second scripted error")}

type replacerReply struct{ limit, failure int }

type replacerWriter struct {
	replies []replacerReply
	calls   []string
}

type replacerStringWriter struct{ *replacerWriter }

func (w *replacerWriter) reply(payload string, mode byte) (int, error) {
	event := []byte{mode}
	for i := range len(payload) {
		event = strconv.AppendUint(append(event, ','), uint64(payload[i]), 10)
	}
	n := len(payload)
	var failure error
	if i := len(w.calls); i < len(w.replies) {
		reply := w.replies[i]
		if reply.limit >= 0 {
			n = min(n, reply.limit)
		}
		failure = writeFailures[reply.failure]
	}
	w.calls = append(w.calls, string(event))
	return n, failure
}

func (w *replacerWriter) Write(payload []byte) (int, error) {
	if payload == nil {
		return w.reply("", 'N')
	}
	return w.reply(string(payload), 'B')
}

func (w replacerStringWriter) WriteString(payload string) (int, error) {
	return w.reply(payload, 'S')
}

func TestGoLeanGenericReplacerWriteString(t *testing.T) {
	file, err := os.Open(os.Getenv("GOLEAN_GENERIC_REPLACER_WRITE_STRING_CASES"))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	rows := json.NewDecoder(file)
	for i := 0; ; i++ {
		var row struct {
			Pairs        [][]byte
			Text         []byte
			StringWriter bool
			Script       [][2]int
			Want         string
		}
		if err := rows.Decode(&row); err == io.EOF {
			if i == 0 {
				t.Fatal("empty comparison fixture")
			}
			t.Logf("%d native writer traces agree", i)
			break
		} else if err != nil {
			t.Fatal(err)
		}
		pairs := make([]string, len(row.Pairs))
		for j, field := range row.Pairs {
			pairs[j] = string(field)
		}
		writer := &replacerWriter{replies: make([]replacerReply, len(row.Script))}
		for j, reply := range row.Script {
			writer.replies[j] = replacerReply{reply[0], reply[1]}
		}
		var native io.Writer = writer
		if row.StringWriter {
			native = replacerStringWriter{writer}
		}
		written, failure := strings.GoLeanGenericReplacerWriteString(pairs, native, string(row.Text))
		code := -1
		for j, sentinel := range writeFailures {
			if failure == sentinel {
				code = j
				break
			}
		}
		got := fmt.Sprintf("%d:%d:%s", written, code, strings.Join(writer.calls, ";"))
		if got != row.Want {
			t.Fatalf("row %d pairs=%q input=%x: Go=%s want=%s", i, pairs, row.Text, got, row.Want)
		}
	}
}
