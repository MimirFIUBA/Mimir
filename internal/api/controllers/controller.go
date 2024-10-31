package controllers

import "mimir/internal/mimir"

var (
	MimirEngine *mimir.MimirEngine
)

func SetMimirEngine(mimirProcessor *mimir.MimirEngine) {
	MimirEngine = mimirProcessor
}
