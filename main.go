package main

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"math/big"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	day = kingpin.Arg("day", "Advent of code day to run").Int()
)

func main() {
	kingpin.Version("0.0.1")
	kingpin.Parse()
	fmt.Printf("Would run day %d\n", *day)
	switch *day {
	case 1:
		day1()
	case 2:
		day2()
	case 3:
		day3()
	case 4:
		day4()
	case 5:
		day5()
	case 6:
		day6()
	case 7:
		day7()
	case 8:
		day8()
	case 9:
		day9()
	case 10:
		day10()
	case 11:
		day11()
	case 12:
		day12()
	case 13:
		day13()
	case 14:
		day14()
	case 15:
		day15()
	case 16:
		day16()
	case 17:
		day17()
	case 18:
		day18()
	case 19:
		day19()
	case 20:
		day20()
	case 21:
		day21()
	case 22:
		day22()
	case 23:
		day23()
	case 24:
		day24()
	case 25:
		day25()
	default:
		fmt.Println("We don't have that day...")
	}
}

func getFuelForFuel(initialFuelAmmount int) (total int) {
	nextAmount := initialFuelAmmount
	for nextAmount > 0 {
		nextAmount = (nextAmount / 3) - 2
		if nextAmount > 0 {
			total += nextAmount
		}
	}
	return total

}

func day1() {
	//data is at the bottom of the file
	total := 0
	dataArray := strings.Split(day1Data, "\n")
	for _, item := range dataArray {
		fuelValue, err := strconv.Atoi(item)
		if err != nil {
			fmt.Println(err)
		}
		total += (fuelValue / 3) - 2
	}
	fmt.Println("Fuel for modules: ", strconv.Itoa(total))
	for _, item := range dataArray {
		moduleMass, _ := strconv.Atoi(item)
		fuelMass := (moduleMass / 3) - 2
		if fuelMass > 0 {
			total += getFuelForFuel(fuelMass)
		}
	}
	fmt.Println("Fuel for modules with the fuel for fuel: ", strconv.Itoa(total))

}

func makeMapForArray(arr []int64) map[int64]int64 {
	theMap := make(map[int64]int64)
	for index, val := range arr {
		theMap[int64(index)] = val
	}
	return theMap
}

func getParamWithMode(data map[int64]int64, mode int64, location int64, relativeBase int64) (paramVal int64) {
	switch mode {
	case 0:
		return data[data[location]]
	case 1:
		return data[location]
	case 2:
		return data[data[location]+relativeBase]
	default:
		fmt.Printf("No mode with value %d when getting location %d", mode, location)
		return -1
	}
}

func setParamWithMode(data map[int64]int64, mode int64, location int64, value int64, relativeBase int64) (paramVal int64) {
	switch mode {
	case 0:
		data[data[location]] = value
	case 1:
		data[location] = value
	case 2:
		data[data[location]+relativeBase] = value
	default:
		fmt.Printf("Cant set with mode of value %d when setting location %d", mode, location)
		return -1
	}
	return 0
}

type inputHandler func() int64
type outputHandler func(int64)

//Intcode computer
func intComp(intCodeData map[int64]int64, getInput inputHandler, sendOutput outputHandler) (int64, map[int64]int64) {
	var relativeBase int64 = 0
	var copiedData = make(map[int64]int64)
	var key, val int64
	for key, val = range intCodeData {
		copiedData[key] = val
	}
	var currentPos int64 = 0
	instruction := copiedData[currentPos]
	running := true

	var firstVal int64
	var secondVal int64

	for running {
		currentOpCode := instruction % 100
		mode1 := (instruction / 100) % 10
		mode2 := (instruction / 1000) % 10
		mode3 := (instruction / 10000) % 10
		var result int64

		if currentOpCode == 1 {
			firstVal = getParamWithMode(copiedData, mode1, currentPos+1, relativeBase)
			secondVal = getParamWithMode(copiedData, mode2, currentPos+2, relativeBase)
			result = setParamWithMode(copiedData, mode3, currentPos+3, firstVal+secondVal, relativeBase)
			if result == -1 {
				running = false
			}
			currentPos += 4
		}

		if currentOpCode == 2 {
			firstVal = getParamWithMode(copiedData, mode1, currentPos+1, relativeBase)
			secondVal = getParamWithMode(copiedData, mode2, currentPos+2, relativeBase)
			result = setParamWithMode(copiedData, mode3, currentPos+3, firstVal*secondVal, relativeBase)
			if result == -1 {
				running = false
			}
			currentPos += 4
		}

		if currentOpCode == 3 {
			var input int64 = getInput()
			result = setParamWithMode(copiedData, mode1, currentPos+1, input, relativeBase)
			if result == -1 {
				running = false
			}
			currentPos += 2
		}

		if currentOpCode == 4 {
			output := getParamWithMode(copiedData, mode1, currentPos+1, relativeBase)
			sendOutput(output)
			currentPos += 2
		}

		if currentOpCode == 5 {
			firstVal = getParamWithMode(copiedData, mode1, currentPos+1, relativeBase)
			secondVal = getParamWithMode(copiedData, mode2, currentPos+2, relativeBase)
			if firstVal != 0 {
				currentPos = secondVal
			} else {
				currentPos += 3
			}
		}

		if currentOpCode == 6 {
			firstVal = getParamWithMode(copiedData, mode1, currentPos+1, relativeBase)
			secondVal = getParamWithMode(copiedData, mode2, currentPos+2, relativeBase)
			if firstVal == 0 {
				currentPos = secondVal
			} else {
				currentPos += 3
			}
		}

		if currentOpCode == 7 {
			firstVal = getParamWithMode(copiedData, mode1, currentPos+1, relativeBase)
			secondVal = getParamWithMode(copiedData, mode2, currentPos+2, relativeBase)
			if firstVal < secondVal {
				result = setParamWithMode(copiedData, mode3, currentPos+3, 1, relativeBase)
			} else {
				result = setParamWithMode(copiedData, mode3, currentPos+3, 0, relativeBase)
			}

			if result == -1 {
				running = false
			}
			currentPos += 4
		}

		if currentOpCode == 8 {
			firstVal = getParamWithMode(copiedData, mode1, currentPos+1, relativeBase)
			secondVal = getParamWithMode(copiedData, mode2, currentPos+2, relativeBase)
			if firstVal == secondVal {
				result = setParamWithMode(copiedData, mode3, currentPos+3, 1, relativeBase)
			} else {
				result = setParamWithMode(copiedData, mode3, currentPos+3, 0, relativeBase)
			}

			if result == -1 {
				running = false
			}
			currentPos += 4
		}

		if currentOpCode == 9 {
			relativeBase += getParamWithMode(copiedData, mode1, currentPos+1, relativeBase)
			currentPos += 2
		}

		if currentOpCode == 99 {
			running = false
		}

		instruction = copiedData[currentPos]
	}
	return copiedData[0], copiedData
}

func getUserInput() int64 {
	fmt.Println("Need number: ")
	var input int64
	_, _ = fmt.Scanf("%d", &input)
	// TODO: Stop people from inputing non-Numbers
	return input
}

func printOutput(output int64) {
	fmt.Println("Output:", output)
}

func day2() {
	var intCodeData = []int64{1, 12, 2, 3, 1, 1, 2, 3, 1, 3, 4, 3, 1, 5, 0, 3, 2, 13, 1, 19, 1, 5, 19, 23, 2, 10, 23, 27, 1, 27, 5, 31, 2, 9, 31, 35, 1, 35, 5, 39, 2, 6, 39, 43, 1, 43, 5, 47, 2, 47, 10, 51, 2, 51, 6, 55, 1, 5, 55, 59, 2, 10, 59, 63, 1, 63, 6, 67, 2, 67, 6, 71, 1, 71, 5, 75, 1, 13, 75, 79, 1, 6, 79, 83, 2, 83, 13, 87, 1, 87, 6, 91, 1, 10, 91, 95, 1, 95, 9, 99, 2, 99, 13, 103, 1, 103, 6, 107, 2, 107, 6, 111, 1, 111, 2, 115, 1, 115, 13, 0, 99, 2, 0, 14, 0}
	var intCodeMap = makeMapForArray(intCodeData)
	//Gets result for initial data
	intCodeDataResultAtZero, _ := intComp(intCodeMap, getUserInput, printOutput)
	fmt.Println("initial result at zero: ", intCodeDataResultAtZero)

	var noun, verb int64
	for noun = 1; noun <= 99; noun++ {
		for verb = 1; verb <= 99; verb++ {
			intCodeData[1] = noun
			intCodeData[2] = verb
			outForVerbNoun, _ := intComp(intCodeMap, getUserInput, printOutput)
			if outForVerbNoun == 19690720 {
				fmt.Println("Noun and verb for result 19690720:")
				fmt.Println("  noun: ", intCodeData[1])
				fmt.Println("  verb: ", intCodeData[2])
			}
		}
	}

}

type Point struct {
	x int64
	y int64
}

func Abs(x int64) int64 {
	if x < 0 {
		return x * -1
	}
	return x
}

func day3() {
	firstWireData := `R1003,U756,L776,U308,R718,D577,R902,D776,R760,U638,R289,D70,L885,U161,R807,D842,R175,D955,R643,U380,R329,U573,L944,D2,L807,D886,L549,U592,R152,D884,L761,D915,L726,D677,L417,D651,L108,D377,L699,D938,R555,D222,L689,D196,L454,U309,L470,D234,R198,U689,L996,U117,R208,D310,R572,D562,L207,U244,L769,U186,R153,D756,R97,D625,R686,U244,R348,U586,L385,D466,R483,U718,L892,D39,R692,U756,L724,U148,R70,U224,L837,D370,L309,U235,R382,D579,R404,D146,R6,U584,L840,D863,R942,U646,R146,D618,L12,U210,R126,U163,R931,D661,L710,D883,L686,D688,L148,D19,R703,U530,R889,U186,R779,D503,R417,U272,R541,U21,L562,D10,L349,U998,R69,D65,R560,D585,L949,D372,L110,D865,R212,U56,L936,U957,L88,U612,R927,U642,R416,U348,L541,D416,L808,D759,R449,D6,L517,D4,R494,D143,L536,U341,R394,U179,L22,D680,L138,U249,L285,U879,L717,U756,L313,U222,R823,D208,L134,U984,R282,U635,R207,D63,L416,U511,L179,D582,L651,U932,R646,U378,R263,U138,L920,U523,L859,D556,L277,D518,R489,U561,L457,D297,R72,U920,L583,U23,L395,D844,R776,D552,L55,D500,R111,U409,R685,D427,R275,U739,R181,U333,L215,U808,R341,D537,R336,U230,R247,U748,R846,U404,R850,D493,R891,U176,L744,U585,L987,D849,R271,D848,L555,U801,R316,U753,L390,U97,L128,U45,R706,U35,L928,U913,R537,D512,R152,D410,R76,D209,R183,U941,R289,U632,L923,D190,R488,D934,R442,D303,R178,D250,R204,U247,R707,U77,R428,D701,R386,U110,R641,U925,R703,D387,L946,U415,R461,D123,L214,U236,L959,U517,R957,D524,R812,D668,R369,U340,L606,D503,R755,U390,R142,D921,L976,D36,L965,D450,L722,D224,L303,U705,L584`
	secondWireData := `L993,U810,L931,D139,R114,D77,L75,U715,R540,D994,L866,U461,R340,D179,R314,D423,R629,D8,L692,U446,L88,D132,L128,U934,L465,D58,L696,D883,L955,D565,R424,U286,R403,U57,L627,D930,R887,D941,L306,D951,R918,U587,R939,U821,L65,D18,L987,D707,L360,D54,L932,U366,R625,U609,R173,D637,R661,U888,L68,U962,R270,U369,R780,U845,L813,U481,R66,D182,R420,U605,R880,D276,L6,D529,R883,U189,R380,D472,R30,U35,L510,D844,L146,U875,R152,U545,R274,U920,R432,U814,R583,D559,L820,U135,L353,U975,L103,U615,R401,U692,L676,D781,R551,D985,L317,U836,R115,D216,L967,U286,R681,U144,L354,U678,L893,D487,R664,D185,R787,D909,L582,D283,L519,D893,L56,U768,L345,D992,L248,U439,R573,D98,L390,D43,L470,D435,R176,U468,R688,U388,L377,U800,R187,U641,L268,U857,L716,D179,R212,U196,L342,U731,R261,D92,R183,D623,L589,D215,L966,U878,L784,U740,R823,D99,L167,D992,R414,U22,L27,U390,R286,D744,L360,U554,L756,U715,R939,D806,R279,U292,L960,U633,L428,U949,R90,D321,R749,U395,L392,U348,L33,D757,R289,D367,L562,D668,L79,D193,L991,D705,L562,U25,R146,D34,R325,U203,R403,D714,R607,U72,L444,D76,R267,U924,R289,U962,L159,U726,L57,D540,R299,U343,R936,U90,L311,U243,L415,D426,L936,D570,L539,D731,R367,D374,L56,D251,L265,U65,L14,D882,L956,U88,R688,D34,R866,U777,R342,D270,L344,D953,L438,D855,L587,U320,L953,D945,L473,U559,L487,D602,R255,U871,L854,U45,R705,D247,R955,U885,R657,D664,L360,D764,L549,D676,R85,U189,L951,D922,R511,D429,R37,U11,R821,U984,R825,U874,R753,D524,L537,U618,L919,D597,L364,D231,L258,U818,R406,D208,R214,U530,R261`

	firstWireDataArray := strings.Split(firstWireData, ",")
	secondWireDataArray := strings.Split(secondWireData, ",")
	var firstPointArray []Point
	var secondPointArray []Point
	var crossesArray []Point
	var crossesCounter int64 = 0
	var currentPosX int64 = 0
	var currentPosY int64 = 0
	var smollestManhattanDistance int64 = math.MaxInt64
	var smollestStepsToIntersection int = math.MaxInt64
	var i int64

	for _, wireString := range firstWireDataArray {
		dir := wireString[0]
		length, _ := strconv.ParseInt(wireString[1:], 10, 64)
		for i = 0; i < length; i++ {
			switch dir {
			case 'R':
				currentPosX++
				firstPointArray = append(firstPointArray, Point{currentPosX, currentPosY})
			case 'L':
				currentPosX--
				firstPointArray = append(firstPointArray, Point{currentPosX, currentPosY})
			case 'D':
				currentPosY++
				firstPointArray = append(firstPointArray, Point{currentPosX, currentPosY})
			case 'U':
				currentPosY--
				firstPointArray = append(firstPointArray, Point{currentPosX, currentPosY})
			}
		}
	}
	currentPosX = 0
	currentPosY = 0
	for _, wireString := range secondWireDataArray {
		dir := wireString[0]
		length, _ := strconv.ParseInt(wireString[1:], 10, 64)
		for i = 0; i < length; i++ {
			switch dir {
			case 'R':
				currentPosX++
				secondPointArray = append(secondPointArray, Point{currentPosX, currentPosY})
			case 'L':
				currentPosX--
				secondPointArray = append(secondPointArray, Point{currentPosX, currentPosY})
			case 'D':
				currentPosY++
				secondPointArray = append(secondPointArray, Point{currentPosX, currentPosY})
			case 'U':
				currentPosY--
				secondPointArray = append(secondPointArray, Point{currentPosX, currentPosY})
			}
		}
	}

	for firstSteps, firstPoint := range firstPointArray {
		for secondSteps, secondPoint := range secondPointArray {
			//fmt.Println("fx: ", firstPoint.x, " fy: ", firstPoint.y, " sx: ", secondPoint.x, " sy: ", secondPoint.y)
			if firstPoint == secondPoint {
				crossesCounter++
				crossesArray = append(crossesArray, firstPoint)
				if smollestStepsToIntersection > (firstSteps + secondSteps) {
					//indexes start at 0, but index 0 is 1 step for both of these, so Im adding two.
					smollestStepsToIntersection = (firstSteps + secondSteps) + 2
				}
			}
		}
	}

	fmt.Println("Cross points found: ", crossesCounter)
	for _, point := range crossesArray {
		//fmt.Println("x: ", point.x, " y: ", point.y)
		currentDistance := (Abs(point.x) + Abs(point.y))
		if smollestManhattanDistance > currentDistance {
			smollestManhattanDistance = currentDistance
		}
	}

	fmt.Println("closest point in manhatten distance", smollestManhattanDistance)
	fmt.Println("Fewest steps to intersection: ", smollestStepsToIntersection)
}

func getPassNumForRuleset(beginingRange int, endingRange int, largerGroupRule bool) int {
	currentNumPass := beginingRange
	passMeetsCriteriaCount := 0
	for currentNumPass <= endingRange {
		neverDecrease := true
		foundDouble := false
		strNumPass := strconv.Itoa(currentNumPass)
		for i := 0; i < len(strNumPass)-1; i++ {
			if strNumPass[i] > strNumPass[i+1] {
				neverDecrease = false
			}
			if largerGroupRule {
				if i < len(strNumPass)-2 {
					if i == 0 {
						if strNumPass[i] == strNumPass[i+1] && strNumPass[i+1] != strNumPass[i+2] {
							foundDouble = true
						}
					} else {
						if strNumPass[i] == strNumPass[i+1] && strNumPass[i+1] != strNumPass[i+2] &&
							strNumPass[i] != strNumPass[i-1] {
							foundDouble = true
						}
					}
				} else {
					if strNumPass[i] == strNumPass[i+1] && strNumPass[i] != strNumPass[i-1] {
						foundDouble = true
					}
				}
			} else {
				if strNumPass[i] == strNumPass[i+1] {
					foundDouble = true
				}
			}
		}

		if neverDecrease && foundDouble {
			passMeetsCriteriaCount++
		}
		currentNumPass++
	}
	return passMeetsCriteriaCount
}

func day4() {
	beginingRange := 256310
	endingRange := 732736
	passNumWithoutLargerGroup := getPassNumForRuleset(beginingRange, endingRange, false)
	fmt.Println("Pass without larger group rule: ", passNumWithoutLargerGroup)

	passNumWithLargerGroup := getPassNumForRuleset(beginingRange, endingRange, true)
	fmt.Println("Pass with larger group rule: ", passNumWithLargerGroup)

}

func day5() {
	inputData := []int64{3, 225, 1, 225, 6, 6, 1100, 1, 238, 225, 104, 0, 1101, 37, 61, 225, 101, 34, 121, 224, 1001, 224, -49, 224, 4, 224, 102, 8, 223, 223, 1001, 224, 6, 224, 1, 224, 223, 223, 1101, 67, 29, 225, 1, 14, 65, 224, 101, -124, 224, 224, 4, 224, 1002, 223, 8, 223, 101, 5, 224, 224, 1, 224, 223, 223, 1102, 63, 20, 225, 1102, 27, 15, 225, 1102, 18, 79, 224, 101, -1422, 224, 224, 4, 224, 102, 8, 223, 223, 1001, 224, 1, 224, 1, 223, 224, 223, 1102, 20, 44, 225, 1001, 69, 5, 224, 101, -32, 224, 224, 4, 224, 1002, 223, 8, 223, 101, 1, 224, 224, 1, 223, 224, 223, 1102, 15, 10, 225, 1101, 6, 70, 225, 102, 86, 40, 224, 101, -2494, 224, 224, 4, 224, 1002, 223, 8, 223, 101, 6, 224, 224, 1, 223, 224, 223, 1102, 25, 15, 225, 1101, 40, 67, 224, 1001, 224, -107, 224, 4, 224, 102, 8, 223, 223, 101, 1, 224, 224, 1, 223, 224, 223, 2, 126, 95, 224, 101, -1400, 224, 224, 4, 224, 1002, 223, 8, 223, 1001, 224, 3, 224, 1, 223, 224, 223, 1002, 151, 84, 224, 101, -2100, 224, 224, 4, 224, 102, 8, 223, 223, 101, 6, 224, 224, 1, 224, 223, 223, 4, 223, 99, 0, 0, 0, 677, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1105, 0, 99999, 1105, 227, 247, 1105, 1, 99999, 1005, 227, 99999, 1005, 0, 256, 1105, 1, 99999, 1106, 227, 99999, 1106, 0, 265, 1105, 1, 99999, 1006, 0, 99999, 1006, 227, 274, 1105, 1, 99999, 1105, 1, 280, 1105, 1, 99999, 1, 225, 225, 225, 1101, 294, 0, 0, 105, 1, 0, 1105, 1, 99999, 1106, 0, 300, 1105, 1, 99999, 1, 225, 225, 225, 1101, 314, 0, 0, 106, 0, 0, 1105, 1, 99999, 108, 677, 677, 224, 1002, 223, 2, 223, 1006, 224, 329, 101, 1, 223, 223, 1107, 677, 226, 224, 102, 2, 223, 223, 1006, 224, 344, 101, 1, 223, 223, 8, 677, 677, 224, 1002, 223, 2, 223, 1006, 224, 359, 101, 1, 223, 223, 1008, 677, 677, 224, 1002, 223, 2, 223, 1006, 224, 374, 101, 1, 223, 223, 7, 226, 677, 224, 1002, 223, 2, 223, 1006, 224, 389, 1001, 223, 1, 223, 1007, 677, 677, 224, 1002, 223, 2, 223, 1006, 224, 404, 1001, 223, 1, 223, 7, 677, 677, 224, 1002, 223, 2, 223, 1006, 224, 419, 1001, 223, 1, 223, 1008, 677, 226, 224, 1002, 223, 2, 223, 1005, 224, 434, 1001, 223, 1, 223, 1107, 226, 677, 224, 102, 2, 223, 223, 1005, 224, 449, 1001, 223, 1, 223, 1008, 226, 226, 224, 1002, 223, 2, 223, 1006, 224, 464, 1001, 223, 1, 223, 1108, 677, 677, 224, 102, 2, 223, 223, 1006, 224, 479, 101, 1, 223, 223, 1108, 226, 677, 224, 1002, 223, 2, 223, 1006, 224, 494, 1001, 223, 1, 223, 107, 226, 226, 224, 1002, 223, 2, 223, 1006, 224, 509, 1001, 223, 1, 223, 8, 226, 677, 224, 102, 2, 223, 223, 1006, 224, 524, 1001, 223, 1, 223, 1007, 226, 226, 224, 1002, 223, 2, 223, 1006, 224, 539, 1001, 223, 1, 223, 107, 677, 677, 224, 1002, 223, 2, 223, 1006, 224, 554, 1001, 223, 1, 223, 1107, 226, 226, 224, 102, 2, 223, 223, 1005, 224, 569, 101, 1, 223, 223, 1108, 677, 226, 224, 1002, 223, 2, 223, 1006, 224, 584, 1001, 223, 1, 223, 1007, 677, 226, 224, 1002, 223, 2, 223, 1005, 224, 599, 101, 1, 223, 223, 107, 226, 677, 224, 102, 2, 223, 223, 1005, 224, 614, 1001, 223, 1, 223, 108, 226, 226, 224, 1002, 223, 2, 223, 1005, 224, 629, 101, 1, 223, 223, 7, 677, 226, 224, 102, 2, 223, 223, 1005, 224, 644, 101, 1, 223, 223, 8, 677, 226, 224, 102, 2, 223, 223, 1006, 224, 659, 1001, 223, 1, 223, 108, 677, 226, 224, 102, 2, 223, 223, 1005, 224, 674, 1001, 223, 1, 223, 4, 223, 99, 226}
	var inputMap = makeMapForArray(inputData)
	// Test data. Prints 999 if input below 8, 1000 if equal, 1001 if greater.
	//inputData := []int64{3, 21, 1008, 21, 8, 20, 1005, 20, 22, 107, 8, 21, 20, 1006, 20, 31, 1106, 0, 36, 98, 0, 0, 1002, 21, 125, 20, 4, 20, 1105, 1, 46, 104, 999, 1105, 1, 46, 1101, 1000, 1, 20, 4, 20, 1105, 1, 46, 98, 99}
	// Test data. if input not 0 it prints 1, if input 0 it prints 0
	//inputData := []int64{3, 12, 6, 12, 15, 1, 13, 14, 13, 4, 13, 99, -1, 0, 1, 9}
	intComp(inputMap, getUserInput, printOutput)
}

type Body struct {
	parent string
	moons  []string
}

func (b *Body) setBodyParent(value string) {
	b.parent = value
}
func (b *Body) setBodyMoons(value string) {
	b.moons = append(b.moons, value)
}

func getStepsToBody(bodyArray []string, bodyToGet string) int {
	for steps, body := range bodyArray {
		if body == bodyToGet {
			return steps
		}
	}
	return -1
}

