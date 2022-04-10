import React from "react";

import { usePlaylist, useSubmitPlaylistLink } from "../services/PlaylistService";
import { PlaylistTable } from './PlaylistTable';

const { ipcRenderer } = (window.require && window.require('electron')) || (window.opener && window.opener.require('electron'));

export function InputBar() {
    const [formData, updateFormData] = React.useState();
    const [downloadPath, updateDownloadPath] = React.useState('');
    const playlist = usePlaylist();
    const submiPlaylistLink = useSubmitPlaylistLink();

    return (
      <>
        <div className="Bar">
          <div>
            <div className="SearchBar">
              <form onSubmit={e => {
                e.preventDefault();
                submiPlaylistLink(formData);
              }} className="inputForm">
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
        { playlist.length > 0 &&
          <PlaylistTable playlist={playlist} downloadPath={downloadPath}/>
        }
      </>
    );
  }
