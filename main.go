package main

import (
	"github.com/gin-gonic/gin"
)

type Answer struct {
	Value    int    `json:"value"`
	Category string `json:"category"`
}

type QuestionDTO struct {
	Answers []Answer `json:"answers"`
}

func calculateTotalScore(questionDTO QuestionDTO) float64 {
	totalSum := 0

	for _, answer := range questionDTO.Answers {
		totalSum += answer.Value
	}

	totalScore := float64(totalSum) / float64(len(questionDTO.Answers))
	return totalScore
}

func calculateScoreByCategory(questionDTO QuestionDTO) map[string]float64 {
	categoryScores := make(map[string]struct{ points, count int })
	averages := make(map[string]float64)

	for _, answer := range questionDTO.Answers {
		categoryScores[answer.Category] = struct{ points, count int }{
			points: categoryScores[answer.Category].points + answer.Value,
			count:  categoryScores[answer.Category].count + 1,
		}
	}

	for category, score := range categoryScores {
		averages[category] = float64(score.points) / float64(score.count)
	}

	return averages
}

func main() {
	router := gin.Default()

	router.POST("/calculate", func(c *gin.Context) {
		var questionDTO QuestionDTO

		if err := c.ShouldBindJSON(&questionDTO); err != nil {
			c.JSON(400, gin.H{"error": "Invalid Request"})
			return
		}

		score := calculateTotalScore(questionDTO)
		averages := calculateScoreByCategory(questionDTO)

		c.JSON(200, gin.H{
			"score":    score,
			"averages": averages,
		})
	})

	router.Run(":8080")
}
