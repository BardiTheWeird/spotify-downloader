import React from "react";

import { IsLoggedInContext, BaseUrlContex } from "../contexts";
import { authorizedFetch } from "../utilities";

import { PlaylistTable } from './PlaylistTable';

const { ipcRenderer } = window.require('electron');

export function InputBar() {
    const [isUserLogged, updateIsUserLogged] = React.useContext(IsLoggedInContext);
    const baseUrl = React.useContext(BaseUrlContex);
    const [formData, updateFormData] = React.useState();
    const [playlist, updatePlaylist] = React.useState();
    const [downloadPath, updateDownloadPath] = React.useState('');
  
    const submitPlaylistLink = async (e) => {
      e.preventDefault();
      if (!isUserLogged) {
        alert('Log in, please');
        return;
      }
      let response = await authorizedFetch(`${baseUrl}/spotify/playlist?link=${formData}`);
      if (response === null) {
        alert("YOU STILL DON'T HANDLE UNAUTHORIZED PLAYLIST SUBMIT (or your (refresh) tokens are ded, idk)");
      }
  
      switch (response.status) {
        case 200:
            let playlist = await response.json();
            updatePlaylist(playlist);
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
          alert('Somethign went wrong');
      }
    }
  
    return (
      <>
        <div className="Bar">
          <div>
            <div className="SearchBar">
              <form onSubmit={submitPlaylistLink} className="inputForm">
                <input placeholder='Spotify Link (https://open.spotify.com/playlist/etc...)' type="text" name='PL-URL' required className="inputForm inputformline" onChange={
                  e => updateFormData(e.target.value.trim())}/>
                <input type="submit" className="uselessButton" value="Submit"/>
              </form>
            </div>
            <div className="SearchBar">
              <form onSubmit={e => e.preventDefault()} className="inputForm">
                <input placeholder='Insert Download Directory' type="text" name='DL-path' required className="inputForm inputformline" onChange={
                  e => updateDownloadPath(e.target.value)}
                    value={downloadPath}
                  />
                <button className="uselessButton" onClick={async e => {
                  e.preventDefault();
                  const path = await ipcRenderer.invoke('openDirectory');
                  updateDownloadPath(path[0]);
                }}>Browse</button>
              </form>
            </div>
          </div>
        </div>
        { playlist &&
          <PlaylistTable playlist={playlist} downloadPath={downloadPath}/>
        }
      </>
    );
  }
