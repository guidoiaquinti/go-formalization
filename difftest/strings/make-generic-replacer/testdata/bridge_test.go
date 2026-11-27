package strings

import "strconv"

func goLeanGenericBytes(raw string) string {
	var out []byte
	for i := range len(raw) {
		if i != 0 {
			out = append(out, ',')
		}
		out = strconv.AppendUint(out, uint64(raw[i]), 10)
	}
	return string(out)
}

func GoLeanMakeGenericReplacer(oldnew []string) string {
	replacer := makeGenericReplacer(oldnew)
	var alphabet [256]byte
	for raw, index := range replacer.mapping {
		if int(index) < replacer.tableSize {
			alphabet[index] = byte(raw)
		}
	}
	var visit func(*trieNode) string
	visit = func(node *trieNode) string {
		terminal := strconv.Itoa(node.priority) + ":" + goLeanGenericBytes(node.value)
		switch {
		case node.table != nil:
			if len(node.table) != replacer.tableSize || node.prefix != "" || node.next != nil {
				panic("invalid table shape")
			}
			var children []string
			for index, child := range node.table {
				if child != nil {
					children = append(children, strconv.Itoa(int(alphabet[index]))+"="+visit(child))
				}
			}
			return "T" + terminal + "{" + Join(children, "|") + "}"
		case node.prefix != "":
			if node.next == nil {
				panic("missing compressed child")
			}
			return "E" + terminal + "[" + goLeanGenericBytes(node.prefix) + "|" + visit(node.next) + "]"
		default:
			if node.next != nil {
				panic("leaf has a child")
			}
			return "L" + terminal
		}
	}
	return strconv.Itoa(replacer.tableSize) + ":" + goLeanGenericBytes(string(replacer.mapping[:])) + "|" + visit(&replacer.root)
}
