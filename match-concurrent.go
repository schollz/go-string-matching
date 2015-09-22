package main

import (
	"bufio"
	"fmt"
	"github.com/arbovm/levenshtein"
	"os"
	"strings"
	"sync"
	"runtime"
	"strconv"
)

//GLOBALS
var findings_matches []string
var findings_leven []int
var wg sync.WaitGroup

func abs(x int) int {
	if x < 0 {
		return -x
	} else if x == 0 {
		return 0 // return correctly abs(-0)
	}
	return x
}

func generateHash(path string) {
	inFile, _ := os.Open(path)
	defer inFile.Close()
	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)
	
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		s := strings.Replace(scanner.Text(), "/", "", -1)
		addToCache("keys.list",s)
		partials, num := getPartials(s)
		for i := 0; i < num; i++ {
			addToCache(partials[i], strconv.Itoa(lineNum))
		}
	}
}

func addToCache(spartial string, s string) {
	f, err := os.OpenFile("cache/"+spartial, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	if _, err = f.WriteString(s + "\n"); err != nil {
		panic(err)
	}
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
			break
		}
	}
	return false
}

func getPartials(s string) ([]string, int) {
	partials := make([]string, 100000)
	num := 0
	s = strings.ToLower(s)
	s = strings.Replace(s, "/", "", -1)
	s = strings.Replace(s, " ", "", -1)
	s = strings.Replace(s, "the", "", -1)
	s = strings.Replace(s, "by", "", -1)
	s = strings.Replace(s, "dr", "", -1)
	s = strings.Replace(s, "of", "", -1)
	slen := len(s)
	if slen <= 6 {
		partials[num] = "asdf"
		num = num + 1
	} else {
		for i := 0; i <= slen-6; i++ {
			partials[num] = s[i : i+6]
			num = num + 1
		}
	}
	return partials, num
}

func getMatch(s string, path string) (string, int) {
	partials, num := getPartials(s)
	matches := make([]string, 1000000)
	numm := 0
	runtime.GOMAXPROCS(8)
	N := 8
	for i := 0; i < num; i++ {

		inFile, _ := os.Open(path + partials[i])
		defer inFile.Close()
		scanner := bufio.NewScanner(inFile)
		scanner.Split(bufio.ScanLines)

		for scanner.Scan() {
			//if stringInSlice(scanner.Text(),matches) == false { ITS NOT WORTH LOOKING THROUGH DUPLICATES
			matches[numm] = scanner.Text()
			numm = numm + 1
			// }
		}

	}

	matches2 := make([]string, numm)
	matches2 = matches[0:numm]
	findings_leven = make([]int, N)
	findings_matches = make([]string, N)

	wg.Add(N)
	for i := 0; i < N; i++ {
		go search(matches2[i*len(matches2)/N:(i+1)*len(matches2)/N], s, i)
	}
	wg.Wait()

	lowest := 100
	best_index := 0
	for i := 0; i < len(findings_leven); i++ {
		if findings_leven[i] < lowest {
			lowest = findings_leven[i]
			best_index = i
		}
	}

	return findings_matches[best_index],lowest
}

func search(matches []string, target string, process int) {
	defer wg.Done()
	match := "No match"
	target = strings.ToLower(target)
	bestLevenshtein := 1000
	for i := 0; i < len(matches); i++ {
		d := levenshtein.Distance(target, strings.ToLower(matches[i]))
		if d < bestLevenshtein {
			bestLevenshtein = d
			match = matches[i]
		}
	}
	findings_matches[process] = match
	findings_leven[process] = bestLevenshtein
}

func main() {
	if strings.EqualFold(os.Args[1], "help") {
		fmt.Println("Version 1.0 - 4-mer tuples\n")
		fmt.Println("./match-concurrent build <NAME OF WORDLIST> - builds cache/ folder in current directory\n")
		fmt.Println("./match-concurrent 'word or words to match' /directions/to/cache/\n")
	} else if strings.EqualFold(os.Args[1], "build") {
		os.Mkdir("cache", 0775)
		generateHash(os.Args[2])
	} else {
		match,lowest := getMatch(os.Args[1],os.Args[2])
		fmt.Printf("%v|||%v\n", match,lowest)
	}
}
