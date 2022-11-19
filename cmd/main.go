package main

import (
	"flag"
	"log"
	"os"
	"runtime"
	"runtime/pprof"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mrmarble/sandbox/pkg/game"
)

var (
	debugCpuprofile     = flag.String("debug_cpuprofile", "", "write CPU profile to file")
	debugMemprofile     = flag.String("debug_memprofile", "", "write memory profile to file")
	debugMemprofileRate = flag.Int("debug_memprofile_rate", runtime.MemProfileRate, "fraction of bytes to be included in -debug_memprofile")
)

func main() {
	flag.Parse()

	if *debugMemprofile != "" {
		// Set the memory profile rate as soon as possible.
		runtime.MemProfileRate = *debugMemprofileRate
	}

	if *debugCpuprofile != "" {
		log.Println("Starting CPU profile")
		f, err := os.Create(*debugCpuprofile)
		if err != nil {
			log.Fatalf("could not create CPU profile: %v", err)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			f.Close()
			log.Fatalf("could not start CPU profile: %v", err)
		}

		defer f.Close()
		defer pprof.StopCPUProfile()
	}

	if *debugMemprofile != "" {
		log.Println("Starting memory profile")
		f, err := os.Create(*debugMemprofile)
		if err != nil {
			log.Fatalf("could not create memory profile: %v", err)
		}
		defer f.Close()
		runtime.GC()
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatalf("could not write memory profile: %v", err)
		}
	}
	game := game.New()

	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
