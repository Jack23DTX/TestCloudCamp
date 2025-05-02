# Simple Go Load Balancer

Небольшой HTTP-балансировщик на Go с round-robin и rate-limiting по IP.

## Требования
- Go 1.23+
- Docker

## Быстрый старт

1. Клонировать репозиторий
   ```bash
   git clone https://github.com/your-user/TestCloudCamp.git
   cd TestCloudCamp
   ```


2. Создать config.yaml рядом с бинарником:

``` yaml
port: ":8080"
backends:
  - "http://host.docker.internal:8081"
  - "http://host.docker.internal:8082"
rate_limit:
  rps: 5
  burst: 10
  ```

3. Собрать и запустить:

- Локально (если запускаете локально нужно изменить настройки config.yaml на "http://localhost:8081" и тд.):
```bash
go build -o proxy-server ./cmd/main.go
./proxy-server
```
- В Docker:

```bash
docker build -t proxy-server .
docker run -p 8080:8080 proxy-server
```

### Тесты

```bash
go test ./... -v
go test ./... -race -bench=.
```

Пример нагрузки

```bash
ab -n 1000 -c 100 http://localhost:8080/
```
