package cut

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Options struct {
	Fields    string
	Delimiter string
	Separated bool
}

func GetResultChannel(opt Options) (<-chan []string, <-chan error) {
	result := make(chan []string, 1)
	errChan := make(chan error, 1)
	defer close(result)

	fieldsSplitted := strings.Split(opt.Fields, ",")
	ranges := make([][]int, len(fieldsSplitted))
	for _, val := range fieldsSplitted {
		if strings.Contains(val, "-") {

			splittedVal := strings.Split(val, "-")
			val1, err1 := strconv.Atoi(splittedVal[0])
			val2, err2 := strconv.Atoi(splittedVal[1])
			if (err1 != nil) || (err2 != nil) {
				errChan <- errors.New("something wrong with ranges")
				return result, errChan
			}
			currRange := make([]int, 2)
			currRange[0] = val1
			currRange[1] = val2
			ranges = append(ranges, currRange)
		} else {
			currRange := make([]int, 1)
			val1, err := strconv.Atoi(val)
			if err != nil {
				errChan <- errors.New("something wrong with ranges")
				return result, errChan
			}
			currRange[0] = val1
			ranges = append(ranges, currRange)
		}
	}

	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			rawLine := strings.TrimRight(scanner.Text(), "\r")
			if opt.Separated {
				if strings.Contains(rawLine, opt.Delimiter) {
					if opt.Fields == "" {
						result <- strings.Split(rawLine, opt.Delimiter)
					} else {
						currentSplittedLine := strings.Split(rawLine, opt.Delimiter)
						resultElem := make([]string, 0)
						for _, val := range ranges {
							if len(val) == 1 {
								if val[0] < 0 {
									errChan <- errors.New("index cannot be negative")
									return
								}
								if val[0] > (len(currentSplittedLine) - 1) {
									errChan <- errors.New("the index is out of bounds: " + fmt.Sprint(val[0]))
									return
								}
								resultElem = append(resultElem, currentSplittedLine[val[0]])
							} else {
								if (val[0] > (len(currentSplittedLine) - 1)) || (val[1] >= (len(currentSplittedLine))) {
									errChan <- errors.New("the index is out of bounds: " + fmt.Sprint(val[0]) + ":" + fmt.Sprint(val[1]))
									return
								}
								if val[0] > val[1] {
									errChan <- errors.New("wrong start index: " + fmt.Sprint(val[0]) + ":" + fmt.Sprint(val[1]))
									return
								}
								if (val[0] < 0) || (val[1] < 0) {
									errChan <- errors.New("index cannot be negative")
									return
								}
								for i := val[0]; i <= val[1]; i++ {
									resultElem = append(resultElem, currentSplittedLine[i])
								}
							}
						}
						result <- resultElem
					}
				}
			} else {
				if opt.Fields == "" {
					result <- strings.Split(rawLine, opt.Delimiter)
				} else {
					currentSplittedLine := strings.Split(rawLine, opt.Delimiter)
					resultElem := make([]string, 0)
					for _, val := range ranges {
						if len(val) == 1 {
							if val[0] < 0 {
								errChan <- errors.New("index cannot be negative")
								return
							}
							if val[0] > (len(currentSplittedLine) - 1) {
								errChan <- errors.New("the index is out of bounds: " + fmt.Sprint(val[0]))
								return
							}
							resultElem = append(resultElem, currentSplittedLine[val[0]])
						} else {
							if (val[0] > (len(currentSplittedLine) - 1)) || (val[1] >= (len(currentSplittedLine))) {
								errChan <- errors.New("the index is out of bounds: " + fmt.Sprint(val[0]) + ":" + fmt.Sprint(val[1]))
								return
							}
							if val[0] > val[1] {
								errChan <- errors.New("wrong start index: " + fmt.Sprint(val[0]) + ":" + fmt.Sprint(val[1]))
								return
							}
							if (val[0] < 0) || (val[1] < 0) {
								errChan <- errors.New("index cannot be negative")
								return
							}
							for i := val[0]; i <= val[1]; i++ {
								resultElem = append(resultElem, currentSplittedLine[i])
							}
						}
					}
					result <- resultElem
				}
			}
		}
	}()

	return result, errChan
}
