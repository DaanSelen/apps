package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

var (
	textFileName   = "config.txt"
	configInputs   [4]string //indexes are linked to keywordsConfig array!
	keywordsConfig [4]string //LINKED TO configInputs by index!
)

const (
	totalKeywords int = len(keywordsConfig)
)

func main() {
	//Here is the section where the keywords must be filled in.
	keywordsConfig[0] = "example"
	keywordsConfig[1] = ""
	keywordsConfig[2] = ""
	keywordsConfig[3] = ""

	checkConfig()

	fmt.Scanln()
}

func checkConfig() {
	f, err := os.Open(textFileName) //Opens the txt file

	if err != nil { //If it can't open the file it errors and quits.
		log.Fatal(err)
	}

	defer f.Close() //Once the program is done with using the txt file it will close it.

	lineByLine := bufio.NewScanner(f)
	a := 0
	for lineByLine.Scan() { //The code is made that if a sentence contains a '#' or the line is completely empty it will ignore the line and continue
		if strings.Contains(lineByLine.Text(), "#") || lineByLine.Text() == "" {
			fmt.Println("Line Ignored")
		} else {
			if strings.Contains(lineByLine.Text(), (keywordsConfig[a] + " = ")) {
				var restant1 = strings.ReplaceAll(lineByLine.Text(), (keywordsConfig[a] + " = "), "")
				dummy, err := strconv.Atoi(restant1)
				if err != nil {
					fmt.Print("Value is not a number: ", restant1+"\n")
					putStringsFromConfigIntoArray(restant1, keywordsConfig[a])
				} else {
					fmt.Println("Value is a number!", dummy)
					putNumberFromConfigIntoArray(restant1, keywordsConfig[a])
				}
			}
			if a < totalKeywords {
				a++
			}
		}
	}
}

func putNumberFromConfigIntoArray(restant1, corrKeyword string) { //This is where the int values go
	for i := 0; i < totalKeywords; i++ {
		if corrKeyword == keywordsConfig[i] {
			configInputs[i] = restant1
		}
	}
}

func putStringsFromConfigIntoArray(restant1, corrKeyword string) { //This is where the string (to be converted to bools) go
	for i := 0; i < totalKeywords; i++ {
		if corrKeyword == keywordsConfig[i] {
			configInputs[i] = strings.ToLower(restant1)
		}
	}
}
