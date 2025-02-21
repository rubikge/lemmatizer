package main

import (
	"encoding/json"
	"fmt"

	"github.com/rubikge/lemmatizer/internal/models"
	"github.com/rubikge/lemmatizer/internal/repository"
	"github.com/rubikge/lemmatizer/internal/services"
	"github.com/rubikge/lemmatizer/internal/utils"
)

func main() {
	jsonData := []byte(`
	{
		"common_positive_keywords": [
			"ищу",
			"нужен",
			"требуется",
			"посоветуйте",
			"помогите",
			"где",
			"как",
			"сколько",
			"цена",
			"стоимость",
			"купить",
			"заказать",
			"ремонт",
			"сервис",
			"мастер",
			"специалист",
			"диагностика",
			"восстановление",
			"наладка",
			"чистка",
			"замена",
			"запчасти",
			"оригинал",
			"качественный",
			"профессиональный",
			"срочный",
			"быстрый",
			"выезд",
			"гарантия"
		],
		"common_negative_keywords": [
			"предлагаю",
			"продаю",
			"услуги",
			"сервис",
			"ремонт",
			"диагностика",
			"восстановление",
			"наладка",
			"чистка",
			"замена",
			"запчасти",
			"оригинал",
			"качественный",
			"профессиональный",
			"срочный",
			"быстрый",
			"выезд",
			"гарантия",
			"недорого",
			"скидки",
			"акция"
		],
		"products": [
			{
				"product_title": "Ремонт iPhone",
				"required_keywords": [
					"ремонт",
					"iphone"
				],
				"min_count_words": 3,
				"waight": 0.7,
				"keywords_with_synonyms": [
					{
						"sinonyms": [
							"ремонт",
							"восстановление",
							"починка",
							"исправление",
							"устранение",
							"наладка",
							"диагностика",
							"чистка",
							"замена",
							"запчасти",
							"комплектующие"
						],
						"waight": 0.6
					},
					{
						"sinonyms": [
							"iphone",
							"айфон",
							"айпад",
							"айпод",
							"apple",
							"эпл",
							"телефон",
							"смартфон"
						],
						"waight": 0.6
					},
					{
						"sinonyms": [
							"разбитый",
							"сломанный",
							"неисправный",
							"поврежденный",
							"завис",
							"не включается",
							"не работает",
							"глючит"
						],
						"waight": 0.5
					},
									{
						"sinonyms": [
							"экран",
							"дисплей",
							"стекло",
							"тачскрин",
							"аккумулятор",
							"батарея",
							"кнопка",
							"разъем",
							"камера",
							"микрофон",
							"динамик"
						],
						"waight": 0.4
					}
				]
			},
			{
				"product_title": "Ремонт iPad",
				"required_keywords": [
					"ремонт",
					"ipad"
				],
				"min_count_words": 3,
				"waight": 0.6,
				"keywords_with_synonyms": [
					{
						"sinonyms": [
							"ремонт",
							"восстановление",
							"починка",
							"исправление",
							"устранение",
							"наладка",
							"диагностика",
							"чистка",
							"замена",
							"запчасти",
							"комплектующие"
						],
						"waight": 0.6
					},
					{
						"sinonyms": [
							"ipad",
							"айпад",
							"apple",
							"эпл",
							"планшет"
						],
						"waight": 0.6
					},
					{
						"sinonyms": [
							"разбитый",
							"сломанный",
							"неисправный",
							"поврежденный",
							"завис",
							"не включается",
							"не работает",
							"глючит"
						],
						"waight": 0.5
					},
									{
						"sinonyms": [
							"экран",
							"дисплей",
							"стекло",
							"тачскрин",
							"аккумулятор",
							"батарея",
							"кнопка",
							"разъем",
							"камера",
							"микрофон",
							"динамик"
						],
						"waight": 0.4
					}
				]
			},
			{
				"product_title": "Ремонт MacBook",
				"required_keywords": [
					"ремонт",
					"macbook"
				],
				"min_count_words": 3,
				"waight": 0.5,
				"keywords_with_synonyms": [
					{
						"sinonyms": [
							"ремонт",
							"восстановление",
							"починка",
							"исправление",
							"устранение",
							"наладка",
							"диагностика",
							"чистка",
							"замена",
							"запчасти",
							"комплектующие"
						],
						"waight": 0.6
					},
					{
						"sinonyms": [
							"macbook",
							"макбук",
							"apple",
							"эпл",
							"ноутбук"
						],
						"waight": 0.6
					},
					{
						"sinonyms": [
							"разбитый",
							"сломанный",
							"неисправный",
							"поврежденный",
							"завис",
							"не включается",
							"не работает",
							"глючит"
						],
						"waight": 0.5
					},
									{
						"sinonyms": [
							"экран",
							"дисплей",
							"стекло",
							"тачскрин",
							"аккумулятор",
							"батарея",
							"кнопка",
							"разъем",
							"клавиатура",
							"трекпад"
						],
						"waight": 0.4
					}
				]
			}
		]
	}
`)
	text := "Подскажите где эппл сервес в тбилиси"

	var product models.Product
	if err := json.Unmarshal(jsonData, &product); err != nil {
		fmt.Println(err)
	}

	r := repository.NewMystemRepository()
	s := services.NewLemmatizerService(r)

	lemmas, err := s.GetLemmasArray(text)
	if err != nil {
		fmt.Println(err)
	}

	words := []string{}
	for _, lemma := range lemmas {
		word := lemma.Lemma
		if word != "" {
			words = append(words, lemma.Lemma)
		}
	}

	searchProducts, err := utils.GetLemmatizedSearchProduct(&product, s)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Result ", services.GetScore(&words, &searchProducts))

}
