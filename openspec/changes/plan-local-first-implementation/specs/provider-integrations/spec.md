## Purpose

Обеспечить заменяемые внешние сервисы с локальными заглушками, контролируемыми сбоями, обезличенными запросами и прозрачной стоимостью.

## ADDED Requirements

### Requirement: Configurable providers and fakes
LLM, Transcriber, Synthesizer, PronunciationAssessor, PaymentProvider, Mailer, ObjectStore MUST заменяться конфигурацией без изменения бизнес-логики, иметь fake и поддержку отмены/таймаутов. Неиспользуемые голосовые адаптеры не требуются в MVP. Основание: ТЗ 5.

#### Scenario: Local offline flow
- **WHEN** выбраны fake-провайдеры
- **THEN** тестовые регистрация, оплата и тьютор работают без ключей и обращений к реальным API.

### Requirement: Bounded resilience
Retry MUST быть не более двух повторов с exponential backoff/jitter только для network/5xx/429; неидемпотентное действие не повторяется без ключа. Circuit breaker и fallback по purpose MUST конфигурироваться. Основание: ТЗ 5.1.

#### Scenario: Client error
- **WHEN** провайдер вернул обычную 4xx, кроме 429
- **THEN** автоматические retry не выполняются, возвращается контролируемая ошибка.

#### Scenario: Canceled request
- **WHEN** истёк общий deadline
- **THEN** retry/fallback не продолжаются за его пределами.

### Requirement: Private and measurable calls
Каждый вызов LLM/ASR/TTS MUST записывать purpose/provider/model/volume/cost/latency/status/trace_id и prompt_version_id для LLM. Содержимое и идентификаторы ученика MUST не уходить провайдеру и не попадать в логи; секреты только env. Основание: ТЗ 2.1, 5.1, 8.3, 9.1.

#### Scenario: Failed fallback chain
- **WHEN** выполнены неуспешные вызовы двух провайдеров
- **THEN** обе попытки видны в учёте затрат без текстов и email, пользователю не раскрываются ключи и технические ответы.

### Requirement: Structured output repair
Структурированный результат MUST проходить строгую схему; после первой ошибки допускается один repair request, затем контролируемый отказ. Основание: ТЗ 5.2.

#### Scenario: Invalid repaired JSON
- **WHEN** исходный и ремонтный ответы невалидны
- **THEN** невалидные данные не записываются как корректный доменный результат, третьего repair нет.

### Requirement: Budget controls and documentation
Тарифы провайдера MUST задаваться конфигурацией; бюджетные пороги переключают/отключают некритичные purposes с алертом. Параметры API MUST подтверждаться официальной документацией с датой проверки. Основание: ТЗ 5.6, 9.8.

#### Scenario: Daily budget exceeded
- **WHEN** расход превышает настроенный порог
- **THEN** применяется заданная дешёвая цепочка либо отключение некритичного purpose, отправляется алерт и не скрывается фактическая стоимость.
