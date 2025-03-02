package ecstypes

const (
	SystemPosition SystemID = 1 << iota
	SystemMotion
	SystemHelm
	SystemSprite
	SystemPlayer
	SystemSyncReceiver
	SystemSyncSender
)
