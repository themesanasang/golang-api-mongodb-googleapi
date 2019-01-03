package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/labstack/echo"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
)

// Retrieve a token, saves the token, then returns the generated client.
func getClient(ctx context.Context, config *oauth2.Config) *http.Client {
	cacheFile, err := tokenCacheFile()
	if err != nil {
		log.Fatalf("Unable to get path to cached credentioal file. %v,", err)
	}
	tok, err := tokenFromFile(cacheFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(cacheFile, tok)
	}
	return config.Client(ctx, tok)
}

func tokenCacheFile() (string, error) {
	tokenCacheDir := filepath.Join("./", "googleDriveToken")
	err := os.MkdirAll(tokenCacheDir, 0700)
	return filepath.Join(tokenCacheDir, url.QueryEscape("drive-go-quickstart.json")), err
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the authorization code: \n%v\n", authURL)

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatalf("Uncable to read authorization code %v", err)
	}

	tok, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatalf("Uncable to retrieve token from web %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)
	defer f.Close()
	return t, err
}

// Saves a token to a file path.
func saveToken(file string, token *oauth2.Token) {
	fmt.Printf("Saving credentioal file to: %s \n", file)
	f, err := os.Create(file)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func (h *Handler) Upload(c echo.Context) (err error) {
	ctx := context.Background()

	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, drive.DriveFileScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(ctx, config)

	srv, err := drive.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %v", err)
	}

	//=>upload file
	filename := "icon-3.png"
	uploadFile, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Cannot find such file %v", err)
	}

	//=> folder id=17qabJANndLbF4tHAIoSRx6ujWZsP6TfK
	var folderIDList []string
	folderIDList = append(folderIDList, "17qabJANndLbF4tHAIoSRx6ujWZsP6TfK")

	f := &drive.File{Name: filename, Parents: folderIDList, Description: "test"}

	r, err := srv.Files.Create(f).Media(uploadFile).Do()
	if err != nil {
		log.Fatalf("Upload Failed %v", err)
	}

	//rr := model.Upload{}
	//rr = r.Id

	/*==== read to db = https://drive.google.com/open?id=xxx ===*/
	//pathimg := "https://drive.google.com/open?id="+r.Id;

	//=>read file ที่ เราสร้างเท่านั้น
	/*r, err := srv.Files.List().PageSize(20).
		Fields("nextPageToken, files(id, name)").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve files: %v", err)
	}*/

	return c.JSON(http.StatusCreated, r.Id)
}
