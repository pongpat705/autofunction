package stuct

type ConfigModel struct {
	Config       []ParamModel `json:"config"`
	WhenPressKey string       `json:"whenPressKey"`
	Mode         string       `json:"mode"`
}
