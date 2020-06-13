package models

type FromFrontend struct {
	X int
	Y int
	RefreshToken string
}

type ForFrontend struct {
	StepSize int
	ColorPalette RGB
	EllipseHeight int
	EllipseWidth int
	X int
	Y int
	SongName string

}