package models

type RGB struct {
	Red int
	Green int
	Blue int
}

type WarmColorsCollection struct {
	WarmOrange RGB
	WarmYellow RGB
	WarmRed RGB
	WarmMagenta RGB
	WarmPink RGB
}

type ConfidentColorsCollection struct {
	ConfidentOrange RGB
	ConfidentYellow RGB
	ConfidentRed RGB
	ConfidentMagenta RGB
	ConfidentPink RGB
}

type FunkyColorsCollection struct {
	FunkyLightGreen RGB
	FunkyDarkGreen RGB
	FunkyLightBlue RGB
	FunkyDarkBlue RGB
	FunkyOrange RGB
}

type CalmColorsCollection struct {
	CalmLightBlue RGB
	CalmDarkBlue RGB
	CalmLightGreen RGB
	CalmDarkGreen RGB
	CalmTurquoise RGB
}

type RelaxedColorsCollection struct {
	RelaxedBlue RGB
	RelaxedOrange RGB
	RelaxedBrown RGB
	RelaxedPink RGB
	RelaxedYellow RGB
}

type SadColorsCollection struct {
	SadBlue RGB
	SadRed RGB
	SadLightOrange RGB
	SadDarkOrange RGB
	SadGreen RGB
	
}

var WarmColors = WarmColorsCollection{
	WarmOrange: RGB {
		Red:   255,
		Green: 114,
		Blue:  81,
	},
	WarmYellow: RGB{
		Red:   247,
		Green: 235,
		Blue:  48,
	},
	WarmRed: RGB{
		Red:   220,
		Green: 72,
		Blue:  8,
	},
	WarmMagenta: RGB{
		Red:   250,
		Green: 154,
		Blue:  251,
	},
	WarmPink: RGB{
		Red:   250,
		Green: 132,
		Blue:  149,
	},
}

var ConfidentColors = ConfidentColorsCollection{
	ConfidentOrange: RGB{
		Red:   246,
		Green: 145,
		Blue:  39,
	},
	ConfidentYellow: RGB{
		Red:   245,
		Green: 241,
		Blue:  24,
	},
	ConfidentRed: RGB{
		Red:   245,
		Green: 105,
		Blue:  16,
	},
	ConfidentMagenta: RGB{
		Red:   245,
		Green: 16,
		Blue:  133,
	},
	ConfidentPink: RGB{
		Red:   249,
		Green: 117,
		Blue:  220,
	},
}

var FunkyColors = FunkyColorsCollection{
	FunkyLightGreen: RGB{
		Red:   77,
		Green: 248,
		Blue:  106,
	},
	FunkyDarkGreen: RGB{
		Red:   6,
		Green: 156,
		Blue:  12,
	},
	FunkyLightBlue: RGB{
		Red:   104,
		Green: 250,
		Blue:  241,
	},
	FunkyDarkBlue: RGB{
		Red:   57,
		Green: 86,
		Blue:  249,
	},
	FunkyOrange: RGB{
		Red:   248,
		Green: 101,
		Blue:  24,
	},
}

var RelaxedColors = RelaxedColorsCollection{
	RelaxedBlue: RGB{
		Red:   19,
		Green: 149,
		Blue:  248,
	},
	RelaxedOrange: RGB{
		Red:   244,
		Green: 160,
		Blue:  8,
	},
	RelaxedBrown: RGB{
		Red:   132,
		Green: 88,
		Blue:  4,
	},
	RelaxedPink: RGB{
		Red:   252,
		Green: 159,
		Blue:  204,
	},
	RelaxedYellow: RGB{
		Red:   252,
		Green: 139,
		Blue:  169,
	},
}

var SadColors = SadColorsCollection{
	SadBlue: RGB{
		Red:   5,
		Green: 77,
		Blue:  126,
	},
	SadRed: RGB{
		Red:   75,
		Green: 13,
		Blue:  3,
	},
	SadLightOrange: RGB{
		Red:   198,
		Green: 129,
		Blue:  8,
	},
	SadDarkOrange: RGB{
		Red:   163,
		Green: 75,
		Blue:  7,
	},
	SadGreen: RGB{
		Red:   2,
		Green: 52,
		Blue:  8,
	},
}






