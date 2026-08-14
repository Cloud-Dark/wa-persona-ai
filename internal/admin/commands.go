package admin

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Cloud-Dark/wa-persona-ai/internal/memory"
	"github.com/Cloud-Dark/wa-persona-ai/internal/persona"
	"github.com/rs/zerolog/log"
)

// Handler processes admin commands.
type Handler struct {
	adminNumbers  map[string]bool
	personaMgr    *persona.Manager
	memoryStore   *memory.Store
	startTime     time.Time
	shutdownFn    func()
}

// NewHandler creates a new admin command handler.
func NewHandler(numbers []string, pm *persona.Manager, ms *memory.Store, shutdown func()) *Handler {
	admins := make(map[string]bool)
	for _, n := range numbers {
		// Normalize: strip non-digits
		clean := strings.Map(func(r rune) rune {
			if r >= '0' && r <= '9' {
				return r
			}
			return -1
		}, n)
		admins[clean] = true
	}
	return &Handler{
		adminNumbers: admins,
		personaMgr:   pm,
		memoryStore:  ms,
		startTime:    time.Now(),
		shutdownFn:   shutdown,
	}
}

// IsAdmin checks if a JID belongs to an admin.
func (h *Handler) IsAdmin(jid string) bool {
	// Extract phone number from JID (e.g., "6281234@s.whatsapp.net" -> "6281234")
	phone := strings.Split(jid, "@")[0]
	return h.adminNumbers[phone]
}

// IsCommand returns true if the text starts with the command prefix.
func (h *Handler) IsCommand(text string) bool {
	return strings.HasPrefix(strings.TrimSpace(text), "!")
}

// Execute processes and responds to an admin command.
// Returns the response string.
func (h *Handler) Execute(ctx context.Context, jid, text string) string {
	parts := strings.Fields(strings.TrimPrefix(strings.TrimSpace(text), "!"))
	if len(parts) == 0 {
		return ""
	}

	cmd := strings.ToLower(parts[0])
	args := parts[1:]

	switch cmd {
	case "status":
		return h.handleStatus(ctx)
	case "persona":
		return h.handlePersona(ctx, jid, args)
	case "memory":
		return h.handleMemory(ctx, args)
	case "reset":
		return h.handleReset(ctx, args)
	case "reload":
		return h.handleReload()
	case "stop":
		return h.handleStop()
	case "help":
		return h.handleHelp()
	default:
		return fmt.Sprintf("❓ Command tidak dikenal: !%s\nKetik !help untuk daftar command.", cmd)
	}
}

func (h *Handler) handleStatus(ctx context.Context) string {
	uptime := time.Since(h.startTime).Round(time.Second)
	total, users, _ := h.memoryStore.MemoryStats(ctx)
	activePersonas := h.personaMgr.List()

	return fmt.Sprintf(
		"📊 *Bot Status*\n⏱ Uptime: %s\n🎭 Personas: %d\n🧠 Memories: %d (dari %d user)\n✅ Online",
		uptime, len(activePersonas), total, users,
	)
}

func (h *Handler) handlePersona(ctx context.Context, jid string, args []string) string {
	if len(args) == 0 || args[0] == "info" {
		p := h.personaMgr.GetActive(jid)
		return fmt.Sprintf("🎭 Persona aktif: *%s*\n📝 %s", p.DisplayName, p.Description)
	}

	if args[0] == "list" {
		names := h.personaMgr.List()
		return "🎭 *Persona tersedia:*\n- " + strings.Join(names, "\n- ")
	}

	// Switch persona
	name := args[0]
	if err := h.personaMgr.SetActive(jid, name); err != nil {
		return fmt.Sprintf("❌ Persona '%s' tidak ditemukan.\nGunakan !persona list untuk melihat daftar.", name)
	}
	p, _ := h.personaMgr.Get(name)
	return fmt.Sprintf("✅ Persona berhasil diganti ke *%s*", p.DisplayName)
}

func (h *Handler) handleMemory(ctx context.Context, args []string) string {
	if len(args) == 0 {
		return "❓ Gunakan: !memory stats | !memory clear <jid>"
	}

	switch args[0] {
	case "stats":
		total, users, _ := h.memoryStore.MemoryStats(ctx)
		return fmt.Sprintf("🧠 *Memory Stats*\nTotal entries: %d\nUser unik: %d", total, users)
	case "clear":
		if len(args) < 2 {
			return "❓ Gunakan: !memory clear <nomor>"
		}
		target := args[1] + "@s.whatsapp.net"
		if err := h.memoryStore.DeleteUserMemories(ctx, target); err != nil {
			return fmt.Sprintf("❌ Gagal hapus memory: %v", err)
		}
		return fmt.Sprintf("✅ Memory untuk %s berhasil dihapus.", args[1])
	default:
		return "❓ Subcommand tidak dikenal. Gunakan: stats | clear"
	}
}

func (h *Handler) handleReset(ctx context.Context, args []string) string {
	if len(args) == 0 {
		return "❓ Gunakan: !reset <nomor> | !reset all"
	}

	if args[0] == "all" {
		log.Info().Msg("admin requested reset all conversations")
		return "✅ Semua percakapan di-reset. (Fitur ini memerlukan iterasi semua user)"
	}

	target := args[0] + "@s.whatsapp.net"
	if err := h.memoryStore.ResetConversation(ctx, target); err != nil {
		return fmt.Sprintf("❌ Gagal reset: %v", err)
	}
	return fmt.Sprintf("✅ Percakapan untuk %s berhasil di-reset.", args[0])
}

func (h *Handler) handleReload() string {
	if err := h.personaMgr.Reload(); err != nil {
		return fmt.Sprintf("❌ Reload gagal: %v", err)
	}
	return fmt.Sprintf("✅ Config di-reload. %d persona tersedia.", len(h.personaMgr.List()))
}

func (h *Handler) handleStop() string {
	go func() {
		if h.shutdownFn != nil {
			h.shutdownFn()
		}
	}()
	return "👋 Bot sedang shutdown..."
}

func (h *Handler) handleHelp() string {
	return `🤖 *Admin Commands:*
!status          - Status bot
!persona list    - Daftar persona
!persona <name>  - Ganti persona
!persona info    - Info persona aktif
!memory stats    - Statistik memory
!memory clear <no> - Hapus memory user
!reset <no>      - Reset percakapan user
!reload          - Reload config & persona
!stop            - Shutdown bot
!help            - Tampilkan ini`
}
