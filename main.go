package main

import (
	"golang.org/x/tools/imports"
	"log"
	"mockif/generator"
	"mockif/ifparser"
	"os"
)

const defaultOutputFilename = "mocks.go"

func main() {
	if len(os.Args) <= 1 {
		log.Fatalln("no input directory specified!")
	}

	inputDir := os.Args[1]

	// Optional param - file name for generated mocks
	outputFilename := defaultOutputFilename
	if len(os.Args) >= 3 {
		outputFilename = os.Args[2]
	}

	parsedPackage := ifparser.ParseDir(inputDir)

	if len(parsedPackage.Interfaces) <= 0 {
		log.Println("no interfaces found in package provided, exiting.")
		os.Exit(0)
	}

	log.Printf("%d interface declarations found, generating mocks...\n", len(parsedPackage.Interfaces))

	output := generator.Generate(parsedPackage)

	filename := inputDir + "/" + outputFilename

	// Remove unused imports
	result, err := imports.Process(filename, []byte(output), nil)
	if err != nil {
		log.Fatalln(err)
	}

	err = os.WriteFile(filename, result, 0644)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Done!")
}
