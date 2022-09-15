package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
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
	finished := make(chan bool, len(ranges))
	timeLimit := time.After(time.Duration(timeout) * time.Second)

	file, err := os.OpenFile(fileName+".txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND|os.O_TRUNC, 0o777)
	if err != nil {
		log.Fatal("File creating error")
	}

	go finishProgram(finishedAll, finished)

	for _, r := range ranges {
		go HandleRange(r, finished, file)
	}
	select {
	case <-timeLimit:
		fmt.Println("Time limit")
		file.Close()
	case <-finishedAll:
		fmt.Println("Strings are ready")
		file.Close()
	}
}

func finishProgram(results chan bool, input chan bool) {
	limit := cap(input)
	for {
		select {
		case <-input:
			limit -= 1
			if limit == 0 {
				results <- true
			}
		default:
			continue
		}
	}
}

func HandleRange(rangeString string, finished chan bool, file *os.File) {
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
	finished <- true
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
