package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/chamzzzzzz/supersimplesoup"
)

type BOC struct {
}

func (o *BOC) Name() string {
	return "BOC"
}

func (o *BOC) GetDepositRate() (*DepositRate, error) {
	urls, err := o.getDateURLs()
	if err != nil {
		return nil, err
	}
	rate := &DepositRate{}
	for _, url := range urls {
		record, err := o.getDateRate(url)
		if err != nil {
			return nil, err
		}
		rate.Records = append(rate.Records, record)
	}
	return rate, nil
}

func (o *BOC) getDateURLs() ([]string, error) {
	client := &http.Client{}
	urls := []string{}
	for i := 0; i < 5; i++ {
		url := "https://www.bankofchina.com/fimarkets/lilv/fd31/"
		if i > 0 {
			url += fmt.Sprintf("index_%d.html", i)
		}
		resp, err := client.Get(url)
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

		div := dom.Query("div", "class", "news")
		if div == nil {
			break
		}

		for _, a := range div.Query("ul", "class", "list").QueryAll("a") {
			url := "https://www.bankofchina.com/fimarkets/lilv/fd31/" + strings.TrimLeft(a.Href(), "./")
			urls = append(urls, url)
		}
	}
	return urls, nil
}

func (o *BOC) getDateRate(url string) (*Record, error) {
	client := &http.Client{}
	resp, err := client.Get(url)
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

	div, err := dom.Find("div", "class", "content con_area")
	if err != nil {
		return nil, err
	}

	h2, err := div.Find("h2", "class", "title")
	if err != nil {
		return nil, err
	}
	date := strings.ReplaceAll(h2.Text(), "人民币存款利率表", "")
	_, err = time.Parse("2006-01-02", date)
	if err != nil {
		return nil, err
	}

	table, err := div.Find("table")
	if err != nil {
		return nil, err
	}
	trs := table.QueryAll("tr")
	if len(trs) < 11 {
		return nil, fmt.Errorf("trs len is less than 11")
	}
	record := &Record{
		Date: date,
	}
	i, j, err := o.getTrRange(trs)
	if err != nil {
		return nil, err
	}
	for _, tr := range trs[i:j] {
		k, v, err := o.getTrKV(tr)
		if err != nil {
			return nil, err
		}
		switch k {
		case "三个月":
			record.ThreeMonths = v
		case "半年":
			record.SixMonths = v
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
	return record, nil
}

func (o *BOC) getTrRange(trs supersimplesoup.Nodes) (int, int, error) {
	m, n := 0, 0
	for i, tr := range trs {
		if i == 0 {
			continue
		}
		k, _, err := o.getTrKV(tr)
		if err != nil {
			return 0, 0, err
		}
		if k == "三个月" {
			m, n = i, i+6
			break
		}
	}
	if m == 0 || n == 0 || n > len(trs) {
		return 0, 0, fmt.Errorf("not found tr range")
	}
	return m, n, nil
}

func (o *BOC) getTrKV(tr *supersimplesoup.Node) (string, string, error) {
	tds := tr.QueryAll("td")
	if len(tds) != 2 {
		return "", "", fmt.Errorf("tr tds len is not 2")
	}
	k := strings.TrimSpace(tds[0].Text())
	v := strings.TrimSpace(tds[1].Text())
	return k, v, nil
}
