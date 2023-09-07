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

type CCB struct {
}

func (c *CCB) Name() string {
	return "CCB"
}

func (c *CCB) GetDepositRate() (*DepositRate, error) {
	urls, err := c.getDateURLs()
	if err != nil {
		return nil, err
	}
	rate := &DepositRate{}
	for _, url := range urls {
		record, err := c.getDateRate(url)
		if err != nil {
			return nil, err
		}
		if record == nil {
			continue
		}
		rate.Records = append(rate.Records, record)
	}
	return rate, nil
}

func (c *CCB) getDateURLs() ([]string, error) {
	client := &http.Client{}
	urls := []string{}
	resp, err := client.Get("http://www3.ccb.com/chn/personal/interestv3/rmbdeposit.shtml")
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

	div, err := dom.Find("div", "class", "list_main w_135 rfloat mB10")
	if div == nil {
		return nil, err
	}

	for _, li := range div.QueryAll("li") {
		url := "http://www3.ccb.com" + li.Attribute("name")
		urls = append(urls, url)
	}
	return urls, nil
}

func (c *CCB) getDateRate(url string) (*Record, error) {
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

	div, err := dom.Find("div", "class", "content f14")
	if err != nil {
		return nil, err
	}

	h2, err := div.Find("h2", "class", "Yahei detail_title")
	if err != nil {
		return nil, err
	}
	date := strings.TrimSpace(h2.Text())
	_, err = time.Parse("2006-01-02", date)
	if err != nil {
		return nil, err
	}

	table, err := div.Find("table")
	if err != nil {
		return nil, err
	}
	trs := table.QueryAll("tr")
	if len(trs) == 8 {
		return nil, nil
	}
	if len(trs) < 11 {
		return nil, fmt.Errorf("trs len is less than 11")
	}
	record := &Record{
		Date: date,
	}
	i, j, err := c.getTrRange(trs)
	if err != nil {
		return nil, err
	}
	for _, tr := range trs[i:j] {
		k, v, err := c.getTrKV(tr)
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

func (c *CCB) getTrRange(trs supersimplesoup.Nodes) (int, int, error) {
	m, n := 0, 0
	for i, tr := range trs {
		if i == 0 {
			continue
		}
		k, _, err := c.getTrKV(tr)
		if err != nil {
			return 0, 0, err
		}
		if k == "三个月" || k == "三年" {
			m, n = i, i+6
			break
		}
	}
	if m == 0 || n == 0 || n > len(trs) {
		return 0, 0, fmt.Errorf("not found tr range")
	}
	return m, n, nil
}

func (o *CCB) getTrKV(tr *supersimplesoup.Node) (string, string, error) {
	tds := tr.QueryAll("td")
	n := len(tds)
	if n != 2 && n != 3 {
		return "", "", fmt.Errorf("tr tds len is not 2 or 3")
	}
	i, j := n-2, n-1
	k := strings.TrimSpace(tds[i].FullText())
	v := strings.TrimSpace(tds[j].FullText())
	return k, v, nil
}
