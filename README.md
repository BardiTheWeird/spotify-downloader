# Spotify Downloader

<a href="https://github.com/BardiTheWeird/spotify-downloader/tree/release-0.1.0">
    <img src="https://raw.githubusercontent.com/BardiTheWeird/spotify-downloader/release-0.1.0/frontend/resources/icon.ico"
         alt="Spotify Downloader logo" title="Spotify Downloader" height="100" width="100" />
</a>

---

An application for downloading of MP3 tracks using Spotify playlist data.

- [Spotify Downloader](#spotify-downloader)
  - [Introduction](#introduction)
  - [How to work with](#how-to-work-with)
    - [Query](#query)
    - [Directory Choice](#directory-choice)
    - [Download](#downloaded)
    - [Cancel Download](#cancel-download)


## Introduction

[Spotify Downloader](https://github.com/BardiTheWeird/spotify-downloader/tree/release-0.1.0) is an application which allows user to download desired musical tracks using [Spotify](https://open.spotify.com) playlist or album data from [YouTubeMusic](https://music.youtube.com/). Using the [Spotify's](https://open.spotify.com) metadata application gets a link to a [YouTubeMusic](https://music.youtube.com/) storage using [SongLink's](https://odesli.co) API. Also authorized user can download liked tracks.

## How to work with

First of all you'll need to install several extra applications:

- [youtube-dl](https://youtube-dl.org). This application is used to get download urls from [YouTube](https://youtube.com/) links.

- [FFMPEG](https://www.ffmpeg.org/download.html). Since downloaded tracks will be in MP4 format we'd need to convert them. Simple and easy tool for it is a FFMPEG, which is used by the [Spotify Downloder](https://github.com/BardiTheWeird/spotify-downloader/tree/release-0.1.0).

User can download tracks from public playlist without authorization but to load private playlist (e.g. Liked Tracks) user will need to log in [Spotify for developers](https://developer.spotify.com/dashboard/), create a formal application and then get a Client ID which will be used for authorization using Spotify OAuth interface. To load Liked tracks user must click on Get All Liked button in the top left corner near User Name.

### Query

To begin Query process User will need to insert a Playlist/Album link into the Query field and click "Submit" button. Window will refresh and show tracks included in the Playlist/Album in the table underneath the Query fields.

If Playlist link will be incorrect, the track-table won't be generated.

### Directory Choice

To chose a Directory for a download User will have 2 options:
 - Insert a directory into the field.
 - Choose a directory using "Browse" button.

### Preview

User can preview tracks availiable for download using Play button on top of the track's cover art.

### Download

To download tracks from the generated table User will need to choose the tracks User needs using the selectors. By default all tracks are chosen. To unchoose or choose all User can click on "All" selector.

To begin download process User will need to click "Download Selected" button. To cancel - "Cancel Download". Quantity of the simultaniously downloading tracks is limited to 10. After a complition of any track from the bunch, next one will start downloading.

In download process User will get notifications on Download Status:
- N/A - Default state, no informaion is acquired before the "Download Button" click.
- Invalid path - Track or Album name uses prohibited symbols (*will be handled later*)
- No YouTube Link - this track does not appear in YouTube database, therefore cannot be downloaded.
- Converting - track is downloaded in MP4 format and the convertation into the MP3 is going.
- Completed - track downloaded and converted.
- Convertation Failed - convertation of track is failed.
- Cancelled - User cancelled the download.

User may need to authorize again after a while due to the authorization expiration.

### Cancel Download

The download can be cancelled on any stage but "Completed" using the "Cancel Download" button.
