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
	defer close(result)

	fieldsSplitted := strings.Split(opt.Fields, ",")
	ranges := make([][]int, len(fieldsSplitted))
	for _, val := range fieldsSplitted {
		if strings.Contains(val, "-") {

			splittedVal := strings.Split(val, "-")
			val1, err1 := strconv.Atoi(splittedVal[0])
			val2, err2 := strconv.Atoi(splittedVal[1])
			if (err1 != nil) || (err2 != nil) {
				return result, errors.New("something wrong with ranges")
			}
			currRange := make([]int, 2)
			currRange[0] = val1
			currRange[1] = val2
			ranges = append(ranges, currRange)
		} else {
			currRange := make([]int, 1)
			val1, err := strconv.Atoi(val)
			if err != nil {
				return result, errors.New("something wrong with ranges")
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
								if val[0] <= (len(currentSplittedLine) - 1) {
									resultElem = append(resultElem, currentSplittedLine[val[0]])
								}

							} else {
								for i := val[0]; i <= val[1]; i++ {
									if i > (len(currentSplittedLine) - 1) {
										break
									} else {
										resultElem = append(resultElem, currentSplittedLine[i])
									}

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
							if val[0] <= (len(currentSplittedLine) - 1) {
								resultElem = append(resultElem, currentSplittedLine[val[0]])
							}

						} else {
							for i := val[0]; i <= val[1]; i++ {
								if i > (len(currentSplittedLine) - 1) {
									break
								} else {
									resultElem = append(resultElem, currentSplittedLine[i])
								}

							}
						}
					}
					result <- resultElem
				}
			}
		}
	}()

	return result, nil
}
