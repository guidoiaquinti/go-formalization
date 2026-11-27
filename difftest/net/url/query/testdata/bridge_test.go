package url

func GoLeanParamsWithinMax(params int) bool { return urlParamsWithinMax(params) }

func GoLeanParseQuery(m Values, s string) error { return parseQuery(m, s) }
