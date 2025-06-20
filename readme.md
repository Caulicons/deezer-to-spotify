# 🎶 Deezer ➡️ Spotify Playlist Migrator 🚀

Easily transfer your public Deezer playlists to Spotify with this Go application!  
Fetch all your Deezer tracks, get their info, and recreate your playlist on Spotify—automatically! 🎧

---

## ✨ Features

- 🔗 **Fetch Deezer Playlist Tracks**  
  Input your Deezer playlist URL and grab all track IDs!

- 🏷️ **Get Track Info**  
  Retrieve ISRC and title for each track from Deezer.

- 💾 **Save as JSON**  
  Store track IDs and info for easy access.

- 🔐 **Spotify OAuth2 Authentication**  
  Securely connect your Spotify account.

- 🆕 **Create Spotify Playlist**  
  Name your new playlist and fill it with your favorite tracks!

- 📦 **Track Results**  
  See which tracks were found and which weren’t.

---

## 🗂️ Project Structure

```
cmd/
  api/
  app/
    main.go
    config/
    dependencies/
config/
data/
  deezer/
    track_id.json
    track_info.json
  spotify/
    track_uri.json
    tracks_not_found.json
internal/
  business/
    deezer/
      usecase/
    spotify/
      usecase/
  domain/
    contract/
    entities/
  infra/
    http/
pkg/
  http/
  utils/
    anyToStruct.go
    json/
```

---

## 🚦 How It Works

1. **🔗 Deezer:**  
   Paste your Deezer playlist URL (e.g. `https://api.deezer.com/user/<user_id>/tracks`)  
   ⬇️  
   All track IDs are saved to `data/deezer/track_id.json`

2. **🏷️ Track Info:**  
   The app fetches ISRC and title for each track  
   ⬇️  
   Info saved to `data/deezer/track_info.json`

3. **🔐 Spotify:**  
   Authenticate via OAuth2 (browser consent)  
   ⬇️  
   Name your new playlist

4. **🆕 Playlist Creation:**  
   The app creates your playlist and adds tracks by ISRC/title  
   ⬇️  
   Spotify track URIs saved to `data/spotify/track_uri.json`  
   Tracks not found: `data/spotify/tracks_not_found.json`

---

## 🛠️ Getting Started

1. **Clone the repo:**  
   ```sh
   git clone https://github.com/yourusername/deezer-to-spotify.git
   cd deezer-to-spotify
   ```

2. **Configure:**  
   Add your Deezer and Spotify credentials to `.env` following the `.env.development` suggests.

3. **Run the app:**  
   ```sh
   go run cmd/app/main.go
   ```

4. **Follow the prompts and enjoy your migrated playlist! 🎉**

---

## 📋 Requirements

- Go 1.18+
- Deezer account with public playlists
- Spotify developer account (for OAuth2 credentials)

---

## 🤝 Contributing

Pull requests are welcome!  
For major changes, please open an issue first.  
Let’s make playlist migration even better! 🚀

---

## 📄 License

MIT

---

_Made with ❤️ for music lovers!_