// this will eventually just be the main package, and it will run all the things here
package main

import (
	"fmt"
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
	dataArray := strings.Split(data, "\n")
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

func getResult(intCodeData []int64) (atZero int64) {
	var copiedData = make([]int64, len(intCodeData))
	copy(copiedData, intCodeData)
	currentPos := 0
	currentOpCode := copiedData[currentPos]
	for currentOpCode != 99 && currentPos < len(copiedData) {
		if currentOpCode == 1 {
			firstVal := copiedData[copiedData[currentPos+1]]
			secondVal := copiedData[copiedData[currentPos+2]]
			placePos := copiedData[currentPos+3]
			copiedData[placePos] = firstVal + secondVal
		}

		if currentOpCode == 2 {
			firstVal := copiedData[copiedData[currentPos+1]]
			secondVal := copiedData[copiedData[currentPos+2]]
			placePos := copiedData[currentPos+3]
			copiedData[placePos] = firstVal * secondVal
		}

		currentPos += 4
		currentOpCode = copiedData[currentPos]
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

const data = `73617
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
