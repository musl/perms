package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	_ "embed"
)

//go:embed words_alpha.txt
var defaultDictionaryString string
var defaultDictionary = map[string]int{}

func init() {
	entries := strings.Split(defaultDictionaryString, "\n")
	for _, entry := range entries {
		defaultDictionary[entry]++
		log.Println(entry)
	}
}

func loadDictionary(path string) (map[string]int, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	wordMap := map[string]int{}
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		wordMap[scanner.Text()]++
	}

	return wordMap, nil
}

func permutations(word string, dictionary map[string]int) []string {
	if len(word) == 0 {
		return []string{word}
	}

	wordMap := map[string]int{}
	runes := []rune(word)

	var rc func(int)

	rc = func(np int) {
		if np == 1 {
			if _, ok := dictionary[string(runes)]; ok {
				wordMap[string(runes)]++
			}
			return
		}

		np1 := np - 1
		pp := len(runes) - np1

		rc(np1)
		for i := pp; i > 0; i-- {
			runes[i], runes[i-1] = runes[i-1], runes[i]
			rc(np1)
		}

		r := runes[0]
		copy(runes, runes[1:pp+1])
		runes[pp] = r
	}
	rc(len(runes))

	words := make([]string, 0, len(wordMap))
	for k := range wordMap {
		words = append(words, k)
	}

	sort.Strings(words)

	return words
}

// Flags holds all of the command line arguments this app provides.
type Flags struct {
	Word           string
	DictionaryPath string
}

func main() {
	log.SetOutput(os.Stderr)
	log.SetPrefix("perms ")
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	f := Flags{}

	flag.StringVar(&f.DictionaryPath, "d", "", "Path to dictionary file (optional, defaults to embedded english dictionary)")
	flag.StringVar(&f.Word, "w", "", "Word to process (required)")
	flag.Parse()

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s <options>\n\nOptions:\n", os.Args[0])
		flag.PrintDefaults()
	}

	f.Word = strings.TrimSpace(f.Word)
	if f.Word == "" {
		flag.Usage()
		return
	}

	dictionary := defaultDictionary

	if f.DictionaryPath != "" {
		loadedDictionary, err := loadDictionary(f.DictionaryPath)
		if err != nil {
			log.Fatal(err)
		}

		dictionary = loadedDictionary
	}

	list := permutations(f.Word, dictionary)
	for _, word := range list {
		fmt.Println(word)
	}
}
