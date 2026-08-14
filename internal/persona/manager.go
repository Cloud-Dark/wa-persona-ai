package persona

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/rs/zerolog/log"
)

// Manager handles loading and switching personas.
type Manager struct {
	mu          sync.RWMutex
	personas    map[string]*Persona
	userPersona map[string]string // userJID -> persona name
	defaultName string
	dir         string
}

// NewManager creates a new persona manager and loads all personas from dir.
func NewManager(dir, defaultName string) (*Manager, error) {
	m := &Manager{
		personas:    make(map[string]*Persona),
		userPersona: make(map[string]string),
		defaultName: defaultName,
		dir:         dir,
	}

	if err := m.LoadAll(); err != nil {
		log.Warn().Err(err).Msg("failed to load personas, using built-in default")
		m.personas["default"] = DefaultPersona()
	}

	if _, ok := m.personas[defaultName]; !ok {
		if _, ok := m.personas["default"]; !ok {
			m.personas["default"] = DefaultPersona()
		}
		m.defaultName = "default"
	}

	return m, nil
}

// LoadAll reads all YAML files from the persona directory.
func (m *Manager) LoadAll() error {
	entries, err := os.ReadDir(m.dir)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("persona directory not found: %s", m.dir)
		}
		return fmt.Errorf("read persona dir: %w", err)
	}

	loaded := 0
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".yaml") || strings.HasPrefix(e.Name(), "_") {
			continue
		}

		path := filepath.Join(m.dir, e.Name())
		p, err := Load(path)
		if err != nil {
			log.Warn().Err(err).Str("file", e.Name()).Msg("skipping invalid persona file")
			continue
		}

		m.mu.Lock()
		m.personas[p.Name] = p
		m.mu.Unlock()
		loaded++
		log.Debug().Str("name", p.Name).Str("file", e.Name()).Msg("persona loaded")
	}

	// Also load from examples/ subdirectory
	examplesDir := filepath.Join(m.dir, "examples")
	if entries2, err := os.ReadDir(examplesDir); err == nil {
		for _, e := range entries2 {
			if e.IsDir() || !strings.HasSuffix(e.Name(), ".yaml") {
				continue
			}
			path := filepath.Join(examplesDir, e.Name())
			p, err := Load(path)
			if err != nil {
				log.Warn().Err(err).Str("file", e.Name()).Msg("skipping invalid example persona")
				continue
			}
			m.mu.Lock()
			m.personas[p.Name] = p
			m.mu.Unlock()
			loaded++
		}
	}

	log.Info().Int("count", loaded).Msg("personas loaded")
	return nil
}

// Get returns a persona by name.
func (m *Manager) Get(name string) (*Persona, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	p, ok := m.personas[name]
	if !ok {
		return nil, fmt.Errorf("persona not found: %s", name)
	}
	return p, nil
}

// GetActive returns the active persona for a user.
func (m *Manager) GetActive(userJID string) *Persona {
	m.mu.RLock()
	name, ok := m.userPersona[userJID]
	m.mu.RUnlock()

	if !ok {
		name = m.defaultName
	}

	p, err := m.Get(name)
	if err != nil {
		p, _ = m.Get(m.defaultName)
		if p == nil {
			return DefaultPersona()
		}
	}
	return p
}

// SetActive assigns a persona to a user.
func (m *Manager) SetActive(userJID, name string) error {
	if _, err := m.Get(name); err != nil {
		return err
	}
	m.mu.Lock()
	m.userPersona[userJID] = name
	m.mu.Unlock()
	return nil
}

// List returns names of all available personas.
func (m *Manager) List() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	names := make([]string, 0, len(m.personas))
	for name := range m.personas {
		names = append(names, name)
	}
	return names
}

// Reload reloads all personas from disk.
func (m *Manager) Reload() error {
	m.mu.Lock()
	m.personas = make(map[string]*Persona)
	m.mu.Unlock()
	return m.LoadAll()
}
