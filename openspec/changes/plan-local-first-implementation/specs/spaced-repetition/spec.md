## Purpose

Поддерживать расписания повторений навыков и слов, сохранять историю и рассчитывать актуальную вероятность вспоминания для интерфейса и планировщика.

## ADDED Requirements

### Requirement: Review states and ratings
Система MUST хранить отдельные FSRS-состояния навыков/слов: stability, difficulty, state, reps, lapses, due_at, last_review_at. Быстрый верный ответ получает Good/Easy по конфигурации, медленный/с подсказкой Hard, неверный Again, flashcard — самооценку. Основание: FR-SRS-1/2.

#### Scenario: Hinted correct answer
- **WHEN** верный ответ получен с подсказкой
- **THEN** применяется Hard, а не Easy, и рассчитывается новое расписание.

### Requirement: Atomic review history
Принятый ответ MUST сохранять согласованные review_log и состояние; повтор попытки MUST не создавать дополнительного review. Основание: FR-SRS-3.

#### Scenario: Concurrent answer retry
- **WHEN** два запроса отвечают на один lesson item
- **THEN** фиксируется одно обновление состояния и один review, повтор не сдвигает due_at.

### Requirement: Current recall and contexts
Retrievability MUST вычисляться на текущий момент; повтор слова MUST выбирать давно не показанный проверенный пример из 3–5. Основание: FR-SRS-4/5.

#### Scenario: Clock advances
- **WHEN** при том же состоянии текущее время сдвинуто вперёд
- **THEN** вероятность пересчитывается, а не возвращается как сохранённая устаревшая оценка.

#### Scenario: Context rotation
- **WHEN** у повторяемого слова несколько опубликованных примеров
- **THEN** выбирается наиболее давно показанный с устойчивым порядком при равенстве.
