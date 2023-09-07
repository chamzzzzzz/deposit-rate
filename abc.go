package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/chamzzzzzz/supersimplesoup"
)

type ABC struct {
}

func (i *ABC) Name() string {
	return "ABC"
}

func (i *ABC) GetDepositRate() (*DepositRate, error) {
	client := &http.Client{}
	resp, err := client.Get("https://www.abchina.com/cn/PersonalServices/Quotation/bwbll/")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	dom, err := supersimplesoup.Parse(bytes.NewReader(b))
	if err != nil {
		return nil, err
	}

	div, err := dom.Find("div", "class", "Custom_UnionStyle")
	if err != nil {
		return nil, err
	}

	center, err := div.Find("center")
	if err != nil {
		return nil, err
	}
	re := regexp.MustCompile(`\d{4}年\d{1,2}月\d{1,2}日`)
	match := re.FindString(center.Text())
	if match == "" {
		return nil, fmt.Errorf("not found date")
	}
	d, err := time.Parse("2006年1月2日", match)
	if err != nil {
		return nil, err
	}
	date := d.Format("2006-01-02")

	table, err := div.Find("table", "class", "DataList")
	if err != nil {
		return nil, err
	}
	trs := table.QueryAll("tr")
	if len(trs) != 20 {
		return nil, fmt.Errorf("trs len is not 20")
	}
	record := &Record{
		Date: date,
	}
	for _, tr := range trs[5:11] {
		tds := tr.QueryAll("td")
		if len(tds) != 2 {
			continue
		}
		k := strings.TrimSpace(tds[0].Text())
		v := strings.TrimSpace(tds[1].Text())
		switch k {
		case "三个月":
			record.ThreeMonths = v
		case "六个月":
			record.SixMonths = v
		case "一年":
			record.OneYear = v
		case "二年":
			record.TwoYears = v
		case "三年":
			record.ThreeYears = v
		case "五年":
			record.FiveYears = v
		default:
			return nil, fmt.Errorf("unknown key: %s", k)
		}
	}
	rate := &DepositRate{
		Records: []*Record{record},
	}
	return rate, nil
}
