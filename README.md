# Тестовое задание: Реализовать API для небольшого интернет магазина.

## Реализовано на базе:
*   Golang & PostgreSQL.
*   Структура модулей (папок) согласно [golang-standards/project-layout](https://github.com/golang-standards/project-layout).
*   Основная логика запросов вынесена в функции (stored procedures) PostgreSQL. Функции сохранены в [бекапе](https://github.com/EwRvp7LV7/48170360shop/configs/sql) БД. 
*   Провайдер БД [jmoiron/sqlx](https://github.com/jmoiron/sqlx).
*   В качестве ключа строк корзины взят тип UUID.
*   Основой для API взят [go-chi/chi](https://github.com/go-chi/chi). Выбран как [быстрый по тестам](https://benhoyt.com/writings/go-routing#benchmarks).
*   Аутенификация токиеном на базе [go-chi/jwtauth](https://github.com/go-chi/jwtauth). Не лучший вариант, достаточно сырая.
*   Валидация при помощи [go-ozzo/ozzo-validation](https://github.com/go-ozzo/ozzo-validation) реализована через методы (не через теги структуры), что снижает возможность ошибки разработчика.
*   API протестировано Postman. Исходники Postman находятся в [test/postman](https://github.com/EwRvp7LV7/48170360shop/test/postman).
*   Начато прописывание комментариев для автоматической генерации swagger.json (находится в папке [docs](https://github.com/EwRvp7LV7/48170360shop/docs)). 

## Возможности API:
Реализовано на связке REST+JSON.
### Для покупателя
*   Смотреть список товаров и цены.
*   Добавлять и убавлять товары в корзине. При вводе отрицательного значения количество товара уменьшается на эту величину. При нулевом количестве товар удаляется. При попытке получить отрицательное количество БД возвращает ошибку.
*   При этом покупателю возвращается итоговый список товаров в его корзине с промежуточной и общей суммой.
```json
[
    {
        "name": "апельсины",
        "price": "170.40 руб",
        "goods_q": 6,
        "sum": "1 022.40 руб"
    },
    {
        "name": "картофель",
        "price": "45.34 руб",
        "goods_q": 2,
        "sum": "90.68 руб"
    },
    {
        "name": "Сумма:",
        "price": "0.00 руб",
        "goods_q": 0,
        "sum": "1 113.08 руб"
    }
]
```
*   *Покупать товар*, то есть уменьшать количество товара на складе на количество товара в корзине с одновременным ее удалением.
### Для менеджера
*   Смотреть все корзины покупателей.
*   Добавлять в список товаров новый товар. При этом название товара является уникальным, нельзя добавить товар с тем же названием.
*   Изменять количество товаров на складе (store) на положительную или отрицательную величину. В ответ БД возвращает текущее количество. При попытке получить отрицательное количество возвращает ошибку.
--------------
Роли покупателя и менеджера ограничены, при попытке покупателя выполнить операции менеджера или наоборот БД вернет ошибку:
```sql
   IF NOT EXISTS (SELECT * FROM b_users JOIN a_user_type USING(type_id) WHERE b_users.account = user_name AND a_user_type.type = 'manager')
   THEN RAISE EXCEPTION 'Denied user --> %', user_name USING HINT = 'Please check your rights';
   END IF;
```

## Сборка и запуск
Для 
```bash
$ psql -V
(PostgreSQL) 13.3
```
1. Создать в PostgreSQL БД smallshop и восстановить в нее backup
```bash
$ psql -h localhost -U postgres -d smallshop -f configs/sql/211024backup.sql
```
2. Добавить параметры соединения с вашей БД в configs/config.toml
``` 
[database]
Host     ="localhost"
Port     ="5432"
User     ="postgres"
Password ="pass"
NameDB   ="smallshop"
``` 
3. 
```bash
$ go run cmd/main.go
```
4. Postman -> Import (Ctrl+O) -> test/postmanshop.postman_collection.json

-------------------------------------------------------------------------------------------------------------
|      API Postman       | description                                                                      |
| --------------------- | --------------------------------------------------------------------------------|
| [user auth]            | Авторизация user (в БД есть user1, user2, user3 с одинаковым паролем userpass)   |
| [user get googs list]  | Возвращает список товаров и цен из таблицы c_goods                               |
| [user put into basket] | Добавляет/убавляет товары в корзине, возвращает содержимое корзины (см. выше) <br /> "goods_add": "1" добавляет к количеству +1 товар,  "goods_add": "5" добавляет 5,  <br />"goods_add": "-1" убавляет -1 товар, и т.д. |                                          
| [user buy]             | Сокращает склад (c_goods.store) на количество товара в корзине и удаляет корзину |
| [user auth check]      | Проверка авторизации пользователя, возвращает user_name                          |
| [user auth logout]     | Logout текущего аккаунта                                           |
|     | **Роль менеджера**                                        |
| [manager auth]         | Авторизация user (в БД есть manager1 с паролем manager1pass)                     |
| [manager get baskets]  | Просмотр всех корзин (содержание таблицы d_basket сгруппированое по user_id)     |
| [manager new goods]    | Добавление нового товара (название, количество, цена) в таблицу c_goods              |
| [manager add to store] | Изменение остатков на складе (c_goods.store), <br />Например, "goods_add": "-6" уменьшает количество на 6  |
-------------------------------------------------------------------------------------------------------------





