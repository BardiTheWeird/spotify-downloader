import React from "react";
import { Howl } from "howler";
import { usePlaylist } from "./PlaylistService";

const PlayerContext = React.createContext();

export function usePlayPause() {
    return React.useContext(PlayerContext)[0];
}

export function PlayerProvider({children}) {
    const [loadedTracks, updateLoadedTracks] = React.useState();
    const [curPlayingIndex, updateCurPlayingIndex] = React.useState(null);
    const spotifyPlaylist = usePlaylist();

    React.useEffect(() => {
        updateLoadedTracks(
            spotifyPlaylist.map(x => {
                return new Howl({
                  src: x.preview_url,
                  html5: true
              })})
        )
    }, [spotifyPlaylist])

    // returns pausedTrackIndex: number | null
    function playPause(index) {
        let pausedTrackIndex = null;
        if (curPlayingIndex !== null) {
            loadedTracks[curPlayingIndex].pause();
            pausedTrackIndex = curPlayingIndex;
        }
    
        if (curPlayingIndex !== index) {
            loadedTracks[index].play();
            updateCurPlayingIndex(index);
        }
        else {
            updateCurPlayingIndex(null);
        }

        return pausedTrackIndex;
    }

    return <PlayerContext.Provider value={[playPause]}>
        {children}
    </PlayerContext.Provider>
}
