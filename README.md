# avito

Решение тестового задания (формулировка задания представлена ниже) на позицию стажера-бэкендера в Авито.


Нужно реализовать HTTP API, которое позволяет взаимодействовать с балансом пользователя: метод получения баланса пользователя `GetBalance()`, метод пополнения и снятия средств со счета пользователя `RefillAndWithdrawMoney()` и метод перевода средств со счета одного пользователя на счет другого `TransferMoney()`.

Чтобы собрать проект, нужно в корневой директории проекта, там, где находится `docker-compose.yml` файл, выполнить следующую команду:
```
docker-compose up -d --build
```
или
```
make up
```

В результате будут созданы три docker-конейнера - первый, `pg` - контейнер с базой данных, второй, `pgadmin` - контейнер с интерактивной платформой для управления базами данных и третий, `app` - само приложение.


Следующая мини-документация объясняет, как использовать HTTP API.

1. Метод GetBalance():
* Входные данные:
  * Content-Type: application/json
  * request body: {"id":id}
  * id - уникальный идентификатор пользователя (число), id > 0
* Выходные данные:
  * Content-Type: application/json
  * response body: {"status":status,"id":id,"balance":balance}
  * status - статус ответа сервера (число)
  * id - идентификатор пользователя (число)
  * balance - баланс пользователя (число, максимум два знака после запятой)
---
2. Метод RefillAndWithdrawMoney():
* Входные данные:
  * Content-Type: application/json
  * request body: {"id":id,"sum":sum}
  * id - уникальный идентификатор пользователя (число), id > 0
  * sum - сумма средств для пополнения или снятия со счета пользователя, sum != 0
* Выходные данные:
  * Content-Type: application/json
  * response body: {"status":status,"id":id,"balance":balance}
  * status - статус ответа сервера (число)
  * id - идентификатор пользователя (число)
  * balance - баланс пользователя (число, максимум два знака после запятой)
---
3. Метод TransferMoney():
* Входные данные:
  * Content-Type: application/json
  * request body: {"from":from,"to":to,"sum":sum}
  * from - уникальный идентификатор пользователя, который переводит деньги на счет пользователя to, from > 0
  * to - уникальный идентификатор пользователя, которому переводит деньги пользователь from, to > 0
  * sum - сумма средств, которая переводится на счет пользователя to, sum > 0
* Выходные данные:
  * Content-Type: application/json
  * response body: {"status":status,"from_id":from_id,"from_balance":from_balance,"to_id":to_id,"to_balance":to_balance}
  * status - статус ответа сервера (число)
  * from_id, to_id - уникальные идентификаторы пользователей (число)
  * from_balance, to_balance - балансы пользователей (число, максимум два знака после запятой)

Статусы ошибок:
1. В случае успеха:
  * status = 0, id > 0, balance >= 0.00
3. В случае фэйла:
* Невалидные данные (вместо числа пришла строка, идентификатор пользователя отрицательный, в методе TransferMoney() при отрицательном значении sum, невалидный Content-Type):
  * status = 1, id = 0, balance = 0.00
* Недостаточно средств (в RefillAndWithdrawMoney() при снятии средств со счета пользователя, в TransferMoney() - при переводе средств с одного счета на другой):
  * status = 2, id = 0, balance = 0.00
* Несуществующий идентификатор пользователя (во всех методах):
  * status = 3, id = 0, balance = 0.00
* Ошибка сервера (во всех методах):
  * status = 4, id = 0, balance = 0.00


### Тестирование


1. Тестирование сервера с помощью тестов:
Запустить следующую команду:
```
DATABASE_URL=postgres://user:password@localhost:5432/database?sslmode=disable go test ./... -v -cover
```

2. Тестирование сервера через терминал:
* метод GetBalance():
```
curl -v --request GET --header "Content-Type: application/json" --data '{"id":1}' localhost:8080/balance
```
* метод RefillAndWithdrawMoney():
* Пополнение баланса пользователя:
```
curl -v --request POST --header "Content-Type: application/json" --data '{"id":2,"sum":5}' localhost:8080/refill
```
* Снятие средств со счета пользователя:
```
curl -v --request POST --header "Content-Type: application/json" --data '{"id":2,"sum":-5}' localhost:8080/withdraw
```

* метод TransferMoney():
* Перевод средств с одного счета на другой:
```
curl -v --request POST --header "Content-Type: application/json" --data '{"from":2,"to":3,"sum":5}' localhost:8080/transfer
```




