package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"regexp"
	"runtime"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	calendar "google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
	people "google.golang.org/api/people/v1"
)

func main() {
	datePtr := flag.String("date", "", "Datum für das die Termine angezeigt werden sollen (Format: 2021-12-31)")
	noTable := flag.Bool("no-table", false, "Keine Tabelle anzeigen")
	yesterday := flag.Bool("yesterday", false, "Termine von gestern anzeigen")
	tomorrow := flag.Bool("tomorrow", false, "Termine von morgen anzeigen")

	flag.Parse()

	if *yesterday && *tomorrow {
		log.Fatalf("Die Optionen -yesterday und -tomorrow können nicht gleichzeitig verwendet werden.")
	}

	if *yesterday && *datePtr != "" {
		log.Fatalf("Die Optionen -yesterday und -date können nicht gleichzeitig verwendet werden.")
	}

	if *tomorrow && *datePtr != "" {
		log.Fatalf("Die Optionen -tomorrow und -date können nicht gleichzeitig verwendet werden.")
	}

	var date time.Time

	if *yesterday {
		date = time.Now().AddDate(0, 0, -1)
	} else if *tomorrow {
		date = date.AddDate(0, 0, 1)
	} else if len(*datePtr) > 0 {
		var err error
		d, err := time.Parse("2006-01-02", *datePtr)
		if err != nil {
			log.Fatalf("Fehler beim Parsen des Datums: %v", err)
		}

		date = time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, time.UTC)
	} else {
		date = time.Now()
	}

	ctx := context.Background()
	config := getConfig()
	client := getClient(config)

	auth := option.WithHTTPClient(client)

	srv, err := calendar.NewService(ctx, auth)
	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}

	pplSrv, err := people.NewService(ctx, auth)
	if err != nil {
		log.Fatalf("Unable to retrieve People client: %v", err)
	}

	person, err := pplSrv.People.Get("people/me").PersonFields("emailAddresses").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve person: %v", err)
	}

	email := ""
	if len(person.EmailAddresses) > 0 {
		email = person.EmailAddresses[0].Value
	} else {
		log.Fatalf("No email address found for the user.")
	}

	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	endOfDay := startOfDay.Add(24 * time.Hour)
	start := startOfDay.Format(time.RFC3339)
	end := endOfDay.Format(time.RFC3339)

	events, err := srv.Events.List("primary").TimeMin(start).TimeMax(end).SingleEvents(true).OrderBy("startTime").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve next ten of the user's events: %v", err)
	}

	items := make([]*calendar.Event, 0)

	for _, item := range events.Items {
		if shouldIgnoreMeeting(item, email) {
			continue
		}

		items = append(items, item)
	}

	if !*noTable {
		printTable(items, date)
	} else {
		printList(items)
	}
}

func shouldIgnoreMeeting(item *calendar.Event, email string) bool {
	isAttendee := false
	for _, attendee := range item.Attendees {
		if attendee.Email == email {
			if attendee.ResponseStatus != "accepted" {
				return true
			}
			isAttendee = true
			break
		}
	}

	if !isAttendee {
		if item.Organizer.Email != email {
			return true
		}
	}

	switch item.Summary {
	case "Außer Haus", "Zuhause", "Zeiten eintragen", "Urlaub", "Krank", "Feiertag", "Krankheit", "Urlaubstag", "Krankheitstag", "Feiertagstag":
		return true
	default:
		return false
	}
}

func getConfigDir() string {
	var configDir string

	if runtime.GOOS == "windows" {
		configDir = os.Getenv("APPDATA")
	} else {
		configHome := os.Getenv("XDG_CONFIG_HOME")
		if configHome == "" {
			configHome = path.Join(os.Getenv("HOME"), ".config")
		}
		configDir = path.Join(configHome, "calendar-export")
	}

	return configDir
}

func getConfig() *oauth2.Config {
	configDir := getConfigDir()
	b, err := os.ReadFile(path.Join(configDir, "credentials.json"))
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, calendar.CalendarReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	return config
}

