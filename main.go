// this will eventually just be the main package, and it will run all the things here
package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
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

func getParamWithMode(data []int64, mode int64, location int64) (paramVal int64) {
	switch mode {
	case 0:
		return data[data[location]]
	case 1:
		return data[location]
	default:
		fmt.Printf("No mode with value %d when getting location %d", mode, location)
		return -1
	}
}

func setParamWithMode(data []int64, mode int64, location int64, value int64) (paramVal int64) {
	switch mode {
	case 0:
		data[data[location]] = value
	case 1:
		data[location] = value
	default:
		fmt.Printf("Cant set with mode of value %d when setting location %d", mode, location)
		return -1
	}
	return 0
}

//Intcode computer
func getResult(intCodeData []int64) (atZero int64) {
	var copiedData = make([]int64, len(intCodeData))
	copy(copiedData, intCodeData)
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
			firstVal = getParamWithMode(copiedData, mode1, currentPos+1)
			secondVal = getParamWithMode(copiedData, mode2, currentPos+2)
			result = setParamWithMode(copiedData, mode3, currentPos+3, firstVal+secondVal)
			if result == -1 {
				running = false
			}
			currentPos += 4
		}

		if currentOpCode == 2 {
			firstVal = getParamWithMode(copiedData, mode1, currentPos+1)
			secondVal = getParamWithMode(copiedData, mode2, currentPos+2)
			result = setParamWithMode(copiedData, mode3, currentPos+3, firstVal*secondVal)
			if result == -1 {
				running = false
			}
			currentPos += 4
		}

		if currentOpCode == 3 {
			fmt.Println("Need number: ")
			var input int64
			_, _ = fmt.Scanf("%d", &input)
			// TODO: Stop people from inputing non-Numbers
			result = setParamWithMode(copiedData, mode1, currentPos+1, input)
			if result == -1 {
				running = false
			}
			fmt.Println("input: ", input)
			currentPos += 2
		}

		if currentOpCode == 4 {
			var output int64
			if mode1 == 0 {
				output = copiedData[copiedData[currentPos+1]]
			} else {
				output = copiedData[currentPos+1]
			}
			fmt.Println("Output: ", output)
			currentPos += 2
		}

		if currentOpCode == 5 {
			firstVal = getParamWithMode(copiedData, mode1, currentPos+1)
			secondVal = getParamWithMode(copiedData, mode2, currentPos+2)
			if firstVal != 0 {
				currentPos = secondVal
			} else {
				currentPos += 3
			}
		}

		if currentOpCode == 6 {
			firstVal = getParamWithMode(copiedData, mode1, currentPos+1)
			secondVal = getParamWithMode(copiedData, mode2, currentPos+2)
			if firstVal == 0 {
				currentPos = secondVal
			} else {
				currentPos += 3
			}
		}

		if currentOpCode == 7 {
			firstVal = getParamWithMode(copiedData, mode1, currentPos+1)
			secondVal = getParamWithMode(copiedData, mode2, currentPos+2)
			if firstVal < secondVal {
				result = setParamWithMode(copiedData, mode3, currentPos+3, 1)
			} else {
				result = setParamWithMode(copiedData, mode3, currentPos+3, 0)
			}

			if result == -1 {
				running = false
			}
			currentPos += 4
		}

		if currentOpCode == 8 {
			firstVal = getParamWithMode(copiedData, mode1, currentPos+1)
			secondVal = getParamWithMode(copiedData, mode2, currentPos+2)
			if firstVal == secondVal {
				result = setParamWithMode(copiedData, mode3, currentPos+3, 1)
			} else {
				result = setParamWithMode(copiedData, mode3, currentPos+3, 0)
			}

			if result == -1 {
				running = false
			}
			currentPos += 4
		}

		if currentOpCode == 99 {
			running = false
		}

		instruction = copiedData[currentPos]
	}

	atZero = copiedData[0]
	return atZero
}

