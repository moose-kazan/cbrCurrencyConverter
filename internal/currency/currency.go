package currency

import (
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/net/html/charset"
)

type CurrencyList struct {
	XMLName xml.Name       `xml:"ValCurs"`
	Date    string         `xml:"Date,attr"`
	Name    string         `xml:"name,attr"`
	Rates   []CurrencyInfo `xml:"Valute"`
}

type CurrencyListInterface interface {
	Convert(src string, dst string) (float64, error)
}

type CurrencyInfo struct {
	XMLName  xml.Name `xml:"Valute"`
	ID       xml.Attr `xml:"ID,attr"`
	NumCode  string   `xml:"NumCode"`
	CharCode string   `xml:"CharCode"`
	Nominal  int      `xml:"Nominal"`
	Name     string   `xml:"Name"`
	Value    string   `xml:"Value"`
}

type CurrencyInfoInterface interface {
	GetNominal() int
	GetValue() float64
	GetISOCode() string
	GetName() string
}

func (cl CurrencyList) Convert(src string, dst string) (float64, error) {
	if src == dst {
		return 1, nil
	}

	var srcRate *CurrencyInfo
	var dstRate *CurrencyInfo
	for _, v := range cl.Rates {
		if v.GetISOCode() == src {
			srcRate = new(CurrencyInfo)
			*srcRate = v
		}
		if v.GetISOCode() == dst {
			dstRate = new(CurrencyInfo)
			*dstRate = v
		}
		if srcRate != nil && dstRate != nil {
			break
		}
	}

	if srcRate == nil || dstRate == nil {
		return 0, errors.New("Unknown currencies!")
	} else if srcRate == nil {
		return 0, errors.New(fmt.Sprintf("Unknown currency: \"%s\"", src))
	} else if dstRate == nil {
		return 0, errors.New(fmt.Sprintf("Unknown currency: \"%s\"", dst))
	}

	rate := dstRate.GetValue() / srcRate.GetValue()
	rate *= float64(srcRate.GetNominal())
	rate /= float64(dstRate.GetNominal())

	return rate, nil
}

func (ci CurrencyInfo) GetNominal() int {
	return ci.Nominal
}

func (ci CurrencyInfo) GetValue() float64 {
	rv, err := strconv.ParseFloat(strings.Replace(ci.Value, ",", ".", 1), 64)
	if err != nil {
		return 0
	}
	return rv
}

func (ci CurrencyInfo) GetISOCode() string {
	return ci.CharCode
}

func (ci CurrencyInfo) GetName() string {
	return ci.Name
}

func New(datestr string) (*CurrencyList, error) {
	var url string = fmt.Sprintf("https://www.cbr.ru/scripts/XML_daily_eng.asp?date_req=%s", datestr)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("User-Agent", "Mozilla/5.0")
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("Bad http code: %d", resp.StatusCode))
	}

	d := xml.NewDecoder(resp.Body)
	d.CharsetReader = charset.NewReaderLabel

	var rv CurrencyList
	err = d.Decode(&rv)
	if err != nil {
		return nil, err
	}

	var rub CurrencyInfo = CurrencyInfo{
		NumCode:  "643",
		CharCode: "RUB",
		Nominal:  1,
		Name:     "Russian Ruble",
		Value:    "1",
	}
	rv.Rates = append(rv.Rates, rub)

	return &rv, nil
}
