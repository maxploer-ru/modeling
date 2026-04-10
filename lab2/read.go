package main

import (
	"log"
	"os"
	"strconv"
	"strings"
)

func readTable1(filename string) ([]float64, []float64, []float64) {
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	var I []float64
	var T0 []float64
	var m []float64
	for i, line := range lines {
		if i == 0 {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) < 3 {
			continue
		}
		iVal, _ := strconv.ParseFloat(parts[0], 64)
		t0Val, _ := strconv.ParseFloat(parts[1], 64)
		mVal, _ := strconv.ParseFloat(parts[2], 64)
		I = append(I, iVal)
		T0 = append(T0, t0Val)
		m = append(m, mVal)
	}
	return I, T0, m
}

func readTable2(filename string) ([]float64, []float64) {
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	var T []float64
	var sigma []float64
	for i, line := range lines {
		if i == 0 {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}
		tVal, _ := strconv.ParseFloat(parts[0], 64)
		sVal, _ := strconv.ParseFloat(parts[1], 64)
		T = append(T, tVal)
		sigma = append(sigma, sVal)
	}
	return T, sigma
}
