package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

type metricsType struct {
	mean, median, sd float32
	mode int
}


func main() {
	nFlag := flag.Int("n", 4, "Choose count of parameters to print\n")
	flag.Parse()
	metric := new(metricsType)
	buf := numberRead()
	nums := atoiSlice(buf)
	metric.paramsCalc(nums)
	metric.paramsPrint(*nFlag)
}

func numberRead() []string {
	fmt.Println("Enter input values:")
	r := bufio.NewReader(os.Stdin)
	buf, err := r.ReadString('\n')
	if err != nil && err != io.EOF {
		fmt.Println("Values reading error.")
		os.Exit(1)
	}
	return strings.Split(buf[:len(buf)-1], " ")
}

func atoiSlice(buf []string) []int {
	nums := make([]int, 0, len(buf))
	for _, elem := range buf {
		num, err := strconv.Atoi(elem)
		if err != nil {
			fmt.Println("Error. Input values must be integer.")
			os.Exit(1)
		}
		nums = append(nums, num)
	}
	sort.Ints(nums)
	return nums
}

func (m *metricsType) paramsCalc(nums []int) {
	var sum int
	for _, elem := range nums {
		sum += elem
	}
	m.mean = float32(sum) / float32(len(nums))
	m.median = float32(nums[len(nums) / 2])
	if len(nums) % 2 == 0 {
		m.median = (m.median + float32(nums[len(nums) / 2 - 1])) / 2
	}
	m.mode = modeCalc(nums)
	m.sd = sdCalc(nums, m.mean)
}

func (m metricsType) paramsPrint(count int) {
	s := fmt.Sprintf("Mean: %.1f\nMedian: %.1f\nMode: %d\nSD: %.2f\n", m.mean, m.median, m.mode, m.sd)
	params := strings.Split(s, "\n")
	for i := 0; i < count; i++ {
		fmt.Println(params[i])
	}
}

func modeCalc(nums []int) int {
	mode := [3]int {nums[0], 0, 0}
	prev := nums[0]
	for _, elem := range nums {
		if elem == mode[0] {
			mode[1]++
		} else if prev == elem{
			mode[2]++
		} else {
			mode[2] = 1
		}
		if mode[1] < mode[2] {
			mode[0] = elem
		}
		prev = elem
	}
	return mode[0]
}

func sdCalc(nums []int, mean float32) float32 {
	var sum float32
	for _, elem := range nums {
		sum += (float32(elem) - mean) * (float32(elem) - mean)
	}
	return float32(math.Sqrt(float64(sum / float32(len(nums)))))
}