> # Тестовое задание на позицию стажера-бекендера
>
> ## Микросервис для работы с балансом пользователей.
>
> **Проблема:**
>
> В нашей компании есть много различных микросервисов. Многие из них так или иначе хотят взаимодействовать с балансом пользователя. На архитектурном комитете приняли решение централизовать работу с балансом пользователя в отдельный сервис. 
>
> **Задача:**
> 
> Необходимо реализовать микросервис для работы с балансом пользователей (зачисление средств, списание средств, перевод средств от пользователя к пользователю, а также метод получения баланса пользователя). Сервис должен предоставлять HTTP API и принимать/отдавать запросы/ответы в формате JSON. 
> 
> **Сценарии использования:**
>
> Далее описаны несколько упрощенных кейсов приближенных к реальности.
> 1. Сервис биллинга с помощью внешних мерчантов (аля через visa/mastercard) обработал зачисление денег на наш счет. Теперь биллингу нужно добавить эти деньги на баланс пользователя. 
> 2. Пользователь хочет купить у нас какую-то услугу. Для этого у нас есть специальный сервис управления услугами, который перед применением услуги проверяет баланс и потом списывает необходимую сумму. 
> 3. В ближайшем будущем планируется дать пользователям возможность перечислять деньги друг-другу внутри нашей платформы. Мы решили заранее предусмотреть такую возможность и заложить ее в архитектуру нашего сервиса. 
>
> **Требования к коду:**
>
> 1. Язык разработки: Go. Мы готовы рассматривать решения на PHP/Python/другом языке, но приоритетом для нас является именно golang.
> 2. Фреймворки и библиотеки можно использовать любые
> 3. Реляционная СУБД: MySQL или PostgreSQL
> 4. Весь код должен быть выложен на Github с Readme файлом с инструкцией по запуску и примерами запросов/ответов (можно просто описать в Readme методы, можно через Postman, можно в Readme curl запросы скопировать, вы поняли идею...)
> 5. Если есть потребность, можно подключить кеши(Redis) и/или очереди(RabbitMQ, Kafka)
> 6. При возникновении вопросов по ТЗ оставляем принятие решения за кандидатом (в таком случае Readme файле к проекту должен быть указан список вопросов с которыми кандидат столкнулся и каким образом он их решил)
> 7. Разработка интерфейса в браузере НЕ ТРЕБУЕТСЯ. Взаимодействие с АПИ предполагается посредством запросов из кода другого сервиса. Для тестирования можно использовать любой удобный инструмент. Например: в терминале через curl или Postman.
>
> **Будет плюсом:**
>
> 1. Использование docker и docker-compose для поднятия и развертывания dev-среды.
> 2. Методы АПИ возвращают человеко-читабельные описания ошибок и соответвующие статус коды при их возникновении.
> 3. Все реализовано на GO, все-же мы собеседуем на GO разработчика. HINT: На собеседовании так или иначе будут вопросы по Go. Кто прочитал, тот молодец :)
> 4. Написаны unit/интеграционные тесты.
>
> **Основное задание (минимум):**
>
> Метод начисления средств на баланс. Принимает id пользователя и сколько средств зачислить.
> 
> Метод списания средств с баланса. Принимает id пользователя и сколько средств списать. 
>
> Метод перевода средств от пользователя к пользователю. Принимает id пользователя с которого нужно списать средства, id пользователя которому должны зачислить средства, а также сумму.
>
> Метод получения текущего баланса пользователя. Принимает id пользователя. Баланс всегда в рублях.
>
> **Детали по заданию:**
>
> 1. Методы начисления и списания можно объединить в один, если это позволяет общая архитектура.
> 2. По умолчанию сервис не содержит в себе никаких данных о балансах (пустая табличка в БД). Данные о балансе появляются при первом зачислении денег. 
> 3. Валидацию данных и обработку ошибок оставляем на усмотрение кандидата. 
> 4. Список полей к методам не фиксированный. Перечислен лишь необходимый минимум. В рамках выполнения доп. заданий возможны дополнительные поля.
> 5. Механизм миграции не нужен. Достаточно предоставить конечный SQL файл с созданием всех необходимых таблиц в БД. 
> 6. Баланс пользователя - очень важные данные в которых недопустимы ошибки (фактически мы работаем тут с реальными деньгами). Необходимо всегда держать баланс в актуальном состоянии и не допускать ситуаций когда баланс может уйти в минус. 
> 7. Валюта баланса по умолчанию всегда рубли.
>
> **Дополнительные задания**
>
> Далее перечислены доп. задания. Они не являются обязательными, но их выполнение даст существенный плюс перед другими кандидатами. 
>
> *Доп. задание 1:*
>
> Эффективные менеджеры захотели добавить в наши приложения товары и услуги в различных от рубля валютах. Необходима возможность вывода баланса пользователя в отличной от рубля валюте.
>
> Задача: добавить к методу получения баланса доп. параметр. Пример: ?currency=USD. 
> Если этот параметр присутствует, то мы должны конвертировать баланс пользователя с рубля на указанную валюту. Данные по текущему курсу валют можно взять отсюда https://exchangeratesapi.io/ или из любого другого открытого источника. 
> 
> Примечание: напоминаем, что базовая валюта которая хранится на балансе у нас всегда рубль. В рамках этой задачи конвертация всегда происходит с базовой валюты.
> 
> *Доп. задание 2:*
> 
> Пользователи жалуются, что не понимают за что были списаны (или зачислены) средства. 
> 
> Задача: необходимо предоставить метод получения списка транзакций с комментариями откуда и зачем были начислены/списаны средства с баланса. Необходимо предусмотреть пагинацию и сортировку по сумме и дате.