func day2() {
	var intCodeData = []int64{1, 12, 2, 3, 1, 1, 2, 3, 1, 3, 4, 3, 1, 5, 0, 3, 2, 13, 1, 19, 1, 5, 19, 23, 2, 10, 23, 27, 1, 27, 5, 31, 2, 9, 31, 35, 1, 35, 5, 39, 2, 6, 39, 43, 1, 43, 5, 47, 2, 47, 10, 51, 2, 51, 6, 55, 1, 5, 55, 59, 2, 10, 59, 63, 1, 63, 6, 67, 2, 67, 6, 71, 1, 71, 5, 75, 1, 13, 75, 79, 1, 6, 79, 83, 2, 83, 13, 87, 1, 87, 6, 91, 1, 10, 91, 95, 1, 95, 9, 99, 2, 99, 13, 103, 1, 103, 6, 107, 2, 107, 6, 111, 1, 111, 2, 115, 1, 115, 13, 0, 99, 2, 0, 14, 0}

	//Gets result for initial data
	intCodeDataResultAtZero := getResult(intCodeData)
	fmt.Println("initial result at zero: ", intCodeDataResultAtZero)

	//Part 2 is what the copiedData is for, so that it can test a whole bunch of nouns and verbs
	//without changing the data permanently
	var noun, verb int64
	for noun = 1; noun <= 99; noun++ {
		for verb = 1; verb <= 99; verb++ {
			intCodeData[1] = noun
			intCodeData[2] = verb
			if getResult(intCodeData) == 19690720 {
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
	// Test data. Prints 999 if input below 8, 1000 if equal, 1001 if greater.
	//inputData := []int64{3, 21, 1008, 21, 8, 20, 1005, 20, 22, 107, 8, 21, 20, 1006, 20, 31, 1106, 0, 36, 98, 0, 0, 1002, 21, 125, 20, 4, 20, 1105, 1, 46, 104, 999, 1105, 1, 46, 1101, 1000, 1, 20, 4, 20, 1105, 1, 46, 98, 99}
	// Test data. if input not 0 it prints 1, if input 0 it prints 0
	//inputData := []int64{3, 12, 6, 12, 15, 1, 13, 14, 13, 4, 13, 99, -1, 0, 1, 9}
	getResult(inputData)
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
	file, err := os.Open("./day2Data.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	day2Data := make(map[string]Body)
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
		if _, ok := day2Data[orbitArr[0]]; !ok {
			day2Data[orbitArr[0]] = Body{"", []string{}}
		}

		if _, ok := day2Data[orbitArr[1]]; !ok {
			day2Data[orbitArr[1]] = Body{"", []string{}}
		}

		mapBody := day2Data[orbitArr[0]]
		mapMoon := day2Data[orbitArr[1]]

		mapBody.setBodyMoons(orbitArr[1])
		mapMoon.setBodyParent(orbitArr[0])

		day2Data[orbitArr[0]] = mapBody
		day2Data[orbitArr[1]] = mapMoon

	}
	breadthFirstQueue = append(breadthFirstQueue, "COM")
	for len(breadthFirstQueue) > 0 {
		currentBodyName, breadthFirstQueue = breadthFirstQueue[0], breadthFirstQueue[1:]
		currentBody := day2Data[currentBodyName]
		//fmt.Println("Name:", currentBodyName, "Parent:", currentBody.parent, "Moons:", currentBody.moons)
		if len(currentBody.moons) == 0 {
			edgeBodies++
		}
		for _, orbits := range currentBody.moons {
			breadthFirstQueue = append(breadthFirstQueue, orbits)
			directOrbits++
		}
		indirectOrbitParent := day2Data[currentBody.parent].parent
		for indirectOrbitParent != "" {
			indirectOrbitParent = day2Data[indirectOrbitParent].parent
			indirectOrbits++
		}
	}

	pathBody = "YOU"
	for pathBody != "" {
		pathBody = day2Data[pathBody].parent
		youOrbitArray = append(youOrbitArray, pathBody)
	}

	pathBody = "SAN"
	for pathBody != "" {
		pathBody = day2Data[pathBody].parent
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
