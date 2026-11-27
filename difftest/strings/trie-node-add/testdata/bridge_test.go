package strings

import "strconv"

func goLeanTrieAddBytes(raw string) string {
	var out []byte
	for i := range len(raw) {
		if i != 0 {
			out = append(out, ',')
		}
		out = strconv.AppendUint(out, uint64(raw[i]), 10)
	}
	return string(out)
}

func GoLeanTrieNodeAdd(oldnew []string, key, value string, priority int) string {
	r := &genericReplacer{tableSize: 256}
	r.root.table = make([]*trieNode, r.tableSize)
	for index := range r.mapping {
		r.mapping[index] = byte(index)
	}
	for i := 0; i < len(oldnew); i += 2 {
		r.root.add(oldnew[i], oldnew[i+1], len(oldnew)-i, r)
	}
	r.root.add(key, value, priority, r)
	var visit func(*trieNode) string
	visit = func(node *trieNode) string {
		terminal := strconv.Itoa(node.priority) + ":" + goLeanTrieAddBytes(node.value)
		switch {
		case node.table != nil:
			if len(node.table) != r.tableSize || node.prefix != "" || node.next != nil {
				panic("invalid table shape")
			}
			var children []string
			for index, child := range node.table {
				if child != nil {
					children = append(children, strconv.Itoa(index)+"="+visit(child))
				}
			}
			return "T" + terminal + "{" + Join(children, "|") + "}"
		case node.prefix != "":
			if node.next == nil {
				panic("missing compressed child")
			}
			return "E" + terminal + "[" + goLeanTrieAddBytes(node.prefix) + "|" + visit(node.next) + "]"
		default:
			if node.next != nil {
				panic("leaf has a child")
			}
			return "L" + terminal
		}
	}
	return visit(&r.root)
}
