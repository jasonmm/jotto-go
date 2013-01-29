package main

import (
	"../../gowc/libgowc"
	jotto "../../libjotto"
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"
)

const APP_VERSION = "0.1"
const APP_NAME = "Jotto"
const APP_VENDOR = "Shady Brook Software"

type Game struct {
	secretWord string
	numGuesses int
}

var (
	game = Game{}
)

// The flag package provides a default help printer via -h switch
var versionFlag *bool = flag.Bool("v", false, "Print the version number.")

// Counts the number of words in the given file descriptor.
func countPossibleWords(fp *os.File) int {
	var cnt int = 0

	reader := bufio.NewReader(fp)
	for {
		if _, err := reader.ReadString('\n'); err != nil {
			break
		}
		cnt++
	}
	return cnt
}

func getWord(lineNum int, fp *os.File) (word string, e error) {
	var cnt int = 0
	var line string = ""
	var err error = nil

	reader := bufio.NewReader(fp)
	for {
		if line, err = reader.ReadString('\n'); err != nil {
			break
		}
		cnt++
		if cnt == lineNum {
			break
		}
	}
	return line, err
}

func chooseSecretWord() error {
	var fp *os.File
	var err error
	var numWords libgowc.Metrics

	fmt.Print("Choosing secret word...")

	if fp, err = os.Open("wordlist.txt"); err != nil {
		return err
	}
	defer fp.Close()

	numWords, err = libgowc.ProcessSingleFile("wordlist.txt")
	if err != nil {
		return err
	}

	rand.Seed(time.Now().UTC().UnixNano())
	selectedIndex := rand.Intn(numWords.Lines)

	game.secretWord, err = getWord(selectedIndex, fp)
	if err != nil {
		return err
	}

	fmt.Println("done.")

	return nil
}

func checkGuess(guess string) (numRight int, numRightPlace int) {
	return jotto.GuessResult(guess, game.secretWord)
}

func main() {
	flag.Parse() // Scan the arguments list 

	if *versionFlag {
		fmt.Println(APP_NAME, "by", APP_VENDOR)
		fmt.Println("Version:", APP_VERSION)
	}

	fmt.Println()
	fmt.Print("Welcome to ", APP_NAME)
	fmt.Println("!")
	fmt.Println()

	if err := chooseSecretWord(); err != nil {
		fmt.Println("Error: ", err)
		return
	}

	fmt.Println("Word: ", game.secretWord)

	var guess string

	for {
		fmt.Println()
		fmt.Print("Enter guess: ")
		if _, err := fmt.Scanln(&guess); err != nil {
			fmt.Println("  - Error! ", err)
			continue
		}

		fmt.Println()
		numRight, numRightPlace := checkGuess(guess)
		if numRight == numRightPlace {
			fmt.Println("Correct!")
			break
		}

		fmt.Println("Guess incorrect: ", numRight, " right; ", numRightPlace, " right place")

	}
}
