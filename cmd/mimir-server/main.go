package main

import (
	"fmt"
	"mimir/internal/mimir"
	"mimir/triggers"
	"os"
	"os/signal"
	"syscall"
)

func setInitialData(mp *mimir.MimirProcessor) {
	// Setup1(mp)
	cond := triggers.NewMatchesCondition(`foo.*`)
	trig := triggers.NewTrigger("test")
	trig.Condition = cond
	print := triggers.NewPrintAction()
	print.Message = "Matchescondition"
	trig.AddAction(print)

	event := triggers.NewEvent()
	event.Data = "seafood"

	trig.Update(*event)

}

func main() {

	fmt.Println("MiMiR starting")

	mimirProcessor := mimir.NewMimirProcessor()
	mimirProcessor.StartGateway()

	setInitialData(mimirProcessor)
	// go mimirProcessor.Run()
	// go api.Start(mimirProcessor.WsChannel)

	fmt.Println("Everything up and running")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	mimir.CloseConnection()

	fmt.Println("Mimir is out of duty, bye!")
}
