package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Mailer sends the magic-link login email.
type Mailer interface {
	SendMagicLink(ctx context.Context, to, link string) error
}

// ResendMailer sends via the Resend HTTP API (no SDK dependency).
type ResendMailer struct {
	apiKey string
	from   string
	client *http.Client
}

func NewResendMailer(apiKey, from string) *ResendMailer {
	return &ResendMailer{apiKey: apiKey, from: from, client: &http.Client{Timeout: 10 * time.Second}}
}

func (m *ResendMailer) SendMagicLink(ctx context.Context, to, link string) error {
	payload, err := json.Marshal(map[string]any{
		"from":    m.from,
		"to":      []string{to},
		"subject": "Your Stanza Bonanza login link",
		"html": fmt.Sprintf(
			`<p>Sign in to Stanza Bonanza:</p><p><a href="%s">Sign in</a></p><p>This link expires in 15 minutes. If you didn't request it, ignore this email.</p>`,
			link,
		),
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.resend.com/emails", bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+m.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := m.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("resend returned %d", resp.StatusCode)
	}
	return nil
}
