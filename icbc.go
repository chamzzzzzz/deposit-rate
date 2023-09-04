package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ICBC struct {
}

func (i *ICBC) Name() string {
	return "ICBC"
}

func (i *ICBC) GetDepositRate() (*DepositRate, error) {
	dates, err := i.getDateList()
	if err != nil {
		return nil, err
	}
	rate := &DepositRate{}
	for _, date := range dates {
		record, err := i.getDateRate(date)
		if err != nil {
			return nil, err
		}
		rate.Records = append(rate.Records, record)
	}
	return rate, nil
}

func (i *ICBC) getDateList() ([]string, error) {
	client := &http.Client{}
	resp, err := client.Get("https://papi.icbc.com.cn/interestRate/deposit/queryRMBDepositDateList?type=CH")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Code    int
		Message string
		Data    []string
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}
	if result.Code != 0 {
		return nil, fmt.Errorf("code is not 0")
	}
	return result.Data, nil
}

func (i *ICBC) getDateRate(date string) (*Record, error) {
	client := &http.Client{}
	resp, err := client.Get(fmt.Sprintf("https://papi.icbc.com.cn/interestRate/deposit/queryRMBDepositInfo?type=CH&date=%s", date))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Code    int
		Message string
		Data    struct {
			ThreeMonths string `json:"zczqThreeMonths"`
			SixMonths   string `json:"zczqHalfYear"`
			OneYear     string `json:"zczqOneYear"`
			TwoYears    string `json:"zczqTwoYears"`
			ThreeYears  string `json:"zczqThreeYears"`
			FiveYears   string `json:"zczqFiveYears"`
		}
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}
	if result.Code != 0 {
		return nil, fmt.Errorf("code is not 0")
	}
	return &Record{
		Date:        date,
		ThreeMonths: result.Data.ThreeMonths,
		SixMonths:   result.Data.SixMonths,
		OneYear:     result.Data.OneYear,
		TwoYears:    result.Data.TwoYears,
		ThreeYears:  result.Data.ThreeYears,
		FiveYears:   result.Data.FiveYears,
	}, nil
}
