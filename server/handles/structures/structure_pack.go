package structures

type Packet struct {
	PortTo   int     `json:"num"`
	Message  string  `json:"message"`
	NumFloat float32 `json:"numFloat"`
	BigMass  []int32 `json:"bigMass"`
}
