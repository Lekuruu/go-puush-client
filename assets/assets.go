package assets

import _ "embed"

//go:embed icons/puush.png
var PuushIconData []byte

//go:embed sounds/success.wav
var SuccessSoundData []byte

//go:embed icons/icon-window.png
var WindowIconData []byte

//go:embed icons/icon-fullscreen.png
var FullscreenIconData []byte

//go:embed icons/icon-upload.png
var UploadIconData []byte

//go:embed icons/icon-selection.png
var SelectionIconData []byte

//go:embed quickstart.png
var QuickstartData []byte // TODO: Remove "windows" text from quickstart asset
