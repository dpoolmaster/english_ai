Все задачи ниже — незавершённая реализация, а не подготовка документов. Выполнять последовательно малыми срезами; у каждой задачи запускать применимые генераторы до тестов. Группы 1–24 локальные, 25–26 — перед внешним пилотом/запуском. Основание и контракты — specs/, архитектура и матрица FR — design.md. Ни одна задача не выполнена фактом создания этого файла.

## 1. L01 — Запускаемые процессы

- [ ] 1.1 Проверить официальные версии Go и совместимость Gin/Templ/pgx/sqlc/River, закрепить toolchain и инструменты; проверка: build и документированная таблица версий.
- [ ] 1.2 Создать Go-модуль, cmd/web и cmd/worker с разделёнными composition roots; проверка: оба бинарника собираются без бизнес-провайдеров.
- [ ] 1.3 Добавить env/YAML config с валидацией и безопасными ошибками; проверка: unit-тесты отсутствующих/невалидных значений и отсутствия секрета в ошибке.
- [ ] 1.4 Добавить slog JSON, request/trace id и graceful shutdown; проверка: тест отмены контекста и ручная остановка обоих процессов.
- [ ] 1.5 Добавить healthz и внутренние метрики; проверка: HTTP-тесты кодов/отсутствия секретов, metrics не опубликован пользовательским router.

## 2. L02/L06 — БД, генераторы и локальная воспроизводимость

- [ ] 2.1 Добавить Compose PostgreSQL, пример env без секретов и команды локального запуска; проверка: БД поднимается на чистом окружении без VPS.
- [ ] 2.2 Подключить pgxpool, clock, transaction helpers и тестовую БД testcontainers; проверка: интеграционный commit/rollback и отмена запроса.
- [ ] 2.3 Настроить migrate/sqlc и первую up/down пару; проверка: fresh/up-repeat/down/up на одноразовой БД и воспроизводимый sqlc generate.
- [ ] 2.4 Подключить River SQL-миграции и worker lifecycle; проверка: постановка/rollback вместе с доменной транзакцией и retry после перезапуска.
- [ ] 2.5 Добавить README и единые команды запуска, миграций, генерации и проверок; проверка: пройти инструкцию с чистого клона.
- [ ] 2.6 Настроить начальный CI tests/vet/lint/govulncheck и generated diff; проверка: намеренно устаревший generated-файл приводит к отказу CI.

## 3. L03 — Пользовательский веб-каркас

- [ ] 3.1 Добавить Gin/Templ layout, локальные версии HTMX/Pico/SSE и app.css; проверка: templ generate и прямое открытие полной страницы.
- [ ] 3.2 Добавить обычные и HTMX формы, responseHandling для 4xx и безопасные 404/500; проверка: HTTP-тесты full/fragment и ошибка формы видна в браузере.
- [ ] 3.3 Добавить security headers и запрет небезопасного вывода; проверка: CSP/headers на 2xx/4xx/5xx, XSS fixture отображается текстом.
- [ ] 3.4 Добавить адаптивную русскую навигацию и PWA-манифест; проверка: desktop/mobile viewport и валидация manifest без внешних ассетов.

## 4. L04 — Минимальный admin auth

- [ ] 4.1 Написать негативные интеграционные сценарии login/session/CSRF/role до реализации; проверка: тесты падают по отсутствующим возможностям, не по окружению.
- [ ] 4.2 Добавить users/sessions миграции и sqlc queries с unique email/hash; проверка: ограничения БД, up/down и ручной SQL review.
- [ ] 4.3 Реализовать argon2id PHC, параметры/rehash/semaphore и admin bootstrap command; проверка: hash vectors, лимит параллельности, повтор bootstrap не повышает роль молча.
- [ ] 4.4 Реализовать login/logout, sid rotation, expiry и middleware ролей; проверка: интеграционные тесты отказов и отзыва старого sid.
- [ ] 4.5 Реализовать pre-auth/session CSRF и PostgreSQL rate limit входа; проверка: отсутствие/чужой токен и конкурентный лимит отклоняют мутацию.
- [ ] 4.6 Зафиксировать локальные cookie/HTTPS настройки и production-ограничения; проверка: production не запускается с ослабленной cookie policy, локальный login работает.

