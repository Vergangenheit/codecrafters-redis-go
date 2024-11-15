package main

import "time"

func expired(res *Resource, currentTime time.Time) bool {
	if res.expired == nil {
		return false
	}
	expiredTs := *res.expired
	if expiredTs.Before(currentTime) {
		return true
	}
	return false
}
