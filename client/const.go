package client

import "time"

const (
	Host               = "ws.xtb.com"
	ModeDemo           = "demo"
	ModeReal           = "real"
	MinReuqestInterval = 200 * time.Millisecond

	modeDemoStream = "demoStream"
	modeRealStream = "realStream"
)

type Quote int64

const (
	Fixed Quote = iota + 1
	Float
	Depth
	Cross
)
