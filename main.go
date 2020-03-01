package main

import (
	"flag"

	colorable "github.com/mattn/go-colorable"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var layerNumber = 0

func main() {
	//CLI Flags
	debug := flag.Bool("v", false, "sets log level to debug")
	inputFilepath := flag.String("i", "input.sl1", "SL1 input file")
	//outputFilepath := flag.String("o", "output.photon", "Phonton ouput file")
	flag.Parse()

	//Setup logs
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out: colorable.NewColorableStdout(),
	})

	//Open a sl1 (zip archive).
	input, err := ReadSL1File(*inputFilepath)
	if err != nil {
		log.Fatal().Err(err).Msgf("Failed to read input file: %s", inputFilepath)
	}
	log.Debug().Msgf("DEBUG input: %#v", input)

	//Load empty .photon file.
	output, err := ReadPhotonFile("./assets/newfile.photon")
	if err != nil {
		log.Fatal().Err(err).Msgf("Failed to generate empty photon file")
	}

	//TODO
	/*
		    photon.exposure_time = float(sl1.config['expTime'])
		    photon.exposure_time_bottom = float(sl1.config['expTimeFirst'])
		    photon.layer_height = float(sl1.config['layerHeight'])
			photon.bottom_layers = int(sl1.config['numFade'])
	*/

	log.Debug().Msgf("DEBUG output: %#v", output)
}
