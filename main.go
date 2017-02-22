package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"runtime/pprof"
	"strings"
	"sync"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

var (
	wordsFile  = flag.String("words_file", "/usr/share/dict/words", "File containing valid words")
	numLetters = flag.Int("num_letters", 3, "Number of letters in resulting puzzles")
	parallel   = flag.Int("parallel", 4, "Number of goroutines to use to match words")

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

	shuffled := make(chan string)
	go shuffle(strings, shuffled)

	// Start the goroutine to consume puzzles and write files, block exit on it
	// completing.
	puzzles := make(chan puzzle)
	var wg2 sync.WaitGroup
	wg2.Add(1)
	go func() {
		defer wg2.Done()
		writePuzzles(puzzles)
	}()

	// Start N goroutines to consume shuffled strings and generate puzzles, block
	// closing puzzles chan until all are done, which will cause writePuzzles to
	// finish.
	var wg sync.WaitGroup
	for i := 0; i < *parallel; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			matchWords(shuffled, puzzles)
		}()
	}
	wg.Wait()
	close(puzzles)

	wg2.Wait()
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

// shuffle emits shuffled versions of the string.
//
// If s is "abcdefg", out will be sent:
// - abcdefg
// - bcdefga
// - cdefgab
// - defgabc
// - efgabcd
// - fgabcde
// - gabcdef
func shuffle(in <-chan string, out chan<- string) {
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

// matchWords emits all words that match in (with spelling bee semantics).
func matchWords(in <-chan string, out chan<- puzzle) {
	// TODO: Don't duplicate this work for each matcher goroutine.
	f, err := os.Open(*wordsFile)
	if err != nil {
		log.Fatalf("Open(%q): %v", *wordsFile, err)
	}
	r := bufio.NewReader(f)
	validRE := regexp.MustCompile("^([a-z]+)$")
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
		if !validRE.MatchString(w) {
			continue
		}

		allWords = append(allWords, w)
	}
	f.Close()
	fmt.Println("Matching", len(allWords), "words")

	for s := range in {
		re := regexp.MustCompile(fmt.Sprintf("^([%s]+)$", s))
		words := []string{}
		for _, word := range allWords {
			// Words must contain the first character.
			if !strings.Contains(word, string(s[0])) {
				continue
			}

			if re.MatchString(word) {
				words = append(words, word)
			}
		}

		// This combination of letters doesn't produce enough answers.
		if len(words) < 10 {
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
	for p := range in {
		fn := p.letters + ".txt"
		f, err := os.Create(fn)
		if err != nil {
			log.Fatalf("Create(%q): %v", fn, err)
		}
		for _, w := range p.words {
			fmt.Fprintln(f, w)
		}
		fmt.Fprintln(f, p.maxPts)
		fmt.Printf("Wrote %s\n", fn)
		f.Close()
	}
}
