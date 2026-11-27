package strings

import (
	"runtime"
	"strconv"
	"sync"
)

func goLeanBuildOnceBytes(raw string) string {
	var out []byte
	for i := range len(raw) {
		if i != 0 {
			out = append(out, ',')
		}
		out = strconv.AppendUint(out, uint64(raw[i]), 10)
	}
	return string(out)
}

func GoLeanReplacerBuildOnce(oldnew []string, present, onceZero, repeat bool) (snapshot string, once *sync.Once) {
	var b *Replacer
	if present {
		b = NewReplacer(oldnew...)
		if !onceZero {
			b.once.Do(func() {})
		}
	}
	defer func() {
		if failure := recover(); failure != nil {
			fault, ok := failure.(runtime.Error)
			if b != nil || !ok || fault.Error() != "runtime error: invalid memory address or nil pointer dereference" {
				panic(failure)
			}
			snapshot = "nil-receiver"
		}
	}()
	b.buildOnce()
	if repeat {
		b.buildOnce()
	}
	pending := "nil:"
	if b.oldnew != nil {
		pending = "slice:"
	}
	return pending + goLeanBuildOnceSnapshot(b.r), &b.once
}
func goLeanBuildOnceSnapshot(actual replacer) string {
	switch result := actual.(type) {
	case *singleStringReplacer:
		return "single:" + goLeanBuildOnceBytes(result.finder.pattern) + "|" + goLeanBuildOnceBytes(result.value)
	case *genericReplacer:
		return "generic:" + strconv.Itoa(result.tableSize) + ":" + goLeanBuildOnceBytes(string(result.mapping[:]))
	case *byteReplacer:
		return "bytes:" + goLeanBuildOnceBytes(string(result[:]))
	case *byteStringReplacer:
		var keys []byte
		for _, key := range result.toReplace {
			if len(key) != 1 {
				panic("invalid toReplace key length")
			}
			keys = append(keys, key[0])
		}
		values := make([]string, len(result.replacements))
		for i, value := range result.replacements {
			if value == nil {
				values[i] = "-"
			} else {
				values[i] = "+" + goLeanBuildOnceBytes(string(value))
			}
		}
		return "byteStrings:" + goLeanBuildOnceBytes(string(keys)) + "|" + Join(values, "|")
	default:
		panic("unrecognized replacement implementation")
	}
}
