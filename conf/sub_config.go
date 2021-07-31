package conf

type VSub struct {
	SubUrl  string
	SubName string
}

// subName <-> VSub
var SubConfigNow map[string]VSub = map[string]VSub{}
