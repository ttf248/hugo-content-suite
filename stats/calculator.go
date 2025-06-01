package stats

import (
	"hugo-content-suite/models"
	"sort"
)

func CalculateCategoryStats(articles []models.Article) []models.CategoryStats {
	categoryMap := make(map[string]int)

	for _, article := range articles {
		category := article.Category
		if category == "" {
			category = "无分类"
		}
		categoryMap[category]++
	}

	var stats []models.CategoryStats
	for name, count := range categoryMap {
		stats = append(stats, models.CategoryStats{
			Name:  name,
			Count: count,
		})
	}

	sort.Slice(stats, func(i, j int) bool {
		return stats[i].Count > stats[j].Count
	})

	return stats
}

func FindNoTagArticles(articles []models.Article) []models.Article {
	var noTagArticles []models.Article

	for _, article := range articles {
		if len(article.Tags) == 0 {
			noTagArticles = append(noTagArticles, article)
		}
	}

	return noTagArticles
}

func GroupTagsByFrequency(tagStats []models.TagStats) (high, medium, low []models.TagStats) {
	for _, stat := range tagStats {
		if stat.Count >= 5 {
			high = append(high, stat)
		} else if stat.Count >= 2 {
			medium = append(medium, stat)
		} else {
			low = append(low, stat)
		}
	}
	return
}
