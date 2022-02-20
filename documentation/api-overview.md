# (Data Transfer) Model
## Error payload
It is sent with some 4xx status codes for addititonal error information
```
type ErrorPayload = {
	status_code: int,
	error_message: string
}
```
## Playlist
Stripped-down Spotify Playlist type
```
type Playlist = {
	id: string,
	name: string,
	owner: PlaylistOwner,
	href: string,
	image: string,
	tracks: Track	
}

type PlaylistOwner = {
	id: string,
	display_name: string,
	href: string
}

type Track = {
	id: string,
	title: string,
	artists: string[],
	album_title: string,
	album_image: string,
	album_href: string,
	href: string,
	preview_url: string
}
```
## DownloadLink
Is sent as a response when found a link to download a Spotify track from
```
type DownloadLink = {
	spotify_id: string,
	link: string
}
```
## DownloadEntry
Is sent as a response when getting a status of a download
```
type DownloadEntry = {
	path: string,
	youtube_link: string,
	total_bytes: int,
	downloaded_bytes: int,
	status: DownloadStatus
}

enum DownloadStatus = {
	DownloadInProgress,
	DownloadFinished,
	DownloadFailed,
	DownloadedCancelled
}
```
# Endpoints
Everything starts with /api/v1
- `GET /playlist?id={spotify_playlist_id}`
	- Returns a Playlist entity for a specified Spotify Id
	- Status codes:
		- 200 + playlist payload
		- 400 + error payload => "id" is empty
		- 401 => not authorized (maybe?)
		- 404 => no playlist with such id
		- 429 => too many requests
		- 500
- `GET /s2y?id={spotify_song_id}`
	- Returns a YouTube link for a given Spotify song Id
	- Status codes:
		- 200 + downloadLink payload
		- 400 => 'id' is empty
		- 404 => (no such id / no yt link) + error payload:
			- 400 => no entry for song with {id}
			- 404 => no YouTube link for song with {id}
		- 500
- `POST /download/start?path={host_path}&link={youtube_link}`
	- Starts a download on a host machine
	- Status codes:
		- 204
		- 400 + error payload => query parameter error
			- "path" is empty
			- "link" is empty
		- 404 => youtube-dl couldn't find a download link
		- 405 => used method other than POST
		- 500 => youtube-dl execution error
- `GET /download/status?path={host_path}`
	- Returns a DownloadEntry of download at {path}
	- Status codes:
		- 200
		- 400 + payload error => path not provided
		- 404 => no download at {path}
		- 500 => can't stat file OR unhandled GetDownloadStatusStatus
- `POST /download/cancel?path={host_path}`
	- Cancels a download at {path}
	- Status codes:
		- 204
		- 400 + payload error => path not provided
		- 404 => no download at path
		- 405 => used method other than POST
		- 409 => not in progress
		- 500 => cancellation status not handled by the server