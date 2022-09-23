package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var fileName string
var timeout int

type rangeArray []string

func (i *rangeArray) String() string {
	return "1:100"
}

func (i *rangeArray) Set(value string) error {
	*i = append(*i, value)
	return nil
}

var ranges rangeArray

// go run . --file res --timeout 15 --range 1:23 --range 90:1050000
func main() {
	flag.StringVar(&fileName, "file", "default.txt", "File name for output")
	flag.IntVar(&timeout, "timeout", 10, "Program allowed time")
	flag.Var(&ranges, "range", "Range elemets for simple numbers search")
	flag.Parse()
	finishedAll := make(chan bool)
	timeLimit := time.After(time.Duration(timeout) * time.Second)

	wg := new(sync.WaitGroup)

	file, err := os.OpenFile(fileName+".txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND|os.O_TRUNC, 0o777)
	if err != nil {
		log.Fatal("File creating error")
	}
	defer file.Close()

	for _, r := range ranges {
		wg.Add(1)
		go HandleRange(wg, r, file)
	}

	go finishProgram(wg, finishedAll)

	select {
	case <-timeLimit:
		fmt.Println("Time limit")
	case <-finishedAll:
		fmt.Println("Strings are ready")
	}
}

func finishProgram(wg *sync.WaitGroup, finishedAll chan bool) {
	wg.Wait()
	finishedAll <- true
}

func HandleRange(wg *sync.WaitGroup, rangeString string, file *os.File) {
	resChan := make(chan []byte)
	rangeArray := strings.Split(rangeString, ":")
	minValue, err := strconv.Atoi(rangeArray[0])
	if err != nil {
		log.Fatal("First number in range isn't a number")
	}
	maxValue, err := strconv.ParseUint(rangeArray[1], 10, 64)
	if err != nil {
		log.Fatal("Second number in range isn't a number")
	}

	go findNumbers(minValue, maxValue, resChan)

	res := <-resChan
	file.WriteString(fmt.Sprintf("For range %v answer is: %v", rangeString, string(res)+"\n"))
	wg.Done()
}

func findNumbers(minValue int, maxValue uint64, resArray chan []byte) {
	simpleNumbers := make([]uint64, 0)
	filteredNumbers := make([]uint64, 0)

	for i := uint64(0); i <= maxValue; i++ {
		simpleNumbers = append(simpleNumbers, i)
	}
	simpleNumbers[1] = 0

	for i := uint64(2); i <= maxValue; i++ {
		if simpleNumbers[i] != 0 {
			j := i + i
			for j <= maxValue {
				simpleNumbers[j] = 0
				j += i
			}
		}
	}
	for _, el := range simpleNumbers {
		if el != 0 && el >= uint64(minValue) {
			filteredNumbers = append(filteredNumbers, el)
		}
	}

	resArray <- []byte(strings.Trim(strings.Replace(fmt.Sprint(filteredNumbers), " ", ", ", -1), "[]"))
}
