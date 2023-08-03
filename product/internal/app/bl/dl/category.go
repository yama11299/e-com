package dl

// category constants
const (
	Budget = iota + 1
	Regular
	Premium
)

// category id to name mapping
var (
	CategoryIDNameMap = map[int]string{
		Budget:  "Budget",
		Regular: "Regular",
		Premium: "Premium",
	}
)