## 5. L05/L06 — Vue и контракт API

- [ ] 5.1 Создать Vue/Vite/TS приложение с библиотеками ТЗ и same-origin dev proxy; проверка: npm ci, dev login и отсутствие обхода CSRF через proxy.
- [ ] 5.2 Описать OpenAPI 3.1, ошибки и текущего admin, настроить TS generation; проверка: типы генерируются без diff, API/контрактный тест совпадают.
- [ ] 5.3 Добавить SPA routes, обработку 401/403 и встроенную production build; проверка: прямой вложенный /admin route и API 404 не перепутаны.
- [ ] 5.4 Добавить audit_log и транзакционный audit helper с redaction; проверка: отказ аудита откатывает тестовую мутацию, секретов в diff нет.
- [ ] 5.5 Подключить frontend build/audit и OpenAPI generation к CI; проверка: чистая сборка Go содержит admin assets без запущенного Vite.

## 6. L07 — Регистрация, email и trial

- [ ] 6.1 Написать тесты регистрации, verification replay и конкурентного trial до кода; проверка: тесты отражают FR-AUTH-1/FR-BIL-2.
- [ ] 6.2 Добавить email_tokens, entitlements и unique trial; проверка: PostgreSQL запрещает две выдачи на user.
- [ ] 6.3 Реализовать email/password validation, consent, нормализацию/anti-abuse и IP/domain limits; проверка: weak password, plus email и превышение лимита.
- [ ] 6.4 Реализовать Mailer fake, безопасный локальный просмотр писем и send_email worker; проверка: письмо доступно локально, токен не попадает в application logs.
- [ ] 6.5 Реализовать verification и trial одной транзакцией, access middleware и paywall; проверка: E2E register/verify, повтор токена, expiry с fake clock и запрет неподтверждённого доступа.

## 7. L08 — Ручной каталог и ревью

- [ ] 7.1 Добавить skills/prereqs/vocab/examples/exercises/error_types и необходимые индексы; проверка: миграции, FK/constraints и полный стабильный справочник ошибок ТЗ.
- [ ] 7.2 Реализовать CRUD навыков и проверку DAG под конкурентными правками; проверка: прямой/косвенный/конкурентный цикл отклоняются.
- [ ] 7.3 Добавить схемы mcq/cloze_choice с schema_version и validator; проверка: валидный fixture принимается, неизвестный/сломанный ключ отвергается при save/publish.
- [ ] 7.4 Добавить draft/review/published/flagged/archived, роли и аудит; проверка: запрещённый переход и публикация без права не меняют данные.
- [ ] 7.5 Добавить Vue CRUD/review/ученический preview и небольшой ручной учебный набор; проверка: E2E create/review/publish и отсутствие draft в ученической выдаче.

## 8. L09 — Первое занятие

- [ ] 8.1 Добавить lessons/items/attempts с одним принятым ответом на item; проверка: constraints и тест конкурентных вставок.
- [ ] 8.2 Реализовать start/show/resume с проверкой владельца и доступа; проверка: чужой id закрыт, reload возвращает первый неотвеченный item.
- [ ] 8.3 Реализовать mcq/cloze_choice checker, транзакционное принятие и готовую обратную связь; проверка: два конфликтующих ответа не дублируют attempt/error/result.
- [ ] 8.4 Реализовать finish/summary и source-key events; проверка: повтор finish не удваивает статистику.
- [ ] 8.5 Добавить Playwright register/verify/lesson/finish на локальной почте; проверка: flow проходит без реальных API и ключей.

## 9. L10/L11 — Типы упражнений и FSRS