func getSummary(item *calendar.Event) string {
	match, err := regexp.MatchString("^(GALAXY|NOTICKET|fix|quality)", item.Summary)
	if err != nil {
		return item.Summary
	}

	if match {
		return item.Summary
	}

	return fmt.Sprintf("Meeting %s", item.Summary)
}

func printList(items []*calendar.Event) {
	for _, item := range items {
		start := parseTime(item.Start)
		end := parseTime(item.End)
		diff := end.Sub(start)
		timeValue := formatDiff(diff)

		fmt.Printf("%s: %s\n", timeValue, getSummary(item))
	}
}

func printTable(items []*calendar.Event, date time.Time) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Zeit", fmt.Sprintf("Zusammenfassung %s", date.Format("2006-01-02"))})
	totalDiff := 0.0

	for _, item := range items {
		start := parseTime(item.Start)
		end := parseTime(item.End)
		diff := end.Sub(start)
		totalDiff += diff.Minutes()
		timeValue := formatDiff(diff)

		t.AppendRow(table.Row{timeValue, getSummary(item)})
	}

	t.AppendSeparator()
	t.AppendFooter(table.Row{"Total", formatDiff(time.Duration(totalDiff * float64(time.Minute)))})
	t.Render()
}

func parseTime(timestamp *calendar.EventDateTime) time.Time {
	if timestamp.DateTime != "" {
		t, err := time.Parse(time.RFC3339, timestamp.DateTime)
		if err != nil {
			log.Fatalf("Fehler beim Parsen des Zeitstempels: %v", err)
		}

		return t
	}

	t, err := time.Parse("2006-01-02", timestamp.Date)
	if err != nil {
		log.Fatalf("Fehler beim Parsen des Zeitstempels: %v", err)
	}

	return t
}

func formatDiff(dur time.Duration) string {
	hours := int(dur / time.Hour)
	minutes := int(dur % time.Hour / time.Minute)
	return fmt.Sprintf("%d:%02d", hours, minutes)
}

func getClient(config *oauth2.Config) *http.Client {
	configDir := getConfigDir()
	tokenFile := path.Join(configDir, "token.json")

	tok, err := tokenFromFile(tokenFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokenFile, tok)
	}
	return config.Client(context.Background(), tok)
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Fatalf("Fehler beim Schließen der Datei: %v", err)
		}
	}()

	var token oauth2.Token
	err = json.NewDecoder(f).Decode(&token)
	return &token, err
}

func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Speichere Token in: %s\n", path)
	f, err := os.Create(path)
	if err != nil {
		log.Fatalf("Unable to save oauth token: %v", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Fatalf("Fehler beim Schließen der Datei: %v", err)
		}
	}()
	err = json.NewEncoder(f).Encode(token)

	if err != nil {
		log.Fatalf("Unable to save oauth token: %v", err)
	}
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	codeCh := make(chan string)

	go func() {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			code := r.URL.Query().Get("code")
			_, err := fmt.Fprintf(w, "Autorisierung erfolgreich. Du kannst dieses Fenster schließen.")
			if err != nil {
				log.Fatalf("Fehler beim Schreiben der Antwort: %v", err)
			}
			codeCh <- code
		})
		log.Fatal(http.ListenAndServe("localhost:8080", nil))
	}()

	config.RedirectURL = "http://localhost:8080"
	config.Scopes = append(config.Scopes, calendar.CalendarReadonlyScope)
	config.Scopes = append(config.Scopes, people.UserinfoEmailScope)
	config.Scopes = append(config.Scopes, people.UserinfoProfileScope)
	authUrl := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

	openURL(authUrl)
	fmt.Println("open the following url in your browse, if not opened automatically:")
	fmt.Println(authUrl)

	code := <-codeCh

	tok, err := config.Exchange(context.Background(), code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}

	return tok
}

func openURL(url string) {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		fmt.Println("Bitte öffne diesen Link manuell: ", url)
	}
	if err != nil {
		log.Println("Fehler beim Öffnen des Browsers:", err)
	}
}
