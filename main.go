package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
)

func main() {
	// parameters
	var csvFile, secretsdumpFile, ntdsFile, systemFile, outDir string
	flag.StringVar(&csvFile, "csv", "", "CSV file with hashes (comma separated username and hash)")
	flag.StringVar(&secretsdumpFile, "secretsdump", "", "Direct output of impacket's secretsdump")
	flag.StringVar(&ntdsFile, "ntds", "", "ntds.dit file")
	flag.StringVar(&systemFile, "system", "", "system registry hive")
	flag.StringVar(&outDir, "output", "", "output directory")

	flag.Parse()

	// load hashes
	ad := NewAD()

	if csvFile != "" {
		err := ad.LoadHashesFromCSV(csvFile)
		if err != nil {
			log.Fatal(err)
		}
	} else if secretsdumpFile != "" {
		err := ad.LoadHashesFromSecretsdump(secretsdumpFile)
		if err != nil {
			log.Fatal(err)
		}
	} else if ntdsFile != "" && systemFile != "" {
		err := ad.LoadHashesFromNTDS(ntdsFile, systemFile)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatal("No input files")
	}

	// analyze
	ad.Analyze()

	// print output
	if outDir != "" {
		// create output directory
		if _, err := os.Stat(outDir); os.IsNotExist(err) {
			err := os.Mkdir(outDir, 0700)
			if err != nil {
				log.Fatal(err)
			}
		}

		// reused passwords CSV
		file, err := os.Create(path.Join(outDir, "reuse.csv"))
		if err != nil {
			log.Fatal(err)
		}

		w := bufio.NewWriter(file)
		for _, r := range ad.ReusedPasswords() {
			line := fmt.Sprintf("%s,%d,%s", r.Hash, r.Count, strings.Join(r.Users, ";"))
			fmt.Fprintln(w, line)
		}
		w.Flush()
		file.Close()

		// reused passwords JSON
		file, err = os.Create(path.Join(outDir, "reuse.json"))
		if err != nil {
			log.Fatal(err)
		}

		w = bufio.NewWriter(file)
		for _, r := range ad.ReusedPasswords() {
			data, err := json.Marshal(r)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Fprintln(w, string(data[:]))
		}
		w.Flush()
		file.Close()

		// reused passwords without hash CSV
		file, err = os.Create(path.Join(outDir, "reuse-without-hash.csv"))
		if err != nil {
			log.Fatal(err)
		}

		w = bufio.NewWriter(file)
		for _, r := range ad.ReusedPasswords() {
			line := fmt.Sprintf("%d,%s", r.Count, strings.Join(r.Users, ";"))
			fmt.Fprintln(w, line)
		}
		w.Flush()
		file.Close()

		// reused passwords without hash JSON
		file, err = os.Create(path.Join(outDir, "reuse-without-hash.json"))
		if err != nil {
			log.Fatal(err)
		}

		w = bufio.NewWriter(file)
		for _, r := range ad.ReusedPasswords() {
			r.Hash = ""
			data, err := json.Marshal(r)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Fprintln(w, string(data[:]))
		}
		w.Flush()
		file.Close()
	} else {
		fmt.Println("Reused passwords:")
		for _, r := range ad.ReusedPasswords() {
			if r.Count > 1 {
				fmt.Println(r)
			}
		}
	}
}
