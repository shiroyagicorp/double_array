package main

import (
	"bufio"
	"errors"
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

func checkDA(names []double_array.Item, da double_array.DoubleArray) error {
	inverse := double_array.ToInverseID(da)
	progress := 0.1
	for i, name := range names {
		ratio := float64(i+1) / float64(len(names))
		if ratio > progress {
			fmt.Printf("%d%%\n", int(ratio*100))
			progress += 0.1
		}

		itemID := da.Lookup(name)
		if itemID == double_array.ItemNotFound {
			return errors.New(
				fmt.Sprintf("failed to lookup %s (%d)", string(name), i))
		}

		deserialized := double_array.Deserialize(da, itemID, inverse)
		if deserialized != string(name) {
			return errors.New(
				fmt.Sprintf("failed to deserialize %s (%d): %s", string(name), i, deserialized))
		}
	}
	return nil
}

func main() {
	flagNames := flag.String("name", "", "name file path")
	flagModel := flag.String("model", "", "model file path")

	flag.Parse()

	if *flagNames == "" || *flagModel == "" {
		flag.Usage()
		os.Exit(1)
	}

	data, err := ioutil.ReadFile(*flagModel)
	if err != nil {
		log.Fatal(err)
	}

	names := readNamesFromFile(*flagNames)
	da, err := double_array.NewDoubleArrayFromBytes(data)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("...done reading DA model...")

	err = checkDA(names, da)
	if err != nil {
		log.Fatal(err)
	}
}
