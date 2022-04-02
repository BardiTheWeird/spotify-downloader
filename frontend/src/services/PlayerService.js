import React from "react";
import { Howl } from "howler";

const PlayerContext = React.createContext();

export function usePlayPause() {
    const [playPause] = React.useContext(PlayerContext);
    return playPause;
}

export function usePlayerLoadTracks() {
    const [, loadTracks] = React.useContext(PlayerContext);
    return loadTracks;
}

export function PlayerProvider({children}) {
    const [loadedTracks, updateLoadedTracks] = React.useState();
    const [curPlayingIndex, updateCurPlayingIndex] = React.useState(null);

    function loadTracks(spotifyPlaylist) {
        updateLoadedTracks(
            spotifyPlaylist.map(x => {
                return new Howl({
                  src: x.preview_url,
                  html5: true
              })})
        )
    }

    function playPause(index) {
        if (curPlayingIndex !== null) {
            loadedTracks[curPlayingIndex].pause();
        }
    
        if (curPlayingIndex !== index) {
            loadedTracks[index].play();
            updateCurPlayingIndex(index);
        }
        else {
            updateCurPlayingIndex(null);
        }
    }

    return <PlayerContext.Provider value={[playPause, loadTracks]}>
        {children}
    </PlayerContext.Provider>
}
