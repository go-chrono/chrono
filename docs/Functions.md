# Functions

The table below shows the functions attached to the primary date and time types of `chrono`. Where a function returns a `chrono` type, that type is specified.

|                                                | `LocalDate` | `LocalTime` |     `LocalDateTime`      | `OffsetTime` |     `OffsetDateTime`      |
| ---------------------------------------------- | :---------: | :---------: | :----------------------: | :----------: | :-----------------------: |
| **`Date() (year int, month Month, day int)`**  |      🗸      |             |                          |              |                           |
| **`IsLeapYear() bool`**                        |      🗸      |             |                          |              |                           |
| **`Weekday() Weekday`**                        |      🗸      |             |                          |              |                           |
| **`YearDay() int`**                            |      🗸      |             |                          |              |                           |
| **`ISOWeek() (isoYear, isoWeek int)`**         |      🗸      |             |                          |              |                           |
| **`Clock() (hour, min, sec int)`**             |             |      🗸      |                          |      🗸       |                           |
| **`Nanosecond() int`**                         |             |      🗸      |                          |      🗸       |                           |
| **`BusinessHour() int`**                       |             |      🗸      |                          |      🗸       |                           |
| **`Offset() Offset`**                          |             |             |                          |      🗸       |             🗸             |
| **`Split() ...`**                              |             |             | `LocalDate`, `LocalTime` |              | `LocalDate`, `OffsetTime` |
| **`Local() ...`**                              |             |             |                          | `LocalTime`  |      `LocalDateTime`      |
| **`In() ...`**                                 |             |             |                          | `LocalTime`  |     `OffsetDateTime`      |
| **`UTC() ...`**                                |             |             |                          | `LocalTime`  |     `OffsetDateTime`      |
| **`Sub() ...`**                                |             |      🗸      |            🗸             |      🗸       |             🗸             |
| **`Add(...) ...`**                             |             | `LocalTime` |     `LocalDateTime`      | `LocalTime`  |     `OffsetDateTime`      |
| **`CanAdd(...) bool`**                         |             |      🗸      |            🗸             |      🗸       |             🗸             |
| **`AddDate(years, months, days int) ...`**     | `LocalDate` |             |     `LocalDateTime`      |              |     `OffsetDateTime`      |
| **`CanAddDate(years, months, days int) bool`** |      🗸      |             |            🗸             |              |             🗸             |
| **`Format(layout string) string`**             |      🗸      |      🗸      |            🗸             |      🗸       |             🗸             |
| **`Parse(layout, value string) error`**        |      🗸      |      🗸      |            🗸             |      🗸       |             🗸             |
