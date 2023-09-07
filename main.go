package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

type Record struct {
	Date        string
	ThreeMonths string
	SixMonths   string
	OneYear     string
	TwoYears    string
	ThreeYears  string
	FiveYears   string
}

type DepositRate struct {
	Records []*Record
}

func (r *DepositRate) has(record *Record) bool {
	for _, r := range r.Records {
		if *r == *record {
			return true
		}
	}
	return false
}

type Bank interface {
	Name() string
	GetDepositRate() (*DepositRate, error)
}

func main() {
	banks := []Bank{
		&ICBC{},
		&ABC{},
		&BOC{},
		&CCB{},
	}
	for _, bank := range banks {
		rate, err := bank.GetDepositRate()
		if err != nil {
			fmt.Printf("get bank deposit rate fail. bank=%s, err='%v'\n", bank.Name(), err)
			return
		}
		fmt.Printf("get bank deposit rate success. bank=%s, record=%d\n", bank.Name(), len(rate.Records))
		file := fmt.Sprintf("%s.json", bank.Name())

		old, err := read(file)
		if err != nil {
			fmt.Printf("read bank deposit file fail. bank=%s, err='%v'\n", bank.Name(), err)
			return
		}
		if old != nil {
			rate = merge(rate, old)
		}

		err = write(file, rate)
		if err != nil {
			fmt.Printf("write bank deposit file fail. bank=%s, err='%v'\n", bank.Name(), err)
			return
		}
		fmt.Printf("write bank deposit file success. bank=%s, file=%s\n", bank.Name(), file)
	}
}

func write(file string, rate *DepositRate) error {
	b, err := json.MarshalIndent(rate, "", "    ")
	if err != nil {
		return err
	}
	return os.WriteFile(file, b, 0644)
}

func read(file string) (*DepositRate, error) {
	b, err := os.ReadFile(file)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var rate DepositRate
	err = json.Unmarshal(b, &rate)
	if err != nil {
		return nil, err
	}
	return &rate, nil
}

func merge(new, old *DepositRate) *DepositRate {
	merged := &DepositRate{}
	for _, r := range new.Records {
		if merged.has(r) {
			continue
		}
		merged.Records = append(merged.Records, r)
	}
	for _, r := range old.Records {
		if merged.has(r) {
			continue
		}
		merged.Records = append(merged.Records, r)
	}

	sort.Slice(merged.Records, func(i, j int) bool {
		return merged.Records[i].Date > merged.Records[j].Date
	})
	return merged
}
