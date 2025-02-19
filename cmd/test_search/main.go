package main

import (
	"fmt"

	"github.com/rubikge/lemmatizer/internal/models"
	"github.com/rubikge/lemmatizer/internal/services"
)

func main() {
	words := []string{
		"учтель", "преподавать",
	}

	searchData := models.SearchData{
		ProductTitle: "Школа",
		Keywords: []models.WeightWord{
			{Word: "учитель", Weight: 1.8}, {Word: "ученик", Weight: 1.7}, {Word: "книга", Weight: 1.4}, {Word: "учебник", Weight: 1.6}, {Word: "класс", Weight: 1.5},
			{Word: "преподавать", Weight: 1.9}, {Word: "изучение", Weight: 1.3}, {Word: "образование", Weight: 2.0}, {Word: "урок", Weight: 1.5}, {Word: "дисциплина", Weight: 1.4},
			{Word: "обучение", Weight: 1.7}, {Word: "учебный", Weight: 1.6}, {Word: "курс", Weight: 1.8}, {Word: "предмет", Weight: 1.2}, {Word: "объяснение", Weight: 1.3},
			{Word: "активность", Weight: 1.3}, {Word: "педагог", Weight: 1.8}, {Word: "план", Weight: 1.7}, {Word: "тема", Weight: 1.6}, {Word: "метод", Weight: 1.5},
			{Word: "группа", Weight: 1.4}, {Word: "занятие", Weight: 1.7}, {Word: "задание", Weight: 1.2}, {Word: "диалог", Weight: 1.4}, {Word: "сессия", Weight: 1.5},
			{Word: "речь", Weight: 1.3}, {Word: "экзамен", Weight: 1.6}, {Word: "дидактика", Weight: 1.5}, {Word: "методика", Weight: 1.7}, {Word: "акцент", Weight: 1.8},
			{Word: "педагогика", Weight: 1.6}, {Word: "контроль", Weight: 1.7}, {Word: "обсуждение", Weight: 1.8}, {Word: "психология", Weight: 1.6}, {Word: "навык", Weight: 1.4},
			{Word: "профессионал", Weight: 1.5}, {Word: "практика", Weight: 1.7}, {Word: "построение", Weight: 1.9}, {Word: "рецензия", Weight: 1.4}, {Word: "оценка", Weight: 1.5},
			{Word: "работа", Weight: 1.3}, {Word: "платформа", Weight: 1.6}, {Word: "студент", Weight: 1.7}, {Word: "тест", Weight: 1.8}, {Word: "консультация", Weight: 1.9},
			{Word: "письменность", Weight: 1.5},
		},
	}

	fmt.Println(services.GetTotalScore(words, &searchData))
}
