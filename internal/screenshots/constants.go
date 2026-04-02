package screenshots

type Quality int

const (
	QualityBest Quality = iota
	QualityHigh
	QualityMedium
)

type FullscreenMode int

const (
	FullscreenModeMouse FullscreenMode = iota
	FullscreenModeAllScreens
	FullscreenModePrimary
)
