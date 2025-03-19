# Calendar Export for Google Calendar

Mit diesem Tool ist es möglich seinen Google Calendar für den heutigen Tag zu exportieren
wobei die gesammelten Stunden berechnet und zusammen gefasst werden.

## Voraussetzungen

Damit die App funktioniert muss eine `credentials.json` Datei im Root Verzeichnis des Projekts
vorhanden sein. Diese Datei kann über die Google Developer Console erstellt werden.

## Installation

```bash
git clone https://github.com/maxischmaxi/calendar-export.git
cd calendar-export
go build
go install
```

## Verwendung

```bash
calendar-export
```

```bash
calendar-export -h
```

```bash
calendar-export -no-table
```

```bash
calendar-export -date=2024-12-24
```

## Lizenz

[MIT](https://choosealicense.com/licenses/mit/)
