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

	cRates = currency.New()
	cRates.SetCache(GetCurrencyCache(a.Preferences()))
	err := cRates.Fetch("")
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
	selectSrc.Options = currencyTitles
	selectDst.Options = currencyTitles

	w.ShowAndRun()
}