func day6() {
	file, err := os.Open("./day6Data.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	day6Data := make(map[string]Body)
	scanner := bufio.NewScanner(file)
	var breadthFirstQueue []string
	var currentBodyName string
	var directOrbits int64
	var indirectOrbits int64
	var edgeBodies int64
	var pathBody string
	var youOrbitArray []string
	var sanOrbitArray []string
	var distanceFromYouAndSan int

	for scanner.Scan() {
		orbitArr := strings.Split(scanner.Text(), ")")
		if _, ok := day6Data[orbitArr[0]]; !ok {
			day6Data[orbitArr[0]] = Body{"", []string{}}
		}

		if _, ok := day6Data[orbitArr[1]]; !ok {
			day6Data[orbitArr[1]] = Body{"", []string{}}
		}

		mapBody := day6Data[orbitArr[0]]
		mapMoon := day6Data[orbitArr[1]]

		mapBody.setBodyMoons(orbitArr[1])
		mapMoon.setBodyParent(orbitArr[0])

		day6Data[orbitArr[0]] = mapBody
		day6Data[orbitArr[1]] = mapMoon

	}
	breadthFirstQueue = append(breadthFirstQueue, "COM")
	for len(breadthFirstQueue) > 0 {
		currentBodyName, breadthFirstQueue = breadthFirstQueue[0], breadthFirstQueue[1:]
		currentBody := day6Data[currentBodyName]
		//fmt.Println("Name:", currentBodyName, "Parent:", currentBody.parent, "Moons:", currentBody.moons)
		if len(currentBody.moons) == 0 {
			edgeBodies++
		}
		for _, orbits := range currentBody.moons {
			breadthFirstQueue = append(breadthFirstQueue, orbits)
			directOrbits++
		}
		indirectOrbitParent := day6Data[currentBody.parent].parent
		for indirectOrbitParent != "" {
			indirectOrbitParent = day6Data[indirectOrbitParent].parent
			indirectOrbits++
		}
	}

	pathBody = "YOU"
	for pathBody != "" {
		pathBody = day6Data[pathBody].parent
		youOrbitArray = append(youOrbitArray, pathBody)
	}

	pathBody = "SAN"
	for pathBody != "" {
		pathBody = day6Data[pathBody].parent
		sanOrbitArray = append(sanOrbitArray, pathBody)
	}

	for youSteps, youBody := range youOrbitArray {
		sanSteps := getStepsToBody(sanOrbitArray, youBody)
		if sanSteps != -1 {
			distanceFromYouAndSan = youSteps + sanSteps
			break
		}
	}
	fmt.Println("Direct orbits and indirect orbits:", directOrbits+indirectOrbits)
	fmt.Println("Distance between you and Santa:", distanceFromYouAndSan)
}

func getValidPhaseSequences(initial []int64) (allCombos [][]int64) {
	var switchAllItemsPastPoint func([]int64, int64)
	switchAllItemsPastPoint = func(arr []int64, index int64) {
		if index == int64(len(arr)) {
			allCombos = append(allCombos, append([]int64{}, arr...))
		} else {
			for i := index; i < int64(len(initial)); i++ {
				arr[index], arr[i] = arr[i], arr[index]
				switchAllItemsPastPoint(arr, index+1)
				arr[index], arr[i] = arr[i], arr[index]
			}
		}
	}
	switchAllItemsPastPoint(initial, 0)

	return allCombos
}

func getIOHandlersForChannel(inChannel chan int64, outChannel chan int64, setting int64) (inputHandler, outputHandler, func() int64) {
	var lastOutput int64
	var sentSetting bool = false
	in := func() int64 {
		if !sentSetting {
			sentSetting = true
			return setting
		}
		toReturn := <-inChannel
		return toReturn
	}
	out := func(output int64) {
		lastOutput = output
		outChannel <- output
	}

	getLastOutput := func() int64 {
		return lastOutput
	}

	return in, out, getLastOutput
}

func getIOHandlersForOneRun(setting int64, input int64) (inputHandler, outputHandler, func() int64) {
	var sentSetting bool = false
	var lastOutput int64 = 0
	in := func() int64 {
		if !sentSetting {
			sentSetting = true
			return setting
		}
		return input
	}

	out := func(output int64) {
		lastOutput = output
	}

	getOutput := func() int64 {
		return lastOutput
	}
	return in, out, getOutput
}

func runAmp(data map[int64]int64, in inputHandler, out outputHandler, finished chan bool) {
	intComp(data, in, out)
	finished <- true
}

func day7() {
	day7Data := []int64{3, 8, 1001, 8, 10, 8, 105, 1, 0, 0, 21, 30, 55, 80, 101, 118, 199, 280, 361, 442, 99999, 3, 9, 101, 4, 9, 9, 4, 9, 99, 3, 9, 101, 4, 9, 9, 1002, 9, 4, 9, 101, 4, 9, 9, 1002, 9, 5, 9, 1001, 9, 2, 9, 4, 9, 99, 3, 9, 101, 5, 9, 9, 1002, 9, 2, 9, 101, 3, 9, 9, 102, 4, 9, 9, 1001, 9, 2, 9, 4, 9, 99, 3, 9, 102, 2, 9, 9, 101, 5, 9, 9, 102, 3, 9, 9, 101, 3, 9, 9, 4, 9, 99, 3, 9, 1001, 9, 2, 9, 102, 4, 9, 9, 1001, 9, 3, 9, 4, 9, 99, 3, 9, 1001, 9, 1, 9, 4, 9, 3, 9, 102, 2, 9, 9, 4, 9, 3, 9, 1002, 9, 2, 9, 4, 9, 3, 9, 102, 2, 9, 9, 4, 9, 3, 9, 1002, 9, 2, 9, 4, 9, 3, 9, 1001, 9, 2, 9, 4, 9, 3, 9, 1002, 9, 2, 9, 4, 9, 3, 9, 101, 2, 9, 9, 4, 9, 3, 9, 101, 2, 9, 9, 4, 9, 3, 9, 101, 2, 9, 9, 4, 9, 99, 3, 9, 101, 2, 9, 9, 4, 9, 3, 9, 102, 2, 9, 9, 4, 9, 3, 9, 1001, 9, 1, 9, 4, 9, 3, 9, 102, 2, 9, 9, 4, 9, 3, 9, 1001, 9, 2, 9, 4, 9, 3, 9, 1001, 9, 2, 9, 4, 9, 3, 9, 1002, 9, 2, 9, 4, 9, 3, 9, 102, 2, 9, 9, 4, 9, 3, 9, 1001, 9, 1, 9, 4, 9, 3, 9, 102, 2, 9, 9, 4, 9, 99, 3, 9, 1001, 9, 1, 9, 4, 9, 3, 9, 101, 1, 9, 9, 4, 9, 3, 9, 1001, 9, 2, 9, 4, 9, 3, 9, 1001, 9, 2, 9, 4, 9, 3, 9, 102, 2, 9, 9, 4, 9, 3, 9, 102, 2, 9, 9, 4, 9, 3, 9, 1001, 9, 2, 9, 4, 9, 3, 9, 102, 2, 9, 9, 4, 9, 3, 9, 1002, 9, 2, 9, 4, 9, 3, 9, 102, 2, 9, 9, 4, 9, 99, 3, 9, 1001, 9, 2, 9, 4, 9, 3, 9, 1001, 9, 2, 9, 4, 9, 3, 9, 102, 2, 9, 9, 4, 9, 3, 9, 101, 1, 9, 9, 4, 9, 3, 9, 101, 2, 9, 9, 4, 9, 3, 9, 102, 2, 9, 9, 4, 9, 3, 9, 102, 2, 9, 9, 4, 9, 3, 9, 102, 2, 9, 9, 4, 9, 3, 9, 101, 1, 9, 9, 4, 9, 3, 9, 101, 2, 9, 9, 4, 9, 99, 3, 9, 1001, 9, 2, 9, 4, 9, 3, 9, 101, 2, 9, 9, 4, 9, 3, 9, 101, 1, 9, 9, 4, 9, 3, 9, 1001, 9, 2, 9, 4, 9, 3, 9, 1002, 9, 2, 9, 4, 9, 3, 9, 101, 1, 9, 9, 4, 9, 3, 9, 1001, 9, 2, 9, 4, 9, 3, 9, 1001, 9, 2, 9, 4, 9, 3, 9, 102, 2, 9, 9, 4, 9, 3, 9, 102, 2, 9, 9, 4, 9, 99}
	//Test data for part 1. Max thruster signal should be 43210 from phase setting sequence 43210
	//day7Data := []int64{3, 15, 3, 16, 1002, 16, 10, 16, 1, 16, 15, 15, 4, 15, 99, 0, 0}
	//Test data for part 1. Max thruster signal should be 54321 from phase setting sequence 01234
	//day7Data := []int64{3, 23, 3, 24, 1002, 24, 10, 24, 1002, 23, -1, 23, 101, 5, 23, 23, 1, 24, 23, 23, 4, 23, 99, 0, 0}
	//Test data for part 2. Max thruster signal should be 18216 from phase setting sequence 97856
	//day7Data := []int64{3, 52, 1001, 52, -5, 52, 3, 53, 1, 52, 56, 54, 1007, 54, 5, 55, 1005, 55, 26, 1001, 54, -5, 54, 1105, 1, 12, 1, 53, 54, 53, 1008, 54, 0, 55, 1001, 55, 1, 55, 2, 53, 55, 53, 4, 53, 1001, 56, -1, 56, 1005, 56, 6, 99, 0, 0, 0, 0, 10}

	day7Map := makeMapForArray(day7Data)

	firstPhaseSettings := []int64{0, 1, 2, 3, 4}
	firstValidSequences := getValidPhaseSequences(firstPhaseSettings)

	secondPhaseSettings := []int64{5, 6, 7, 8, 9}
	secondValidSequences := getValidPhaseSequences(secondPhaseSettings)

	var lastOutput int64 = 0
	var lastOutputOfE func() int64

	var largestOutputFirst int64 = 0
	var bestSequenceFirst []int64
	var largestOutputSecond int64 = 0
	var bestSequenceSecond []int64

	var channels = make([]chan int64, 5)

	for _, sequence := range firstValidSequences {
		lastOutput = 0
		for _, setting := range sequence {
			in, out, getOutput := getIOHandlersForOneRun(setting, lastOutput)
			intComp(day7Map, in, out)
			lastOutput = getOutput()
		}
		if lastOutput > largestOutputFirst {
			largestOutputFirst = lastOutput
			bestSequenceFirst = sequence
		}
	}

	for _, sequence := range secondValidSequences {
		lastOutput = 0
		var inputHandlers = make([]inputHandler, 5)
		var outputHandlers = make([]outputHandler, 5)

		for i := 0; i < 5; i++ {
			channels[i] = make(chan int64, 10)
		}
		finished := make(chan bool, 5)

		for i, setting := range sequence {

			outChannel := channels[(i+1)%5]
			in, out, getLastOutput := getIOHandlersForChannel(channels[i], outChannel, setting)
			inputHandlers[i] = in
			outputHandlers[i] = out
			if i == 4 {
				lastOutputOfE = getLastOutput
			}
		}
		for i := range sequence {
			go runAmp(day7Map, inputHandlers[i], outputHandlers[i], finished)
		}
		channels[0] <- 0

		for i := 0; i < 5; i++ {
			<-finished
		}
		lastOutput = lastOutputOfE()

		fmt.Println(lastOutput)
		if lastOutput > largestOutputSecond {
			largestOutputSecond = lastOutput
			bestSequenceSecond = sequence
		}
	}

	fmt.Println("largestOutput first run:", largestOutputFirst)
	fmt.Println("bestSequence first run:", bestSequenceFirst)
	fmt.Println("largestOutput continous:", largestOutputSecond)
	fmt.Println("bestSequence continous:", bestSequenceSecond)

}

func day8() {
	const width = 25
	const height = 6
	var data string
	var layers []string
	var leastZeros = math.MaxInt64
	var leastZerosLayer int
	currentOnes := 0
	currentTwos := 0

	file, err := os.Open("./day8ImageData.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		data = scanner.Text()
	}

	for i := 0; i < len(data); i += width * height {
		layer := data[i : i+(width*height)]
		layers = append(layers, layer)
	}

	for index, layer := range layers {
		currentZeros := 0
		for i := 0; i < len(layer); i++ {
			if layer[i] == '0' {
				currentZeros++
			}
		}
		if currentZeros < leastZeros {
			leastZeros = currentZeros
			leastZerosLayer = index
		}
	}

	for _, str := range layers[leastZerosLayer] {
		if str == '1' {
			currentOnes++
		}
		if str == '2' {
			currentTwos++
		}
	}
	fmt.Println("least zeros layer:", leastZerosLayer)
	fmt.Println("1s and 2s multiplied in that layer:", currentOnes*currentTwos)

	img := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{width, height}})
	var datImg [height][width]string
	for index := range layers {
		usedIndex := (len(layers) - 1) - index
		layer := layers[usedIndex]

		for point, clr := range layer {
			x := point % width
			y := point / width
			switch clr {
			case '0':
				img.Set(x, y, color.Black)
				datImg[y][x] = " "
			case '1':
				img.Set(x, y, color.White)
				datImg[y][x] = "█"
			case '2':
			}
		}
	}
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			fmt.Print(datImg[y][x])
		}
		fmt.Print("\n")
	}
	f, _ := os.Create("day8Image.png")
	png.Encode(f, img)

}

func day9() {
	boostIntCodeData := []int64{1102, 34463338, 34463338, 63, 1007, 63, 34463338, 63, 1005, 63, 53, 1102, 3, 1, 1000, 109, 988, 209, 12, 9, 1000, 209, 6, 209, 3, 203, 0, 1008, 1000, 1, 63, 1005, 63, 65, 1008, 1000, 2, 63, 1005, 63, 902, 1008, 1000, 0, 63, 1005, 63, 58, 4, 25, 104, 0, 99, 4, 0, 104, 0, 99, 4, 17, 104, 0, 99, 0, 0, 1101, 309, 0, 1024, 1101, 0, 24, 1002, 1102, 388, 1, 1029, 1102, 1, 21, 1019, 1101, 0, 33, 1015, 1102, 1, 304, 1025, 1101, 344, 0, 1027, 1101, 25, 0, 1003, 1102, 1, 1, 1021, 1101, 29, 0, 1012, 1101, 0, 23, 1005, 1102, 1, 32, 1007, 1102, 38, 1, 1000, 1101, 30, 0, 1016, 1102, 1, 347, 1026, 1101, 0, 26, 1010, 1101, 0, 39, 1004, 1102, 1, 36, 1011, 1101, 0, 393, 1028, 1101, 0, 37, 1013, 1101, 0, 35, 1008, 1101, 34, 0, 1001, 1101, 0, 495, 1022, 1102, 1, 28, 1018, 1101, 0, 0, 1020, 1102, 1, 22, 1006, 1101, 488, 0, 1023, 1102, 31, 1, 1009, 1102, 1, 20, 1017, 1101, 0, 27, 1014, 109, 10, 21102, 40, 1, 4, 1008, 1014, 37, 63, 1005, 63, 205, 1001, 64, 1, 64, 1106, 0, 207, 4, 187, 1002, 64, 2, 64, 109, -18, 1207, 8, 37, 63, 1005, 63, 227, 1001, 64, 1, 64, 1106, 0, 229, 4, 213, 1002, 64, 2, 64, 109, 17, 1207, -7, 25, 63, 1005, 63, 247, 4, 235, 1106, 0, 251, 1001, 64, 1, 64, 1002, 64, 2, 64, 109, -8, 1202, 6, 1, 63, 1008, 63, 29, 63, 1005, 63, 275, 1001, 64, 1, 64, 1106, 0, 277, 4, 257, 1002, 64, 2, 64, 109, 25, 1205, -6, 293, 1001, 64, 1, 64, 1105, 1, 295, 4, 283, 1002, 64, 2, 64, 109, -4, 2105, 1, 2, 4, 301, 1106, 0, 313, 1001, 64, 1, 64, 1002, 64, 2, 64, 109, -9, 1208, -4, 31, 63, 1005, 63, 335, 4, 319, 1001, 64, 1, 64, 1105, 1, 335, 1002, 64, 2, 64, 109, 16, 2106, 0, -2, 1106, 0, 353, 4, 341, 1001, 64, 1, 64, 1002, 64, 2, 64, 109, -13, 2102, 1, -8, 63, 1008, 63, 38, 63, 1005, 63, 373, 1105, 1, 379, 4, 359, 1001, 64, 1, 64, 1002, 64, 2, 64, 109, 9, 2106, 0, 3, 4, 385, 1105, 1, 397, 1001, 64, 1, 64, 1002, 64, 2, 64, 109, -11, 21107, 41, 42, 0, 1005, 1014, 415, 4, 403, 1106, 0, 419, 1001, 64, 1, 64, 1002, 64, 2, 64, 109, 14, 1206, -7, 431, 1106, 0, 437, 4, 425, 1001, 64, 1, 64, 1002, 64, 2, 64, 109, -23, 2107, 37, -5, 63, 1005, 63, 455, 4, 443, 1105, 1, 459, 1001, 64, 1, 64, 1002, 64, 2, 64, 109, 10, 21107, 42, 41, -2, 1005, 1013, 475, 1105, 1, 481, 4, 465, 1001, 64, 1, 64, 1002, 64, 2, 64, 2105, 1, 8, 1001, 64, 1, 64, 1106, 0, 497, 4, 485, 1002, 64, 2, 64, 109, -6, 21108, 43, 41, 8, 1005, 1017, 517, 1001, 64, 1, 64, 1106, 0, 519, 4, 503, 1002, 64, 2, 64, 109, 5, 2101, 0, -9, 63, 1008, 63, 23, 63, 1005, 63, 541, 4, 525, 1106, 0, 545, 1001, 64, 1, 64, 1002, 64, 2, 64, 109, -13, 1201, 5, 0, 63, 1008, 63, 20, 63, 1005, 63, 565, 1105, 1, 571, 4, 551, 1001, 64, 1, 64, 1002, 64, 2, 64, 109, 16, 1205, 4, 589, 4, 577, 1001, 64, 1, 64, 1106, 0, 589, 1002, 64, 2, 64, 109, -16, 1202, 4, 1, 63, 1008, 63, 23, 63, 1005, 63, 615, 4, 595, 1001, 64, 1, 64, 1106, 0, 615, 1002, 64, 2, 64, 109, 1, 2101, 0, 6, 63, 1008, 63, 33, 63, 1005, 63, 639, 1001, 64, 1, 64, 1105, 1, 641, 4, 621, 1002, 64, 2, 64, 109, 8, 21101, 44, 0, 8, 1008, 1018, 44, 63, 1005, 63, 667, 4, 647, 1001, 64, 1, 64, 1105, 1, 667, 1002, 64, 2, 64, 109, -7, 1201, 1, 0, 63, 1008, 63, 39, 63, 1005, 63, 689, 4, 673, 1106, 0, 693, 1001, 64, 1, 64, 1002, 64, 2, 64, 109, 7, 2102, 1, -8, 63, 1008, 63, 24, 63, 1005, 63, 715, 4, 699, 1105, 1, 719, 1001, 64, 1, 64, 1002, 64, 2, 64, 109, 5, 2108, 34, -7, 63, 1005, 63, 739, 1001, 64, 1, 64, 1105, 1, 741, 4, 725, 1002, 64, 2, 64, 109, -22, 2108, 25, 10, 63, 1005, 63, 763, 4, 747, 1001, 64, 1, 64, 1106, 0, 763, 1002, 64, 2, 64, 109, 31, 1206, -4, 781, 4, 769, 1001, 64, 1, 64, 1105, 1, 781, 1002, 64, 2, 64, 109, -10, 21101, 45, 0, 5, 1008, 1019, 47, 63, 1005, 63, 805, 1001, 64, 1, 64, 1105, 1, 807, 4, 787, 1002, 64, 2, 64, 109, 2, 21108, 46, 46, -3, 1005, 1013, 825, 4, 813, 1106, 0, 829, 1001, 64, 1, 64, 1002, 64, 2, 64, 109, -22, 2107, 40, 10, 63, 1005, 63, 845, 1105, 1, 851, 4, 835, 1001, 64, 1, 64, 1002, 64, 2, 64, 109, 17, 1208, -7, 36, 63, 1005, 63, 871, 1001, 64, 1, 64, 1105, 1, 873, 4, 857, 1002, 64, 2, 64, 109, 16, 21102, 47, 1, -9, 1008, 1018, 47, 63, 1005, 63, 899, 4, 879, 1001, 64, 1, 64, 1106, 0, 899, 4, 64, 99, 21102, 1, 27, 1, 21101, 0, 913, 0, 1105, 1, 920, 21201, 1, 39657, 1, 204, 1, 99, 109, 3, 1207, -2, 3, 63, 1005, 63, 962, 21201, -2, -1, 1, 21102, 1, 940, 0, 1105, 1, 920, 21201, 1, 0, -1, 21201, -2, -3, 1, 21101, 955, 0, 0, 1105, 1, 920, 22201, 1, -1, -2, 1106, 0, 966, 21202, -2, 1, -2, 109, -3, 2105, 1, 0}
	//Test data. Produces a copy of itself
	//boostIntCodeData := []int64{109, 1, 204, -1, 1001, 100, 1, 100, 1008, 100, 16, 101, 1006, 101, 0, 99}
	boostIntCodeMap := makeMapForArray(boostIntCodeData)
	intComp(boostIntCodeMap, getUserInput, printOutput)
}

func getAngle(x1, y1, x2, y2 int64) float64 {
	//Subracting 90 puts the start of the circle at the top
	//Adding 360 if less than 0 turns things into positive degrees for easy sorting
	angle := (math.Atan2(float64(y1-y2), float64(x1-x2)) * 180) / math.Pi
	angle -= 90
	if angle < 0 {
		angle += 360
	}
	return angle
}

func day10() {
	//To use test data. Best asteroid should be 11,13, with most visible being 210
	//asteroidsData := asteroidsTestData
	var asteroids = strings.Split(asteroidsData, "\n")
	var asteroidsArr []Point
	var asteroidsVisible = make(map[Point]int64)
	// var vaporized bool = false
	// var vaporizedCount int64 = 0
	var twoHundrethVaporized Point
	var bestAsteroidPos Point
	var mostVisible int64

	for y, line := range asteroids {
		for x, point := range line {
			if point == '#' {
				asteroidsArr = append(asteroidsArr, Point{int64(x), int64(y)})
				//Asteroids can see themselves
				asteroidsVisible[Point{int64(x), int64(y)}] = 1
			}
		}
	}

	for _, asteroid := range asteroidsArr {
		var slopesUsed = make(map[float64]Point)
		for _, compareAsteroid := range asteroidsArr {
			if asteroid != compareAsteroid {
				anglePoint := Point{1, 1}
				slope := 0.0
				if (compareAsteroid.x - asteroid.x) != 0 {
					slope = float64(compareAsteroid.y-asteroid.y) / float64(compareAsteroid.x-asteroid.x)
				}
				if compareAsteroid.x < asteroid.x {
					anglePoint.x = -1
				}
				if compareAsteroid.y < asteroid.y {
					anglePoint.y = -1
				}
				if slopesUsed[slope] != anglePoint {
					asteroidsVisible[asteroid]++
					slopesUsed[slope] = anglePoint
				}

				if mostVisible < asteroidsVisible[asteroid] {
					bestAsteroidPos = asteroid
					mostVisible = asteroidsVisible[asteroid]
				}
			}
		}
	}
	/*
		Get angles for every asteroid
		Add first astroid per angle to temp array
		Sort temp array by angle
		Add to end array in that order
	*/
	var anglesForPoints = make(map[Point]float64)
	var pointsForAngles = make(map[float64][]Point)
	for _, asteroid := range asteroidsArr {
		if asteroid != bestAsteroidPos {
			angle := getAngle(bestAsteroidPos.x, bestAsteroidPos.y, asteroid.x, asteroid.y)
			anglesForPoints[asteroid] = angle
			pointsForAngles[angle] = append(pointsForAngles[angle], asteroid)
		}
	}
	var vaporizedArray []Point
	for len(vaporizedArray) < len(asteroidsArr) {
		var tempVaporizedArray []Point
		for _, list := range pointsForAngles {
			if len(list) > 0 {
				nextAsteroid := list[0]
				list = list[1:]
				tempVaporizedArray = append(tempVaporizedArray, nextAsteroid)
			}
		}
		sort.Slice(tempVaporizedArray, func(a, b int) bool {
			return anglesForPoints[tempVaporizedArray[a]] < anglesForPoints[tempVaporizedArray[b]]
		})
		for _, point := range tempVaporizedArray {
			vaporizedArray = append(vaporizedArray, point)
		}
	}

	twoHundrethVaporized = vaporizedArray[199]
	fmt.Println("bestPos:", bestAsteroidPos)
	fmt.Println("mostVisible:", mostVisible)
	fmt.Println("200th vaporized:", (twoHundrethVaporized.x*100)+twoHundrethVaporized.y)
}

type Panel struct {
	clr      int64
	traveled bool
}

func testPanelBot() {
	//This mimics the example on the Advent of code site. Print should be 6
	_, out, getpanels := panelBotHandlers(false)
	out(1)
	out(0)
	out(0)
	out(0)
	out(1)
	out(0)
	out(1)
	out(0)
	out(0)
	out(1)
	out(1)
	out(0)
	out(1)
	out(0)
	fmt.Println(len(getpanels()))
}