- [ ] 9.1 Реализовать normalizer с флагами/опечатками; проверка: эталонные и граничные тесты пробелов/апострофов/регистра/сокращений/длины слова.
- [ ] 9.2 Добавить cloze_input схемы, checker и UI; проверка: допустимые ответы и отклонённые опечатки в интеграционном занятии.
- [ ] 9.3 Добавить match/reorder схемы, частичное оценивание и UI; проверка: fixtures неверных пар, перестановок и повторяющихся слов.
- [ ] 9.4 Добавить flashcard и детерминированный translate_to_en; проверка: самооценка и listed translations не вызывают LLM.
- [ ] 9.5 Добавить skill/vocab states и review_log, обёртку go-fsrs и time thresholds; проверка: фиксированные clock/ratings дают ожидаемые состояния.
- [ ] 9.6 Включить review/error/state в answer transaction; проверка: rollback и параллельный replay не дублируют FSRS.
- [ ] 9.7 Реализовать retrievability и ротацию 3–5 контекстов; проверка: изменение времени и least-recent context со стабильными ties.

## 10. L12 — Полный planner

- [ ] 10.1 Реализовать candidates review/weak/new и SQL индексы; проверка: due/14-day errors/prerequisites/goal fixtures и EXPLAIN на наполненной базе.
- [ ] 10.2 Реализовать доли/округление/пустые категории/размер и фиксированный seed; проверка: golden очередь и property tests границ 6–26.
- [ ] 10.3 Сериализовать выбор новых навыков по user/day; проверка: конкурентный start не обходит 1/lesson и 2/day.
- [ ] 10.4 Реализовать explain+3 easy, ограничение типов и адаптацию; проверка: две ошибки вставляют одно задание даже при retry.
- [ ] 10.5 Добавить content gap и безопасную замену flagged item; проверка: опубликованность меняется после start, draft не показывается, недостаток виден в админке.

## 11. L13 — Онбординг и placement

- [ ] 11.1 Добавить настройки цели/минут/timezone/reminders и start A1; проверка: формы сохраняют допустимые значения без фиктивного теста.
- [ ] 11.2 Реализовать Rasch/MAP/SE/выбор/остановку отдельно от HTTP; проверка: fixtures, no-repeat и <=25 вопросов.
- [ ] 11.3 Выполнить фиксированную симуляцию 1000 учеников; проверка: отчёт MAE <= половины CEFR; при провале остановить срез и согласовать уточнение, не менять формулы молча.
- [ ] 11.4 Добавить placement sessions/items, answer UI и идемпотентное завершение/засев; проверка: integration replay и E2E тест с картой знаний.
- [ ] 11.5 Добавить retest 30 дней и weekly recalibration >=200 ответов; проверка: clock boundaries, 199/200 ответов и сохранение истории.

## 12. L14 — Полный auth и восстановление

- [ ] 12.1 Реализовать sliding sessions/devices/logout everywhere; проверка: продление не чаще суток и отзыв всех устройств.
- [ ] 12.2 Реализовать password reset 1h, rehash и смену пароля; проверка: одинаковые ответы unknown email, token replay и отзыв всех старых сессий.
- [ ] 12.3 Реализовать подтверждённую смену email с unique conflict; проверка: до подтверждения старый адрес действует, чужой занятый email не присваивается.
- [ ] 12.4 Подключить captcha по официальной документации и конфигурируемому fake; проверка: порог login и CSP auth-страниц, production captcha не подменяется fake.

## 13. L15 — OAuth

- [ ] 13.1 Проверить документацию Яндекс/VK, PKCE/email semantics и записать дату; проверка: документированный контракт без предположений по API.
- [ ] 13.2 Добавить одноразовый state и identities, реализовать Yandex adapter; проверка: httptest callback replay/чужая cookie и linking existing account.
- [ ] 13.3 Реализовать VK adapter и email confirmation при отсутствии verified email; проверка: нельзя автоматически захватить существующий аккаунт.
- [ ] 13.4 Добавить UI входа/привязки и реальные тестовые callback; проверка: integration scenarios FR-AUTH-4 и отдельный ручной account-takeover report.

## 14. L16 — Экспорт, удаление и retention

- [ ] 14.1 Реализовать JSON export и реестр данных готовых модулей; проверка: полный own-data fixture без чужих записей/password/session secrets.
- [ ] 14.2 Реализовать account delete/access revoke и delete_accounts; проверка: повтор/частичный сбой безопасны, персональные данные удалены <=30 дней.
- [ ] 14.3 Реализовать обезличивание платежей и FK-политики; проверка: необходимые учётные записи сохранены без идентификаторов ученика.
- [ ] 14.4 Добавить cleanup сессий/tokens/auth_attempts/provider logs и агрегаты; проверка: clock boundaries 30 дней/12 месяцев и отсутствие удаления актуальных записей.

