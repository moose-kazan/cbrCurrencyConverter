package main

import (
	"cbr/internal/currency"
	"fmt"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
)

const cache_ttl = 7200

func GetCurrencyCache(pref fyne.Preferences) *currency.CurrencyCache {

	var cCache currency.CurrencyCache = currency.CurrencyCache{
		CacheGet: func(date string) *string {
			var cache_key_value string = fmt.Sprintf("cache_value_%s", date)
			var cache_key_time string = fmt.Sprintf("cache_time_%s", date)
			var cacheTimeStr string = pref.String(cache_key_time)
			var cacheTIme int64
			var err error
			cacheTIme, err = strconv.ParseInt(cacheTimeStr, 10, 64)
			if err != nil {
				return nil
			}
			if cacheTIme+cache_ttl < time.Now().Unix() {
				// Cleanup
				pref.RemoveValue(cache_key_time)
				pref.RemoveValue(cache_key_value)
				return nil
			}

			var rv = pref.String(cache_key_value)
			if rv == "" {
				return nil
			}
			return &rv
		},
		CacheSet: func(date string, data string) {
			var cache_key_value string = fmt.Sprintf("cache_value_%s", date)
			var cache_key_time string = fmt.Sprintf("cache_time_%s", date)
			pref.SetString(cache_key_time, fmt.Sprintf("%d", time.Now().Unix()))
			pref.SetString(cache_key_value, data)
		},
	}
	return &cCache
}
