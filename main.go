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
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

var (
	wordsFile  = flag.String("words_file", "/usr/share/dict/words", "File containing words")
	numLetters = flag.Int("num_letters", 3, "Number of letters")

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

	strings := make(chan string, 1000)
	go genAllStrings(*numLetters, strings)

	shuffled := make(chan string, 1000)
	go shuffle(strings, shuffled)

	puzzles := make(chan puzzle, 1000)
	go matchWords(shuffled, puzzles)

	writePuzzles(puzzles)
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
	close(out)
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