func panelBotHandlers(shouldStartOnWhite bool) (inputHandler, outputHandler, func() map[Point]Panel) {
	var panelPoints = make(map[Point]Panel)
	var firstOutput = true
	var currentPos = Point{0, 0}
	var dir int64 = 0 //Up 0 Right 1 Down 2 Left 3
	if shouldStartOnWhite {
		panel := panelPoints[currentPos]
		panel.clr = 1
		panelPoints[currentPos] = panel
	}

	in := func() int64 {
		return panelPoints[currentPos].clr
	}

	out := func(output int64) {
		if firstOutput {
			firstOutput = false
			// fmt.Print(output)
			panel := panelPoints[currentPos]
			panel.clr = output
			panel.traveled = true
			panelPoints[currentPos] = panel
		} else {
			firstOutput = true
			// fmt.Print(" ", output, "\n")
			if output == 0 {
				dir--
			} else {
				dir++
			}
			dir = (dir + 4) % 4

			switch dir {
			case 0:
				currentPos.y--
			case 1:
				currentPos.x++
			case 2:
				currentPos.y++
			case 3:
				currentPos.x--
			}
		}
	}

	getPanelPoints := func() map[Point]Panel {
		return panelPoints
	}

	return in, out, getPanelPoints
}

func day11() {
	//This is for testing. It should be the only thing run. Makes for easy debugging that way.
	//testPanelBot()
	var panelBotIntCodesArr = []int64{3, 8, 1005, 8, 318, 1106, 0, 11, 0, 0, 0, 104, 1, 104, 0, 3, 8, 102, -1, 8, 10, 1001, 10, 1, 10, 4, 10, 1008, 8, 1, 10, 4, 10, 101, 0, 8, 29, 1, 107, 12, 10, 2, 1003, 8, 10, 3, 8, 102, -1, 8, 10, 1001, 10, 1, 10, 4, 10, 1008, 8, 0, 10, 4, 10, 1002, 8, 1, 59, 1, 108, 18, 10, 2, 6, 7, 10, 2, 1006, 3, 10, 3, 8, 1002, 8, -1, 10, 1001, 10, 1, 10, 4, 10, 1008, 8, 1, 10, 4, 10, 1002, 8, 1, 93, 1, 1102, 11, 10, 3, 8, 102, -1, 8, 10, 1001, 10, 1, 10, 4, 10, 108, 1, 8, 10, 4, 10, 101, 0, 8, 118, 2, 1102, 10, 10, 3, 8, 102, -1, 8, 10, 101, 1, 10, 10, 4, 10, 1008, 8, 0, 10, 4, 10, 101, 0, 8, 145, 1006, 0, 17, 1006, 0, 67, 3, 8, 1002, 8, -1, 10, 101, 1, 10, 10, 4, 10, 1008, 8, 0, 10, 4, 10, 101, 0, 8, 173, 2, 1109, 4, 10, 1006, 0, 20, 3, 8, 102, -1, 8, 10, 1001, 10, 1, 10, 4, 10, 108, 0, 8, 10, 4, 10, 102, 1, 8, 201, 3, 8, 1002, 8, -1, 10, 1001, 10, 1, 10, 4, 10, 1008, 8, 0, 10, 4, 10, 1002, 8, 1, 224, 1006, 0, 6, 1, 1008, 17, 10, 2, 101, 5, 10, 3, 8, 1002, 8, -1, 10, 1001, 10, 1, 10, 4, 10, 108, 1, 8, 10, 4, 10, 1001, 8, 0, 256, 2, 1107, 7, 10, 1, 2, 4, 10, 2, 2, 12, 10, 1006, 0, 82, 3, 8, 1002, 8, -1, 10, 1001, 10, 1, 10, 4, 10, 1008, 8, 1, 10, 4, 10, 1002, 8, 1, 294, 2, 1107, 2, 10, 101, 1, 9, 9, 1007, 9, 988, 10, 1005, 10, 15, 99, 109, 640, 104, 0, 104, 1, 21102, 1, 837548352256, 1, 21102, 335, 1, 0, 1105, 1, 439, 21102, 1, 47677543180, 1, 21102, 346, 1, 0, 1106, 0, 439, 3, 10, 104, 0, 104, 1, 3, 10, 104, 0, 104, 0, 3, 10, 104, 0, 104, 1, 3, 10, 104, 0, 104, 1, 3, 10, 104, 0, 104, 0, 3, 10, 104, 0, 104, 1, 21102, 1, 235190374592, 1, 21101, 393, 0, 0, 1105, 1, 439, 21102, 3451060455, 1, 1, 21102, 404, 1, 0, 1105, 1, 439, 3, 10, 104, 0, 104, 0, 3, 10, 104, 0, 104, 0, 21102, 837896909668, 1, 1, 21102, 1, 427, 0, 1105, 1, 439, 21102, 1, 709580555020, 1, 21102, 438, 1, 0, 1105, 1, 439, 99, 109, 2, 21201, -1, 0, 1, 21102, 1, 40, 2, 21102, 1, 470, 3, 21102, 460, 1, 0, 1106, 0, 503, 109, -2, 2105, 1, 0, 0, 1, 0, 0, 1, 109, 2, 3, 10, 204, -1, 1001, 465, 466, 481, 4, 0, 1001, 465, 1, 465, 108, 4, 465, 10, 1006, 10, 497, 1101, 0, 0, 465, 109, -2, 2105, 1, 0, 0, 109, 4, 1201, -1, 0, 502, 1207, -3, 0, 10, 1006, 10, 520, 21101, 0, 0, -3, 21202, -3, 1, 1, 22101, 0, -2, 2, 21101, 1, 0, 3, 21101, 0, 539, 0, 1106, 0, 544, 109, -4, 2105, 1, 0, 109, 5, 1207, -3, 1, 10, 1006, 10, 567, 2207, -4, -2, 10, 1006, 10, 567, 21202, -4, 1, -4, 1105, 1, 635, 22101, 0, -4, 1, 21201, -3, -1, 2, 21202, -2, 2, 3, 21101, 0, 586, 0, 1105, 1, 544, 22102, 1, 1, -4, 21102, 1, 1, -1, 2207, -4, -2, 10, 1006, 10, 605, 21102, 1, 0, -1, 22202, -2, -1, -2, 2107, 0, -3, 10, 1006, 10, 627, 21202, -1, 1, 1, 21101, 627, 0, 0, 105, 1, 502, 21202, -2, -1, -2, 22201, -4, -2, -4, 109, -5, 2105, 1, 0}
	var panelBotIntCodes = makeMapForArray(panelBotIntCodesArr)
	var lowestX, lowestY, heighestX, heighestY int64

	//Part 1 gets all points traveled, or all panels painted at least once.
	inHandle, outHandle, getPanelPoints := panelBotHandlers(false)
	intComp(panelBotIntCodes, inHandle, outHandle)
	panelPoints := getPanelPoints()
	fmt.Println("First run on black panel:", len(panelPoints))

	//Run part 2 starting on a white panel
	inHandle, outHandle, getPanelPoints = panelBotHandlers(true)
	intComp(panelBotIntCodes, inHandle, outHandle)
	panelPoints = getPanelPoints()

	for point := range panelPoints {
		if point.x < lowestX {
			lowestX = point.x
		}
		if point.y < lowestY {
			lowestY = point.y
		}
		if point.x > heighestX {
			heighestX = point.x
		}
		if point.y > heighestY {
			heighestY = point.y
		}
	}

	fmt.Println("Second run registration identifier:")
	var x, y int64
	for y = lowestY; y <= heighestY; y++ {
		for x = lowestX; x <= heighestX; x++ {
			if panelPoints[Point{x, y}].clr == 0 {
				fmt.Print(" ")
			} else {
				fmt.Print("█")
			}
		}
		fmt.Print("\n")
	}
}

type xyzPoint struct {
	x int64
	y int64
	z int64
}

func getMoonCombinations(moonPositions []xyzPoint) [][]int {
	var combinations [][]int
	var x, y int
	for x = 0; x < len(moonPositions); x++ {
		for y = x + 1; y < len(moonPositions); y++ {
			combinations = append(combinations, []int{x, y})
		}
	}
	return combinations
}

func xyzPointsInArray(a []xyzPoint, b [][]xyzPoint) int {
	for x := 0; x < len(b); x++ {
		found := 0
		for y := 0; y < len(b[x]); y++ {
			if b[x][y] == a[y] {
				found++
			}
		}
		if found == len(a) {
			return x
		}
	}
	return -1
}

func greatestCommonDivisor(a, b int64) int64 {
	for b != 0 {
		temp := b
		b = a % b
		a = temp
	}
	return a
}

func leastCommonMultiple(a, b int64) int64 {
	// I dont do (a*b) / gcd(a,b) becuase gcd(a,b) is a divisor of both a and b, so this is more efficient
	return (a / greatestCommonDivisor(a, b)) * b
}

func day12() {
	var moonPositions = []xyzPoint{(xyzPoint{3, 15, 8}), (xyzPoint{5, -1, -2}), (xyzPoint{-10, 8, 2}), (xyzPoint{8, 4, -5})}
	//Test data. After 10 steps total energy should be 179 and 2772 steps to get to a past position
	//var moonPositions = []xyzPoint{(xyzPoint{-1, 0, 2}), (xyzPoint{2, -10, -7}), (xyzPoint{4, -8, 8}), (xyzPoint{3, 5, -1})}
	//Test data. After 100 steps total energy should be 1940
	//var moonPositions = []xyzPoint{(xyzPoint{-8, -10, 0}), (xyzPoint{5, 5, 10}), (xyzPoint{2, -7, 3}), (xyzPoint{9, -8, -3})}
	var copiedPositions = []xyzPoint{}
	copy(copiedPositions, moonPositions)
	var prevMoonPositions [][]xyzPoint
	var prevMoonVelocities [][]xyzPoint
	var repeatIntervals = xyzPoint{0, 0, 0}
	var repeatCount = 0
	var moonVelocities = make([]xyzPoint, len(moonPositions))
	var total int64 = 0
	var steps int64 = 0
	var afterSteps int64 = 1000
	moonCombos := getMoonCombinations(moonPositions)
	for repeatCount < 3 || steps < afterSteps {
		steps++
		//Copy positions and velocities so we dont append a pointer
		copiedPositions := make([]xyzPoint, len(moonPositions))
		copiedVelocities := make([]xyzPoint, len(moonPositions))
		copy(copiedPositions, moonPositions)
		copy(copiedVelocities, moonVelocities)
		prevMoonPositions = append(prevMoonPositions, copiedPositions)
		prevMoonVelocities = append(prevMoonVelocities, copiedVelocities)
		for _, moonCombo := range moonCombos {
			//Pull out moons and velocities to change their data
			moon0 := moonPositions[moonCombo[0]]
			moon1 := moonPositions[moonCombo[1]]

			moon0Velocity := moonVelocities[moonCombo[0]]
			moon1Velocity := moonVelocities[moonCombo[1]]

			//Change velocities per axis. Apply gravity
			if moon0.x < moon1.x {
				moon0Velocity.x++
				moon1Velocity.x--
			} else if moon0.x > moon1.x {
				moon0Velocity.x--
				moon1Velocity.x++
			}

			if moon0.y < moon1.y {
				moon0Velocity.y++
				moon1Velocity.y--
			} else if moon0.y > moon1.y {
				moon0Velocity.y--
				moon1Velocity.y++
			}

			if moon0.z < moon1.z {
				moon0Velocity.z++
				moon1Velocity.z--
			} else if moon0.z > moon1.z {
				moon0Velocity.z--
				moon1Velocity.z++
			}

			//Set moons and velocities with changed data
			moonPositions[moonCombo[0]] = moon0
			moonPositions[moonCombo[1]] = moon1

			moonVelocities[moonCombo[0]] = moon0Velocity
			moonVelocities[moonCombo[1]] = moon1Velocity
		}

		for moonIndex := range moonPositions {
			//Change positions for velocities
			moon := moonPositions[moonIndex]
			moonVelocity := moonVelocities[moonIndex]
			moon.x += moonVelocity.x
			moon.y += moonVelocity.y
			moon.z += moonVelocity.z
			moonPositions[moonIndex] = moon
			moonVelocities[moonIndex] = moonVelocity
		}

		// Get intervals
		if repeatIntervals.x == 0 {
			var allMatch bool = true
			for _, velocity := range moonVelocities {
				if velocity.x != 0 {
					allMatch = false
					break
				}
			}
			if allMatch {
				repeatIntervals.x = steps
				fmt.Println("Found X")
				repeatCount++
			}
		}

		if repeatIntervals.y == 0 {
			var allMatch bool = true
			for _, velocity := range moonVelocities {
				if velocity.y != 0 {
					allMatch = false
					break
				}
			}
			if allMatch {
				repeatIntervals.y = steps
				fmt.Println("Found Y")
				repeatCount++
			}
		}

		if repeatIntervals.z == 0 {
			var allMatch bool = true
			for _, velocity := range moonVelocities {
				if velocity.z != 0 {
					allMatch = false
					break
				}
			}
			if allMatch {
				repeatIntervals.z = steps
				fmt.Println("Found Z")
				repeatCount++
			}
		}

		if steps == afterSteps {
			//Get total energy
			for moonIndex := range moonPositions {
				var totalPotential int64
				var totalKinetic int64
				totalPotential += Abs(moonPositions[moonIndex].x)
				totalPotential += Abs(moonPositions[moonIndex].y)
				totalPotential += Abs(moonPositions[moonIndex].z)
				totalKinetic += Abs(moonVelocities[moonIndex].x)
				totalKinetic += Abs(moonVelocities[moonIndex].y)
				totalKinetic += Abs(moonVelocities[moonIndex].z)
				total += totalPotential * totalKinetic
			}
			fmt.Println("total energy:", total)
		}
	}

	pastPointSteps := leastCommonMultiple(repeatIntervals.x, repeatIntervals.y)
	pastPointSteps = leastCommonMultiple(pastPointSteps, repeatIntervals.z) * 2
	fmt.Println("Steps to past point in time:", pastPointSteps)
}