## 15. L17 — Полный каталог, жалобы и мониторинг

- [ ] 15.1 Добавить проверяемый YAML/CSV import и отчёт ошибок; проверка: cycle/duplicate/missing reference не оставляют повреждённый каталог.
- [ ] 15.2 Добавить reports form/API/queue и поддержку email; проверка: комментарий экранируется, чужой закрытый объект не раскрывается.
- [ ] 15.3 Реализовать distinct users/7-day flagged и действия reviewer; проверка: повторы одного автора и конкурентная третья жалоба.
- [ ] 15.4 Добавить flag_exercises при >=50 attempts и порогах 15/98%; проверка: boundary fixtures и задача идемпотентна.
- [ ] 15.5 Добавить hotkeys/bulk/preview и полный review workflow; проверка: ручной замер целевых 60 упражнений/час и аудит каждой изменённой записи.

## 16. L18 — Биллинг на fake

- [ ] 16.1 Написать order/access/webhook race tests до реализации; проверка: сценарии mismatched amount/order/status и duplicate success сформулированы явно.
- [ ] 16.2 Добавить plans/orders/payments/inbox/financial operations и unique entitlement per order; проверка: constraints/migrations и allowed transitions.
- [ ] 16.3 Реализовать fake provider, checkout/return polling и rate limit; проверка: return без оплаты не выдаёт доступ, replay checkout безопасен.
- [ ] 16.4 Реализовать durable webhook + verified processing transaction; проверка: forged body/repeated/out-of-order events и две оплаты одного user.
- [ ] 16.5 Реализовать reconcile 10min/5min и одно paid email; проверка: lost webhook восстанавливается без второй выдачи.
- [ ] 16.6 Реализовать fake refunds/manual grants с durable audit/reason; проверка: timeout retry не создаёт второй возврат и доступ меняется только по выбранному действию.

## 17. L19 — Конкурентные квоты

- [ ] 17.1 Написать concurrent last-unit/failure/replay/crash тесты; проверка: используются настоящие PostgreSQL транзакции.
- [ ] 17.2 Добавить quota operations/reservations/counters и UTC period ADR; проверка: сумма consumption+active reservations не превышает limit под нагрузкой.
- [ ] 17.3 Реализовать admit/succeed/fail с атомарным UPSERT и domain result; проверка: повтор success списывает один раз, оба provider failures не списывают.
- [ ] 17.4 Реализовать leases/deadlines/recovery и period rollover; проверка: crash до/после result commit и смена суток не блокируют bucket навсегда.
- [ ] 17.5 Добавить тарифные defaults, остаток и 80% warning; проверка: trial/paid boundaries без включения голосовых/письменных функций MVP.

## 18. L20 — Провайдеры и LLM

- [ ] 18.1 Описать интерфейсы ТЗ, fake и purpose routing, проверить выбранный реальный API; проверка: бизнес-пакеты не импортируют реализации, дата документации записана.
- [ ] 18.2 Реализовать OpenAI-compatible complete/stream и cancellation; проверка: httptest success, timeout, truncated stream и malformed response.
- [ ] 18.3 Реализовать bounded retry/breaker/fallback; проверка: 4xx не повторяется, network/429/5xx ограничены, общий deadline соблюдается.
- [ ] 18.4 Реализовать Structured schema+one repair и prompt version references; проверка: второй invalid JSON даёт контролируемую ошибку.
- [ ] 18.5 Добавить calls journal/pricing/trace/redaction; проверка: каждая попытка учтена, PII/секреты/тексты отсутствуют в логах и provider metadata.
- [ ] 18.6 Добавить budget routing и cost_alerts; проверка: fake clock/budget переводит некритичный purpose в configured fallback с алертом.
- [ ] 18.7 Завершить translate_to_en LLM fallback через quota operation; проверка: listed answer не вызывает модель, failure/limit не выдаёт ложный успех.

