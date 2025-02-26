# lemmatizer

Run service: 
```
go run ./cmd/lemmatizer/.
```

Test search: 
```
go run ./cmd/test_search/.
```

Docker:
```
docker build -t search-app .
```

Варианты лемматизатора русского языка:

MyStem от Яндекса
(Либо через CLI, либо   pymystem3, и наверное так- но с 2019 не обновлялся)
pymystem3 - wrapper

ChatGPT предложил
spaCy + pymorphy2 - тоже 5 лет назад последнее обновление
(github.com/dveselov/morph pymorphy2 на Go)

github.com/amarin/gomorph - еще вариант на Go

Natasha - на питоне, 4 месяца назад последняя контрибуция
pymorphy3 -https://github.com/no-plagiarism/pymorphy3

UDPipe - не русский

https://habr.com/ru/articles/503420/ - статья 2020 года про сравнение pymorphy2 и pymystem3. Mystem эффективен на больших текстах, так как долго запускается exe. Может, в параллельных потоках запускать, может формировать большой текст перед обработкой.

https://habr.com/ru/companies/otus/articles/808435/ - статаья 2024 года, 5 NLP инструментов на python.
из нового DeepPavlov - но там про создание ai ассистентов больше, для лемматизации он pymorphy2 вероятно использует
