# rest-jwt


## Требования

- Go 1.18+
- PostgreSQL

## Установка

1. Клонируйте репозиторий:

```bash
git clone https://github.com/MaximInnopolis/rest-jwt.git
cd rest-jwt
```

2. Соберите докер-билд:
```bash
make up-all
```

3. Проведите миграцию:
```bash
make migrate
```

4. Запустите тесты:
```bash
make test
```