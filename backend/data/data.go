package data

import (
	"strings"

	"razorsh4rk.github.io/jobsite/db"
)

func GetSkillHeatmap(listings []db.InterpretedListing) map[string]int {
	hMap := make(map[string]int)

	for _, l := range listings {
		skills := l.Techstack

		for _, s := range skills {
			key := strings.ToLower(s)
			hMap[key]++
		}
	}

	return hMap
}