## 19. L21 — Реальная оплата

- [ ] 19.1 Проверить документацию ЮKassa, методы/чеки/IP ranges и доступы магазина; проверка: записанные ссылки/дата и подтверждённые условия без секретов.
- [ ] 19.2 Реализовать create/get/refund/webhook adapter и source restrictions; проверка: httptest fixtures, идемпотентность, RUB/amount/order metadata mismatch.
- [ ] 19.3 Проверить тестовый checkout/return/reconcile через доступный HTTPS webhook endpoint; проверка: отчёт тестового полного цикла, endpoint отдельно от локальной обязательной среды.

## 20. L22 — Генерация и промпты

- [ ] 20.1 Добавить immutable prompt versions, unique active и eval fixtures; проверка: concurrent activation оставляет одну активную версию и audit.
- [ ] 20.2 Реализовать generate_exercises с прогрессом/model/version/cost; проверка: retry после части batch не дублирует контент.
- [ ] 20.3 Реализовать schema/second solver/unique answer проверки; проверка: disagreement не попадает незаметно в published.
- [ ] 20.4 Реализовать level/length/stoplist/hash/pg_trgm и >=2 translation variants; проверка: проблемные fixtures направляются на review с причиной.
- [ ] 20.5 Добавить Vue generation/prompts/eval comparison и ручную активацию; проверка: regression блокирует автоматическое включение, откат переключает версию.

## 21. L23 — Текстовый тьютор

- [ ] 21.1 Добавить conversations/messages/scenarios и operation ids; проверка: owner/role/access и один активный user turn на диалог.
- [ ] 21.2 Реализовать context builder/profile snapshot/history summary и prompt; проверка: token/message budgets, отсутствие PII, длины и 40-reply limit.
- [ ] 21.3 Реализовать POST turn и SSE read-only subscription; проверка: reconnect/repeated POST не запускает второй LLM-вызов.
- [ ] 21.4 Реализовать response states/token escape/commit-before-done; проверка: HTML fixture безопасен, lost done восстанавливает completed result без второго списания.
- [ ] 21.5 Реализовать fallback до первого token и partial-failure после; проверка: fault-injection до/после token, timeout 20с и no-charge failure.
- [ ] 21.6 Добавить hint(3)/translate, rate limit и chat UI; проверка: браузерные лимиты и доступ к чужому stream отклоняются сервером.
- [ ] 21.7 Реализовать analyze_errors JSON/repair/idempotency и finish panel; проверка: invalid result пропускается, duplicate job не дублирует error events.
- [ ] 21.8 Добавить injection/toxicity rules, disclaimer/report и tutor eval; проверка: эталонные recast/length/language/injection сценарии и регрессионный отчёт.

## 22. L24 — Дашборд, коучинг и события

- [ ] 22.1 Реализовать streak/freeze/minutes/words/mastery/due; проверка: timezone boundaries, повтор finish и незавершённое занятие не удваивают день.
- [ ] 22.2 Добавить notifications/dedupe и scheduler каждые 15min; проверка: параллельные scheduler создают одно уведомление на локальный день.
- [ ] 22.3 Реализовать daily/3-day/7-day/access expiry письма; проверка: lesson completed before delivery отменяет неактуальное daily reminder.
- [ ] 22.4 Добавить signed unsubscribe/settings и SMTP adapter; проверка: forged signature отклоняется, queued coaching suppresses after opt-out, security mail не отключается.
- [ ] 22.5 Добавить product events из specs/release-readiness с source dedupe; проверка: повтор запросов не удваивает paid/lesson_finished/reminder_sent и внешних trackers нет.

## 23. L24 — Завершение админки

- [ ] 23.1 Добавить users search/profile/progress/errors/sessions/block/grants; проверка: block не обходится старой сессией или повторным входом.
- [ ] 23.2 Добавить placement/scenarios/settings/plans/quotas/refunds screens; проверка: API CRUD соответствует OpenAPI и server validation.
- [ ] 23.3 Добавить dashboard/conversion/revenue/costs/top20/content gaps; проверка: агрегаты совпадают с фиксированным набором событий и затрат.
- [ ] 23.4 Добавить audit viewer/cursor pagination и полную матрицу reviewer denial; проверка: все мутации дают audit, reviewer не читает users/payments/settings.

