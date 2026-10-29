## Purpose

Предоставить единую административную панель с контрактным API, серверными полномочиями, проверяемым аудитом и доступными средствами управления MVP.

## ADDED Requirements

### Requirement: Contract and SPA behavior
API MUST использовать /admin/api/v1, OpenAPI 3.1, JSON, курсорную пагинацию, ошибки code/message/details и CSRF для мутаций; frontend types генерируются из контракта. SPA MUST отдаваться без роли, API MUST возвращать 403 без требуемых полномочий. Основание: FR-ADM-1/2.

#### Scenario: Direct SPA and denied API
- **WHEN** пользователь без роли открывает вложенный маршрут /admin и запрашивает API
- **THEN** оболочка загружается, данные API не выдаются, экран корректно отображает отсутствие доступа.

### Requirement: Reviewer isolation
Reviewer MUST видеть только контент и жалобы, admin — разрешённые административные функции; API MUST не доверять скрытию кнопок. Основание: FR-ADM-1, FR-AUTH-8.

#### Scenario: Reviewer refund attempt
- **WHEN** reviewer вручную вызывает refund/users/settings API
- **THEN** запрос отклоняется без раскрытия или изменения защищённых данных.

### Requirement: Complete MVP administration
Панель MUST включать dashboard, skills/prereqs/exercises, review, generation jobs, prompts, placement, users/sessions, orders/refunds, reports, scenarios, settings/quotas/trial, costs и audit log. Пользователи поддерживают поиск по email, блокировку, принудительный выход и ручной доступ. Основание: ТЗ 4.12.

#### Scenario: User blocked
- **WHEN** admin блокирует пользователя и завершает его сессии
- **THEN** дальнейшие защищённые запросы и повторный вход не обходят блокировку, действие аудируется.

#### Scenario: Cursor pagination
- **WHEN** администратор запрашивает следующую страницу журнала по курсору
- **THEN** возвращаются записи в устойчивом порядке без раскрытия чужих областей reviewer и без offset-зависимого пропуска на неизменном снимке.

### Requirement: Atomic audit and admin limits
Каждая успешная мутация MUST записывать actor/action/entity/time/IP и безопасный diff; секреты в diff запрещены. При невозможности долговечно записать аудит мутация MUST не считаться успешной. API ограничивается 600/мин. Основание: FR-ADM-1, ТЗ 9.1.

#### Scenario: Audit write failure
- **WHEN** бизнес-изменение возможно, но запись audit log падает
- **THEN** локальная мутация откатывается; внешний финансовый эффект сохраняет восстанавливаемый след операции и не повторяется вслепую.
