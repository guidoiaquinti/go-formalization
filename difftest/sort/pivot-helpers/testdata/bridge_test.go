package sort

type goLeanPivotData struct{ less func(int, int) bool }

func (d goLeanPivotData) Len() int           { panic("pivot helper unexpectedly called Len") }
func (d goLeanPivotData) Less(i, j int) bool { return d.less(i, j) }
func (d goLeanPivotData) Swap(i, j int)      { panic("pivot helper unexpectedly called Swap") }

func GoLeanPivotHelpers(functionVariant bool, kind, a, b, c, initial int, less func(int, int) bool) (int, int, int) {
	data := goLeanPivotData{less}
	if functionVariant {
		callbacks := lessSwap{Less: data.Less, Swap: data.Swap}
		switch kind {
		case 0:
			x, y := order2_func(callbacks, a, b, &initial)
			return x, y, initial
		case 1:
			x := median_func(callbacks, a, b, c, &initial)
			return x, -1, initial
		case 2:
			x := medianAdjacent_func(callbacks, a, &initial)
			return x, -1, initial
		default:
			x, hint := choosePivot_func(callbacks, a, b)
			return x, int(hint), 0
		}
	}
	switch kind {
	case 0:
		x, y := order2(data, a, b, &initial)
		return x, y, initial
	case 1:
		x := median(data, a, b, c, &initial)
		return x, -1, initial
	case 2:
		x := medianAdjacent(data, a, &initial)
		return x, -1, initial
	default:
		x, hint := choosePivot(data, a, b)
		return x, int(hint), 0
	}
}
