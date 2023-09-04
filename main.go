package main

import (
	"encoding/json"
	"fmt"
	"os"
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

type Bank interface {
	Name() string
	GetDepositRate() (*DepositRate, error)
}

func main() {
	banks := []Bank{
		&ABC{},
	}
	for _, bank := range banks {
		rate, err := bank.GetDepositRate()
		if err != nil {
			fmt.Printf("get bank deposit rate fail. bank=%s, err='%v'\n", bank.Name(), err)
			return
		}
		file := fmt.Sprintf("%s.json", bank.Name())
		err = write(file, rate)
		if err != nil {
			fmt.Printf("write bank deposit file faile. bank=%s, err='%v'\n", bank.Name(), err)
			return
		}
	}
}

func write(file string, rate *DepositRate) error {
	b, err := json.MarshalIndent(rate, "", "    ")
	if err != nil {
		return err
	}
	return os.WriteFile(file, b, 0644)
}
