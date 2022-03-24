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
	tracks: Track[]	
}

type Track = {
	id: string,
	title: string,
	artists: string[],
	album_title: string,
	album_image: string,
}
```
# Endpoints
Everything starts with /api/v1
- `GET /spotify/playlist?id={spotify_playlist_id}`
	- `GET /spotify/playlist?link={spotify_playlist_link}`
	- Requires an `Authorization` header with a Spotify access token.
	- Returns a Playlist entity for a specified Spotify Playlist Id or Spotify Playlist Link. When both `id` and `link` are provided, `id` takes precedence.
	- Status codes:
		- 200 + playlist payload
		- 400 + error payload => invalid link or "id" is empty
		- 401 => not authorized (maybe?)
		- 404 => no playlist with such id
		- 429 => too many requests
		- 500
- `POST /download/start`
	- Request Body:
	```
	type DownloadRequest {
		id: string,
		folder: string,
		filename: string,
		
		title: string,
		artist: string,
		album: string,
		image: string
	}
	```
	- Starts a download on a host machine at `folder_path/file_name`
		- If ffmpeg is detected, it also converts the downloaded file to .mp3 with provided metadata
	- Status codes:
		- 204
		- 400 => no songlink entry for song with such id (most likely, the id sent was wrong)
		- 403 => can't create a file at filepath
		- 404 =>
			- no youtube link for song with such id
			- no download link for youtube link (youtube api and/or youtube-dl weirdness)
		- 429 => songlink too many requests
		- 500 => songlink/download error sending request
- `GET /download/status?id={trackId}`
	- Response Model:
```
	type DownloadEntry = {
		total_bytes: int,
		downloaded_bytes: int?,
		status: DownloadStatus
	}
	
	enum DownloadStatus = {
		DownloadInProgress,
		DownloadConvertationInProgress,
		DownloadFinished,
		DownloadErrorConverting,
		DownloadFailed,
		DownloadedCancelled
	}
```
	- Status codes:
		- 200
		- 400 + payload error => path not provided
		- 404 => no download at {path}
		- 500 => can't stat file OR unhandled GetDownloadStatusStatus
- `POST /download/cancel?id={trackId}`
	- Cancels a download at `folder_path/file_name`
	- Status codes:
		- 204
		- 400 + payload error => path not provided
		- 404 => no download at path
		- 405 => used method other than POST
		- 409 => not in progress
		- 500 => cancellation status not handled by the server
- `GET /features`
	- Returns available/installed features
	- Status code 200
	- Response Body:
	```
	type FeaturesAvailable{
		youtube_dl: bool,
		ffmpeg: bool
	}
	```