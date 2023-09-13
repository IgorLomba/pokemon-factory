package services

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
)

func AppendCapabilities(current []string, capabilities []string) []string {
	for i := 0; i < 4; i++ {
		capability := capabilities[rand.Intn(len(capabilities)-1)]
		if !contains(current, capability) {
			current = append(current, capability)
		}
	}
	return current
}

func GetNameAndAdjFromFiles(namesFile, adjFileName, capabilitiesFile string) ([]string, []string, []string) {
	var names []string
	var adj []string
	var capabilities []string
	fn, err := os.Open(namesFile)
	if err != nil {
		fmt.Println("please try later, pokemon factory unavailable")
		return nil, nil, nil
	}
	defer func(fn *os.File) {
		err := fn.Close()
		if err != nil {
			panic(err)
		}
	}(fn)
	var scanner = bufio.NewScanner(fn)
	for scanner.Scan() {
		names = append(names, scanner.Text())
	}
	fa, err := os.Open(adjFileName)
	if err != nil {
		panic(err)
	}
	defer func(fa *os.File) {
		err := fa.Close()
		if err != nil {
			panic(err)
		}
	}(fa)
	scanner = bufio.NewScanner(fa)
	for scanner.Scan() {
		adj = append(adj, scanner.Text())
	}
	capFile, err := os.Open(adjFileName)
	if err != nil {
		panic(err)
	}
	defer func(capFile *os.File) {
		err := capFile.Close()
		if err != nil {
			panic(err)
		}
	}(capFile)
	scanner = bufio.NewScanner(capFile)
	for scanner.Scan() {
		capabilities = append(capabilities, scanner.Text())
	}
	return names, adj, capabilities
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if strings.ToLower(s) == strings.ToLower(item) {
			return true
		}
	}
	return false
}
