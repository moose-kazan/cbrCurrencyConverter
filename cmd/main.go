package main

import (
	"fmt"
	"runtime"

	"cbr/internal/currency"

	"bitbucket.org/rj/goey"
	"bitbucket.org/rj/goey/base"
	"bitbucket.org/rj/goey/loop"
)

var (
	mainWindow    *goey.Window
	convertSrc    *goey.SelectInput = new(goey.SelectInput)
	convertDst    *goey.SelectInput = new(goey.SelectInput)
	convertResult string
	cRates        currency.CurrencyList
)

func main() {
	err := loop.Run(createMainWindow)
	if err != nil {
		fmt.Println("Error: ", err)
	}
}

func runConvert() {
	if convertSrc.Value < 0 {
		return
	}
	if convertDst.Value < 0 {
		return
	}
	valSrc := convertSrc.Items[convertSrc.Value][0:3]
	valDst := convertDst.Items[convertDst.Value][0:3]

	rate, _ := cRates.Convert(valSrc, valDst)

	resultFormat := "1 %s = %f %s\n1 %s = %f %s\n"
	if runtime.GOOS == "windows" {
		resultFormat = "1 %s = %f %s\r\n1 %s = %f %s\r\n"
	} else if runtime.GOOS == "darwin" {
		resultFormat = "1 %s = %f %s\r1 %s = %f %s\r"
	}
	convertResult = fmt.Sprintf(
		resultFormat,
		valSrc, 1/rate, valDst,
		valDst, rate, valSrc,
	)

}

func createMainWindow() error {
	cRates = *currency.New()
	err := cRates.Fetch("")

	if err != nil {
		panic(err)
	}

	var currencyTitles []string
	for _, v := range cRates.Rates {
		currencyTitles = append(
			currencyTitles,
			fmt.Sprintf("%s (%s)", v.GetISOCode(), v.GetName()))
	}

	convertSrc.Items = currencyTitles
	convertSrc.OnChange = func(c int) {
		convertSrc.Value = c
		runConvert()
		updateMainWindow()
	}

	convertDst.Items = currencyTitles
	convertDst.OnChange = func(c int) {
		convertDst.Value = c
		runConvert()
		updateMainWindow()
	}

	mw, err := goey.NewWindow("CBR Currency Converter", renderMainWindow())

	if err != nil {
		return err
	}

	mw.SetScroll(false, true)
	mainWindow = mw

	return nil
}

func renderMainWindow() base.Widget {
	return &goey.VBox{
		Children: []base.Widget{
			&goey.P{
				Text:  "Currency converter based on data from\nCentral Bank of Russia",
				Align: goey.JustifyCenter,
			},

			convertSrc,
			convertDst,
			&goey.Expand{
				Child: &goey.TextArea{
					ReadOnly: true,
					Value:    convertResult,
					MinLines: 10,
				},
			},
		},
	}
}

func updateMainWindow() {
	err := mainWindow.SetChild(renderMainWindow())

	if err != nil {
		fmt.Println(err)
	}
}
