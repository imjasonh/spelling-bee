package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"runtime/pprof"
	"strings"
	"sync"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

var (
	wordsFile  = flag.String("words_file", "/usr/share/dict/words", "File containing valid words")
	numLetters = flag.Int("num_letters", 3, "Number of letters in resulting puzzles")
	parallel   = flag.Int("parallel", 3, "Number of goroutines to use to generate puzzles")
	v          = flag.Bool("v", false, "verbose logging")

	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
)

func main() {
	flag.Parse()

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	strings := make(chan string)
	go genAllStrings(*numLetters, strings)

	rotated := make(chan string)
	go rotate(strings, rotated)

	puzzles := make(chan puzzle)

	// Consume puzzles and write files.
	var wg2 sync.WaitGroup
	wg2.Add(1)
	go func() {
		defer wg2.Done()
		writePuzzles(puzzles)
	}()

	allWords := genAllWords()

	// Consume rotated words and generate puzzles.
	var wg sync.WaitGroup
	for i := 0; i < *parallel; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			matchWords(allWords, rotated, puzzles)
		}()
	}
	wg.Wait()
	// When puzzle generators are done, close puzzles. This will cause
	// writePuzzles to finish, and the program to exit.
	close(puzzles)

	wg2.Wait()
}

func genAllWords() []string {
	f, err := os.Open(*wordsFile)
	if err != nil {
		log.Fatalf("Open(%q): %v", *wordsFile, err)
	}
	r := bufio.NewReader(f)
	allWords := []string{}
	for {
		l, err := r.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("ReadBytes: %v", err)
		}
		w := string(l)
		w = strings.TrimSpace(w)
		// Words must be >5 letters.
		if len(w) < 5 {
			continue
		}
		// Words must be lowercase, no punctuation.
		if !containsOnly(w, alphabet) {
			continue
		}

		allWords = append(allWords, w)
	}
	f.Close()
	fmt.Println("Matching", len(allWords), "words")
	return allWords
}

// genAllStrings generates all unique strings of length n and sends them to
// out.
func genAllStrings(n int, out chan<- string) {
	for _, c := range alphabet {
		if n == 1 {
			out <- string(c)
			continue
		}

		ch := make(chan string, 1000)
		go genAllStrings(n-1, ch)
		for rest := range ch {
			if rest[0] > byte(c) {
				out <- string(c) + rest
			}
		}
	}
	close(out)
}

// rotate emits rotated versions of the string.
//
// If s is "abcdefg", out will be sent:
// - abcdefg
// - bcdefga
// - cdefgab
// - defgabc
// - efgabcd
// - fgabcde
// - gabcdef
func rotate(in <-chan string, out chan<- string) {
	for s := range in {
		for i := 0; i < len(s); i++ {
			first, rest := s[:i], s[i:]
			out <- rest + first
		}
	}
	close(out)
}

type puzzle struct {
	letters string
	words   []string
	maxPts  int
}

func containsOnly(s, target string) bool {
	rs := map[rune]struct{}{}
	for _, r := range target {
		rs[r] = struct{}{}
	}
	for _, r := range s {
		if _, found := rs[r]; !found {
			return false
		}
	}
	return true
}

// matchWords emits all words that match in (with spelling bee semantics).
func matchWords(allWords []string, in <-chan string, out chan<- puzzle) {

	for s := range in {
		runes := map[rune]struct{}{}
		for _, c := range s {
			runes[c] = struct{}{}
		}

		words := []string{}
		for _, word := range allWords {
			// Words must contain the first character.
			if !strings.Contains(word, string(s[0])) {
				continue
			}

			// Words must contain only letters in this set.
			if containsOnly(word, s) {
				words = append(words, word)
			}
		}

		// This combination of letters doesn't produce enough answers.
		if len(words) < 10 {
			if *v {
				fmt.Print("1")
			}
			continue
		}

		// Score the puzzle and ensure at least one answer uses all letters.
		someContainsAll := false
		maxPts := 0
		for _, w := range words {
			containsAll := true
			for _, let := range s {
				if !strings.ContainsRune(w, let) {
					containsAll = false
				}
			}
			if containsAll {
				maxPts += 3
				someContainsAll = true
			} else {
				maxPts += 1
			}
		}
		if !someContainsAll {
			if *v {
				fmt.Print("2")
			}
			continue
		}

		out <- puzzle{
			letters: s,
			words:   words,
			maxPts:  maxPts,
		}
	}
}

func writePuzzles(in <-chan puzzle) {
	t := time.Tick(time.Second)
	for {
		select {
		case p, ok := <-in:
			if !ok {
				return
			}
			fn := p.letters + ".txt"
			f, err := os.Create(fn)
			if err != nil {
				log.Fatalf("Create(%q): %v", fn, err)
			}
			for _, w := range p.words {
				fmt.Fprintln(f, w)
			}
			fmt.Fprintln(f, p.maxPts)
			if *v {
				fmt.Println("wrote", p.letters)
			}
			f.Close()
		case <-t:
			if *v {
				fmt.Print("%")
			}
		}
	}
}
