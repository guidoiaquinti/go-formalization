package strings

import "strconv"

func goLeanReplacerBuildBytes(raw string) string {
	var out []byte
	for i := range len(raw) {
		if i != 0 {
			out = append(out, ',')
		}
		out = strconv.AppendUint(out, uint64(raw[i]), 10)
	}
	return string(out)
}

func GoLeanReplacerBuild(oldnew []string) (string, any, any) {
	b := NewReplacer(oldnew...)
	actual := b.build()
	switch result := actual.(type) {
	case *singleStringReplacer:
		return "single:" + goLeanReplacerBuildBytes(result.finder.pattern) + "|" + goLeanReplacerBuildBytes(result.value),
			result, makeSingleStringReplacer(oldnew[0], oldnew[1])
	case *genericReplacer:
		return "generic:" + strconv.Itoa(result.tableSize) + ":" + goLeanReplacerBuildBytes(string(result.mapping[:])),
			result, makeGenericReplacer(oldnew)
	case *byteReplacer:
		return "bytes:" + goLeanReplacerBuildBytes(string(result[:])), nil, nil
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
				values[i] = "+" + goLeanReplacerBuildBytes(string(value))
			}
		}
		return "byteStrings:" + goLeanReplacerBuildBytes(string(keys)) + "|" + Join(values, "|"), nil, nil
	default:
		panic("unrecognized replacement implementation")
	}
}
