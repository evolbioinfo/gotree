package main

import (
	"fmt"
	"github.com/fredericlemoine/gotree/io"
	"github.com/fredericlemoine/gotree/io/utils"
	"os"
	"os/signal"
	"time"
)

func main() {
	//fmt.Fprintf(os.Stderr, "Started Quartets\n")
	quartet, err := utils.ReadRefTree("tests/data/quartets.nw.gz")
	nbquartets := 0
	t := time.Now()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			fmt.Println(sig)
			fmt.Println(time.Now().Sub(t), nbquartets)
		}
	}()

	if err != nil {
		io.ExitWithMessage(err)
	}
	quartet.Quartets(func(tb1, tb2, tb3, tb4 uint) {
		nbquartets++
		//fmt.Fprintf(os.Stderr, "(%d,%d)(%d,%d)\n", tb1, tb2, tb3, tb4)
	})
	//fmt.Fprintf(os.Stderr, "End Quartets\n")
	// Total quartets     : 40 693 092 081 640
	// Total quartets ref :     40 842 660 378
}