func clearConsole(OS string) {

	switch OS {
	case "linux":
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	case "windows":
		cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func printScreen(screen map[Point]int64, width, height int64) {
	var x, y int64
	fmt.Println("Score:", screen[Point{-1, 0}])
	for x = 0; y <= height; y++ {
		for x = 0; x <= width; x++ {
			switch screen[Point{x, y}] {
			case 0:
				fmt.Print(" ")
			case 1:
				fmt.Print("█")
			case 2:
				fmt.Print("■")
			case 3:
				fmt.Print("—")
			case 4:
				fmt.Print("•")
			}
		}
		fmt.Print("\n")
	}
}

func getScreenIOHandlers(showScreen bool, width, height int64) (inputHandler, outputHandler, func() map[Point]int64) {
	var outputNum int64 = 0
	var setX, setY int64
	var drawAfter int64 = 5
	var frame int64 = 0
	ballPos := Point{0, 0}
	paddlePos := Point{0, 0}
	var screen = make(map[Point]int64)
	in := func() int64 {
		if paddlePos.x < ballPos.x {
			return 1
		} else if paddlePos.x > ballPos.x {
			return -1
		}
		return 0
	}

	out := func(output int64) {
		switch outputNum {
		case 0:
			setX = output
			outputNum++
		case 1:
			setY = output
			outputNum++
		case 2:
			screen[Point{setX, setY}] = output
			frame++
			if output == 4 {
				ballPos = Point{setX, setY}
			}
			if output == 3 {
				paddlePos = Point{setX, setY}
			}
			if frame == drawAfter && showScreen {
				frame = 0
				clearConsole("windows")
				printScreen(screen, width, height)
			}
			outputNum = 0
		}
	}

	getScreen := func() map[Point]int64 {
		return screen
	}

	return in, out, getScreen
}

func day13() {
	intCodesArr := []int64{1, 380, 379, 385, 1008, 2663, 456801, 381, 1005, 381, 12, 99, 109, 2664, 1101, 0, 0, 383, 1101, 0, 0, 382, 20101, 0, 382, 1, 20102, 1, 383, 2, 21102, 37, 1, 0, 1105, 1, 578, 4, 382, 4, 383, 204, 1, 1001, 382, 1, 382, 1007, 382, 44, 381, 1005, 381, 22, 1001, 383, 1, 383, 1007, 383, 23, 381, 1005, 381, 18, 1006, 385, 69, 99, 104, -1, 104, 0, 4, 386, 3, 384, 1007, 384, 0, 381, 1005, 381, 94, 107, 0, 384, 381, 1005, 381, 108, 1106, 0, 161, 107, 1, 392, 381, 1006, 381, 161, 1101, -1, 0, 384, 1106, 0, 119, 1007, 392, 42, 381, 1006, 381, 161, 1102, 1, 1, 384, 21001, 392, 0, 1, 21102, 1, 21, 2, 21102, 1, 0, 3, 21102, 138, 1, 0, 1105, 1, 549, 1, 392, 384, 392, 21001, 392, 0, 1, 21101, 0, 21, 2, 21102, 1, 3, 3, 21101, 0, 161, 0, 1106, 0, 549, 1101, 0, 0, 384, 20001, 388, 390, 1, 20101, 0, 389, 2, 21102, 1, 180, 0, 1105, 1, 578, 1206, 1, 213, 1208, 1, 2, 381, 1006, 381, 205, 20001, 388, 390, 1, 21001, 389, 0, 2, 21101, 0, 205, 0, 1106, 0, 393, 1002, 390, -1, 390, 1102, 1, 1, 384, 20102, 1, 388, 1, 20001, 389, 391, 2, 21102, 1, 228, 0, 1105, 1, 578, 1206, 1, 261, 1208, 1, 2, 381, 1006, 381, 253, 20101, 0, 388, 1, 20001, 389, 391, 2, 21102, 1, 253, 0, 1105, 1, 393, 1002, 391, -1, 391, 1101, 1, 0, 384, 1005, 384, 161, 20001, 388, 390, 1, 20001, 389, 391, 2, 21102, 1, 279, 0, 1105, 1, 578, 1206, 1, 316, 1208, 1, 2, 381, 1006, 381, 304, 20001, 388, 390, 1, 20001, 389, 391, 2, 21101, 0, 304, 0, 1105, 1, 393, 1002, 390, -1, 390, 1002, 391, -1, 391, 1102, 1, 1, 384, 1005, 384, 161, 21001, 388, 0, 1, 20102, 1, 389, 2, 21102, 1, 0, 3, 21101, 338, 0, 0, 1105, 1, 549, 1, 388, 390, 388, 1, 389, 391, 389, 20102, 1, 388, 1, 20102, 1, 389, 2, 21102, 1, 4, 3, 21101, 0, 365, 0, 1106, 0, 549, 1007, 389, 22, 381, 1005, 381, 75, 104, -1, 104, 0, 104, 0, 99, 0, 1, 0, 0, 0, 0, 0, 0, 315, 20, 18, 1, 1, 22, 109, 3, 22101, 0, -2, 1, 21202, -1, 1, 2, 21102, 0, 1, 3, 21101, 0, 414, 0, 1106, 0, 549, 22102, 1, -2, 1, 21202, -1, 1, 2, 21101, 429, 0, 0, 1106, 0, 601, 1202, 1, 1, 435, 1, 386, 0, 386, 104, -1, 104, 0, 4, 386, 1001, 387, -1, 387, 1005, 387, 451, 99, 109, -3, 2105, 1, 0, 109, 8, 22202, -7, -6, -3, 22201, -3, -5, -3, 21202, -4, 64, -2, 2207, -3, -2, 381, 1005, 381, 492, 21202, -2, -1, -1, 22201, -3, -1, -3, 2207, -3, -2, 381, 1006, 381, 481, 21202, -4, 8, -2, 2207, -3, -2, 381, 1005, 381, 518, 21202, -2, -1, -1, 22201, -3, -1, -3, 2207, -3, -2, 381, 1006, 381, 507, 2207, -3, -4, 381, 1005, 381, 540, 21202, -4, -1, -1, 22201, -3, -1, -3, 2207, -3, -4, 381, 1006, 381, 529, 21202, -3, 1, -7, 109, -8, 2106, 0, 0, 109, 4, 1202, -2, 44, 566, 201, -3, 566, 566, 101, 639, 566, 566, 2102, 1, -1, 0, 204, -3, 204, -2, 204, -1, 109, -4, 2106, 0, 0, 109, 3, 1202, -1, 44, 594, 201, -2, 594, 594, 101, 639, 594, 594, 20101, 0, 0, -2, 109, -3, 2106, 0, 0, 109, 3, 22102, 23, -2, 1, 22201, 1, -1, 1, 21102, 1, 509, 2, 21101, 264, 0, 3, 21102, 1012, 1, 4, 21102, 1, 630, 0, 1106, 0, 456, 21201, 1, 1651, -2, 109, -3, 2106, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 2, 0, 0, 2, 0, 2, 0, 0, 2, 0, 2, 0, 0, 2, 2, 0, 2, 0, 2, 2, 0, 0, 0, 0, 2, 0, 0, 2, 2, 2, 0, 2, 0, 2, 0, 2, 2, 2, 0, 1, 1, 0, 0, 0, 2, 2, 0, 0, 0, 0, 2, 2, 2, 0, 0, 2, 0, 2, 0, 0, 0, 0, 2, 2, 2, 2, 0, 0, 2, 0, 2, 2, 2, 2, 0, 0, 0, 2, 0, 2, 0, 0, 0, 1, 1, 0, 2, 0, 0, 0, 2, 0, 2, 2, 2, 0, 0, 2, 0, 2, 0, 0, 0, 0, 2, 2, 0, 0, 2, 2, 0, 0, 2, 0, 0, 0, 0, 2, 2, 0, 2, 0, 0, 2, 0, 0, 0, 1, 1, 0, 0, 2, 0, 2, 2, 0, 0, 0, 2, 0, 2, 2, 2, 0, 2, 0, 2, 2, 2, 2, 0, 2, 2, 2, 0, 0, 0, 0, 0, 2, 0, 2, 0, 2, 0, 2, 2, 0, 0, 2, 0, 1, 1, 0, 0, 0, 0, 0, 2, 0, 2, 2, 0, 0, 2, 0, 2, 0, 2, 2, 0, 0, 0, 0, 2, 2, 0, 2, 2, 2, 0, 2, 2, 0, 2, 2, 0, 0, 2, 2, 0, 0, 0, 2, 0, 1, 1, 0, 0, 2, 2, 2, 0, 0, 0, 0, 0, 2, 2, 0, 0, 2, 0, 0, 0, 2, 2, 0, 2, 2, 0, 0, 2, 0, 0, 2, 0, 0, 0, 2, 2, 0, 2, 2, 2, 0, 2, 0, 0, 1, 1, 0, 0, 2, 0, 0, 2, 0, 2, 2, 2, 2, 2, 0, 2, 0, 0, 0, 0, 2, 2, 2, 2, 2, 2, 2, 0, 2, 2, 2, 2, 2, 2, 2, 0, 2, 0, 0, 0, 2, 2, 2, 0, 1, 1, 0, 2, 0, 0, 2, 2, 2, 0, 2, 2, 2, 0, 2, 2, 0, 0, 0, 2, 0, 2, 2, 0, 2, 0, 2, 2, 2, 0, 2, 0, 0, 0, 0, 0, 0, 2, 0, 2, 2, 2, 2, 0, 1, 1, 0, 0, 2, 2, 2, 2, 0, 2, 2, 2, 2, 0, 2, 2, 2, 0, 0, 0, 2, 0, 0, 0, 0, 0, 2, 2, 2, 2, 2, 2, 0, 0, 2, 0, 2, 2, 2, 2, 2, 0, 0, 0, 1, 1, 0, 0, 2, 2, 2, 2, 2, 0, 0, 2, 2, 0, 0, 2, 2, 0, 2, 0, 0, 0, 0, 2, 2, 0, 2, 0, 2, 2, 2, 2, 0, 2, 0, 0, 2, 2, 0, 2, 2, 0, 2, 0, 1, 1, 0, 2, 0, 2, 2, 0, 0, 2, 0, 2, 0, 2, 2, 0, 2, 0, 2, 0, 2, 2, 0, 0, 2, 2, 2, 2, 0, 2, 2, 2, 0, 2, 2, 0, 0, 0, 2, 2, 0, 0, 2, 0, 1, 1, 0, 2, 0, 0, 2, 2, 2, 2, 2, 0, 0, 2, 0, 0, 2, 2, 2, 0, 2, 0, 2, 0, 2, 0, 0, 0, 2, 0, 0, 0, 0, 2, 2, 2, 2, 0, 2, 2, 0, 0, 0, 0, 1, 1, 0, 0, 2, 0, 0, 2, 2, 0, 2, 0, 2, 0, 2, 2, 0, 0, 2, 0, 2, 0, 0, 2, 2, 2, 2, 2, 2, 0, 0, 0, 0, 2, 2, 2, 2, 2, 0, 2, 0, 2, 2, 0, 1, 1, 0, 0, 0, 2, 2, 2, 2, 2, 2, 2, 0, 0, 2, 0, 0, 0, 0, 2, 0, 2, 2, 2, 2, 0, 2, 0, 2, 0, 0, 2, 0, 2, 2, 2, 2, 0, 0, 0, 2, 2, 2, 0, 1, 1, 0, 0, 0, 2, 2, 0, 2, 0, 2, 2, 0, 2, 0, 2, 0, 0, 0, 2, 0, 2, 0, 2, 2, 0, 2, 0, 0, 2, 0, 0, 2, 0, 2, 0, 2, 2, 2, 0, 0, 2, 2, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 78, 97, 10, 89, 31, 40, 53, 97, 63, 60, 92, 10, 54, 27, 53, 42, 36, 34, 79, 30, 8, 70, 22, 20, 18, 67, 79, 30, 81, 50, 67, 46, 39, 15, 72, 26, 35, 61, 6, 36, 2, 26, 65, 94, 82, 27, 37, 6, 71, 66, 84, 19, 69, 5, 62, 89, 57, 49, 1, 9, 59, 67, 30, 74, 71, 37, 66, 77, 43, 4, 59, 42, 85, 4, 87, 1, 24, 64, 85, 25, 29, 67, 97, 15, 22, 6, 34, 97, 97, 47, 22, 19, 40, 89, 45, 36, 93, 77, 26, 85, 30, 40, 65, 21, 45, 91, 18, 77, 45, 13, 74, 18, 47, 67, 79, 1, 31, 22, 1, 96, 94, 60, 44, 56, 79, 64, 74, 56, 91, 79, 41, 23, 9, 57, 9, 86, 63, 82, 55, 92, 63, 63, 94, 73, 76, 40, 88, 18, 26, 66, 29, 27, 20, 1, 94, 90, 43, 11, 67, 33, 27, 47, 34, 73, 65, 67, 77, 54, 92, 84, 6, 29, 41, 8, 8, 38, 83, 36, 74, 29, 26, 70, 68, 57, 54, 38, 75, 37, 24, 64, 30, 89, 43, 61, 6, 4, 65, 81, 39, 85, 91, 22, 28, 17, 47, 95, 52, 40, 76, 77, 81, 52, 59, 19, 37, 90, 23, 33, 5, 82, 3, 64, 46, 70, 22, 24, 9, 96, 97, 69, 48, 66, 58, 97, 51, 15, 86, 6, 23, 7, 35, 52, 57, 3, 79, 82, 71, 87, 64, 91, 93, 69, 77, 95, 1, 57, 5, 2, 65, 35, 57, 14, 35, 12, 14, 60, 45, 52, 67, 32, 26, 93, 63, 54, 45, 8, 48, 83, 5, 49, 95, 60, 78, 98, 54, 62, 9, 1, 39, 57, 63, 82, 52, 90, 64, 38, 95, 8, 12, 72, 22, 53, 78, 63, 72, 65, 59, 1, 87, 95, 81, 79, 38, 92, 61, 60, 59, 3, 39, 31, 47, 69, 70, 6, 55, 44, 49, 54, 49, 50, 11, 87, 85, 89, 15, 70, 58, 5, 87, 65, 79, 86, 92, 98, 49, 73, 8, 79, 30, 55, 4, 30, 11, 55, 80, 28, 63, 28, 33, 9, 49, 70, 34, 83, 29, 97, 67, 65, 89, 50, 88, 29, 40, 5, 3, 11, 87, 85, 43, 2, 51, 18, 58, 39, 81, 8, 15, 2, 42, 95, 64, 8, 76, 60, 73, 67, 30, 28, 11, 84, 56, 73, 14, 66, 43, 21, 40, 31, 48, 11, 65, 27, 9, 37, 60, 91, 34, 11, 83, 45, 9, 77, 70, 97, 9, 13, 68, 20, 17, 15, 6, 13, 44, 59, 51, 91, 73, 60, 37, 40, 18, 69, 48, 14, 44, 96, 71, 21, 27, 90, 9, 91, 14, 80, 38, 69, 69, 52, 28, 15, 54, 63, 46, 32, 78, 54, 79, 95, 83, 16, 44, 29, 26, 92, 31, 51, 66, 14, 94, 49, 1, 93, 43, 57, 50, 82, 45, 95, 83, 74, 50, 87, 47, 55, 62, 31, 1, 88, 1, 77, 59, 64, 26, 48, 22, 61, 56, 20, 54, 59, 62, 3, 59, 28, 98, 45, 53, 47, 72, 73, 72, 43, 30, 23, 94, 10, 76, 63, 63, 8, 30, 92, 25, 61, 61, 32, 64, 25, 57, 61, 95, 81, 23, 67, 28, 59, 48, 68, 21, 85, 48, 32, 93, 98, 50, 89, 27, 46, 38, 63, 38, 87, 76, 76, 10, 71, 36, 91, 2, 47, 2, 36, 37, 90, 25, 97, 27, 71, 67, 77, 4, 11, 57, 68, 87, 94, 12, 83, 91, 94, 92, 35, 49, 46, 4, 31, 64, 39, 12, 92, 26, 12, 75, 29, 11, 5, 83, 8, 23, 73, 62, 74, 55, 75, 38, 40, 90, 73, 71, 38, 15, 75, 10, 38, 55, 74, 82, 13, 32, 55, 90, 47, 6, 25, 65, 88, 85, 40, 13, 66, 54, 39, 82, 19, 15, 18, 74, 19, 54, 70, 30, 56, 28, 2, 20, 50, 44, 51, 7, 4, 79, 97, 90, 71, 97, 5, 25, 95, 22, 36, 61, 30, 16, 68, 61, 23, 22, 60, 93, 9, 92, 98, 40, 41, 11, 47, 7, 57, 15, 51, 51, 77, 22, 32, 4, 27, 10, 76, 76, 50, 81, 96, 46, 28, 38, 69, 41, 43, 47, 86, 66, 54, 22, 33, 45, 75, 75, 51, 37, 62, 62, 25, 71, 35, 49, 93, 44, 18, 92, 39, 32, 11, 31, 96, 2, 33, 94, 45, 14, 82, 57, 79, 81, 57, 6, 19, 63, 35, 11, 55, 18, 38, 22, 43, 82, 76, 35, 7, 21, 74, 50, 83, 7, 55, 94, 23, 79, 85, 20, 4, 65, 18, 12, 62, 35, 74, 23, 20, 96, 71, 25, 95, 45, 95, 4, 18, 82, 71, 79, 4, 12, 41, 44, 23, 8, 86, 6, 78, 5, 54, 68, 60, 12, 73, 18, 95, 31, 86, 23, 5, 36, 40, 97, 35, 48, 28, 15, 9, 27, 54, 14, 22, 97, 63, 41, 37, 12, 20, 38, 41, 27, 70, 35, 10, 89, 31, 90, 44, 46, 44, 49, 66, 71, 58, 74, 7, 24, 6, 96, 68, 27, 16, 89, 80, 1, 38, 26, 88, 60, 47, 27, 46, 32, 34, 44, 74, 51, 70, 13, 57, 14, 31, 40, 71, 55, 22, 87, 23, 9, 37, 38, 18, 17, 34, 84, 84, 49, 74, 81, 31, 4, 45, 11, 71, 89, 16, 56, 91, 61, 61, 67, 92, 14, 88, 89, 10, 11, 77, 38, 40, 89, 76, 7, 5, 74, 54, 64, 97, 25, 20, 1, 41, 9, 41, 97, 1, 31, 21, 96, 98, 88, 52, 71, 25, 62, 42, 8, 91, 84, 43, 75, 37, 22, 32, 58, 87, 22, 6, 13, 62, 48, 85, 81, 48, 70, 3, 13, 93, 88, 52, 7, 66, 84, 27, 37, 21, 62, 72, 40, 30, 28, 12, 88, 48, 47, 96, 98, 47, 76, 80, 98, 42, 25, 72, 13, 15, 31, 81, 40, 16, 85, 77, 82, 41, 67, 93, 73, 58, 86, 68, 85, 28, 60, 13, 87, 9, 12, 40, 20, 4, 92, 51, 456801}
	intCodes := makeMapForArray(intCodesArr)
	in, out, getScreen := getScreenIOHandlers(false, 0, 0)
	intComp(intCodes, in, out)
	screen := getScreen()
	var blockTileCount int64 = 0
	var width, height int64
	for point, tile := range screen {
		if point.x > width {
			width = point.x
		}
		if point.y > height {
			height = point.y
		}
		if tile == 2 {
			blockTileCount++
		}
	}
	fmt.Println("Block tiles after exit:", blockTileCount)

	// "Insert" quarters. Memory address 0 is the number of quarters that have been inserted.
	intCodes[0] = 2
	in, out, getScreen = getScreenIOHandlers(false, width, height)
	intComp(intCodes, in, out)
	screen = getScreen()
	printScreen(screen, width, height)
	print("score:", screen[Point{-1, 0}])
}

type reactionChem struct {
	amount int
	name   string
}

func getNeededChemForChemOfAmount(neededChem string, resultChem string, amount int,
	neededForResult map[string][]reactionChem, producedAmount map[string]int, excess map[string]int) int {
	var gottenAmount int = 0
	var scale int = 1

	//ORE requires nothing to produce it, so return 0.
	if resultChem == "ORE" || amount == 0 {
		return 0
	}

	neededChemicals := neededForResult[resultChem]
	if len(neededChemicals) != 0 {
		if amount > producedAmount[resultChem] {
			scale = int(math.Ceil(float64(amount) / float64(producedAmount[resultChem])))

		}
		for _, chem := range neededChemicals {

			if chem.name == neededChem {
				gottenAmount += chem.amount * scale
				continue
			}

			// Golang % operator is the remainder function, and so a divide by 0 can happen
			// This is why I have this here
			if chem.name == "ORE" {
				continue
			}

			// If amountNeeded(4) - excess(5) is less than zero,
			// amountNeeded is set to 0 and excess would be  excess(5) -= amountNeeded(0).
			// So this temp is here to keep the amountNeeded value for later
			// so excess(5) -= tempAmount(4)
			tempAmount := (chem.amount * scale)
			amountNeeded := (chem.amount * scale) - excess[chem.name]
			if amountNeeded < 0 {
				amountNeeded = 0
			}

			excess[chem.name] -= tempAmount
			if excess[chem.name] < 0 {
				excess[chem.name] = 0
			}

			if (amountNeeded % producedAmount[chem.name]) != 0 {
				excess[chem.name] += producedAmount[chem.name] - (amountNeeded % producedAmount[chem.name])
			}

			gottenAmount += getNeededChemForChemOfAmount(neededChem, chem.name, amountNeeded,
				neededForResult, producedAmount, excess)
		}
		return gottenAmount
	}
	return 0
}

func getAmountOfChemProducedFromChemOfAmount(fromChem string, fromChemAmount int, producedChem string,
	neededForResult map[string][]reactionChem, producedAmount map[string]int) int {
	var amountNeeded = 0
	var stepAmount = 1000000
	var amountGuess = 0
	var found bool = false
	for !found {
		var excess = make(map[string]int)
		amountGuess += stepAmount
		amountNeeded = getNeededChemForChemOfAmount(fromChem, producedChem, amountGuess,
			neededForResult, producedAmount, excess)

		if amountNeeded > fromChemAmount {
			fmt.Print(amountGuess, " ")
			amountGuess -= stepAmount
			if stepAmount == 1 {
				found = true
			}
			stepAmount = int(math.Ceil(float64(stepAmount) / 2))
			fmt.Print(stepAmount, "\n")
		}

	}

	return amountGuess

}

func day14() {
	file, err := os.Open("day14Data.txt")
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	var chemicalsNeededForResult = make(map[string][]reactionChem)
	var producedAmount = make(map[string]int)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		reactionString := scanner.Text()
		//Get inputs and output for reaction
		reactionArray := strings.Split(reactionString, " => ")

		//Get result data
		resultArray := strings.Split(reactionArray[1], " ")
		resultAmount, _ := strconv.Atoi(resultArray[0])
		resultName := resultArray[1]
		producedAmount[resultName] = resultAmount

		//Get chemicals needed in structs
		var completeInputChemArray []reactionChem
		inputChemStringArray := strings.Split(reactionArray[0], ", ")
		for _, inputChem := range inputChemStringArray {
			if inputChem == "ORE" {
				completeInputChemArray = append(completeInputChemArray, reactionChem{0, "ORE"})
				continue
			}
			inputChemArray := strings.Split(inputChem, " ")
			inputChemAmount, _ := strconv.Atoi(inputChemArray[0])
			parsedInputChem := reactionChem{inputChemAmount, inputChemArray[1]}
			completeInputChemArray = append(completeInputChemArray, parsedInputChem)
		}

		chemicalsNeededForResult[resultName] = completeInputChemArray
	}
	var excessChemicals = make(map[string]int)
	oreNeededForOneFuel := getNeededChemForChemOfAmount("ORE", "FUEL", 1,
		chemicalsNeededForResult, producedAmount, excessChemicals)
	fuelProducedFromTrillionOre := getAmountOfChemProducedFromChemOfAmount("ORE", 1000000000000, "FUEL",
		chemicalsNeededForResult, producedAmount)
	fmt.Println("Ore needed for one fuel:", oreNeededForOneFuel)
	fmt.Println("Fuel produced from one trillion ore:", fuelProducedFromTrillionOre)
}

type droidTile struct {
	tile     string
	traveled bool
}

const (
	NORTH int64 = 1
	SOUTH int64 = 2
	WEST  int64 = 3
	EAST  int64 = 4
)

func getPointForDirection(dir int64, pos Point) Point {
	var tilePoint Point
	switch dir {
	case NORTH:
		tilePoint = Point{pos.x, pos.y - 1}
	case SOUTH:
		tilePoint = Point{pos.x, pos.y + 1}
	case WEST:
		tilePoint = Point{pos.x - 1, pos.y}
	case EAST:
		tilePoint = Point{pos.x + 1, pos.y}
	}
	return tilePoint
}

func drawAreaAroundDroid(area map[Point]droidTile, topLeft, bottomRight Point) {
	for y := topLeft.y; y <= bottomRight.y; y++ {
		for x := topLeft.x; x <= bottomRight.x; x++ {
			if area[Point{x, y}].tile == "" {
				fmt.Print(" ")
			} else {
				fmt.Print(area[Point{x, y}].tile)
			}
		}
		fmt.Println()
	}
}

func getDroidIOHandlers(drawSetting int64) (inputHandler, outputHandler, func() []Point, func() int64) {
	var shortestPath []Point
	rand.Seed(time.Now().UnixNano())
	var areaAroundDroid = make(map[Point]droidTile)
	var stack []Point
	var topLeft Point
	var bottomRight Point
	var OxygenTankPoint Point

	var found bool = false
	var gotMap bool = false
	var backtracking bool = false
	var availableDirs = []int64{NORTH, SOUTH, WEST, EAST}
	var dir int64 = NORTH
	var currentPos = Point{0, 0}
	areaAroundDroid[currentPos] = droidTile{"D", true}

	in := func() int64 {
		for {
			if len(availableDirs) == 0 {
				backtracking = true
				if currentPos == (Point{0, 0}) {
					gotMap = true
					if drawSetting == 0 {
						drawAreaAroundDroid(areaAroundDroid, topLeft, bottomRight)
					}
					return 0
				}
				switch drawSetting {
				case 1:
					fmt.Print("OUT OF DIRS, BACKTRACKING.")
				case 2:
					if found {
						fmt.Print("OUT OF DIRS, BACKTRACKING.")
					}
				}
				if currentPos.x < stack[len(stack)-1].x {
					dir = EAST
				}
				if currentPos.x > stack[len(stack)-1].x {
					dir = WEST
				}
				if currentPos.y < stack[len(stack)-1].y {
					dir = SOUTH
				}
				if currentPos.y > stack[len(stack)-1].y {
					dir = NORTH
				}
				stack = stack[:len(stack)-1]
				return dir
			}
			backtracking = false

			dirIndex := int64(rand.Intn(len(availableDirs)) % 4)
			dir = availableDirs[dirIndex]
			tilePoint := getPointForDirection(dir, currentPos)
			tile := areaAroundDroid[tilePoint]

			if tile.traveled {
				availableDirs = append(availableDirs[:dirIndex], availableDirs[dirIndex+1:]...)
				continue
			}
			return dir
		}
	}

	out := func(output int64) {
		availableDirs = []int64{NORTH, SOUTH, WEST, EAST}

		tilePoint := getPointForDirection(dir, currentPos)
		tile := areaAroundDroid[tilePoint]
		tile.traveled = true
		switch output {
		case 0:
			tile.tile = "#"
			areaAroundDroid[tilePoint] = tile
		case 1:
			// For drawing droid
			tile.tile = "D"
			currentTile := areaAroundDroid[currentPos]
			currentTile.tile = "-"
			if currentPos == OxygenTankPoint && currentPos != (Point{0, 0}) {
				currentTile.tile = "O"
			}
			areaAroundDroid[currentPos] = currentTile

			if tilePoint.x < topLeft.x {
				topLeft.x = tilePoint.x
			}
			if tilePoint.y < topLeft.y {
				topLeft.y = tilePoint.y
			}
			if tilePoint.x > bottomRight.x {
				bottomRight.x = tilePoint.x
			}
			if tilePoint.y > bottomRight.y {
				bottomRight.y = tilePoint.y
			}
			if !backtracking {
				stack = append(stack, currentPos)
			}
			currentPos = tilePoint
		case 2:
			// Copy stack now that the droid has found the oxygen system
			tile.tile = "O"
			stack = append(stack, currentPos)
			shortestPath = make([]Point, len(stack))
			currentPos = tilePoint
			OxygenTankPoint = tilePoint
			copy(shortestPath, stack)
			found = true
		}
		areaAroundDroid[tilePoint] = tile

		switch drawSetting {
		case 1:
			clearConsole("windows")
			drawAreaAroundDroid(areaAroundDroid, topLeft, bottomRight)
			if found {
				fmt.Println("FOUND")
			}
		case 2:
			if found {
				clearConsole("windows")
				drawAreaAroundDroid(areaAroundDroid, topLeft, bottomRight)
				fmt.Println("FOUND")
			}
		}
	}

	getShortestPath := func() []Point {
		return shortestPath
	}

	getTimeForOxygenSpread := func() int64 {
		var time int64 = 0
		var canGoToQueue []Point
		var perMinuteQueue []Point
		canGoToQueue = append(canGoToQueue, OxygenTankPoint)
		availableDirs = []int64{NORTH, SOUTH, WEST, EAST}

		for len(canGoToQueue) != 0 {
			perMinuteQueue = make([]Point, len(canGoToQueue))
			copy(perMinuteQueue, canGoToQueue)
			canGoToQueue = []Point{}
			for _, point := range perMinuteQueue {

				for _, dir := range availableDirs {
					dirPoint := getPointForDirection(dir, point)
					dirTile := areaAroundDroid[dirPoint]
					if dirTile.tile != "#" && dirTile.tile != "O" {
						dirTile.tile = "O"
						areaAroundDroid[dirPoint] = dirTile
						canGoToQueue = append(canGoToQueue, dirPoint)
					}
				}
			}

			time++
			if drawSetting > 0 && drawSetting < 4 {
				clearConsole("windows")
				drawAreaAroundDroid(areaAroundDroid, topLeft, bottomRight)
				if found {
					fmt.Println("FOUND")
				}
			}
		}
		//Oxygen tank is the first point, so subtract one to account for that.
		time -= 1

		return time
	}

	return in, out, getShortestPath, getTimeForOxygenSpread
}

func day15() {
	var intCodeDataArr = []int64{3, 1033, 1008, 1033, 1, 1032, 1005, 1032, 31, 1008, 1033, 2, 1032, 1005, 1032, 58, 1008, 1033, 3, 1032, 1005, 1032, 81, 1008, 1033, 4, 1032, 1005, 1032, 104, 99, 1002, 1034, 1, 1039, 1002, 1036, 1, 1041, 1001, 1035, -1, 1040, 1008, 1038, 0, 1043, 102, -1, 1043, 1032, 1, 1037, 1032, 1042, 1106, 0, 124, 1001, 1034, 0, 1039, 1002, 1036, 1, 1041, 1001, 1035, 1, 1040, 1008, 1038, 0, 1043, 1, 1037, 1038, 1042, 1106, 0, 124, 1001, 1034, -1, 1039, 1008, 1036, 0, 1041, 1002, 1035, 1, 1040, 1001, 1038, 0, 1043, 101, 0, 1037, 1042, 1105, 1, 124, 1001, 1034, 1, 1039, 1008, 1036, 0, 1041, 102, 1, 1035, 1040, 1001, 1038, 0, 1043, 101, 0, 1037, 1042, 1006, 1039, 217, 1006, 1040, 217, 1008, 1039, 40, 1032, 1005, 1032, 217, 1008, 1040, 40, 1032, 1005, 1032, 217, 1008, 1039, 39, 1032, 1006, 1032, 165, 1008, 1040, 3, 1032, 1006, 1032, 165, 1102, 1, 2, 1044, 1106, 0, 224, 2, 1041, 1043, 1032, 1006, 1032, 179, 1102, 1, 1, 1044, 1106, 0, 224, 1, 1041, 1043, 1032, 1006, 1032, 217, 1, 1042, 1043, 1032, 1001, 1032, -1, 1032, 1002, 1032, 39, 1032, 1, 1032, 1039, 1032, 101, -1, 1032, 1032, 101, 252, 1032, 211, 1007, 0, 59, 1044, 1105, 1, 224, 1102, 1, 0, 1044, 1105, 1, 224, 1006, 1044, 247, 101, 0, 1039, 1034, 1001, 1040, 0, 1035, 101, 0, 1041, 1036, 1002, 1043, 1, 1038, 1002, 1042, 1, 1037, 4, 1044, 1105, 1, 0, 93, 27, 71, 56, 88, 17, 30, 78, 5, 57, 79, 56, 3, 82, 62, 58, 16, 2, 21, 89, 95, 33, 12, 32, 90, 12, 7, 76, 83, 31, 8, 13, 27, 89, 60, 33, 7, 40, 22, 50, 8, 63, 35, 45, 57, 94, 81, 4, 65, 33, 47, 73, 28, 98, 11, 70, 95, 17, 82, 39, 19, 73, 62, 56, 80, 85, 23, 91, 39, 86, 91, 82, 50, 37, 86, 4, 90, 83, 8, 65, 56, 63, 15, 99, 51, 3, 60, 60, 77, 58, 90, 82, 5, 52, 14, 87, 37, 74, 85, 43, 17, 61, 91, 35, 31, 81, 19, 12, 34, 54, 9, 66, 34, 69, 67, 21, 4, 14, 87, 22, 76, 26, 82, 79, 4, 69, 48, 73, 8, 73, 57, 61, 83, 23, 83, 60, 3, 41, 75, 67, 53, 44, 91, 27, 52, 84, 66, 13, 65, 95, 81, 83, 30, 26, 60, 12, 33, 92, 81, 46, 78, 25, 13, 72, 87, 26, 63, 57, 35, 2, 60, 96, 63, 26, 2, 76, 95, 21, 38, 60, 5, 79, 86, 89, 47, 42, 12, 91, 30, 52, 69, 55, 67, 73, 47, 44, 5, 86, 8, 52, 69, 81, 23, 70, 3, 38, 41, 89, 88, 58, 41, 9, 96, 27, 67, 21, 14, 68, 67, 35, 84, 23, 20, 91, 63, 47, 75, 34, 70, 57, 13, 54, 82, 33, 61, 27, 97, 88, 46, 44, 56, 74, 14, 5, 96, 71, 16, 40, 86, 61, 84, 41, 81, 81, 16, 88, 51, 41, 96, 76, 28, 97, 44, 41, 65, 87, 50, 73, 58, 71, 46, 73, 51, 43, 18, 46, 99, 74, 65, 9, 89, 3, 77, 22, 34, 93, 94, 39, 54, 96, 12, 35, 62, 87, 56, 69, 64, 9, 34, 91, 64, 71, 28, 10, 94, 1, 96, 20, 67, 92, 39, 37, 26, 79, 68, 16, 76, 57, 83, 92, 46, 75, 99, 26, 64, 39, 72, 65, 37, 93, 65, 5, 53, 62, 36, 13, 97, 14, 38, 85, 33, 76, 56, 99, 29, 64, 84, 28, 19, 91, 92, 55, 33, 88, 32, 70, 38, 53, 76, 1, 76, 35, 26, 75, 18, 18, 7, 88, 19, 53, 65, 22, 91, 20, 85, 15, 13, 72, 82, 13, 31, 75, 62, 68, 4, 56, 91, 89, 56, 10, 46, 63, 7, 74, 50, 15, 85, 87, 64, 77, 12, 95, 10, 66, 77, 51, 6, 61, 75, 91, 75, 85, 61, 78, 4, 97, 99, 4, 90, 34, 89, 44, 44, 68, 89, 30, 20, 70, 24, 22, 81, 22, 77, 61, 33, 89, 2, 11, 75, 50, 85, 13, 43, 56, 78, 73, 49, 27, 38, 78, 56, 90, 17, 94, 72, 51, 5, 55, 67, 32, 19, 81, 81, 45, 83, 18, 96, 33, 75, 53, 4, 29, 87, 80, 33, 57, 78, 80, 43, 68, 57, 71, 83, 10, 18, 98, 70, 36, 61, 31, 73, 33, 69, 24, 78, 76, 43, 88, 96, 16, 14, 91, 43, 66, 15, 98, 44, 48, 68, 57, 72, 48, 49, 89, 62, 31, 55, 83, 68, 86, 97, 16, 25, 87, 13, 74, 40, 82, 43, 48, 85, 40, 45, 72, 33, 60, 84, 4, 47, 96, 19, 92, 75, 73, 46, 6, 69, 4, 81, 98, 89, 48, 55, 89, 24, 64, 31, 47, 50, 93, 72, 47, 72, 36, 79, 7, 24, 66, 60, 65, 18, 81, 93, 40, 37, 36, 62, 94, 48, 8, 77, 21, 82, 22, 65, 20, 46, 85, 47, 52, 70, 55, 74, 19, 65, 15, 72, 81, 57, 67, 46, 94, 21, 16, 94, 84, 36, 43, 62, 82, 48, 47, 79, 5, 96, 39, 58, 85, 80, 31, 7, 98, 23, 69, 22, 99, 37, 69, 35, 66, 36, 70, 3, 69, 47, 6, 64, 38, 69, 42, 57, 91, 89, 21, 89, 13, 42, 78, 24, 44, 79, 74, 65, 63, 85, 10, 50, 71, 94, 26, 78, 55, 5, 26, 71, 46, 20, 83, 96, 51, 87, 2, 99, 83, 5, 38, 86, 8, 13, 94, 61, 93, 39, 67, 23, 60, 74, 87, 57, 30, 72, 23, 19, 95, 57, 93, 83, 58, 34, 83, 35, 4, 47, 81, 88, 24, 87, 34, 93, 79, 70, 18, 24, 73, 98, 76, 77, 24, 93, 18, 66, 56, 87, 25, 29, 7, 7, 97, 40, 61, 56, 96, 96, 1, 42, 21, 92, 73, 11, 10, 97, 69, 58, 93, 2, 82, 27, 96, 7, 84, 44, 67, 57, 63, 13, 79, 56, 72, 34, 89, 26, 94, 24, 86, 99, 71, 73, 98, 26, 89, 10, 98, 5, 64, 70, 85, 32, 61, 35, 67, 0, 0, 21, 21, 1, 10, 1, 0, 0, 0, 0, 0, 0}
	var intCodeData = makeMapForArray(intCodeDataArr)
	in, out, getShortestPath, getTimeForOxygenSpread := getDroidIOHandlers(0)
	intComp(intCodeData, in, out)

	shortestPath := getShortestPath()
	timeForOxygenSpread := getTimeForOxygenSpread()
	fmt.Println("Length of shortest path:", len(shortestPath))
	fmt.Println("Time for oxygen spread:", timeForOxygenSpread)
}

func messageOfLenAtPoint(messageLen, point int, sequence []string) []string {
	var phases int = 100
	var basePattern = []int64{0, 1, 0, -1}
	var replaceNum int64 = 0
	var addedNumsForPattern int64 = 0
	for i := 0; i < phases; i++ {
		var newSequence = make([]string, len(sequence))
		if point > len(sequence)/2 {
			replaceNum = 0
			addedNumsForPattern = 0
		}
		for replaceIndex := 0; replaceIndex < len(sequence); replaceIndex++ {
			patternIndex := 0
			count := 1
			if point > len(sequence)/2 {
				num, _ := strconv.ParseInt(sequence[(len(sequence)-1)-replaceIndex], 10, 64)
				replaceNum += num
				newSequence[(len(sequence)-1)-replaceIndex] = strconv.Itoa(int(replaceNum % 10))
			} else {
				replaceNum = 0
				addedNumsForPattern = 0
				for _, numString := range sequence {
					num, _ := strconv.ParseInt(numString, 10, 64)
					if count == replaceIndex+1 {
						count = 0
						replaceNum += addedNumsForPattern * basePattern[patternIndex]
						addedNumsForPattern = 0
						patternIndex = (patternIndex + 1) % len(basePattern)
					}

					if patternIndex == 1 || patternIndex == 3 {
						addedNumsForPattern += num
					}
					count++
				}
				// Once more because the loop ends
				replaceNum += addedNumsForPattern * basePattern[patternIndex]
				newSequence[replaceIndex] = strconv.Itoa(int(Abs(replaceNum) % 10))
			}
		}
		sequence = newSequence
	}

	return sequence[point : point+messageLen]
}

func day16() {
	var sequenceString string = "59713137269801099632654181286233935219811755500455380934770765569131734596763695509279561685788856471420060118738307712184666979727705799202164390635688439701763288535574113283975613430058332890215685102656193056939765590473237031584326028162831872694742473094498692690926378560215065112055277042957192884484736885085776095601258138827407479864966595805684283736114104361200511149403415264005242802552220930514486188661282691447267079869746222193563352374541269431531666903127492467446100184447658357579189070698707540721959527692466414290626633017164810627099243281653139996025661993610763947987942741831185002756364249992028050315704531567916821944"
	var sequence = strings.Split(sequenceString, "")

	// Part 1
	fmt.Println("First eight digits after 100 phases:", messageOfLenAtPoint(8, 0, sequence))

	// Part 2
	fmt.Println("Repeating sequence 10,000 times")
	var expandedSequenceString string
	for i := 0; i < 10000; i++ {
		expandedSequenceString += sequenceString
	}
	var expandedSequence = strings.Split(expandedSequenceString, "")
	skipNum, _ := strconv.Atoi(sequenceString[0:7])
	fmt.Println("Running part 2")
	trueMessage := messageOfLenAtPoint(8, skipNum, expandedSequence)

	fmt.Println("First eight digits after skipping to point", skipNum, ":", trueMessage)
}

//Ascii convert that makes reading easier
const (
	SCAFFOLD int64 = 35
	EMPTY    int64 = 46
	NEWLINE  int64 = 10
	A        int64 = 65
	B        int64 = 66
	C        int64 = 67
	L        int64 = 76
	R        int64 = 82
	COMMA    int64 = 44
	yChar    int64 = 121
	nChar    int64 = 110
)

func scaffoldBotIOHandlers() (inputHandler, outputHandler, func() (map[Point]int64, int64, int64)) {
	var screen = make(map[Point]int64)
	var x, y, width, height int64

	//A,C,C,B,B,A,A,C,C,B,NEWLINE
	mainProgram := []int64{A, COMMA, C, COMMA, C, COMMA, B, COMMA, B, COMMA, A, COMMA, A, COMMA, C, COMMA, C, COMMA, B, NEWLINE}

	//L,12,R,4,R,4,NEWLINE
	aProgram := []int64{L, COMMA, 49, 50, COMMA, R, COMMA, 52, COMMA, R, COMMA, 52, NEWLINE}
	//R,12,R,4,L,6,L,8,L,8,NEWLINE
	bProgram := []int64{R, COMMA, 49, 50, COMMA, R, COMMA, 52, COMMA, L, COMMA, 54, COMMA, L, COMMA, 56, COMMA, L, COMMA, 56, NEWLINE}
	//R,12,R,4,L,12,NEWLINE
	cProgram := []int64{R, COMMA, 49, 50, COMMA, R, COMMA, 52, COMMA, L, COMMA, 49, 50, NEWLINE}

	yesNoProgram := []int64{nChar, NEWLINE}
	currentProgram := 0
	currentOperation := 0
	in := func() int64 {
		var returnCode int64
		switch currentProgram {
		case 0:
			returnCode = mainProgram[currentOperation]
		case 1:
			returnCode = aProgram[currentOperation]
		case 2:
			returnCode = bProgram[currentOperation]
		case 3:
			returnCode = cProgram[currentOperation]
		case 4:
			returnCode = yesNoProgram[currentOperation]
		}
		currentOperation++

		if returnCode == NEWLINE {
			currentOperation = 0
			currentProgram++
		}

		return returnCode
	}

	out := func(output int64) {
		screen[Point{x, y}] = output
		if x > width {
			width = x
		}

		if y > height {
			height = y
		}

		if output == NEWLINE {
			y++
			x = 0
		} else {
			x++
		}

		if output > 128 {
			fmt.Println("Dust collected:", output)
		}
	}

	getScreen := func() (map[Point]int64, int64, int64) {
		return screen, width, height
	}

	return in, out, getScreen
}

func day17() {
	var scaffoldBotProgramArr = []int64{1, 330, 331, 332, 109, 3752, 1101, 1182, 0, 16, 1102, 1451, 1, 24, 102, 1, 0, 570, 1006, 570, 36, 101, 0, 571, 0, 1001, 570, -1, 570, 1001, 24, 1, 24, 1106, 0, 18, 1008, 571, 0, 571, 1001, 16, 1, 16, 1008, 16, 1451, 570, 1006, 570, 14, 21101, 0, 58, 0, 1105, 1, 786, 1006, 332, 62, 99, 21102, 333, 1, 1, 21102, 73, 1, 0, 1105, 1, 579, 1102, 1, 0, 572, 1101, 0, 0, 573, 3, 574, 101, 1, 573, 573, 1007, 574, 65, 570, 1005, 570, 151, 107, 67, 574, 570, 1005, 570, 151, 1001, 574, -64, 574, 1002, 574, -1, 574, 1001, 572, 1, 572, 1007, 572, 11, 570, 1006, 570, 165, 101, 1182, 572, 127, 1002, 574, 1, 0, 3, 574, 101, 1, 573, 573, 1008, 574, 10, 570, 1005, 570, 189, 1008, 574, 44, 570, 1006, 570, 158, 1106, 0, 81, 21102, 340, 1, 1, 1106, 0, 177, 21102, 1, 477, 1, 1105, 1, 177, 21102, 514, 1, 1, 21101, 0, 176, 0, 1106, 0, 579, 99, 21101, 184, 0, 0, 1105, 1, 579, 4, 574, 104, 10, 99, 1007, 573, 22, 570, 1006, 570, 165, 1002, 572, 1, 1182, 21101, 0, 375, 1, 21101, 0, 211, 0, 1106, 0, 579, 21101, 1182, 11, 1, 21101, 0, 222, 0, 1105, 1, 979, 21101, 0, 388, 1, 21101, 233, 0, 0, 1106, 0, 579, 21101, 1182, 22, 1, 21102, 1, 244, 0, 1106, 0, 979, 21102, 401, 1, 1, 21101, 0, 255, 0, 1106, 0, 579, 21101, 1182, 33, 1, 21102, 1, 266, 0, 1105, 1, 979, 21101, 414, 0, 1, 21102, 1, 277, 0, 1106, 0, 579, 3, 575, 1008, 575, 89, 570, 1008, 575, 121, 575, 1, 575, 570, 575, 3, 574, 1008, 574, 10, 570, 1006, 570, 291, 104, 10, 21102, 1, 1182, 1, 21102, 1, 313, 0, 1105, 1, 622, 1005, 575, 327, 1102, 1, 1, 575, 21101, 327, 0, 0, 1105, 1, 786, 4, 438, 99, 0, 1, 1, 6, 77, 97, 105, 110, 58, 10, 33, 10, 69, 120, 112, 101, 99, 116, 101, 100, 32, 102, 117, 110, 99, 116, 105, 111, 110, 32, 110, 97, 109, 101, 32, 98, 117, 116, 32, 103, 111, 116, 58, 32, 0, 12, 70, 117, 110, 99, 116, 105, 111, 110, 32, 65, 58, 10, 12, 70, 117, 110, 99, 116, 105, 111, 110, 32, 66, 58, 10, 12, 70, 117, 110, 99, 116, 105, 111, 110, 32, 67, 58, 10, 23, 67, 111, 110, 116, 105, 110, 117, 111, 117, 115, 32, 118, 105, 100, 101, 111, 32, 102, 101, 101, 100, 63, 10, 0, 37, 10, 69, 120, 112, 101, 99, 116, 101, 100, 32, 82, 44, 32, 76, 44, 32, 111, 114, 32, 100, 105, 115, 116, 97, 110, 99, 101, 32, 98, 117, 116, 32, 103, 111, 116, 58, 32, 36, 10, 69, 120, 112, 101, 99, 116, 101, 100, 32, 99, 111, 109, 109, 97, 32, 111, 114, 32, 110, 101, 119, 108, 105, 110, 101, 32, 98, 117, 116, 32, 103, 111, 116, 58, 32, 43, 10, 68, 101, 102, 105, 110, 105, 116, 105, 111, 110, 115, 32, 109, 97, 121, 32, 98, 101, 32, 97, 116, 32, 109, 111, 115, 116, 32, 50, 48, 32, 99, 104, 97, 114, 97, 99, 116, 101, 114, 115, 33, 10, 94, 62, 118, 60, 0, 1, 0, -1, -1, 0, 1, 0, 0, 0, 0, 0, 0, 1, 58, 18, 0, 109, 4, 2102, 1, -3, 587, 20101, 0, 0, -1, 22101, 1, -3, -3, 21102, 1, 0, -2, 2208, -2, -1, 570, 1005, 570, 617, 2201, -3, -2, 609, 4, 0, 21201, -2, 1, -2, 1105, 1, 597, 109, -4, 2105, 1, 0, 109, 5, 2102, 1, -4, 630, 20101, 0, 0, -2, 22101, 1, -4, -4, 21101, 0, 0, -3, 2208, -3, -2, 570, 1005, 570, 781, 2201, -4, -3, 653, 20101, 0, 0, -1, 1208, -1, -4, 570, 1005, 570, 709, 1208, -1, -5, 570, 1005, 570, 734, 1207, -1, 0, 570, 1005, 570, 759, 1206, -1, 774, 1001, 578, 562, 684, 1, 0, 576, 576, 1001, 578, 566, 692, 1, 0, 577, 577, 21102, 1, 702, 0, 1106, 0, 786, 21201, -1, -1, -1, 1106, 0, 676, 1001, 578, 1, 578, 1008, 578, 4, 570, 1006, 570, 724, 1001, 578, -4, 578, 21102, 731, 1, 0, 1105, 1, 786, 1105, 1, 774, 1001, 578, -1, 578, 1008, 578, -1, 570, 1006, 570, 749, 1001, 578, 4, 578, 21101, 0, 756, 0, 1106, 0, 786, 1106, 0, 774, 21202, -1, -11, 1, 22101, 1182, 1, 1, 21101, 0, 774, 0, 1106, 0, 622, 21201, -3, 1, -3, 1106, 0, 640, 109, -5, 2105, 1, 0, 109, 7, 1005, 575, 802, 20101, 0, 576, -6, 21001, 577, 0, -5, 1106, 0, 814, 21102, 0, 1, -1, 21102, 0, 1, -5, 21101, 0, 0, -6, 20208, -6, 576, -2, 208, -5, 577, 570, 22002, 570, -2, -2, 21202, -5, 59, -3, 22201, -6, -3, -3, 22101, 1451, -3, -3, 2102, 1, -3, 843, 1005, 0, 863, 21202, -2, 42, -4, 22101, 46, -4, -4, 1206, -2, 924, 21102, 1, 1, -1, 1105, 1, 924, 1205, -2, 873, 21101, 0, 35, -4, 1105, 1, 924, 1202, -3, 1, 878, 1008, 0, 1, 570, 1006, 570, 916, 1001, 374, 1, 374, 2101, 0, -3, 895, 1101, 0, 2, 0, 2101, 0, -3, 902, 1001, 438, 0, 438, 2202, -6, -5, 570, 1, 570, 374, 570, 1, 570, 438, 438, 1001, 578, 558, 922, 20102, 1, 0, -4, 1006, 575, 959, 204, -4, 22101, 1, -6, -6, 1208, -6, 59, 570, 1006, 570, 814, 104, 10, 22101, 1, -5, -5, 1208, -5, 39, 570, 1006, 570, 810, 104, 10, 1206, -1, 974, 99, 1206, -1, 974, 1102, 1, 1, 575, 21102, 973, 1, 0, 1106, 0, 786, 99, 109, -7, 2105, 1, 0, 109, 6, 21101, 0, 0, -4, 21102, 0, 1, -3, 203, -2, 22101, 1, -3, -3, 21208, -2, 82, -1, 1205, -1, 1030, 21208, -2, 76, -1, 1205, -1, 1037, 21207, -2, 48, -1, 1205, -1, 1124, 22107, 57, -2, -1, 1205, -1, 1124, 21201, -2, -48, -2, 1106, 0, 1041, 21101, 0, -4, -2, 1106, 0, 1041, 21101, -5, 0, -2, 21201, -4, 1, -4, 21207, -4, 11, -1, 1206, -1, 1138, 2201, -5, -4, 1059, 2102, 1, -2, 0, 203, -2, 22101, 1, -3, -3, 21207, -2, 48, -1, 1205, -1, 1107, 22107, 57, -2, -1, 1205, -1, 1107, 21201, -2, -48, -2, 2201, -5, -4, 1090, 20102, 10, 0, -1, 22201, -2, -1, -2, 2201, -5, -4, 1103, 1202, -2, 1, 0, 1106, 0, 1060, 21208, -2, 10, -1, 1205, -1, 1162, 21208, -2, 44, -1, 1206, -1, 1131, 1105, 1, 989, 21102, 439, 1, 1, 1105, 1, 1150, 21101, 477, 0, 1, 1105, 1, 1150, 21101, 0, 514, 1, 21102, 1, 1149, 0, 1106, 0, 579, 99, 21102, 1, 1157, 0, 1106, 0, 579, 204, -2, 104, 10, 99, 21207, -3, 22, -1, 1206, -1, 1138, 2102, 1, -5, 1176, 2101, 0, -4, 0, 109, -6, 2105, 1, 0, 24, 13, 46, 1, 11, 1, 46, 1, 11, 1, 46, 1, 11, 1, 46, 1, 11, 13, 34, 1, 23, 1, 34, 1, 23, 1, 34, 1, 23, 1, 34, 1, 23, 1, 34, 1, 23, 1, 34, 1, 23, 1, 34, 1, 23, 1, 30, 5, 23, 1, 30, 1, 27, 1, 30, 1, 25, 5, 1, 1, 26, 1, 25, 1, 1, 1, 1, 1, 1, 1, 14, 5, 5, 9, 17, 5, 1, 1, 1, 1, 14, 1, 3, 1, 5, 1, 1, 1, 5, 1, 17, 1, 1, 1, 3, 1, 1, 1, 14, 1, 3, 1, 5, 1, 1, 1, 5, 1, 17, 1, 1, 13, 8, 1, 3, 1, 5, 1, 1, 1, 5, 1, 17, 1, 5, 1, 1, 1, 6, 7, 1, 1, 3, 13, 1, 1, 17, 1, 5, 1, 1, 1, 6, 1, 5, 1, 1, 1, 9, 1, 1, 1, 3, 1, 1, 1, 17, 1, 5, 1, 1, 1, 6, 1, 5, 1, 1, 1, 9, 1, 1, 1, 1, 5, 17, 9, 6, 1, 5, 1, 1, 1, 9, 1, 1, 1, 1, 1, 1, 1, 25, 1, 8, 1, 5, 13, 1, 5, 25, 1, 8, 1, 7, 1, 13, 1, 27, 1, 8, 1, 7, 1, 13, 1, 23, 5, 8, 1, 7, 1, 13, 1, 23, 1, 12, 9, 13, 1, 23, 1, 34, 1, 23, 1, 34, 1, 23, 1, 34, 1, 23, 1, 34, 1, 23, 1, 34, 1, 23, 1, 34, 13, 11, 1, 46, 1, 11, 1, 46, 1, 11, 1, 46, 1, 11, 1, 46, 13, 12}
	var scaffoldBotProgram = makeMapForArray(scaffoldBotProgramArr)
	var alignParamSum int64 = 0

	in, out, getScreen := scaffoldBotIOHandlers()
	intComp(scaffoldBotProgram, in, out)

	screen, width, height := getScreen()
	var x, y int64
	for y = 0; y <= height; y++ {
		for x = 0; x <= width; x++ {
			current := string(screen[Point{x, y}])
			l := string(screen[Point{x - 1, y}])
			r := string(screen[Point{x + 1, y}])
			u := string(screen[Point{x, y + 1}])
			d := string(screen[Point{x, y - 1}])

			if l == "#" && r == "#" && u == "#" && d == "#" && current == "#" {
				alignParamSum += x * y
				fmt.Print("O")
				continue
			}
			fmt.Print(current)
		}
	}
	fmt.Println("\nAlign parameter sum:", alignParamSum)
	scaffoldBotProgram[0] = 2
	intComp(scaffoldBotProgram, in, out)
}

type Route struct {
	routeName string
	pos       Point
	length    int
}

func distancesFrom(source Point, mazeMap map[Point]droidTile, keys, doors map[string]bool) map[string]Route {
	availableDirs := []int64{NORTH, SOUTH, WEST, EAST}
	var mazeCopy = make(map[Point]droidTile)
	for key, val := range mazeMap {
		mazeCopy[key] = val
	}
	routeInfo := make(map[string]Route)

	var queue = []Route{Route{"", source, 0}}
	for routeIndex := 0; routeIndex < len(queue); routeIndex++ {
		route := queue[routeIndex]
		tileAtPos := mazeCopy[route.pos]
		if (keys[tileAtPos.tile] || doors[tileAtPos.tile]) && route.length != 0 {
			// fmt.Println(tileAtPos.tile, "found at distance", route.length, "with route:", route.routeName)
			routeInfo[tileAtPos.tile] = route
			route.routeName += tileAtPos.tile
		}
		tileAtPos.traveled = true
		mazeCopy[route.pos] = tileAtPos

		for _, dir := range availableDirs {
			dirPoint := getPointForDirection(dir, route.pos)
			dirTile := mazeCopy[dirPoint]
			if dirTile.tile != "#" && !dirTile.traveled {
				queue = append(queue, Route{route.routeName, dirPoint, route.length + 1})
			}
		}
	}
	return routeInfo
}

func findRouteInfo(width, height int64, mazeMap map[Point]droidTile, keys, doors map[string]bool) map[string]map[string]Route {
	routeInfo := make(map[string]map[string]Route)
	var x, y int64
	for y = 0; y < height; y++ {
		for x = 0; x < width; x++ {
			tileAtPos := mazeMap[Point{x, y}]
			if keys[tileAtPos.tile] || strings.Contains("@1234", tileAtPos.tile) {
				routeInfo[tileAtPos.tile] = distancesFrom(Point{x, y}, mazeMap, keys, doors)
			}
		}
	}
	return routeInfo
}

func canReachKey(currentKeys string, route Route) bool {
	for _, routeItem := range route.routeName {
		if !strings.Contains(currentKeys, string(routeItem)) && !strings.Contains(currentKeys, strings.ToLower(string(routeItem))) {
			return false
		}
	}
	return true
}

func day18() {
	file, err := os.Open("day18Tunnels.txt")
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	//Yoinking the droidTile struct because it has what I need
	var mazeMap = make(map[Point]droidTile)
	var doors = make(map[string]bool)
	var keys = make(map[string]bool)
	var startPos Point
	var x, y, width, height int64
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		for index := range line {
			if line[index] == '@' {
				startPos = Point{x, y}
			}
			if line[index] != '#' && line[index] != '@' && line[index] != '.' {
				stringChar := string(line[index])
				if stringChar == strings.ToUpper(stringChar) {
					doors[stringChar] = true
				}
				if stringChar == strings.ToLower(stringChar) {
					keys[stringChar] = true
				}
			}

			if mazeMap[Point{x, y}].tile == "" {
				mazeMap[Point{x, y}] = droidTile{string(line[index]), false}
			}

			// x and y start at 0, so add 1 to account for that.
			if x >= width {
				width = x + 1
			}
			if y >= height {
				height = y + 1
			}
			x++
		}
		x = 0
		y++
	}

	var part1ShortestSteps = math.MaxInt64
	routeInfo := findRouteInfo(width, height, mazeMap, keys, doors)
	// I cant use a struct with a map or an array in it as a key,
	// so Im using a formatted string
	infoStartString := "@:"
	info := make(map[string]int)
	info[infoStartString] = 0

	for range keys {
		newInfo := make(map[string]int)
		for dataString, currentDist := range info {
			currentLocation := string(dataString[0])
			currentKeys := dataString[2:]
			for newKey := range keys {
				if !strings.Contains(currentKeys, newKey) {
					route := routeInfo[currentLocation][newKey]
					reachable := canReachKey(currentKeys, route)
					if reachable {
						newDistance := currentDist + route.length
						newKeys := currentKeys + newKey
						newKeysList := strings.Split(newKeys, "")
						sort.Strings(newKeysList)
						newKeys = strings.Join(newKeysList, "")
						newDataString := newKey + ":" + newKeys

						if newInfo[newDataString] == 0 || newDistance < newInfo[newDataString] {
							newInfo[newDataString] = newDistance
						}
					}
				}
			}
		}
		info = newInfo
	}

	for _, distance := range info {
		if distance < part1ShortestSteps {
			part1ShortestSteps = distance
		}
	}
	part1FinalPositions := len(info)

	// Going to put these here because both parts take a smol bit,
	// and its good to know its doing something.
	fmt.Println("There are", part1FinalPositions, "final positions for part 1")
	fmt.Println("Shortest steps for part 1:", part1ShortestSteps)

	// Set mazeMap up for part 2 with its four robots.
	mazeMap[startPos] = droidTile{"#", false}
	mazeMap[Point{startPos.x - 1, startPos.y}] = droidTile{"#", false}
	mazeMap[Point{startPos.x + 1, startPos.y}] = droidTile{"#", false}
	mazeMap[Point{startPos.x, startPos.y - 1}] = droidTile{"#", false}
	mazeMap[Point{startPos.x, startPos.y + 1}] = droidTile{"#", false}

	mazeMap[Point{startPos.x - 1, startPos.y - 1}] = droidTile{"1", false}
	mazeMap[Point{startPos.x + 1, startPos.y - 1}] = droidTile{"2", false}
	mazeMap[Point{startPos.x - 1, startPos.y + 1}] = droidTile{"3", false}
	mazeMap[Point{startPos.x + 1, startPos.y + 1}] = droidTile{"4", false}
	var part2ShortestSteps = math.MaxInt64
	routeInfo = findRouteInfo(width, height, mazeMap, keys, doors)
	// A more differenter formatted string
	infoStartString = "1,2,3,4:"
	info = make(map[string]int)
	info[infoStartString] = 0
	for range keys {
		newInfo := make(map[string]int)
		for dataString, currentDist := range info {
			currentLocationsString := string(dataString[:7])
			currentLocations := strings.Split(currentLocationsString, ",")
			currentKeys := dataString[8:]
			for newKey := range keys {
				if !strings.Contains(currentKeys, newKey) {
					for robot, currentLocation := range currentLocations {
						if routeInfo[currentLocation][newKey].length != 0 {
							route := routeInfo[currentLocation][newKey]
							reachable := canReachKey(currentKeys, route)

							if reachable {
								newDistance := currentDist + route.length
								// Get new sorted keys
								newKeys := currentKeys + newKey
								newKeysList := strings.Split(newKeys, "")
								sort.Strings(newKeysList)
								newKeys = strings.Join(newKeysList, "")

								// Get new locations. Copy locations to a new list so its not a pointer.
								newLocations := make([]string, len(currentLocations))
								copy(newLocations, currentLocations)
								newLocations[robot] = newKey
								newLocationsString := strings.Join(newLocations, ",")
								newDataString := newLocationsString + ":" + newKeys

								if newInfo[newDataString] == 0 || newDistance < newInfo[newDataString] {
									newInfo[newDataString] = newDistance
								}
							}
						}
					}
				}
			}
		}
		info = newInfo
	}
	for _, distance := range info {
		if distance < part2ShortestSteps {
			part2ShortestSteps = distance
		}
	}
	part2FinalPositions := len(info)

	fmt.Println("There are", part2FinalPositions, "final positions for part 2")
	fmt.Println("Shortest steps for part 2:", part2ShortestSteps)
}

