# Тестовое задание: Реализовать API для небольшого интернет магазина.

## Реализовано на базе:
*   Golang & PostgreSQL.
*   Основная логика запросов вынесена в функции (stored procedures) PostgreSQL. Функции сохранены в [бекапе](https://github.com/EwRvp7LV7/48170360shop/configs/sql) БД. 
*   Провайдер БД [jmoiron/sqlx](https://github.com/jmoiron/sqlx).
*   В качестве ключа строк корзины взят тип UUID.
*   Основой для API взят [go-chi/chi](https://github.com/go-chi/chi). Выбран как [быстрый по тестам](https://benhoyt.com/writings/go-routing#benchmarks).
*   Аутенификация токиеном на базе [go-chi/jwtauth](https://github.com/go-chi/jwtauth). Не лучший вариант, достаточно сырая.
*   Валидация при помощи [go-ozzo/ozzo-validation](https://github.com/go-ozzo/ozzo-validation) реализована через методы (не через теги структуры), что снижает возможность ошибки разработчика.
*   API протестировано Postmen. Исходники Postmen находятся в [test/postman](https://github.com/EwRvp7LV7/48170360shop/test/postman).
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
*   "Покупать товар" то есть уменьшать количество товара на складе на количество товара в корзине с одновременным ее удалением.
### Для менеджера
*   Смотреть все корзины покупателей.
*   Добавлять в список товаров новый товар. При этом название товара является уникальным, нельзя добавить товар с тем же названием.
*   Изменять количество товаров на складе (store) на положительную или отрицательную величину. БД возвращает текущее значение. При попытке получить отрицательное количество возвращает ошибку.
*   Роли покупателя и менеджера ограничены, при попытке покупателя выполнить операции менеджера или наоборот БД вернет ошибку.