## 24. L25 — Комплексная локальная приёмка

- [ ] 24.1 Расширить export/delete/retention на все добавленные модули; проверка: полный user fixture удалён/экспортирован без потери обязательных учётных записей.
- [ ] 24.2 Выполнить Playwright register/verify/onboard/placement/lesson/fake payment/chat/admin review; проверка: все сценарии проходят с fake без реальных аккаунтов.
- [ ] 24.3 Выполнить security/concurrency integration suite auth/CSRF/roles/quota/webhook/lesson; проверка: независимые PostgreSQL транзакции и fault-injection не нарушают инварианты.
- [ ] 24.4 Выполнить генераторы, go test ./..., go vet ./..., lint/govulncheck, frontend build/audit; проверка: актуальные generated files и >=70% coverage ключевых пакетов.
- [ ] 24.5 Сверить каждый FR этапа 1 с тестом/ручным доказательством; проверка: матрица приёмки явно разделяет passed/local, pending/external и deferred этапы 2–4.
- [ ] 24.6 Выполнить 100 конкурентных логинов с согласованным бюджетом памяти; проверка: измеренный peak укладывается в бюджет, semaphore не допускает более 4 хэшей одновременно.

## 25. Внешняя готовность — только после локального MVP

- [ ] 25.1 Обследовать общий VPS/Nginx и согласовать топологию Caddy без влияния на соседние проекты; проверка: документированные порты/сети/resource budget и безопасный deployment plan.
- [ ] 25.2 Настроить образы/Compose/deploy by tag/TLS/secrets и rollback; проверка: smoke после deploy и предыдущий образ работает с расширенной схемой.
- [ ] 25.3 Проверить production SSE, metrics isolation, non-root и алерты; проверка: первый token не буферизуется proxy, DB/metrics недоступны извне.
- [ ] 25.4 Настроить encrypted backup другому российскому provider/retention и restore; проверка: чистая машина, RTO <=4ч/RPO <=24ч и тест неуспешного backup alert.
- [ ] 25.5 Проверить реальные SMTP/SPF/DKIM/DMARC, OAuth и LLM fallback; проверка: отчёты доставки/входа/отказа с датой и окружением.
- [ ] 25.6 Проверить боевую оплату, доступ, чек и возврат; проверка: подтверждённый полный цикл без публикации финансовых/персональных данных в репозитории.
- [ ] 25.7 Выполнить ASVS L1 и ручной OAuth takeover review; проверка: закрытые критичные findings и отчёт, не заменённый unit tests.
- [ ] 25.8 Выполнить k6 150 sessions/10min и real tutor latency; проверка: пороги specs/release-readiness, отчёт с hardware/provider configuration.

## 26. Контент, юридические предпосылки и пилот

- [ ] 26.1 Подготовить программу A1–A2 и подтвердить лицензии источников до импорта; проверка: утверждённый DAG и реестр коммерческих прав.
- [ ] 26.2 Подготовить/проверить >=700 published упражнений, словарь MVP и >=120 placement items (30/уровень); проверка: DB report учитывает published и человеческое ревью, не draft count.
- [ ] 26.3 Подтвердить хранение РФ/уведомление РКН/политику/согласие/оферту/возвраты/контакты с владельцем и юристом до реальных пользователей; проверка: зафиксированный статус внешних подтверждений и опубликованные документы.
- [ ] 26.4 Провести пилот 20–50 пользователей минимум 14 дней; проверка: события активации/завершения/регулярности/конверсии и список проблем собраны.
- [ ] 26.5 Устранить критичные проблемы и сверить приложение D/приёмку MVP; проверка: итоговый отчёт с доказательствами, оставшимися ограничениями и решением владельца о платном запуске.
- [ ] 26.6 Рассчитать обычную/предельную себестоимость по реальным ценам и настроить метрики надёжности; проверка: переменные затраты <=25% выручки, полный выбор квот не убыточен, availability/provider errors/жалобы измеряются без заявления о недоступных долгосрочных результатах.
