# Installation Guide — WA Persona AI

> Status: Draft
> Terakhir diperbarui: 2026-08-14
> Pemilik: Cloud-Dark

## 1. System Requirements

### Minimum

| Resource | Requirement |
|---|---|
| CPU | 1 vCPU |
| RAM | 2 GB |
| Disk | 10 GB |
| OS | Linux (Ubuntu 22.04+), macOS, Windows |
| Network | Stable internet connection |

### Recommended

| Resource | Requirement |
|---|---|
| CPU | 2 vCPU |
| RAM | 4 GB |
| Disk | 20 GB SSD |
| OS | Ubuntu 22.04 LTS |

## 2. Installation Methods

### Method A: Binary Release (Recommended)

```bash
# Download latest release
curl -L https://github.com/Cloud-Dark/wa-persona-ai/releases/latest/download/wa-persona-ai-linux-amd64 -o wa-persona-ai
chmod +x wa-persona-ai

# Create directories
mkdir -p config persona data/session data/memory/vectors logs

# Download example config
curl -L https://raw.githubusercontent.com/Cloud-Dark/wa-persona-ai/main/config/config.example.yaml -o config/config.yaml
curl -L https://raw.githubusercontent.com/Cloud-Dark/wa-persona-ai/main/persona/default.yaml -o persona/default.yaml

# Set API key
export CLAUDE_API_KEY="your-key-here"

# Run
./wa-persona-ai
```

### Method B: Docker (Recommended for Production)

```bash
# Clone repository
git clone https://github.com/Cloud-Dark/wa-persona-ai.git
cd wa-persona-ai

# Copy and edit config
cp config/config.example.yaml config/config.yaml
cp .env.example .env
# Edit .env with your API keys

# Start with Docker Compose
docker-compose up -d

# View logs
docker-compose logs -f wa-persona-ai

# Stop
docker-compose down
```

### Method C: Build from Source

```bash
# Prerequisites: Go 1.22+, Git, GCC (for SQLite)

# Clone
git clone https://github.com/Cloud-Dark/wa-persona-ai.git
cd wa-persona-ai

# Build
go build -o bin/wa-persona-ai ./cmd/wa-persona-ai/

# Setup config
cp config/config.example.yaml config/config.yaml
cp .env.example .env

# Run
./bin/wa-persona-ai
```

## 3. First-Time Setup

1. **Run the application** — QR code akan muncul di terminal
2. **Scan QR code** menggunakan WhatsApp di ponsel (Linked Devices → Link a Device)
3. **Verify connection** — terminal menampilkan "Connected successfully"
4. **Test** — kirim pesan ke nomor bot dari WhatsApp lain
5. **Configure persona** — edit `persona/default.yaml` sesuai kebutuhan

## 4. Systemd Service (Linux Production)

```ini
# /etc/systemd/system/wa-persona-ai.service
[Unit]
Description=WA Persona AI
After=network.target

[Service]
Type=simple
User=wabot
WorkingDirectory=/opt/wa-persona-ai
ExecStart=/opt/wa-persona-ai/wa-persona-ai
Restart=always
RestartSec=10
EnvironmentFile=/opt/wa-persona-ai/.env

[Install]
WantedBy=multi-user.target
```

```bash
# Enable and start
sudo systemctl enable wa-persona-ai
sudo systemctl start wa-persona-ai

# Check status
sudo systemctl status wa-persona-ai

# View logs
sudo journalctl -u wa-persona-ai -f
```

## 5. Backup & Restore

### Backup

```bash
# Backup session, memory, and config
tar -czf backup-$(date +%Y%m%d).tar.gz \
  data/session/ \
  data/memory/ \
  config/config.yaml \
  persona/
```

### Restore

```bash
# Restore from backup
tar -xzf backup-YYYYMMDD.tar.gz
```

## 6. Updating

```bash
# Binary: download new release and replace
# Docker: pull and restart
docker-compose pull
docker-compose up -d

# Source: pull and rebuild
git pull
go build -o bin/wa-persona-ai ./cmd/wa-persona-ai/
```

## 7. Referensi

- Lihat [16_CONFIG_REFERENCE.md](16_CONFIG_REFERENCE.md) untuk konfigurasi
- Lihat [17_DEVELOPER_SETUP.md](17_DEVELOPER_SETUP.md) untuk development
