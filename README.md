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
$ calendar-export -no-table
```

```bash
$ calendar-export -date=2024-12-24
```

## Example Output

```bash
+-------+---------------------------------------------------+-------+
| ZEIT  | ZUSAMMENFASSUNG                                   | TOTAL |
+-------+---------------------------------------------------+-------+
| 0:45  | Meeting GALAXY-1241 review                        | 0:45  |
| 0:30  | Meeting GALAXY-1235 fixed font-size               | 1:15  |
| 0:15  | Meeting Galaxy Daily Standup (DEVs)               | 1:30  |
| 0:45  | Meeting NOTICKET general quality improvement talk | 2:15  |
| 0:15  | Meeting WebDev Daily                              | 2:30  |
| 1:30  | Meeting Refinement Content Components             | 4:00  |
| 1:45  | Meeting FW: Bi-Weekly | DevOps                    | 5:45  |
+-------+---------------------------------------------------+-------+
| TOTAL |                                                   | 5:45  |
+-------+---------------------------------------------------+-------+
```

## Lizenz

[MIT](https://choosealicense.com/licenses/mit/)
