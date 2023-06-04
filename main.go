package main

import (
	"cbr/internal/apptheme"
	"cbr/internal/currency"
	"fmt"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var (
	a          fyne.App
	w          fyne.Window
	cRates     *currency.CurrencyList
	selectSrc  *widget.Select
	selectDst  *widget.Select
	resultItem *widget.Entry
)

func runConvert(_ string) {
	if selectSrc.SelectedIndex() == -1 {
		return
	}
	if selectDst.SelectedIndex() == -1 {
		return
	}
	var src = selectSrc.Selected[0:3]
	var dst = selectDst.Selected[0:3]
	rate, _ := cRates.Convert(src, dst)
	resultItem.SetText(
		fmt.Sprintf(
			"1 %s = %f %s\n1 %s = %f %s\n",
			dst, rate, src,
			src, 1/rate, dst,
		),
	)
}

func main() {
	os.Setenv("FYNE_THEME", "light")
	a = app.NewWithID("cbrRates")
	a.Settings().SetTheme(apptheme.New())
	w = a.NewWindow("CBR Currency Converter")
	w.Resize(fyne.NewSize(480, 320))
	content := container.NewVBox()
	//content.Add(widget.NewLabel("Currency converter based on cbr.ru API"))
	appHeader := widget.NewLabel(
		"Currency converter based on data from\nCentral Bank of Russia",
	)
	appHeader.Alignment = fyne.TextAlignCenter

	content.Add(appHeader)

	selectSrc = widget.NewSelect([]string{}, runConvert)
	selectDst = widget.NewSelect([]string{}, runConvert)
	resultItem = widget.NewEntry()
	resultItem.Disable()
	resultItem.MultiLine = true
	content.Add(selectSrc)
	content.Add(selectDst)
	content.Add(resultItem)

	w.SetContent(content)

	var err error
	cRates, err = currency.New("")
	if err != nil {
		fmt.Println(err)
		a.Quit()
	}
	var currencyTitles []string
	for _, v := range cRates.Rates {
		currencyTitles = append(
			currencyTitles,
			fmt.Sprintf("%s (%s)", v.GetISOCode(), v.GetName()))
	}
	for i := 1; i < len(currencyTitles); i++ {
		if currencyTitles[i] < currencyTitles[i-1] {
			for j := i; j > 0 && currencyTitles[j] < currencyTitles[j-1]; j-- {
				currencyTitles[j], currencyTitles[j-1] = currencyTitles[j-1], currencyTitles[j]
			}
		}
	}
	selectSrc.Options = currencyTitles
	selectDst.Options = currencyTitles

	w.ShowAndRun()
}

/*
func main() {
	data, err := currency.Get("30/05/2023")

	if err != nil {
		fmt.Println(err)
		return
	} else {
	}

	//fmt.Println(data)

	fmt.Printf("Exchange rates on %s:\n", data.Date)
	for _, v := range data.Rates {
		fmt.Printf(
			"% 6d %s = %f RUB\n",
			v.GetNominal(),
			v.GetISOCode(),
			v.GetValue(),
		)
	}

	var currencies []string = []string{
		"USD",
		"EUR",
		"AZN",
		"BYN",
		"KZT",
		"RUB",
		"GBP",
		"VND",
	}

	for _, src := range currencies {
		for _, dst := range currencies {
			r, e := data.Convert(src, dst)
			if e != nil {
				fmt.Println(e)
			} else {
				fmt.Printf("%s to %s = %f\n", src, dst, r)
			}
		}
	}

}
*/
