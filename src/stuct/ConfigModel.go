package stuct

type ConfigModel struct {
	StartPositionX int    `json:"startPositionX"`
	StartPositionY int    `json:"startPositionY"`
	EndPositionX   int    `json:"endPositionX"`
	EndPositionY   int    `json:"endPositionY"`
	KeyPressDelay  int    `json:"keyPressDelay"`
	KeyToPress     string `json:"keyToPress"`
}
