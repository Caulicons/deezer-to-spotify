

# 🎵 Deezer to Spotify Playlist Converter 🎵

An interactive CLI tool that helps you transfer your favorite tracks from **Deezer** to **Spotify** with ease!

---

## ✨ Features

- 🖥️ **Interactive Terminal Interface** – Step-by-step CLI to guide you through the migration
- 🎵 **Deezer Playlist Support** – Import from "Loved" tracks or any public playlist
- 🔍 **Smart Track Matching** – Uses ISRC codes for accurate Spotify track mapping
- 🔐 **Seamless OAuth Authentication** – Spotify login via browser
- 📋 **Playlist Management** – List, create, and modify Spotify playlists
- 🚀 **Batch Processing** – Transfer entire playlists in one command
- 📊 **Transfer Statistics** – See what was successfully transferred and what wasn't

---

## 🚀 Getting Started

### ✅ Prerequisites

- Go 1.18 or higher
- A [Spotify Developer Account](https://developer.spotify.com/)
- Spotify API credentials (Client ID & Secret)
- A Deezer account with playlists

---

### 📦 Installation

```bash
# 1. Clone the repository
git clone https://github.com/your-username/deezer-to-spotify.git
cd deezer-to-spotify

# 2. Create your .env file with your Spotify credentials
touch .env
# Add SPOTIFY_CLIENT_ID and SPOTIFY_CLIENT_SECRET to .env

# 3. Build the application
go build -o deezer2spotify ./cmd/shell

# 4. Run the application
./deezer2spotify
````

---

## 🔧 How It Works

1. **Enter Deezer Playlist Details**

   * Paste a playlist URL or use your "Loved" tracks
   * The app fetches track metadata via the Deezer API

2. **Authenticate with Spotify**

   * A browser tab opens automatically
   * Log in and authorize the app securely with Spotify

3. **Manage Your Music**

   * List your Spotify playlists
   * Create a new one
   * Choose a playlist to receive the Deezer tracks

4. **Transfer Tracks**

   * The app searches Spotify for each Deezer track (using ISRC when possible)
   * Shows what was found or missed
   * Transfers all matched tracks to the selected playlist

---

## 📁 Project Structure

```text
cmd/
  shell/           # CLI main app
  api/             # (optional) API integration
internal/
  business/        # Use cases for Spotify and Deezer
  domain/entities/ # Core domain models (e.g., Token, Track)
  infra/http/      # HTTP layer (Spotify auth callback etc.)
pkg/
  jsonUtils/       # JSON file helpers
.env               # Your Spotify credentials
```

---

## 🤔 Troubleshooting

* **Authentication Issues**: Make sure your Spotify `CLIENT_ID` and `CLIENT_SECRET` are correctly set in `.env`
* **Missing Tracks**: Some tracks may not be available on Spotify or have slightly different metadata
* **API Rate Limits**: Spotify may rate-limit you if transferring very large playlists quickly

---

## 🤝 Contributing

Contributions are welcome!

```bash
# Fork the repository
# Create your feature branch
git checkout -b feature/amazing-feature

# Commit your changes
git commit -m 'Add some amazing feature'

# Push to your branch
git push origin feature/amazing-feature

# Open a Pull Request
```

---

## 📄 License

This project is licensed under the [MIT License](LICENSE).


