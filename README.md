## GoMockServer

### 1. Создайте конфигурационный файл

Создайте файл `conf.yaml` в папке вашего проекта:

```yaml
port: 9876
log_level: debug
log_format: text
endpoints:
  - type: "OK"
    method: "GET"
    status: 200
    path: "/scripts/XML_daily.asp?date_req=02/03/2002"
    response_path: "./example/ok.xml"

  - type: "internal"
    method: "GET"
    status: 500
    path: "/scripts/XML_daily.asp?date_req=03/03/2002"
    response_path: "./example/error.xml"
```

### 2. Создайте response файлы

Создайте файлы указанные в yaml файле

**ok.xml**
```xml
<ValCurs Date="02.03.2002" name="Foreign Currency Market">
<Valute ID="R01010">
<NumCode>036</NumCode>
<CharCode>AUD</CharCode>
<Nominal>1</Nominal>
<Name>Австралийский доллар</Name>
<Value>16,0102</Value>
<VunitRate>16,0102</VunitRate>
</Valute>
<Valute ID="R01035">
<NumCode>826</NumCode>
<CharCode>GBP</CharCode>
<Nominal>1</Nominal>
<Name>Фунт стерлингов</Name>
<Value>43,8254</Value>
<VunitRate>43,8254</VunitRate>
</Valute>
```

**error.xml**
```xml
<ValCurs>Internal server error</ValCurs>
```

### 3. Запустите сервер

```bash
go run .
```

или с своей конфигурацией:

```bash
go run . --config=./custom-config.yaml
```

### Server Configuration

| Опция | Описание |
|--------|-------------|
| `port` | порт сервера | 
| `log_level` | уровень логгирования (debug, info) | 
| `log_format` | формат логов (text, json) | 
| `endpoints` | массив конфигурацией эндпоинтов | 

| Опция | Описание |
|--------|-------------|
| `method` | HTTP метод (GET, POST, PUT, DELETE, etc.) | 
| `status` | HTTP response status code | 
| `path` | путь до эндпоинта | 
| `response_path` | путь до response файла | 
