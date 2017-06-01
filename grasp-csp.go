
package main

import (
	"fmt"
	"sort"
	"math/rand"
	"os"
	"bufio"
	"log"
	"strconv"
	"math"

	"github.com/xrash/smetrics"
)

type CountAndSortResult struct {
	char byte
	count int
}

type LocalSearchResult struct {
	str string
	cost int
}

/*
res: ccc

abc
acc
cdd
dda

[a] = 2
[c] = 1
[d] = 1

reverse
[2] = [ a ]
[1] = [c, d]
[3] = [c, d]

sorting
[2, 1, 3]

sorted
[3, 2, 1]

most common chars
[ {a, 2}, {c, 1}, {d, 1} ]

aab
aac
*/
func CountAndSort(strings []string, index int) []CountAndSortResult {

	charFreq := make(map[byte]int)

	for _, s := range strings {
		charFreq[s[index]] += 1
	}

	reverseMap := make(map[int][]byte)
	var sortingArray []int
	for char, value := range charFreq {
		reverseMap[value] = append(reverseMap[value], char)
	}

	for value := range reverseMap {
		sortingArray = append(sortingArray, value)
	}

	sort.Sort(sort.Reverse(sort.IntSlice(sortingArray)))

	var mostCommonChars []CountAndSortResult
	for _, k := range sortingArray {
		for _, c := range reverseMap[k] {
			mostCommonChars = append(mostCommonChars, CountAndSortResult{c, k})
		}
	}

	// fmt.Println("mostComm:", mostCommonChars)
	return mostCommonChars
}

func Cost(value string, target string) int {
	v, _ := smetrics.Hamming(value, target)
	return v
}

func CostSum(value string, strings []string, size int) int {

	c := 0

	for _, target := range strings {
		v, _ := smetrics.Hamming(value, target[0:size])
		c += v
	}

	return c

}

func Construct(strings []string, alphabet []string, stringSize int, alpha float64) string {

	x := ""

	for i:=0; len(x) < stringSize; i++ {

		allCandidates := CountAndSort(strings, i)
		min := allCandidates[len(allCandidates)-1].count
		max := allCandidates[0].count
		threshold := float64(max - min) * alpha

		var rcl []byte
		for _,element := range allCandidates {

			if float64(element.count) > threshold {
				rcl = append(rcl, element.char)
			}

		}
		chosen := rand.Intn(len(rcl))
		x = x + string(rcl[chosen])

	}

	return x
}

func NeighborhoodRandom(value string, alphabet []string) string {
	pos := rand.Intn(len(value))
	newChar := alphabet[ rand.Intn(len(alphabet)) ]

	chars := []byte(value)
	chars[pos] = newChar[0]

	return string(chars)
}

func LocalSearchRandom(value string, currentCost int, strings []string, stringSize int, alphabet []string) (string,int) {

	NhbMax := 10
	for i:=0; i<NhbMax; i++ {
		nhbValue := NeighborhoodRandom(value, alphabet)
		nhbCost := CostSum(nhbValue, strings, stringSize)

		// cost check
		if nhbCost < currentCost {
			value = nhbValue
			currentCost = nhbCost
		}
	}

	return value,currentCost
}

func IncrementAlphabetValue(value string, position int, alphabet []string) string {

	valueChars := []byte(value)

	if position > len(value) {

		return string(valueChars)

	}

	alphabetPos := -1
	for i, char := range alphabet {
		if char[0] == valueChars[position] {
			alphabetPos = i
			break
		}
	}

	valueChars[position] = alphabet[ (alphabetPos + 1) % len(alphabet) ][0] // <- convert to byte

	return string(valueChars)
}

func LocalSearchIterative(value string, currentCost int, strings []string, stringSize int, alphabet []string) (string,int) {

	res := value
	cost := currentCost

	stop := false

	for stop == false {

		// initial: abc
		// neighbors:
		// bca
		// bcb
		// baa
		// bab
		// ...
		neighbor := IncrementAlphabetValue(res, 0, alphabet)
		for i:=1; i<stringSize; i++ {
			neighbor = IncrementAlphabetValue(neighbor, i, alphabet)
		}

		stop = true

		// total neighbors: (alphabetSize-1) ^ stringSize
		for i:=0; i<stringSize; i++ {

			for j:=0; j<len(alphabet)-1; j++ {
				newCost := CostSum(neighbor, strings, stringSize)
				if newCost < cost {
					res = neighbor
					cost = newCost
					stop = false
				}
				neighbor = IncrementAlphabetValue(neighbor, i, alphabet)
			}

		}

	}

	return res, cost

}

// Return closest string, lower bound, upper bound
func CSP(strings []string, alphabet []string, stringSize int, maxIterations int, alpha float64) (string,int,int) {

	cost := 0
	var res string
	for i:=0; i<maxIterations; i++ {
		// construção
		res = Construct(strings, alphabet, stringSize, alpha)
		cost = CostSum(res, strings, stringSize)

		// busca local
		// res,cost = LocalSearchRandom(res, cost, strings, stringSize, alphabet)
		res,cost = LocalSearchIterative(res, cost, strings, stringSize, alphabet)

	}

	lower := math.MaxInt32
	upper := 0
	var dist int

	for _, target := range strings {
		dist = Cost(res, target)

		if dist < lower {
			lower = dist
		}

		if lower > upper {
			upper = dist
		}
	}

	return res,lower,upper
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s [inputfile]\n", os.Args[0])
	os.Exit(2)
}

func main() {

	maxIterations := 10
	greedinessFactor := 0.5

	if len(os.Args) < 2 {
		usage()
	}
	inputFilename := os.Args[1]
	inputFile, err := os.Open(inputFilename)
    if err != nil {
        log.Fatal(err)
    }
    defer inputFile.Close()

	var alphabetSize int
	scanner := bufio.NewScanner(inputFile)
    if scanner.Scan() {
        alphabetSize, _ = strconv.Atoi(scanner.Text())
		// fmt.Println("Size:", alphabetSize)
    }

	var stringsSize int
    if scanner.Scan() {
        stringsSize, _ = strconv.Atoi(scanner.Text())
		// fmt.Println("Size:", stringsSize)
    }

	var stringsLength int
    if scanner.Scan() {
        stringsLength, _ = strconv.Atoi(scanner.Text())
		// fmt.Println("Size:", stringsLength)
    }

	alphabet := make([]string, alphabetSize)
	for i:=0; i<alphabetSize; i++ {
		if scanner.Scan() {
			alphabet[i] = scanner.Text()
		}
	}
	// fmt.Println("Alphabet:", alphabet)

	strings := make([]string, stringsSize)
	for i:=0; i<stringsSize; i++ {
		if scanner.Scan() {
			strings[i] = scanner.Text()
		}
	}
	// fmt.Println("Strings:", strings)

	_, lower, upper := CSP(strings, alphabet, stringsLength, maxIterations, greedinessFactor)

	// fmt.Println("Closest:", closest, "lower:", lower, "upper:", upper)
	fmt.Println(lower, upper)

}