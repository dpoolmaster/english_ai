## Purpose

Управлять разовым пробным и оплаченным доступом, безопасно обрабатывать платежи и возвраты и восстанавливать состояние после потерянных уведомлений.

## ADDED Requirements

### Requirement: Plans and single trial
Тарифы MUST задавать 1/3/12 календарных месяцев, RUB, цену в копейках, активность и квоты без автопродления. Trial MUST выдаваться один раз после подтверждения email на TRIAL_DAYS=3 по умолчанию; защита включает нормализованный email без plus, IP limits и список одноразовых доменов. Основание: FR-BIL-1/2.

#### Scenario: Concurrent email confirmation
- **WHEN** подтверждение и его повтор выполняются параллельно
- **THEN** создаётся один trial, его срок не продлевается повтором.

### Requirement: Access and extension
Доступ MUST определяться now < max(entitlements.ends_at). Без доступа кабинет/вход/оплата работают, занятия закрыты. Оплата MUST продлевать от max(now,current_end). Основание: FR-BIL-3.

#### Scenario: Access boundary
- **WHEN** текущее время равно ends_at
- **THEN** учебный доступ закрыт, покупка остаётся доступной.

#### Scenario: Two successful orders
- **WHEN** два разных подтверждённых заказа оплачены конкурентно
- **THEN** оба периода добавляются без потери продления и без двойной выдачи одного заказа.

### Requirement: Checkout and return
Checkout MUST создавать pending order с ключом идемпотентности и фиксировать сумму/условия покупки. Создание платежа передаёт order_id и данные чека; return page опрашивает каждые 2 секунды до 60 секунд и MUST не выдавать доступ по самому факту возврата. Лимит checkout 10/час. Основание: FR-BIL-4, ТЗ 9.1.

#### Scenario: Browser return without payment
- **WHEN** пользователь открыл return_url неоплаченного заказа
- **THEN** доступ не выдаётся, отображается проверка либо актуальное состояние заказа.

### Requirement: Verified idempotent payment events
Webhook MUST быстро отвечать 200 после долговечной записи события, проверять источник/rate limit и повторно получать платёж у провайдера. Статус, RUB/сумма/order_id MUST совпасть; order/payment/entitlement обновляются атомарно, provider payment уникален. Основание: FR-BIL-4, ТЗ 7.4.

#### Scenario: Forged successful body
- **WHEN** тело сообщает success, но API провайдера возвращает другую сумму, заказ или статус
- **THEN** оплаченный доступ не выдаётся, ошибка наблюдаема и требует обработки.

#### Scenario: Duplicate and out-of-order event
- **WHEN** success повторяется, затем приходит старое canceled
- **THEN** выдаётся один entitlement и подтверждённая оплата не отменяется устаревшим событием.

### Requirement: Reconciliation and administrative changes
Каждые 10 минут система MUST сверять pending старше 5 минут; подтверждённая оплата вызывает одно письмо. Возврат MUST вызывать API провайдера, admin выбирает сокращение доступа; ручной доступ требует причины и аудита. Основание: FR-BIL-4/5.

#### Scenario: Lost webhook
- **WHEN** успешная оплата осталась pending из-за недоставленного webhook
- **THEN** reconciliation подтверждает её через API и выдаёт доступ ровно один раз.

#### Scenario: Refund retry
- **WHEN** запрос возврата повторён после timeout
- **THEN** не создаётся второй финансовый возврат, итог сверяется с провайдером и решение о доступе аудируется.
