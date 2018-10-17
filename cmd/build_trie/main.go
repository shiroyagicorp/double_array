package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/shiroyagicorp/double_array"
)

func readNamesFromFile(path string) []double_array.Item {
	fmt.Printf("reading %s", path)
	fp, err := os.Open(path)
	defer fp.Close()
	if err != nil {
		log.Fatalf("failed to open file: %s", path)
	}

	scanner := bufio.NewScanner(fp)
	names := make([]double_array.Item, 0)
	i := 0
	for scanner.Scan() {
		name := scanner.Text()
		names = append(names, double_array.Item([]rune(name)))
		i++
		if i%10000 == 0 {
			fmt.Print(".")
		}
	}
	fmt.Printf("done [count=%d, cap=%d]\n", len(names), cap(names))

	return names
}

func main() {
	flagIn := flag.String("in", "", "input file path")
	flagOut := flag.String("out", "", "output file path")

	flag.Parse()

	if *flagIn == "" || *flagOut == "" {
		flag.Usage()
		os.Exit(1)
	}

	names := readNamesFromFile(*flagIn)
	da, err := double_array.NewDoubleArray(names)
	if err != nil {
		log.Fatal(err)
	}

	data, err := da.Serialize()
	if err != nil {
		log.Fatal(err)
	}

	ioutil.WriteFile(*flagOut, data, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