func tractorBeamDroneIOHandlers(x, y int64) (inputHandler, outputHandler, func() bool) {
	var sentX bool = false
	var pulled = false
	in := func() int64 {
		if !sentX {
			sentX = true
			return x
		}
		return y
	}
	out := func(output int64) {
		if output == 1 {
			pulled = true
		}
	}

	getPulled := func() bool {
		return pulled
	}
	return in, out, getPulled
}

func day19() {
	var tractorBeamDroneProgramArr = []int64{109, 424, 203, 1, 21101, 11, 0, 0, 1105, 1, 282, 21102, 18, 1, 0, 1105, 1, 259, 2102, 1, 1, 221, 203, 1, 21102, 1, 31, 0, 1106, 0, 282, 21101, 38, 0, 0, 1105, 1, 259, 21001, 23, 0, 2, 21201, 1, 0, 3, 21101, 0, 1, 1, 21101, 0, 57, 0, 1105, 1, 303, 1201, 1, 0, 222, 20102, 1, 221, 3, 20101, 0, 221, 2, 21101, 259, 0, 1, 21102, 80, 1, 0, 1106, 0, 225, 21101, 127, 0, 2, 21102, 91, 1, 0, 1106, 0, 303, 1201, 1, 0, 223, 20102, 1, 222, 4, 21101, 259, 0, 3, 21101, 0, 225, 2, 21102, 225, 1, 1, 21102, 1, 118, 0, 1106, 0, 225, 21001, 222, 0, 3, 21101, 0, 89, 2, 21101, 133, 0, 0, 1105, 1, 303, 21202, 1, -1, 1, 22001, 223, 1, 1, 21101, 0, 148, 0, 1105, 1, 259, 2102, 1, 1, 223, 21002, 221, 1, 4, 21001, 222, 0, 3, 21101, 0, 21, 2, 1001, 132, -2, 224, 1002, 224, 2, 224, 1001, 224, 3, 224, 1002, 132, -1, 132, 1, 224, 132, 224, 21001, 224, 1, 1, 21102, 195, 1, 0, 106, 0, 108, 20207, 1, 223, 2, 20102, 1, 23, 1, 21102, 1, -1, 3, 21101, 0, 214, 0, 1105, 1, 303, 22101, 1, 1, 1, 204, 1, 99, 0, 0, 0, 0, 109, 5, 1201, -4, 0, 249, 22102, 1, -3, 1, 21201, -2, 0, 2, 22101, 0, -1, 3, 21102, 250, 1, 0, 1105, 1, 225, 21202, 1, 1, -4, 109, -5, 2105, 1, 0, 109, 3, 22107, 0, -2, -1, 21202, -1, 2, -1, 21201, -1, -1, -1, 22202, -1, -2, -2, 109, -3, 2106, 0, 0, 109, 3, 21207, -2, 0, -1, 1206, -1, 294, 104, 0, 99, 22101, 0, -2, -2, 109, -3, 2106, 0, 0, 109, 5, 22207, -3, -4, -1, 1206, -1, 346, 22201, -4, -3, -4, 21202, -3, -1, -1, 22201, -4, -1, 2, 21202, 2, -1, -1, 22201, -4, -1, 1, 21201, -2, 0, 3, 21101, 0, 343, 0, 1106, 0, 303, 1105, 1, 415, 22207, -2, -3, -1, 1206, -1, 387, 22201, -3, -2, -3, 21202, -2, -1, -1, 22201, -3, -1, 3, 21202, 3, -1, -1, 22201, -3, -1, 2, 22101, 0, -4, 1, 21101, 384, 0, 0, 1106, 0, 303, 1105, 1, 415, 21202, -4, -1, -4, 22201, -4, -3, -4, 22202, -3, -2, -2, 22202, -2, -4, -4, 22202, -3, -2, -3, 21202, -4, -1, -2, 22201, -3, -2, 1, 21201, 1, 0, -4, 109, -5, 2105, 1, 0}
	var tractorBeamDroneProgram = makeMapForArray(tractorBeamDroneProgramArr)
	pointsAffected := 0
	var x, y int64
	var sizeX int64 = 100
	var sizeY int64 = 100
	var width int64 = 50
	var height int64 = 50
	var wantedXCoord int64
	var wantedYCoord int64
	for y = 0; y < height; y++ {
		for x = 0; x < width; x++ {
			in, out, getPulled := tractorBeamDroneIOHandlers(x, y)
			intComp(tractorBeamDroneProgram, in, out)
			isPulled := getPulled()
			if isPulled {
				fmt.Print("#")
				pointsAffected++
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
	fmt.Println("Points affected in 50x50 square:", pointsAffected)
	var gottenXYClosestCoords bool = false
	for !gottenXYClosestCoords {

		//Subtract 1 because the coordinates start at 0
		in, out, getPulled := tractorBeamDroneIOHandlers(x, y+sizeY-1)
		intComp(tractorBeamDroneProgram, in, out)
		isPulledBottomLeft := getPulled()

		in, out, getPulled = tractorBeamDroneIOHandlers(x+sizeX-1, y)
		intComp(tractorBeamDroneProgram, in, out)
		isPulledTopRight := getPulled()

		if isPulledBottomLeft && isPulledTopRight {
			wantedXCoord = x
			wantedYCoord = y
			gottenXYClosestCoords = true
		}

		if isPulledBottomLeft {
			x -= sizeX + 1
			y++
		} else {
			x++
		}
	}
	fmt.Println("Wanted coords to fit size of 100x100:", (wantedXCoord*10000)+wantedYCoord)
}

type teleporter struct {
	pointA Point
	pointB Point
	name   string
}

func day20() {
	file, err := os.Open("day20DonutMaze.txt")
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	var mazeMap = make(map[Point]string)
	var teleNameMap = make(map[Point]string)
	var teleMap = make(map[string]teleporter)
	scanner := bufio.NewScanner(file)

	var x, y, width, height int64
	for scanner.Scan() {
		line := scanner.Text()
		for _, char := range line {
			if char != ' ' && char != '#' && char != '.' {
				northPoint := mazeMap[Point{x, y - 1}]
				westPoint := mazeMap[Point{x - 1, y}]
				var telePoint Point
				var teleName string
				if northPoint != "" && northPoint != " " && northPoint != "#" && northPoint != "." {
					//Look up two, to check for a period.
					if mazeMap[Point{x, y - 2}] == "." {
						telePoint = Point{x, y - 2}
					} else {
						telePoint = Point{x, y + 1}
					}
					teleName = northPoint + string(char)
				}

				if westPoint != "" && westPoint != " " && westPoint != "#" && westPoint != "." {
					//Look back two, to check for a period.
					if mazeMap[Point{x - 2, y}] == "." {
						telePoint = Point{x - 2, y}
					} else {
						telePoint = Point{x + 1, y}
					}
					teleName = westPoint + string(char)
				}

				//Set the teleNameMap and teleMap so I can get the telporters later
				if telePoint != (Point{0, 0}) {
					teleNameMap[telePoint] = teleName
					tele := teleMap[teleName]
					if tele.pointA == (Point{0, 0}) {
						tele.pointA = telePoint
					} else {
						tele.pointB = telePoint
					}
					teleMap[teleName] = tele
				}
			}

			mazeMap[Point{x, y}] = string(char)
			width = x
			height = y
			x++
		}
		x = 0
		y++
	}

	fmt.Println("Done parsing maze...")

	var shortestPathWithoutLevels = []Point{}
	availableDirs := []int64{NORTH, SOUTH, WEST, EAST}

	// AA and ZZ are treated like teleporters,
	// so I can loop them up that way.
	// They only have a pointA.
	var startPos = teleMap["AA"].pointA

	//Take point, return parent
	var breadthFirstMap = make(map[Point]Point)
	breadthFirstMap[startPos] = Point{-1, -1}

	//Make it so startPos is traveled so it doesnt try to move over it and overwrite the parent
	var traveledMap = make(map[Point]bool)
	traveledMap[startPos] = true

	var canGotoQueue = []Point{}
	canGotoQueue = append(canGotoQueue, startPos)

	// This is the part 1 breadth first search.
	// I was going to reduce the amount of code,
	// and try and do both parts at the same time,
	// but I wanted to make it easy to read and understand.
	for len(canGotoQueue) != 0 {
		var localQueue = make([]Point, len(canGotoQueue))
		copy(localQueue, canGotoQueue)
		canGotoQueue = []Point{}
		for _, point := range localQueue {
			for _, dir := range availableDirs {
				dirPoint := getPointForDirection(dir, point)
				dirTile := mazeMap[dirPoint]

				if teleNameMap[dirPoint] != "" && !traveledMap[dirPoint] {
					breadthFirstMap[dirPoint] = point
					teleName := teleNameMap[dirPoint]
					if teleName == "ZZ" {
						fmt.Println("Found end of maze without levels...")
						currentParent := breadthFirstMap[dirPoint]
						for currentParent != (Point{-1, -1}) {
							shortestPathWithoutLevels = append(shortestPathWithoutLevels, currentParent)
							currentParent = breadthFirstMap[currentParent]
						}
					} else {
						// fmt.Println("Using teleporter: ", teleName)
						tele := teleMap[teleName]
						if dirPoint != tele.pointA {
							breadthFirstMap[tele.pointA] = dirPoint
							traveledMap[tele.pointA] = true
							canGotoQueue = append(canGotoQueue, tele.pointA)
						} else {
							breadthFirstMap[tele.pointB] = dirPoint
							traveledMap[tele.pointB] = true
							canGotoQueue = append(canGotoQueue, tele.pointB)
						}
						traveledMap[dirPoint] = true
						mazeMap[dirPoint] = dirTile
					}
				}

				if dirTile == "." && !traveledMap[dirPoint] {
					traveledMap[dirPoint] = true
					breadthFirstMap[dirPoint] = point
					canGotoQueue = append(canGotoQueue, dirPoint)
				}
			}
		}
	}

	var shortestPathWithLevels = []xyzPoint{}
	var leftAndTopEdgeOffset int64 = 2
	var rightEdgeOffset int64 = width - 2
	var bottomEdgeOffset int64 = height - 2
	var lowestLevel int64 = 0

	var xyzStartPos = xyzPoint{teleMap["AA"].pointA.x, teleMap["AA"].pointA.y, 0}

	//Take point, return parent
	var xyzBreadthFirstMap = make(map[xyzPoint]xyzPoint)
	xyzBreadthFirstMap[xyzStartPos] = xyzPoint{-1, -1, 0}

	//Make it so startPos is traveled so it doesnt try to move over it and overwrite the parent
	var xyzTraveledMap = make(map[xyzPoint]bool)
	xyzTraveledMap[xyzStartPos] = true

	var xyzCanGotoQueue = []xyzPoint{}
	xyzCanGotoQueue = append(xyzCanGotoQueue, xyzStartPos)

	for len(xyzCanGotoQueue) != 0 && len(shortestPathWithLevels) == 0 {
		var localQueue = make([]xyzPoint, len(xyzCanGotoQueue))
		copy(localQueue, xyzCanGotoQueue)
		xyzCanGotoQueue = []xyzPoint{}
		for _, point := range localQueue {
			for _, dir := range availableDirs {
				dirXYPoint := getPointForDirection(dir, Point{point.x, point.y})
				dirPoint := xyzPoint{dirXYPoint.x, dirXYPoint.y, point.z}
				dirTile := mazeMap[dirXYPoint]

				if teleNameMap[dirXYPoint] != "" && dirTile == "." && !xyzTraveledMap[dirPoint] {
					xyzBreadthFirstMap[dirPoint] = point
					teleName := teleNameMap[dirXYPoint]
					if teleName == "ZZ" {
						// fmt.Println("Found zz at level:", dirPoint.z)
						if point.z == 0 {
							fmt.Println("Found end of maze with levels...")
							currentParent := xyzBreadthFirstMap[dirPoint]
							for currentParent != (xyzPoint{-1, -1, 0}) {
								if currentParent.z > lowestLevel {
									lowestLevel = currentParent.z
								}
								shortestPathWithLevels = append(shortestPathWithLevels, currentParent)

								currentParent = xyzBreadthFirstMap[currentParent]
							}
						}
					} else if teleName != "AA" {
						tele := teleMap[teleName]
						var teleXYZPoint xyzPoint
						var level = point.z
						if dirXYPoint == tele.pointB {
							// Because pointA is not equal, we are at pointB.
							// So if pointA is on the outer edge,
							// then we are going one level deeper at pointB on the inner edge,
							// so add one to level.
							if tele.pointA.x == leftAndTopEdgeOffset || tele.pointA.y == leftAndTopEdgeOffset ||
								tele.pointA.x == rightEdgeOffset || tele.pointA.y == bottomEdgeOffset {
								level++
							} else {
								level--
							}

							if level >= 0 && level < int64(len(teleNameMap)) {
								teleXYZPoint = xyzPoint{tele.pointA.x, tele.pointA.y, level}
							}
						} else {
							if tele.pointB.x == leftAndTopEdgeOffset || tele.pointB.y == leftAndTopEdgeOffset ||
								tele.pointB.x == rightEdgeOffset || tele.pointB.y == bottomEdgeOffset {
								level++
							} else {
								level--
							}

							if level >= 0 && level < int64(len(teleNameMap)) {
								teleXYZPoint = xyzPoint{tele.pointB.x, tele.pointB.y, level}
							}
						}
						//Dont check z axis because level can be 0
						if teleXYZPoint.x != 0 && teleXYZPoint.y != 0 {
							// fmt.Println("Using teleporter: ", teleName, "from level:", point.z, "to go to level:", level)
							xyzTraveledMap[dirPoint] = true
							xyzBreadthFirstMap[teleXYZPoint] = dirPoint
							xyzTraveledMap[teleXYZPoint] = true
							xyzCanGotoQueue = append(xyzCanGotoQueue, teleXYZPoint)
						}
					}
				}

				if dirTile == "." && !xyzTraveledMap[dirPoint] {
					xyzTraveledMap[dirPoint] = true
					xyzBreadthFirstMap[dirPoint] = point
					xyzCanGotoQueue = append(xyzCanGotoQueue, dirPoint)
				}
			}
		}
	}

	fmt.Println("Shortest path steps without levels:", len(shortestPathWithoutLevels))
	fmt.Println("Shortest path steps with levels:", len(shortestPathWithLevels))
	fmt.Println("Deepest level shortest path goes to:", lowestLevel)
	fmt.Println("Number of teleporters:", len(teleMap))

}

const (
	D     int64 = 68
	E     int64 = 69
	F     int64 = 70
	G     int64 = 71
	H     int64 = 72
	I     int64 = 73
	J     int64 = 74
	K     int64 = 75
	N     int64 = 78
	O     int64 = 79
	T     int64 = 84
	U     int64 = 85
	W     int64 = 87
	SPACE int64 = 32
)

func springDroidIOHandlers(shouldRun bool) (inputHandler, outputHandler) {
	var currentChar = 0
	var droidProgram []int64

	if !shouldRun {
		// Check if spaces A, B, and C, exist.
		// Only if they all exist T will be set to true.
		// If T is false, one of those three is gone.
		// Wait until we can land on space D to jump.
		droidProgram = []int64{
			O, R, SPACE, A, SPACE, T, NEWLINE,
			A, N, D, SPACE, B, SPACE, T, NEWLINE,
			A, N, D, SPACE, C, SPACE, T, NEWLINE,
			N, O, T, SPACE, T, SPACE, J, NEWLINE,
			A, N, D, SPACE, D, SPACE, J, NEWLINE,
			W, A, L, K, NEWLINE}
	} else {
		droidProgram = []int64{
			// E lines up jumps that are off a little bit,
			// and H prevents E from killing the bot
			// AND D J could be moved above the E H lines,
			// because the droid has memory from its last state,
			// but this is easier to read.
			O, R, SPACE, A, SPACE, T, NEWLINE,
			A, N, D, SPACE, B, SPACE, T, NEWLINE,
			A, N, D, SPACE, C, SPACE, T, NEWLINE,
			N, O, T, SPACE, T, SPACE, J, NEWLINE,

			O, R, SPACE, E, SPACE, T, NEWLINE,
			O, R, SPACE, H, SPACE, T, NEWLINE,
			A, N, D, SPACE, T, SPACE, J, NEWLINE,
			A, N, D, SPACE, D, SPACE, J, NEWLINE,
			R, U, N, NEWLINE}

	}

	in := func() int64 {
		toReturn := droidProgram[currentChar]
		fmt.Print(string(toReturn))
		currentChar++
		return toReturn
	}
	out := func(output int64) {
		if output > 128 {
			fmt.Println(output)
		} else {
			fmt.Print(string(output))
		}

	}

	return in, out
}

func day21() {
	var springDroidProgramArr = []int64{109, 2050, 21101, 0, 966, 1, 21102, 1, 13, 0, 1105, 1, 1378, 21101, 0, 20, 0, 1106, 0, 1337, 21101, 0, 27, 0, 1106, 0, 1279, 1208, 1, 65, 748, 1005, 748, 73, 1208, 1, 79, 748, 1005, 748, 110, 1208, 1, 78, 748, 1005, 748, 132, 1208, 1, 87, 748, 1005, 748, 169, 1208, 1, 82, 748, 1005, 748, 239, 21102, 1041, 1, 1, 21101, 73, 0, 0, 1105, 1, 1421, 21101, 78, 0, 1, 21102, 1041, 1, 2, 21102, 1, 88, 0, 1105, 1, 1301, 21101, 0, 68, 1, 21101, 1041, 0, 2, 21102, 1, 103, 0, 1106, 0, 1301, 1101, 0, 1, 750, 1106, 0, 298, 21102, 1, 82, 1, 21101, 1041, 0, 2, 21101, 0, 125, 0, 1106, 0, 1301, 1101, 0, 2, 750, 1105, 1, 298, 21101, 0, 79, 1, 21101, 1041, 0, 2, 21102, 147, 1, 0, 1106, 0, 1301, 21101, 84, 0, 1, 21102, 1041, 1, 2, 21102, 162, 1, 0, 1105, 1, 1301, 1101, 3, 0, 750, 1106, 0, 298, 21101, 0, 65, 1, 21101, 1041, 0, 2, 21101, 184, 0, 0, 1106, 0, 1301, 21101, 76, 0, 1, 21101, 1041, 0, 2, 21102, 199, 1, 0, 1106, 0, 1301, 21102, 1, 75, 1, 21101, 0, 1041, 2, 21101, 0, 214, 0, 1106, 0, 1301, 21102, 221, 1, 0, 1106, 0, 1337, 21101, 10, 0, 1, 21101, 0, 1041, 2, 21101, 0, 236, 0, 1105, 1, 1301, 1106, 0, 553, 21101, 0, 85, 1, 21102, 1, 1041, 2, 21102, 1, 254, 0, 1105, 1, 1301, 21101, 0, 78, 1, 21101, 1041, 0, 2, 21102, 1, 269, 0, 1106, 0, 1301, 21101, 276, 0, 0, 1106, 0, 1337, 21101, 10, 0, 1, 21101, 1041, 0, 2, 21101, 291, 0, 0, 1106, 0, 1301, 1102, 1, 1, 755, 1106, 0, 553, 21102, 1, 32, 1, 21101, 0, 1041, 2, 21102, 313, 1, 0, 1106, 0, 1301, 21101, 320, 0, 0, 1106, 0, 1337, 21102, 1, 327, 0, 1106, 0, 1279, 2102, 1, 1, 749, 21101, 65, 0, 2, 21101, 0, 73, 3, 21102, 346, 1, 0, 1105, 1, 1889, 1206, 1, 367, 1007, 749, 69, 748, 1005, 748, 360, 1102, 1, 1, 756, 1001, 749, -64, 751, 1106, 0, 406, 1008, 749, 74, 748, 1006, 748, 381, 1101, 0, -1, 751, 1105, 1, 406, 1008, 749, 84, 748, 1006, 748, 395, 1102, 1, -2, 751, 1105, 1, 406, 21102, 1100, 1, 1, 21102, 406, 1, 0, 1105, 1, 1421, 21101, 32, 0, 1, 21101, 1100, 0, 2, 21101, 0, 421, 0, 1105, 1, 1301, 21101, 428, 0, 0, 1105, 1, 1337, 21101, 435, 0, 0, 1105, 1, 1279, 2101, 0, 1, 749, 1008, 749, 74, 748, 1006, 748, 453, 1101, 0, -1, 752, 1105, 1, 478, 1008, 749, 84, 748, 1006, 748, 467, 1101, 0, -2, 752, 1105, 1, 478, 21102, 1168, 1, 1, 21101, 478, 0, 0, 1105, 1, 1421, 21102, 1, 485, 0, 1106, 0, 1337, 21102, 1, 10, 1, 21101, 1168, 0, 2, 21102, 1, 500, 0, 1106, 0, 1301, 1007, 920, 15, 748, 1005, 748, 518, 21101, 0, 1209, 1, 21102, 1, 518, 0, 1105, 1, 1421, 1002, 920, 3, 529, 1001, 529, 921, 529, 1001, 750, 0, 0, 1001, 529, 1, 537, 1001, 751, 0, 0, 1001, 537, 1, 545, 102, 1, 752, 0, 1001, 920, 1, 920, 1105, 1, 13, 1005, 755, 577, 1006, 756, 570, 21101, 1100, 0, 1, 21101, 0, 570, 0, 1106, 0, 1421, 21101, 0, 987, 1, 1105, 1, 581, 21102, 1001, 1, 1, 21102, 1, 588, 0, 1106, 0, 1378, 1102, 1, 758, 594, 102, 1, 0, 753, 1006, 753, 654, 21002, 753, 1, 1, 21102, 1, 610, 0, 1106, 0, 667, 21101, 0, 0, 1, 21101, 0, 621, 0, 1105, 1, 1463, 1205, 1, 647, 21101, 0, 1015, 1, 21101, 635, 0, 0, 1105, 1, 1378, 21102, 1, 1, 1, 21101, 646, 0, 0, 1105, 1, 1463, 99, 1001, 594, 1, 594, 1105, 1, 592, 1006, 755, 664, 1102, 0, 1, 755, 1106, 0, 647, 4, 754, 99, 109, 2, 1102, 1, 726, 757, 22102, 1, -1, 1, 21102, 1, 9, 2, 21101, 697, 0, 3, 21101, 0, 692, 0, 1105, 1, 1913, 109, -2, 2105, 1, 0, 109, 2, 101, 0, 757, 706, 2101, 0, -1, 0, 1001, 757, 1, 757, 109, -2, 2106, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 63, 95, 159, 191, 223, 127, 0, 109, 152, 253, 182, 236, 231, 212, 42, 69, 154, 142, 196, 86, 214, 120, 166, 139, 175, 117, 102, 178, 124, 204, 189, 59, 156, 244, 92, 110, 184, 241, 87, 220, 245, 122, 167, 227, 53, 50, 197, 93, 85, 234, 190, 252, 248, 54, 163, 217, 94, 207, 58, 115, 70, 99, 71, 188, 247, 168, 221, 116, 239, 155, 186, 232, 49, 213, 218, 137, 103, 56, 170, 119, 34, 242, 76, 169, 173, 46, 172, 187, 171, 141, 153, 238, 61, 125, 121, 222, 199, 243, 174, 229, 235, 226, 201, 215, 108, 138, 126, 78, 249, 62, 51, 79, 57, 118, 181, 38, 39, 84, 228, 55, 77, 113, 179, 107, 136, 198, 140, 35, 246, 205, 162, 219, 43, 68, 185, 111, 237, 100, 183, 47, 157, 98, 230, 114, 158, 202, 216, 203, 200, 206, 106, 233, 254, 143, 60, 251, 101, 250, 177, 123, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 20, 73, 110, 112, 117, 116, 32, 105, 110, 115, 116, 114, 117, 99, 116, 105, 111, 110, 115, 58, 10, 13, 10, 87, 97, 108, 107, 105, 110, 103, 46, 46, 46, 10, 10, 13, 10, 82, 117, 110, 110, 105, 110, 103, 46, 46, 46, 10, 10, 25, 10, 68, 105, 100, 110, 39, 116, 32, 109, 97, 107, 101, 32, 105, 116, 32, 97, 99, 114, 111, 115, 115, 58, 10, 10, 58, 73, 110, 118, 97, 108, 105, 100, 32, 111, 112, 101, 114, 97, 116, 105, 111, 110, 59, 32, 101, 120, 112, 101, 99, 116, 101, 100, 32, 115, 111, 109, 101, 116, 104, 105, 110, 103, 32, 108, 105, 107, 101, 32, 65, 78, 68, 44, 32, 79, 82, 44, 32, 111, 114, 32, 78, 79, 84, 67, 73, 110, 118, 97, 108, 105, 100, 32, 102, 105, 114, 115, 116, 32, 97, 114, 103, 117, 109, 101, 110, 116, 59, 32, 101, 120, 112, 101, 99, 116, 101, 100, 32, 115, 111, 109, 101, 116, 104, 105, 110, 103, 32, 108, 105, 107, 101, 32, 65, 44, 32, 66, 44, 32, 67, 44, 32, 68, 44, 32, 74, 44, 32, 111, 114, 32, 84, 40, 73, 110, 118, 97, 108, 105, 100, 32, 115, 101, 99, 111, 110, 100, 32, 97, 114, 103, 117, 109, 101, 110, 116, 59, 32, 101, 120, 112, 101, 99, 116, 101, 100, 32, 74, 32, 111, 114, 32, 84, 52, 79, 117, 116, 32, 111, 102, 32, 109, 101, 109, 111, 114, 121, 59, 32, 97, 116, 32, 109, 111, 115, 116, 32, 49, 53, 32, 105, 110, 115, 116, 114, 117, 99, 116, 105, 111, 110, 115, 32, 99, 97, 110, 32, 98, 101, 32, 115, 116, 111, 114, 101, 100, 0, 109, 1, 1005, 1262, 1270, 3, 1262, 20102, 1, 1262, 0, 109, -1, 2106, 0, 0, 109, 1, 21102, 1288, 1, 0, 1105, 1, 1263, 20101, 0, 1262, 0, 1101, 0, 0, 1262, 109, -1, 2106, 0, 0, 109, 5, 21102, 1, 1310, 0, 1106, 0, 1279, 22102, 1, 1, -2, 22208, -2, -4, -1, 1205, -1, 1332, 21201, -3, 0, 1, 21101, 0, 1332, 0, 1106, 0, 1421, 109, -5, 2106, 0, 0, 109, 2, 21101, 0, 1346, 0, 1105, 1, 1263, 21208, 1, 32, -1, 1205, -1, 1363, 21208, 1, 9, -1, 1205, -1, 1363, 1106, 0, 1373, 21101, 1370, 0, 0, 1106, 0, 1279, 1105, 1, 1339, 109, -2, 2106, 0, 0, 109, 5, 1201, -4, 0, 1385, 21001, 0, 0, -2, 22101, 1, -4, -4, 21102, 0, 1, -3, 22208, -3, -2, -1, 1205, -1, 1416, 2201, -4, -3, 1408, 4, 0, 21201, -3, 1, -3, 1106, 0, 1396, 109, -5, 2106, 0, 0, 109, 2, 104, 10, 22102, 1, -1, 1, 21102, 1, 1436, 0, 1106, 0, 1378, 104, 10, 99, 109, -2, 2105, 1, 0, 109, 3, 20002, 594, 753, -1, 22202, -1, -2, -1, 201, -1, 754, 754, 109, -3, 2106, 0, 0, 109, 10, 21101, 5, 0, -5, 21102, 1, 1, -4, 21102, 0, 1, -3, 1206, -9, 1555, 21101, 3, 0, -6, 21102, 1, 5, -7, 22208, -7, -5, -8, 1206, -8, 1507, 22208, -6, -4, -8, 1206, -8, 1507, 104, 64, 1105, 1, 1529, 1205, -6, 1527, 1201, -7, 716, 1515, 21002, 0, -11, -8, 21201, -8, 46, -8, 204, -8, 1105, 1, 1529, 104, 46, 21201, -7, 1, -7, 21207, -7, 22, -8, 1205, -8, 1488, 104, 10, 21201, -6, -1, -6, 21207, -6, 0, -8, 1206, -8, 1484, 104, 10, 21207, -4, 1, -8, 1206, -8, 1569, 21102, 0, 1, -9, 1106, 0, 1689, 21208, -5, 21, -8, 1206, -8, 1583, 21102, 1, 1, -9, 1105, 1, 1689, 1201, -5, 716, 1588, 21001, 0, 0, -2, 21208, -4, 1, -1, 22202, -2, -1, -1, 1205, -2, 1613, 22101, 0, -5, 1, 21101, 1613, 0, 0, 1105, 1, 1444, 1206, -1, 1634, 22101, 0, -5, 1, 21101, 1627, 0, 0, 1105, 1, 1694, 1206, 1, 1634, 21102, 1, 2, -3, 22107, 1, -4, -8, 22201, -1, -8, -8, 1206, -8, 1649, 21201, -5, 1, -5, 1206, -3, 1663, 21201, -3, -1, -3, 21201, -4, 1, -4, 1106, 0, 1667, 21201, -4, -1, -4, 21208, -4, 0, -1, 1201, -5, 716, 1676, 22002, 0, -1, -1, 1206, -1, 1686, 21101, 1, 0, -4, 1106, 0, 1477, 109, -10, 2106, 0, 0, 109, 11, 21101, 0, 0, -6, 21102, 1, 0, -8, 21101, 0, 0, -7, 20208, -6, 920, -9, 1205, -9, 1880, 21202, -6, 3, -9, 1201, -9, 921, 1724, 21002, 0, 1, -5, 1001, 1724, 1, 1732, 21001, 0, 0, -4, 21201, -4, 0, 1, 21102, 1, 1, 2, 21102, 9, 1, 3, 21102, 1754, 1, 0, 1105, 1, 1889, 1206, 1, 1772, 2201, -10, -4, 1766, 1001, 1766, 716, 1766, 21001, 0, 0, -3, 1106, 0, 1790, 21208, -4, -1, -9, 1206, -9, 1786, 22101, 0, -8, -3, 1106, 0, 1790, 22102, 1, -7, -3, 1001, 1732, 1, 1796, 20102, 1, 0, -2, 21208, -2, -1, -9, 1206, -9, 1812, 22102, 1, -8, -1, 1105, 1, 1816, 21202, -7, 1, -1, 21208, -5, 1, -9, 1205, -9, 1837, 21208, -5, 2, -9, 1205, -9, 1844, 21208, -3, 0, -1, 1105, 1, 1855, 22202, -3, -1, -1, 1106, 0, 1855, 22201, -3, -1, -1, 22107, 0, -1, -1, 1106, 0, 1855, 21208, -2, -1, -9, 1206, -9, 1869, 21202, -1, 1, -8, 1105, 1, 1873, 21202, -1, 1, -7, 21201, -6, 1, -6, 1106, 0, 1708, 21201, -8, 0, -10, 109, -11, 2105, 1, 0, 109, 7, 22207, -6, -5, -3, 22207, -4, -6, -2, 22201, -3, -2, -1, 21208, -1, 0, -6, 109, -7, 2106, 0, 0, 0, 109, 5, 1202, -2, 1, 1912, 21207, -4, 0, -1, 1206, -1, 1930, 21101, 0, 0, -4, 22101, 0, -4, 1, 21202, -3, 1, 2, 21102, 1, 1, 3, 21102, 1949, 1, 0, 1106, 0, 1954, 109, -5, 2105, 1, 0, 109, 6, 21207, -4, 1, -1, 1206, -1, 1977, 22207, -5, -3, -1, 1206, -1, 1977, 21201, -5, 0, -5, 1105, 1, 2045, 21202, -5, 1, 1, 21201, -4, -1, 2, 21202, -3, 2, 3, 21102, 1, 1996, 0, 1105, 1, 1954, 21201, 1, 0, -5, 21102, 1, 1, -2, 22207, -5, -3, -1, 1206, -1, 2015, 21101, 0, 0, -2, 22202, -3, -2, -3, 22107, 0, -4, -1, 1206, -1, 2037, 22102, 1, -2, 1, 21101, 2037, 0, 0, 106, 0, 1912, 21202, -3, -1, -3, 22201, -5, -3, -5, 109, -6, 2106, 0, 0}
	var springDroidProgram = makeMapForArray(springDroidProgramArr)

	// Part 1 with walking
	in, out := springDroidIOHandlers(false)
	intComp(springDroidProgram, in, out)

	// Part 2 with running
	in, out = springDroidIOHandlers(true)
	intComp(springDroidProgram, in, out)
}

type cardShuffleInstruction struct {
	instructionType string
	value           int64
}

func cutDeck(cardsArr []int64, numToCut int64) []int64 {
	var newCardsArr []int64
	if numToCut > 0 {
		newCardsArr = append(newCardsArr, cardsArr[numToCut:]...)
		newCardsArr = append(newCardsArr, cardsArr[:numToCut]...)
	} else {
		numToCut = Abs(numToCut)
		newCardsArr = append(newCardsArr, cardsArr[int64(len(cardsArr))-numToCut:]...)
		newCardsArr = append(newCardsArr, cardsArr[:int64(len(cardsArr))-numToCut]...)
	}
	return newCardsArr
}

func dealWithIncrement(cardsArr []int64, increment int64) []int64 {
	var newCardsArr = make([]int64, len(cardsArr))
	var cardsDealt = 0
	var dealFromIndex = 0
	var dealToIndex = 0

	for cardsDealt != len(cardsArr) {
		if dealToIndex > len(cardsArr)-1 {
			dealToIndex = dealToIndex % len(cardsArr)
		}

		newCardsArr[dealToIndex] = cardsArr[dealFromIndex]

		dealToIndex += int(increment)
		dealFromIndex++
		cardsDealt++
	}

	return newCardsArr
}

func dealIntoNewStack(cardsArr []int64) []int64 {
	var newCardsArr []int64
	for cardIndex := range cardsArr {
		cardFromBack := cardsArr[(len(cardsArr)-1)-cardIndex]
		newCardsArr = append(newCardsArr, cardFromBack)
	}

	return newCardsArr
}

func runInstructionsForCards(cardsArr []int64, cardShuffleInstructions []cardShuffleInstruction) []int64 {
	for _, instruction := range cardShuffleInstructions {
		switch instruction.instructionType {
		case "cut":
			cardsArr = cutDeck(cardsArr, instruction.value)
		case "deal with increment":
			cardsArr = dealWithIncrement(cardsArr, instruction.value)
		case "deal into new stack":
			cardsArr = dealIntoNewStack(cardsArr)
		}
	}
	return cardsArr
}

func fillCards(cardsArr []int64) []int64 {
	for i := range cardsArr {
		cardsArr[i] = int64(i)
	}
	return cardsArr
}

func invBigInt(x, c *big.Int) *big.Int {
	exponent := big.NewInt(0).Sub(c, big.NewInt(2))
	inverseValue := big.NewInt(0).Exp(x, exponent, c)
	return inverseValue
}

func runInstructionsForLargeAmountsOfCards(incrementRate, offset, numberOfCards *big.Int,
	cardShuffleInstructions []cardShuffleInstruction) (*big.Int, *big.Int) {
	for _, instruction := range cardShuffleInstructions {
		switch instruction.instructionType {
		case "cut":
			offset.Add(offset, big.NewInt(0).Mul(incrementRate, big.NewInt(instruction.value)))
		case "deal with increment":
			incrementRate.Mul(incrementRate, invBigInt(big.NewInt(instruction.value), numberOfCards))
		case "deal into new stack":
			offset.Sub(offset, incrementRate)
			incrementRate.Neg(incrementRate)
		}
	}
	return incrementRate, offset
}

func day22() {
	file, err := os.Open("day22CardShuffleInstructions.txt")
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	var cardShuffleInstructions []cardShuffleInstruction
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		instructinoString := scanner.Text()
		instructionArr := strings.Split(instructinoString, " ")
		lastElementIndex := len(instructionArr) - 1
		var newCardShuffleInstruction cardShuffleInstruction

		if instructionArr[lastElementIndex] == "stack" {
			newCardShuffleInstruction.instructionType = instructinoString
			newCardShuffleInstruction.value = -1
		} else {
			newCardShuffleInstruction.instructionType = strings.Join(instructionArr[:lastElementIndex], " ")
			value, _ := strconv.ParseInt(instructionArr[lastElementIndex], 10, 64)
			newCardShuffleInstruction.value = value
		}
		cardShuffleInstructions = append(cardShuffleInstructions, newCardShuffleInstruction)
	}

	var part1Deck = make([]int64, 10007)
	part1Deck = fillCards(part1Deck)
	var card2019Pos = 0

	part1Deck = runInstructionsForCards(part1Deck, cardShuffleInstructions)
	for index, card := range part1Deck {
		if card == 2019 {
			card2019Pos = index
		}
	}

	fmt.Println("Positon of card 2019:", card2019Pos)

	var part2DeckSize = big.NewInt(119315717514047)
	var repeatInstructionAmount = big.NewInt(101741582076661)
	var incrementRate = big.NewInt(1)
	var offset = big.NewInt(0)
	var cardAtPos = big.NewInt(2020)

	incrementRate, offset = runInstructionsForLargeAmountsOfCards(incrementRate, offset, part2DeckSize,
		cardShuffleInstructions)

	negativeIncrementRate := big.NewInt(0).Sub(big.NewInt(1), incrementRate)
	inveseNegativeIncrementRate := invBigInt(negativeIncrementRate, part2DeckSize)
	offset.Mul(offset, inveseNegativeIncrementRate)
	incrementRate.Exp(incrementRate, repeatInstructionAmount, part2DeckSize)
	cardAtPos.Mul(cardAtPos, incrementRate)
	negativeIncrementRate = big.NewInt(0).Sub(big.NewInt(1), incrementRate)
	offset.Mul(offset, negativeIncrementRate)
	cardAtPos.Add(cardAtPos, offset)
	cardAtPos.Mod(cardAtPos, part2DeckSize)
	var cardAt2020 = cardAtPos.String()

	fmt.Println("Card at position 2020 with extended instructions:", cardAt2020)
}

type Int64Value struct {
	Value int64
}

// NewInt64Queue returns a new queue with the given initial size.
func NewInt64Queue(size int) *Int64Queue {
	return &Int64Queue{
		values: make([]*Int64Value, size),
		size:   size,
	}
}

// Queue is a basic FIFO queue based on a circular list that resizes as needed.
type Int64Queue struct {
	values []*Int64Value
	size   int
	head   int
	tail   int
	count  int
}

// Push adds a node to the queue.
func (q *Int64Queue) Push(n *Int64Value) {
	if q.head == q.tail && q.count > 0 {
		nodes := make([]*Int64Value, len(q.values)+q.size)
		copy(nodes, q.values[q.head:])
		copy(nodes[len(q.values)-q.head:], q.values[:q.head])
		q.head = 0
		q.tail = len(q.values)
		q.values = nodes
	}
	q.values[q.tail] = n
	q.tail = (q.tail + 1) % len(q.values)
	q.count++
}

// Pop removes and returns a node from the queue in first to last order.
func (q *Int64Queue) Pop() *Int64Value {
	if q.count == 0 {
		return nil
	}
	node := q.values[q.head]
	q.head = (q.head + 1) % len(q.values)
	q.count--
	return node
}
func networkIOHandlers(networkID int64, network []Int64Queue, localReceive chan int64) (inputHandler, outputHandler) {
	var returnedNetworkID bool = false
	var sendPacketVal int64 = 0
	var sendToNetworkID int64 = -1
	var sentFirstYPacket bool = false
	var lastXReceived int64 = -1
	var lastYReceived int64 = -1

	NATsystem := func(value int64) {
		// fmt.Println("----------------------------------Received at NAT system")
		var allIdle = true
		for _, networkQueue := range network {
			if networkQueue.count > 0 {
				allIdle = false
				// fmt.Println("NOT IDLE--------------------", checkingNetworkID)
			}
		}

		if allIdle {
			// fmt.Println("ALL IDLE", networkID)
			switch sendPacketVal {
			case 1:
				lastXReceived = value
			case 2:
				// Packets only come from computer 10.
				// If packets came from multiple computers,
				// you would need sentFirstYPacket to not be a local var.
				if value == lastYReceived || !sentFirstYPacket {
					localReceive <- value
					sentFirstYPacket = true
				}
				lastYReceived = value
				network[0].Push(&Int64Value{lastXReceived})
				network[0].Push(&Int64Value{lastYReceived})
			}
		}
	}

	in := func() int64 {
		// fmt.Println("Running on comp:", networkID)
		if !returnedNetworkID {
			returnedNetworkID = true
			// fmt.Printf("returning networkID %d\n", networkID)
			return networkID
		}
		packet := network[networkID].Pop()
		if packet != nil {
			// fmt.Println("Received packet at networkID:", networkID)
			return packet.Value
		}
		return -1
	}

	out := func(output int64) {
		switch sendPacketVal {
		case 0:
			// fmt.Println("Set send to network ID")
			sendToNetworkID = output
			sendPacketVal++
		case 1:
			if sendToNetworkID == 255 {
				NATsystem(output)
			} else {
				network[sendToNetworkID].Push(&Int64Value{output})
				// fmt.Println("Sent X at networkID:", networkID, "to networkID:", sendToNetworkID)
			}
			sendPacketVal++
		case 2:
			if sendToNetworkID == 255 {
				NATsystem(output)
			} else {
				network[sendToNetworkID].Push(&Int64Value{output})
				// fmt.Println("Sent Y at networkID:", networkID, "to networkID:", sendToNetworkID)
			}
			sendPacketVal = 0
		}
	}

	return in, out
}

func day23() {
	var networkProgramArr = []int64{3, 62, 1001, 62, 11, 10, 109, 2243, 105, 1, 0, 1555, 2097, 668, 1728, 833, 1425, 2029, 2060, 1798, 1631, 1136, 864, 1988, 1330, 938, 1957, 1169, 969, 2140, 2208, 571, 899, 1660, 1293, 1454, 1485, 1858, 633, 1390, 1691, 802, 1829, 1031, 1596, 998, 1262, 730, 1761, 1520, 2171, 1231, 699, 1062, 765, 604, 1893, 1105, 1200, 1361, 1922, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 3, 64, 1008, 64, -1, 62, 1006, 62, 88, 1006, 61, 170, 1106, 0, 73, 3, 65, 20101, 0, 64, 1, 20102, 1, 66, 2, 21102, 1, 105, 0, 1106, 0, 436, 1201, 1, -1, 64, 1007, 64, 0, 62, 1005, 62, 73, 7, 64, 67, 62, 1006, 62, 73, 1002, 64, 2, 133, 1, 133, 68, 133, 102, 1, 0, 62, 1001, 133, 1, 140, 8, 0, 65, 63, 2, 63, 62, 62, 1005, 62, 73, 1002, 64, 2, 161, 1, 161, 68, 161, 1102, 1, 1, 0, 1001, 161, 1, 169, 1001, 65, 0, 0, 1102, 1, 1, 61, 1102, 1, 0, 63, 7, 63, 67, 62, 1006, 62, 203, 1002, 63, 2, 194, 1, 68, 194, 194, 1006, 0, 73, 1001, 63, 1, 63, 1105, 1, 178, 21102, 1, 210, 0, 106, 0, 69, 1201, 1, 0, 70, 1102, 1, 0, 63, 7, 63, 71, 62, 1006, 62, 250, 1002, 63, 2, 234, 1, 72, 234, 234, 4, 0, 101, 1, 234, 240, 4, 0, 4, 70, 1001, 63, 1, 63, 1105, 1, 218, 1105, 1, 73, 109, 4, 21102, 0, 1, -3, 21101, 0, 0, -2, 20207, -2, 67, -1, 1206, -1, 293, 1202, -2, 2, 283, 101, 1, 283, 283, 1, 68, 283, 283, 22001, 0, -3, -3, 21201, -2, 1, -2, 1106, 0, 263, 22102, 1, -3, -3, 109, -4, 2106, 0, 0, 109, 4, 21102, 1, 1, -3, 21101, 0, 0, -2, 20207, -2, 67, -1, 1206, -1, 342, 1202, -2, 2, 332, 101, 1, 332, 332, 1, 68, 332, 332, 22002, 0, -3, -3, 21201, -2, 1, -2, 1105, 1, 312, 21201, -3, 0, -3, 109, -4, 2105, 1, 0, 109, 1, 101, 1, 68, 358, 21002, 0, 1, 1, 101, 3, 68, 367, 20101, 0, 0, 2, 21102, 1, 376, 0, 1105, 1, 436, 22102, 1, 1, 0, 109, -1, 2105, 1, 0, 1, 2, 4, 8, 16, 32, 64, 128, 256, 512, 1024, 2048, 4096, 8192, 16384, 32768, 65536, 131072, 262144, 524288, 1048576, 2097152, 4194304, 8388608, 16777216, 33554432, 67108864, 134217728, 268435456, 536870912, 1073741824, 2147483648, 4294967296, 8589934592, 17179869184, 34359738368, 68719476736, 137438953472, 274877906944, 549755813888, 1099511627776, 2199023255552, 4398046511104, 8796093022208, 17592186044416, 35184372088832, 70368744177664, 140737488355328, 281474976710656, 562949953421312, 1125899906842624, 109, 8, 21202, -6, 10, -5, 22207, -7, -5, -5, 1205, -5, 521, 21102, 1, 0, -4, 21102, 0, 1, -3, 21101, 0, 51, -2, 21201, -2, -1, -2, 1201, -2, 385, 470, 21002, 0, 1, -1, 21202, -3, 2, -3, 22207, -7, -1, -5, 1205, -5, 496, 21201, -3, 1, -3, 22102, -1, -1, -5, 22201, -7, -5, -7, 22207, -3, -6, -5, 1205, -5, 515, 22102, -1, -6, -5, 22201, -3, -5, -3, 22201, -1, -4, -4, 1205, -2, 461, 1105, 1, 547, 21101, 0, -1, -4, 21202, -6, -1, -6, 21207, -7, 0, -5, 1205, -5, 547, 22201, -7, -6, -7, 21201, -4, 1, -4, 1106, 0, 529, 22101, 0, -4, -7, 109, -8, 2105, 1, 0, 109, 1, 101, 1, 68, 563, 21001, 0, 0, 0, 109, -1, 2106, 0, 0, 1101, 85199, 0, 66, 1102, 1, 2, 67, 1102, 1, 598, 68, 1102, 302, 1, 69, 1101, 0, 1, 71, 1101, 602, 0, 72, 1105, 1, 73, 0, 0, 0, 0, 49, 96146, 1102, 1, 79357, 66, 1102, 1, 1, 67, 1101, 631, 0, 68, 1102, 556, 1, 69, 1102, 1, 0, 71, 1101, 0, 633, 72, 1106, 0, 73, 1, 1115, 1102, 102593, 1, 66, 1102, 1, 3, 67, 1101, 0, 660, 68, 1101, 302, 0, 69, 1102, 1, 1, 71, 1102, 1, 666, 72, 1105, 1, 73, 0, 0, 0, 0, 0, 0, 21, 69404, 1102, 1, 54751, 66, 1102, 1, 1, 67, 1102, 1, 695, 68, 1102, 1, 556, 69, 1102, 1, 1, 71, 1101, 697, 0, 72, 1106, 0, 73, 1, -176, 23, 44753, 1102, 84163, 1, 66, 1101, 0, 1, 67, 1101, 726, 0, 68, 1101, 556, 0, 69, 1102, 1, 1, 71, 1102, 1, 728, 72, 1105, 1, 73, 1, 12, 39, 237092, 1101, 0, 10159, 66, 1102, 1, 1, 67, 1101, 0, 757, 68, 1101, 556, 0, 69, 1102, 1, 3, 71, 1102, 759, 1, 72, 1105, 1, 73, 1, 10, 33, 2293, 37, 135068, 12, 125182, 1101, 42197, 0, 66, 1102, 1, 4, 67, 1102, 1, 792, 68, 1102, 253, 1, 69, 1102, 1, 1, 71, 1102, 1, 800, 72, 1105, 1, 73, 0, 0, 0, 0, 0, 0, 0, 0, 10, 91297, 1101, 0, 36373, 66, 1101, 1, 0, 67, 1102, 829, 1, 68, 1102, 1, 556, 69, 1102, 1, 1, 71, 1101, 831, 0, 72, 1105, 1, 73, 1, 1613, 33, 4586, 1101, 0, 78787, 66, 1101, 1, 0, 67, 1102, 860, 1, 68, 1102, 1, 556, 69, 1101, 1, 0, 71, 1102, 862, 1, 72, 1105, 1, 73, 1, 13, 42, 572971, 1102, 1, 63079, 66, 1101, 0, 1, 67, 1102, 891, 1, 68, 1102, 1, 556, 69, 1101, 3, 0, 71, 1101, 893, 0, 72, 1106, 0, 73, 1, 3, 42, 327412, 7, 79873, 23, 134259, 1102, 1, 17351, 66, 1102, 1, 5, 67, 1102, 926, 1, 68, 1101, 253, 0, 69, 1102, 1, 1, 71, 1102, 1, 936, 72, 1106, 0, 73, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 3, 47797, 1101, 1399, 0, 66, 1102, 1, 1, 67, 1102, 1, 965, 68, 1101, 0, 556, 69, 1102, 1, 1, 71, 1102, 1, 967, 72, 1105, 1, 73, 1, 8, 39, 118546, 1101, 0, 28751, 66, 1101, 1, 0, 67, 1102, 1, 996, 68, 1101, 0, 556, 69, 1102, 1, 0, 71, 1101, 0, 998, 72, 1105, 1, 73, 1, 1683, 1101, 0, 79279, 66, 1101, 0, 2, 67, 1102, 1025, 1, 68, 1102, 302, 1, 69, 1101, 0, 1, 71, 1102, 1029, 1, 72, 1105, 1, 73, 0, 0, 0, 0, 43, 84394, 1101, 37087, 0, 66, 1102, 1, 1, 67, 1102, 1, 1058, 68, 1102, 1, 556, 69, 1102, 1, 1, 71, 1101, 0, 1060, 72, 1105, 1, 73, 1, 107, 7, 239619, 1101, 0, 81853, 66, 1102, 1, 7, 67, 1101, 0, 1089, 68, 1101, 0, 302, 69, 1102, 1, 1, 71, 1102, 1103, 1, 72, 1105, 1, 73, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 43, 168788, 1102, 103079, 1, 66, 1102, 1, 1, 67, 1101, 0, 1132, 68, 1101, 0, 556, 69, 1101, 0, 1, 71, 1101, 1134, 0, 72, 1105, 1, 73, 1, 9, 23, 89506, 1102, 1, 91297, 66, 1101, 2, 0, 67, 1101, 0, 1163, 68, 1101, 0, 351, 69, 1101, 1, 0, 71, 1101, 0, 1167, 72, 1105, 1, 73, 0, 0, 0, 0, 255, 29569, 1102, 1, 26687, 66, 1101, 1, 0, 67, 1102, 1196, 1, 68, 1102, 1, 556, 69, 1102, 1, 1, 71, 1101, 1198, 0, 72, 1105, 1, 73, 1, 8867, 42, 245559, 1101, 81869, 0, 66, 1102, 1, 1, 67, 1101, 0, 1227, 68, 1101, 556, 0, 69, 1101, 1, 0, 71, 1101, 1229, 0, 72, 1105, 1, 73, 1, 233, 7, 159746, 1101, 59753, 0, 66, 1102, 1, 1, 67, 1102, 1, 1258, 68, 1101, 0, 556, 69, 1101, 1, 0, 71, 1101, 1260, 0, 72, 1106, 0, 73, 1, 11, 42, 163706, 1102, 1, 78467, 66, 1101, 1, 0, 67, 1102, 1, 1289, 68, 1102, 1, 556, 69, 1101, 1, 0, 71, 1102, 1291, 1, 72, 1105, 1, 73, 1, 131, 27, 102593, 1102, 1, 44753, 66, 1101, 4, 0, 67, 1101, 1320, 0, 68, 1101, 0, 302, 69, 1101, 1, 0, 71, 1102, 1, 1328, 72, 1106, 0, 73, 0, 0, 0, 0, 0, 0, 0, 0, 19, 1373, 1102, 28051, 1, 66, 1101, 1, 0, 67, 1101, 0, 1357, 68, 1101, 0, 556, 69, 1101, 1, 0, 71, 1101, 0, 1359, 72, 1106, 0, 73, 1, 160, 12, 187773, 1102, 1, 50359, 66, 1102, 1, 1, 67, 1102, 1, 1388, 68, 1102, 556, 1, 69, 1101, 0, 0, 71, 1101, 0, 1390, 72, 1105, 1, 73, 1, 1232, 1101, 93251, 0, 66, 1102, 1, 1, 67, 1102, 1, 1417, 68, 1101, 0, 556, 69, 1102, 3, 1, 71, 1102, 1, 1419, 72, 1106, 0, 73, 1, 7, 42, 81853, 38, 278583, 23, 179012, 1102, 1, 53759, 66, 1102, 1, 1, 67, 1101, 1452, 0, 68, 1101, 556, 0, 69, 1102, 0, 1, 71, 1101, 0, 1454, 72, 1106, 0, 73, 1, 1370, 1101, 19597, 0, 66, 1101, 1, 0, 67, 1101, 1481, 0, 68, 1102, 556, 1, 69, 1102, 1, 1, 71, 1102, 1, 1483, 72, 1105, 1, 73, 1, 2311, 27, 205186, 1102, 18253, 1, 66, 1101, 3, 0, 67, 1101, 1512, 0, 68, 1102, 1, 302, 69, 1101, 0, 1, 71, 1102, 1518, 1, 72, 1105, 1, 73, 0, 0, 0, 0, 0, 0, 43, 126591, 1101, 0, 92861, 66, 1102, 1, 3, 67, 1101, 1547, 0, 68, 1102, 1, 302, 69, 1102, 1, 1, 71, 1102, 1, 1553, 72, 1105, 1, 73, 0, 0, 0, 0, 0, 0, 34, 158558, 1102, 29569, 1, 66, 1101, 1, 0, 67, 1101, 1582, 0, 68, 1101, 556, 0, 69, 1102, 6, 1, 71, 1101, 1584, 0, 72, 1106, 0, 73, 1, 20982, 34, 79279, 19, 2746, 19, 4119, 25, 18253, 25, 36506, 25, 54759, 1102, 2293, 1, 66, 1102, 1, 3, 67, 1101, 0, 1623, 68, 1101, 302, 0, 69, 1101, 0, 1, 71, 1101, 0, 1629, 72, 1105, 1, 73, 0, 0, 0, 0, 0, 0, 21, 52053, 1101, 88259, 0, 66, 1101, 0, 1, 67, 1102, 1, 1658, 68, 1101, 556, 0, 69, 1101, 0, 0, 71, 1101, 1660, 0, 72, 1106, 0, 73, 1, 1672, 1101, 0, 12379, 66, 1101, 0, 1, 67, 1102, 1, 1687, 68, 1101, 556, 0, 69, 1102, 1, 1, 71, 1102, 1689, 1, 72, 1105, 1, 73, 1, -3333, 21, 34702, 1102, 39569, 1, 66, 1101, 1, 0, 67, 1102, 1718, 1, 68, 1101, 556, 0, 69, 1102, 4, 1, 71, 1101, 1720, 0, 72, 1106, 0, 73, 1, 1, 7, 319492, 33, 6879, 20, 170398, 27, 307779, 1102, 1, 47797, 66, 1102, 2, 1, 67, 1101, 1755, 0, 68, 1101, 0, 302, 69, 1102, 1, 1, 71, 1102, 1, 1759, 72, 1105, 1, 73, 0, 0, 0, 0, 38, 185722, 1102, 33767, 1, 66, 1102, 1, 4, 67, 1101, 0, 1788, 68, 1101, 0, 302, 69, 1101, 0, 1, 71, 1101, 0, 1796, 72, 1105, 1, 73, 0, 0, 0, 0, 0, 0, 0, 0, 12, 312955, 1102, 1, 55381, 66, 1102, 1, 1, 67, 1101, 1825, 0, 68, 1102, 556, 1, 69, 1102, 1, 1, 71, 1102, 1827, 1, 72, 1105, 1, 73, 1, 192, 39, 59273, 1101, 80953, 0, 66, 1101, 0, 1, 67, 1101, 0, 1856, 68, 1101, 556, 0, 69, 1101, 0, 0, 71, 1102, 1, 1858, 72, 1106, 0, 73, 1, 1204, 1102, 1, 96451, 66, 1101, 0, 1, 67, 1101, 0, 1885, 68, 1101, 556, 0, 69, 1101, 0, 3, 71, 1102, 1, 1887, 72, 1105, 1, 73, 1, 5, 37, 67534, 37, 101301, 12, 62591, 1101, 0, 54361, 66, 1102, 1, 1, 67, 1101, 0, 1920, 68, 1101, 0, 556, 69, 1102, 1, 0, 71, 1101, 0, 1922, 72, 1106, 0, 73, 1, 1692, 1101, 48073, 0, 66, 1101, 0, 3, 67, 1102, 1949, 1, 68, 1102, 302, 1, 69, 1102, 1, 1, 71, 1102, 1955, 1, 72, 1105, 1, 73, 0, 0, 0, 0, 0, 0, 21, 86755, 1102, 23671, 1, 66, 1102, 1, 1, 67, 1102, 1984, 1, 68, 1102, 556, 1, 69, 1101, 1, 0, 71, 1101, 0, 1986, 72, 1105, 1, 73, 1, 125, 37, 33767, 1102, 1, 62591, 66, 1101, 0, 6, 67, 1101, 0, 2015, 68, 1101, 0, 302, 69, 1102, 1, 1, 71, 1101, 2027, 0, 72, 1106, 0, 73, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 10, 182594, 1102, 57493, 1, 66, 1102, 1, 1, 67, 1101, 2056, 0, 68, 1101, 556, 0, 69, 1101, 0, 1, 71, 1101, 2058, 0, 72, 1106, 0, 73, 1, -23027, 20, 85199, 1102, 1, 79873, 66, 1102, 1, 4, 67, 1102, 1, 2087, 68, 1101, 0, 302, 69, 1102, 1, 1, 71, 1102, 2095, 1, 72, 1105, 1, 73, 0, 0, 0, 0, 0, 0, 0, 0, 21, 17351, 1102, 63487, 1, 66, 1102, 1, 1, 67, 1102, 1, 2124, 68, 1101, 0, 556, 69, 1101, 0, 7, 71, 1102, 1, 2126, 72, 1105, 1, 73, 1, 2, 39, 177819, 42, 409265, 49, 48073, 49, 144219, 38, 92861, 12, 250364, 12, 375546, 1101, 36691, 0, 66, 1102, 1, 1, 67, 1101, 2167, 0, 68, 1102, 1, 556, 69, 1102, 1, 1, 71, 1102, 1, 2169, 72, 1106, 0, 73, 1, 256, 3, 95594, 1101, 59273, 0, 66, 1102, 4, 1, 67, 1102, 1, 2198, 68, 1101, 302, 0, 69, 1102, 1, 1, 71, 1101, 2206, 0, 72, 1105, 1, 73, 0, 0, 0, 0, 0, 0, 0, 0, 42, 491118, 1101, 1373, 0, 66, 1101, 0, 3, 67, 1101, 0, 2235, 68, 1102, 302, 1, 69, 1101, 0, 1, 71, 1101, 2241, 0, 72, 1106, 0, 73, 0, 0, 0, 0, 0, 0, 43, 42197}
	var networkProgram = makeMapForArray(networkProgramArr)
	var firstYReceived int64 = 0
	var firstYRepeatedTwiceInARow int64 = 0

	var numOfComputers int64 = 50
	var network = make([]Int64Queue, numOfComputers)
	var comp int64
	for comp = 0; comp < numOfComputers; comp++ {
		network[comp] = *NewInt64Queue(int(numOfComputers))
	}

	// This is where we will get the output we want.
	var localReceive = make(chan int64, numOfComputers)

	// Some race conditions, but they are rare-ish.
	// firstYReceived is always correct though.
	for comp = 0; comp < numOfComputers; comp++ {
		in, out := networkIOHandlers(comp, network, localReceive)
		go intComp(networkProgram, in, out)
	}

	firstYReceived = <-localReceive
	firstYRepeatedTwiceInARow = <-localReceive
	fmt.Println("First Y value received at networkID 255:", firstYReceived)
	fmt.Println("First Y value sent to networkID 0 from NAT system twice in a row:", firstYRepeatedTwiceInARow)
}

type BugTile struct {
	biodiversityRating int64
	adjacentBugCount   int64
	tile               string
	hasBug             bool
}

func day24() {
	file, err := os.Open("day24BugData.txt")
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	var bugMap = make(map[Point]BugTile)
	var levelMap = make(map[int64]map[Point]BugTile)
	levelMap[0] = make(map[Point]BugTile)
	scanner := bufio.NewScanner(file)

	var x, y, width, height, currentBiodiversityRatingCount int64
	for scanner.Scan() {
		line := scanner.Text()
		for _, char := range line {
			if x > width {
				width = x
			}
			if y > height {
				height = y
			}
			var newBugTile BugTile
			newBugTile.biodiversityRating = int64(math.Pow(2, float64(currentBiodiversityRatingCount)))
			newBugTile.tile = string(char)
			if char == '#' {
				newBugTile.hasBug = true
			}
			bugMap[Point{x, y}] = newBugTile
			levelMap[0][Point{x, y}] = newBugTile
			currentBiodiversityRatingCount++
			x++
		}
		x = 0
		y++
	}

	var availableDirs = []int64{NORTH, SOUTH, WEST, EAST}
	var previousLayoutsMap = make(map[string]bool)

	var biodiversityRatingOfFirstDoubleLayout int64
	var bugsPressentAfter200Minutes int64

	var foundDoubleBioRating = false
	for !foundDoubleBioRating {
		var layoutString = ""
		var totalBioDiversityRating int64 = 0
		for y = 0; y <= height; y++ {
			for x = 0; x <= width; x++ {
				bugTile := bugMap[Point{x, y}]
				layoutString += bugTile.tile

				if bugTile.hasBug {
					totalBioDiversityRating += bugTile.biodiversityRating
				}

				bugTile.adjacentBugCount = 0
				for _, dir := range availableDirs {
					checkBugPoint := getPointForDirection(dir, Point{x, y})
					if bugMap[checkBugPoint].hasBug {
						bugTile.adjacentBugCount++
					}
				}
				bugMap[Point{x, y}] = bugTile
			}
		}

		if previousLayoutsMap[layoutString] == true {
			biodiversityRatingOfFirstDoubleLayout = totalBioDiversityRating
			foundDoubleBioRating = true
		} else {
			previousLayoutsMap[layoutString] = true
			for y = 0; y <= height; y++ {
				for x = 0; x <= width; x++ {
					bugTile := bugMap[Point{x, y}]
					if bugTile.hasBug {
						if bugTile.adjacentBugCount != 1 {
							bugTile.tile = "."
							bugTile.hasBug = false
						}
					} else {
						if bugTile.adjacentBugCount == 1 || bugTile.adjacentBugCount == 2 {
							bugTile.tile = "#"
							bugTile.hasBug = true
						}
					}
					bugMap[Point{x, y}] = bugTile
				}
			}
		}

	}

	// Part 2 stuff
	// I make these maps so that I can get the adjacent bugs for them,
	// and the other levels should be handled from there
	levelMap[-1] = make(map[Point]BugTile)
	levelMap[1] = make(map[Point]BugTile)

	// Width and height are the edges, not a count of the items on each axis,
	// and the grid starts at 0, so I dont need to subtract 1 for an offset.
	hole := Point{width / 2, height / 2}
	var minutesPassed = 0
	for minutesPassed <= 200 {
		for level, levelBugMap := range levelMap {
			if levelMap[level-1] == nil {
				levelMap[level-1] = make(map[Point]BugTile)
			}
			if levelMap[level+1] == nil {
				levelMap[level+1] = make(map[Point]BugTile)
			}
			for y = 0; y <= height; y++ {
				for x = 0; x <= width; x++ {
					bugTile := levelBugMap[Point{x, y}]
					bugPoint := Point{x, y}
					bugTile.adjacentBugCount = 0

					if bugPoint != hole {
						for _, dir := range availableDirs {
							checkBugPoint := getPointForDirection(dir, bugPoint)
							// Handle outside edges
							if checkBugPoint.x == -1 {
								if levelMap[level-1][Point{hole.x - 1, hole.y}].hasBug {
									bugTile.adjacentBugCount++
								}
							} else if checkBugPoint.y == -1 {
								if levelMap[level-1][Point{hole.x, hole.y - 1}].hasBug {
									bugTile.adjacentBugCount++
								}

							} else if checkBugPoint.x == width+1 {
								if levelMap[level-1][Point{hole.x + 1, hole.y}].hasBug {
									bugTile.adjacentBugCount++
								}
							} else if checkBugPoint.y == width+1 {
								if levelMap[level-1][Point{hole.x, hole.y + 1}].hasBug {
									bugTile.adjacentBugCount++
								}
							}

							// Handle inside edges.
							var checkX, checkY int64
							if bugPoint == (Point{hole.x - 1, hole.y}) && checkBugPoint == hole {
								for checkY = 0; checkY <= height; checkY++ {
									if levelMap[level+1][Point{0, checkY}].hasBug {
										bugTile.adjacentBugCount++
									}
								}
							} else if bugPoint == (Point{hole.x, hole.y - 1}) && checkBugPoint == hole {
								for checkX = 0; checkX <= width; checkX++ {
									if levelMap[level+1][Point{checkX, 0}].hasBug {
										bugTile.adjacentBugCount++
									}
								}
							} else if bugPoint == (Point{hole.x + 1, hole.y}) && checkBugPoint == hole {
								for checkY = 0; checkY <= height; checkY++ {
									if levelMap[level+1][Point{width, checkY}].hasBug {
										bugTile.adjacentBugCount++
									}
								}
							} else if bugPoint == (Point{hole.x, hole.y + 1}) && checkBugPoint == hole {
								for checkX = 0; checkX <= width; checkX++ {
									if levelMap[level+1][Point{checkX, height}].hasBug {
										bugTile.adjacentBugCount++
									}
								}
							}

							// Handle everything else
							if levelBugMap[checkBugPoint].hasBug {
								bugTile.adjacentBugCount++
							}
						}
					}
					levelBugMap[Point{x, y}] = bugTile
				}
			}
		}

		if minutesPassed == 200 {
			for _, levelBugMap := range levelMap {
				for y = 0; y <= height; y++ {
					for x = 0; x <= width; x++ {
						bugTile := levelBugMap[Point{x, y}]
						if bugTile.hasBug {
							bugsPressentAfter200Minutes++
						}
					}
				}
			}
		} else {
			for _, levelBugMap := range levelMap {
				for y = 0; y <= height; y++ {
					for x = 0; x <= width; x++ {
						bugTile := levelBugMap[Point{x, y}]
						if bugTile.hasBug {
							if bugTile.adjacentBugCount != 1 {
								bugTile.tile = "."
								bugTile.hasBug = false
							}
						} else {
							if bugTile.adjacentBugCount == 1 || bugTile.adjacentBugCount == 2 {
								bugTile.tile = "#"
								bugTile.hasBug = true
							}
						}
						levelBugMap[Point{x, y}] = bugTile
					}
				}
			}
		}
		minutesPassed++
	}

	fmt.Println("Biodiversity rating of first layout that appears twice:", biodiversityRatingOfFirstDoubleLayout)
	fmt.Println("Bugs present after 200 minutes:", bugsPressentAfter200Minutes)
}

func adventureGameIOHandlers() (inputHandler, outputHandler) {
	var command string
	var currentCommandIndex int = 0
	var haveCommandToSend bool = false
	in := func() int64 {
		for {
			if !haveCommandToSend {
				var input string = ""
				scanner := bufio.NewScanner(os.Stdin)
				scanner.Scan()
				input = scanner.Text()
				input = strings.ToLower(input)
				if input == "-h" || input == "help" || input == "?" {
					fmt.Println("You can use commands 'north', 'south', 'east', 'west' to move the droid")
					fmt.Println("You can use commands 'take (item)' and 'drop (item)' to add and remove items from your inventory")
					fmt.Println("You can use 'inv' to list everything in your inventory")
				} else if input != "" {
					command = input
					haveCommandToSend = true
				}
			} else {
				if currentCommandIndex == len(command) {
					currentCommandIndex = 0
					haveCommandToSend = false
					return NEWLINE
				}
				if string(command[currentCommandIndex]) == " " {
					currentCommandIndex++
					return SPACE
				}
				asciiNumChar := int64(command[currentCommandIndex])
				currentCommandIndex++

				return asciiNumChar
			}
		}
	}

	out := func(output int64) {
		fmt.Print(string(output))
	}

	return in, out
}

func day25() {
	file, err := os.Open("day25IntCodeAdventureGameInput.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	adventureGameIntCodeString := scanner.Text()
	adventureGameIntCodeStringArr := strings.Split(adventureGameIntCodeString, ",")
	var adventureGameIntCodeArr []int64
	for _, stringIntCode := range adventureGameIntCodeStringArr {
		intCode, _ := strconv.ParseInt(stringIntCode, 10, 64)
		adventureGameIntCodeArr = append(adventureGameIntCodeArr, intCode)
	}
	adventureGameIntCodeProgram := makeMapForArray(adventureGameIntCodeArr)
	in, out := adventureGameIOHandlers()
	intComp(adventureGameIntCodeProgram, in, out)

	/*
		Needed Items:
		ornament
		easter egg
		hypercube
		monolith
	*/

}

var asteroidsData = `.#......##.#..#.......#####...#..
...#.....##......###....#.##.....
..#...#....#....#............###.
.....#......#.##......#.#..###.#.
#.#..........##.#.#...#.##.#.#.#.
..#.##.#...#.......#..##.......##
..#....#.....#..##.#..####.#.....
#.............#..#.........#.#...
........#.##..#..#..#.#.....#.#..
.........#...#..##......###.....#
##.#.###..#..#.#.....#.........#.
.#.###.##..##......#####..#..##..
.........#.......#.#......#......
..#...#...#...#.#....###.#.......
#..#.#....#...#.......#..#.#.##..
#.....##...#.###..#..#......#..##
...........#...#......#..#....#..
#.#.#......#....#..#.....##....##
..###...#.#.##..#...#.....#...#.#
.......#..##.#..#.............##.
..###........##.#................
###.#..#...#......###.#........#.
.......#....#.#.#..#..#....#..#..
.#...#..#...#......#....#.#..#...
#.#.........#.....#....#.#.#.....
.#....#......##.##....#........#.
....#..#..#...#..##.#.#......#.#.
..###.##.#.....#....#.#......#...
#.##...#............#..#.....#..#
.#....##....##...#......#........
...#...##...#.......#....##.#....
.#....#.#...#.#...##....#..##.#.#
.#.#....##.......#.....##.##.#.##`

var asteroidsTestData = `.#..##.###...#######
##.############..##.
.#.######.########.#
.###.#######.####.#.
#####.##.#.##.###.##
..#####..#.#########
####################
#.####....###.#.#.##
##.#################
#####.##.###..####..
..######..##.#######
####.##.####...##..#
.#####..#.######.###
##...#.##########...
#.##########.#######
.####.#.###.###.#.##
....##.##.###..#####
.#.#.###########.###
#.#.#.#####.####.###
###.##.####.##.#..##`

const day1Data = `73617
104372
131825
85022
105514
78478
87420
118553
97680
89479
146989
79746
108085
117895
143811
102509
102382
92975
72978
94208
130521
83042
138605
107566
63374
71176
129487
118408
115425
63967
98282
121829
92834
61084
70122
87221
132773
141347
133225
81199
94994
60881
110074
63499
143107
76618
86818
135394
106908
96085
99801
112903
51751
56002
70924
62180
133025
68025
122660
64898
77339
62109
133891
134460
84224
54836
59748
125540
67796
71845
92899
130103
74612
136820
96212
132002
97405
82629
63717
62805
112693
147810
139827
116220
69711
50236
137833
103743
147456
112098
84867
75615
132738
81072
89444
58443
94465
112494
82127
132533`
