package google

import (
	"context"
	"fmt"
	"github.com/SamoylikV/LocaleParse/internal/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func Parse(cfg *config.Config, readRange string) (map[string]string, error) {
	var tokenSource oauth2.TokenSource
	var err error
	tokenSource, err = google.DefaultTokenSource(context.Background(), sheets.SpreadsheetsReadonlyScope)
	if err != nil {
		var jwtConfig *jwt.Config
		jwtConfig, err = google.JWTConfigFromJSON([]byte(cfg.GoogleCredentialsJSON), sheets.SpreadsheetsReadonlyScope)
		if err != nil {
			return nil, fmt.Errorf("unable to create token source or JWT config: %w", err)
		}

		tokenSource = jwtConfig.TokenSource(context.Background())
	}

	client := oauth2.NewClient(context.Background(), tokenSource)

	srv, err := sheets.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve Sheets client: %w", err)
	}
	resp, err := srv.Spreadsheets.Values.Get(cfg.SpreadsheetID, readRange).Do()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve data from sheet: %w", err)
	}

	dataMap := make(map[string]string)
	if len(resp.Values) == 0 {
		return dataMap, nil
	}

	for _, row := range resp.Values {
		if len(row) >= 2 {
			key := fmt.Sprintf("%v", row[0])
			value := fmt.Sprintf("%v", row[1])

			if key != "" {
				dataMap[key] = value
			}
		}
	}
	return dataMap, nil
}
