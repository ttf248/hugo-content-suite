package stats

import (
	"hugo-content-suite/models"
)

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
