import React from "react";

import { useBaseUrl } from "../services/BaseUrlService";
import { authorizedFetch } from "../utilities";

const PlaylistContext = React.createContext([]);

export function usePlaylist() {
    return React.useContext(PlaylistContext)[0];
}

export function useSubmitPlaylistLink() {
    return React.useContext(PlaylistContext)[1];
}

export function PlaylistProvider({children}) {
    const [playlist, updatePlaylist] = React.useState([]);
    const baseUrl = useBaseUrl();

    async function submitPlaylistLink(spotifyLink) {
        let response = await authorizedFetch(`${baseUrl}/spotify/playlist?link=${spotifyLink}`);
        switch (response.status) {
          case 200:
              let playlist = await response.json();
              updatePlaylist(playlist.tracks);
            break;
          case 400:
            alert('Bad Spotify link');
            break;
          case 401:
            alert('Log in, please');
            break;
          case 404:
            alert('No playlist or album with such id');
            break;
          case 429:
          case 500:
            alert('Something went wrong');
        }
    }
    return (
        <PlaylistContext.Provider value={[playlist, submitPlaylistLink]}>
            {children}
        </PlaylistContext.Provider>
    )
}