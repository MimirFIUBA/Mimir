package controllers

import "mimir/internal/mimir"

var (
	MimirProcessor *mimir.MimirProcessor
)

func SetMimirProcessor(mimirProcessor *mimir.MimirProcessor) {
	MimirProcessor = mimirProcessor
}
