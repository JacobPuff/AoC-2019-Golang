package main

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"

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
func intComp(intCodeData map[int64]int64, getInput inputHandler, sendOutput outputHandler) int64 {
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
	return copiedData[0]
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
	intCodeDataResultAtZero := intComp(intCodeMap, getUserInput, printOutput)
	fmt.Println("initial result at zero: ", intCodeDataResultAtZero)

	var noun, verb int64
	for noun = 1; noun <= 99; noun++ {
		for verb = 1; verb <= 99; verb++ {
			intCodeData[1] = noun
			intCodeData[2] = verb
			if intComp(intCodeMap, getUserInput, printOutput) == 19690720 {
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
