package main

import (
	"bufio"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/C-Sto/gosecretsdump/pkg/ditreader"
)

// AD stores the users and password hashes from Active Directory
type AD struct {
	hashes map[string]string

	reuse []*ReusedPassword
}

// NewAD returns a new AD instance
func NewAD() *AD {
	ad := AD{}

	ad.hashes = make(map[string]string)

	return &ad
}

// LoadHashesFromCSV loads the password hashes from a CSV file (comma separated)
func (ad *AD) LoadHashesFromCSV(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// split at ;
		parts := strings.Split(line, ",")
		if len(parts) != 2 {
			log.Print("Ignoring line: " + line)
			continue
		}

		ad.hashes[parts[0]] = strings.ToLower(parts[1])
	}

	err = scanner.Err()
	if err != nil {
		return err
	}

	return nil
}

// LoadHashesFromSecretsdump loads the password hashes from a text file with the secretsdump output
func (ad *AD) LoadHashesFromSecretsdump(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// split at :
		parts := strings.Split(line, ":")
		if len(parts) != 7 {
			continue
		}

		ad.hashes[parts[0]] = strings.ToLower(parts[3])
	}

	err = scanner.Err()
	if err != nil {
		return err
	}

	return nil
}

// LoadHashesFromNTDS loads the password hashes from a file
func (ad *AD) LoadHashesFromNTDS(ntdsFile string, systemFile string) error {
	dr, err := ditreader.New(systemFile, ntdsFile)
	if err != nil {
		return err
	}

	data := dr.GetOutChan()

	for d := range data {
		// check dh.UAC.AccountDisable ?
		line := d.HashString()
		parts := strings.Split(line, ":")
		if len(parts) != 7 {
			log.Print("Ignoring line: " + line)
			continue
		}

		ad.hashes[parts[0]] = strings.ToLower(parts[3])
	}

	return nil
}

// Analyze runs the analysis of the collected password hashs
func (ad *AD) Analyze() {
	ad.findReusedPasswords()
}

// ReusedPasswords returns the reused passwords and according users
func (ad *AD) ReusedPasswords() []*ReusedPassword {
	return ad.reuse
}

// FindReusedPasswords searches for accounts that have the same password hash
func (ad *AD) findReusedPasswords() {
	// a hashtable is used to easily check, whether the hash was found already or not
	store := make(map[string]*ReusedPassword)

	// search
	for user, hash := range ad.hashes {
		// ignore empty ntlm hash
		if hash == "31d6cfe0d16ae931b73c59d7e0c089c0" {
			continue
		}

		// store
		_, reused := store[hash]
		if reused {
			store[hash].Add(user)
		} else {
			store[hash] = NewReusedPassword(user, hash)
		}
	}

	// count how many passwords are used more than once (so we know how big the result list is)
	numReuses := 0
	for _, reusedpassword := range store {
		if reusedpassword.Count > 1 {
			numReuses++
		}
	}

	// filter and convert to output format (list instead of hashmap)
	result := make([]*ReusedPassword, numReuses)
	counter := 0
	for _, reusedpassword := range store {
		if reusedpassword.Count > 1 {
			result[counter] = reusedpassword
			counter++
		}
	}

	// sort
	sort.Slice(result, func(i, j int) bool {
		return result[i].Count > result[j].Count
	})

	// done
	ad.reuse = result
}
