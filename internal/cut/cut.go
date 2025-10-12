package cut

import (
	"bufio"
	"errors"
	"os"
	"strconv"
	"strings"
)

type Options struct {
	Fields    string
	Delimiter string
	Separated bool
}

func GetResultChannel(opt Options) (<-chan []string, error) {
	result := make(chan []string, 1)

	fieldsSplitted := strings.Split(opt.Fields, ",")
	ranges := make([][]int, 0)
	for _, val := range fieldsSplitted {
		val = strings.TrimSpace(val)
		if val == "" {
			continue
		}

		if strings.Contains(val, "-") {
			splittedVal := strings.Split(val, "-")
			if len(splittedVal) != 2 {
				return result, errors.New("invalid field range")
			}
			val1, err1 := strconv.Atoi(splittedVal[0])
			val2, err2 := strconv.Atoi(splittedVal[1])
			if (err1 != nil) || (err2 != nil) || val1 <= 0 || val2 < val1 {
				return result, errors.New("something wrong with ranges")
			}
			ranges = append(ranges, []int{val1, val2})
		} else {
			val1, err := strconv.Atoi(val)
			if err != nil || val1 <= 0 {
				return result, errors.New("something wrong with ranges")
			}
			ranges = append(ranges, []int{val1})
		}
	}

	go func() {
		defer close(result)
		scanner := bufio.NewScanner(os.Stdin)

		for scanner.Scan() {
			rawLine := strings.TrimRight(scanner.Text(), "\r")
			hasDelimiter := strings.Contains(rawLine, opt.Delimiter)

			if !hasDelimiter {
				if opt.Separated {
					continue
				}
				result <- []string{rawLine}
				continue
			}

			currentSplittedLine := strings.Split(rawLine, opt.Delimiter)
			resultElem := make([]string, 0)

			if opt.Fields == "" {
				result <- currentSplittedLine
				continue
			}

			for _, val := range ranges {
				if len(val) == 1 {
					idx := val[0] - 1
					if idx >= 0 && idx < len(currentSplittedLine) {
						resultElem = append(resultElem, currentSplittedLine[idx])
					}
				} else {
					for i := val[0] - 1; i < val[1]; i++ {
						if i >= len(currentSplittedLine) {
							break
						}
						resultElem = append(resultElem, currentSplittedLine[i])
					}
				}
			}

			if len(resultElem) > 0 {
				result <- resultElem
			}
		}
	}()

	return result, nil
}
