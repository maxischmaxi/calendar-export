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

## Lizenz

[MIT](https://choosealicense.com/licenses/mit/)
