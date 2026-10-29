# Закреплённые версии инструментов

Проверено: 2026-10-29. Версии взяты из официальных источников (go.dev, GitHub Releases,
proxy.golang.org); совместимость проверена сборкой `cmd/web` и `cmd/worker`.

## Toolchain

| Инструмент | Версия | Источник проверки |
| --- | --- | --- |
| Go (stable) | 1.27.3 | https://go.dev/dl/?mode=json |
| Go (закреплён в `go.mod`) | 1.27.3 | директива `go` |
| templ | v0.3.1100 | https://github.com/a-h/templ/releases/latest |
| sqlc | v1.31.1 | https://github.com/sqlc-dev/sqlc/releases/latest |
| golang-migrate | v4.20.1 | https://github.com/golang-migrate/migrate/releases/latest |
| golangci-lint | v2.15.0 | https://github.com/golangci/golangci-lint/releases/latest |
| govulncheck | v1.9.0 | https://proxy.golang.org/golang.org/x/vuln/@latest |

`templ` и `sqlc` закреплены директивами `tool` в `go.mod` и запускаются как
`go tool templ` / `go tool sqlc`, поэтому версия не зависит от глобально
установленных бинарников. Библиотеки `River` и `golang-migrate` закреплены
импортами с build-тегом `tools` в `internal/platform/tools/tools.go`; они
подключаются к рабочим пакетам в срезах L02/L06. `golangci-lint` и
`govulncheck` фиксируются этой таблицей и будут установлены в CI-срезе (L06).

## Библиотеки

| Библиотека | Версия | Требование Go из `go.mod` |
| --- | --- | --- |
| Gin | v1.12.0 | 1.25.0 |
| Templ | v0.3.1100 | 1.26.0 |
| pgx | v5.11.0 | 1.25.0 |
| River | v0.50.0 | 1.26.0 |
| sqlc | v1.31.1 | 1.26.0 |
| golang-migrate | v4.20.1 | проверено сборкой |
| testcontainers-go | v0.44.0 | проверено сборкой |
| Prometheus client_golang | v1.25.0 | проверено сборкой |
| yaml.v3 | v3.0.1 | проверено сборкой |
| x/sync | v0.24.0 | проверено сборкой |

Все перечисленные модули требуют Go не выше закреплённого 1.27.3; конфликтов
toolchain нет. HTMX, Pico.css и Vue/Vite фиксируются в своих срезах (L03/L05).

## Проверка

```bash
go version
go build ./cmd/web ./cmd/worker
go tool templ version
go tool sqlc version
```

Обновление версий выполняется отдельной задачей с прогоном всех тестов и
повторной фиксацией таблицы.
