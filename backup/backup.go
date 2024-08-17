package backup

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v2"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

// GetDriveService creates and returns a Google Drive service client
func GetDriveService() (*drive.Service, error) {
	ctx := context.Background()

	// Read the credentials file
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		return nil, fmt.Errorf("unable to read credentials file: %v", err)
	}

	// Parse the OAuth2 config
	config, err := google.ConfigFromJSON(b, drive.DriveFileScope)
	if err != nil {
		return nil, fmt.Errorf("unable to parse credentials file: %v", err)
	}

	// Retrieve the token from file or get it from the web
	tokenFile := "token.json"
	tok, err := getTokenFromFile(tokenFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokenFile, tok)
	}

	// Create the Drive client
	client := config.Client(ctx, tok)
	srv, err := drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("unable to create Drive client: %v", err)
	}

	return srv, nil
}

// getTokenFromFile retrieves a token from the token file
func getTokenFromFile(tokenFile string) (*oauth2.Token, error) {
	f, err := os.Open(tokenFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	tok := &oauth2.Token{}
	if err := json.NewDecoder(f).Decode(tok); err != nil {
		return nil, err
	}
	return tok, nil
}

// getTokenFromWeb obtains a token from the web
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser:\n%v\n", authURL)
	fmt.Print("Enter the authorization code: ")

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.Background(), code)
	if err != nil {
		log.Fatalf("Unable to retrieve token: %v", err)
	}
	return tok
}

// saveToken saves the token to a file
func saveToken(tokenFile string, token *oauth2.Token) {
	f, err := os.Create(tokenFile)
	if err != nil {
		log.Fatalf("Unable to create token file: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

// UploadFile uploads a file to Google Drive
func UploadFile(filePath string, mimeType string) error {
	srv, err := GetDriveService()
	if err != nil {
		return fmt.Errorf("unable to get Drive service: %v", err)
	}

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("unable to open file: %v", err)
	}
	defer file.Close()

	fileMetadata := &drive.File{
		// Name: file.Name(),
	}

	_, err = srv.Files.Insert(fileMetadata).
		Media(file, googleapi.ContentType(mimeType)).
		Do()
	if err != nil {
		return fmt.Errorf("unable to upload file: %v", err)
	}

	fmt.Println("File uploaded successfully")
	return nil
}

func backupPostgres(dbName, user, password, backupFilePath string) error {
	// Set the PGPASSWORD environment variable
	os.Setenv("PGPASSWORD", password)

	// Prepare the pg_dump command
	cmd := exec.Command("pg_dump", "-U", user, "-F", "c", "-b", "-v", "-f", backupFilePath, dbName)

	// Run the command
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("pg_dump failed: %s\n%s", err, output)
	}

	fmt.Println("Backup completed successfully.")
	return nil
}
