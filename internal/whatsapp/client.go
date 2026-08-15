package whatsapp

import (
	"context"
	"fmt"
	"os"

	"github.com/mdp/qrterminal/v3"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	_ "modernc.org/sqlite"
)

// Client wraps the whatsmeow client.
type Client struct {
	WA         *whatsmeow.Client
	sessionDir string
}

// NewClient creates a new WhatsApp client.
func NewClient(sessionDir string) (*Client, error) {
	if err := os.MkdirAll(sessionDir, 0755); err != nil {
		return nil, fmt.Errorf("create session dir: %w", err)
	}

	dbPath := sessionDir + "/whatsapp.db"
	container, err := sqlstore.New(context.Background(), "sqlite", fmt.Sprintf("file:%s?_foreign_keys=on", dbPath), waLog.Noop)
	if err != nil {
		return nil, fmt.Errorf("create sqlstore: %w", err)
	}

	deviceStore, err := container.GetFirstDevice(context.Background())
	if err != nil {
		return nil, fmt.Errorf("get device store: %w", err)
	}

	waClient := whatsmeow.NewClient(deviceStore, waLog.Noop)

	return &Client{
		WA:         waClient,
		sessionDir: sessionDir,
	}, nil
}

// Connect connects to WhatsApp, showing QR code if needed.
func (c *Client) Connect() error {
	if c.WA.Store.ID == nil {
		// First time — need QR code
		qrChan, err := c.WA.GetQRChannel(context.Background())
		if err != nil {
			return fmt.Errorf("get qr channel: %w", err)
		}

		if err := c.WA.Connect(); err != nil {
			return fmt.Errorf("connect: %w", err)
		}

		for evt := range qrChan {
			switch evt.Event {
			case "code":
				printQR(evt.Code)
			case "success":
				log.Info().Msg("QR code scanned successfully")
			case "timeout":
				return fmt.Errorf("QR code timeout")
			case "error":
				return fmt.Errorf("QR code error: %v", evt.Error)
			}
		}
	} else {
		if err := c.WA.Connect(); err != nil {
			return fmt.Errorf("connect: %w", err)
		}
	}
	return nil
}

// AddEventHandler registers a WhatsApp event handler.
func (c *Client) AddEventHandler(handler func(evt interface{})) uint32 {
	return c.WA.AddEventHandler(handler)
}

// IsConnected returns true if the client is connected.
func (c *Client) IsConnected() bool {
	return c.WA.IsConnected()
}

// Disconnect disconnects from WhatsApp.
func (c *Client) Disconnect() {
	c.WA.Disconnect()
}

// OnConnected registers a handler called when connection is established.
func (c *Client) OnConnected(fn func()) {
	c.WA.AddEventHandler(func(evt interface{}) {
		if _, ok := evt.(*events.Connected); ok {
			fn()
		}
	})
}

// setupLogger creates a zerolog-based WhatsApp logger (unused but kept for future).
func setupLogger() zerolog.Logger {
	return log.With().Str("component", "whatsapp").Logger()
}

// printQR renders a scannable QR code in the terminal.
func printQR(code string) {
	fmt.Println("\n╔═══════════════════════════════════════╗")
	fmt.Println("║       WA PERSONA AI — SCAN QR CODE   ║")
	fmt.Println("╚═══════════════════════════════════════╝")
	qrterminal.GenerateWithConfig(code, qrterminal.Config{
		Level:     qrterminal.M,
		Writer:    os.Stdout,
		BlackChar: qrterminal.BLACK,
		WhiteChar: qrterminal.WHITE,
		QuietZone: 1,
	})
	fmt.Println("Open WhatsApp → Linked Devices → Link a Device → scan kode di atas")
}
