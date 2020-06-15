package models

type FromFrontend struct {
	X int `json:"x"`
	Y int `json:"y"`
	RefreshToken string `json:"refresh_token"`
}

type ForFrontend struct {
	StepSize int `json:"step_size"`
	ColorPalette RGB `json:"color_palette"`
	EllipseHeight int `json:"ellipse_height"`
	EllipseWidth int `json:"ellipse_width"`
	X int `json:"x"`
	Y int `json:"y"`
	SongName string `json:"song_name"`

}