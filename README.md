# Calendar Export for Google Calendar

Mit diesem Tool ist es möglich seinen Google Calendar für den heutigen Tag zu exportieren
wobei die gesammelten Stunden berechnet und zusammen gefasst werden.

## Voraussetzungen

Damit die App funktioniert muss eine `credentials.json` Datei vorhanden sein.
Diese Datei kann über die Google Developer Console erstellt werden.

### Unix

```bash
~/.config/calendar-export/credentials.json
```

### Windows

```bash
%APPDATA%\calendar-export\credentials.json
```

## Installation

```bash
$ curl -sSfL https://raw.githubusercontent.com/maxischmaxi/calendar-export/main/install.sh | sh
```

```bash
$ git clone https://github.com/maxischmaxi/calendar-export.git
$ cd calendar-export
$ go build
$ go install
```

## Verwendung

```bash
$ calendar-export
```

```bash
$ calendar-export -h
```

```bash
$ calendar-export -yesterday
```

```bash
$ calendar-export -tomorrow
```

```bash
$ calendar-export -no-table
```

```bash
$ calendar-export -date=2024-12-24
```

## Example Output

```bash
+-------+-----------------------------------------------------------------+
| ZEIT  | ZUSAMMENFASSUNG                                                 |
+-------+-----------------------------------------------------------------+
| 1:00  | XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX                     |
| 0:15  | XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX                             |
| 0:45  | XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX              |
| 0:30  | XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX                         |
| 0:15  | XXXXXXXXXXXXXXXXXXXX                                            |
| 1:00  | XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX                                |
| 1:00  | XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX |
+-------+-----------------------------------------------------------------+
| TOTAL | 4:45                                                            |
+-------+-----------------------------------------------------------------+
```

## New Version

```bash
$ git tag v1.0.0
$ git push origin v1.0.0
```

## Lizenz

[MIT](https://choosealicense.com/licenses/mit/)
