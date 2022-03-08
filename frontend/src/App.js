import logo from './logo.svg';
import './App.css';
import React from 'react';

export function App() {
  return (
    <div className="App">
      <header className="App-header">
        <div>Please make sure you installed</div>
        <a href='https://youtube-dl.org' className="App-link">"YOUTUBE downloader"</a>
        <div>Prior to beginning the search</div>
        <div>Enter Spotify Playlist URL:</div>
      </header>
    </div>
    
  );
}

export function InputBar() {
  const [formData, updateFormData] = React.useState();
  const [playlist, updatePlaylist] = React.useState();

  const handleChange = (e) => {
    updateFormData(e.target.value.trim());
  }

  const handleSubmit = (e) => {
    e.preventDefault();
    fetch("http://localhost:8080/api/v1/spotify/playlist?link=" + formData)
      .then(async response => {
        updatePlaylist(await response.json())
      })
  }
  
  return (
    <>
      <div className="Bar">
        <div>
          <div className="SearchBar">
            <form onSubmit={handleSubmit} class="inputForm">
              <input type="text" name='PL-URL' required class="inputForm" onChange={handleChange}/>
              <input type="submit"/>
            </form>
          </div>
        </div>
      </div>
      { playlist &&
        <PlaylistTable playlist={playlist}/>
      }
    </>
  );
}

export default App;

export function PlaylistTable({playlist}) {
  const [checkedTracks, updateCheckedTracks] = 
    React.useState(Array(playlist.tracks.length).fill(true));

  return (
    <>
      <button className='DonwloadButton' onClick={() => {
          const tracksToDownload = playlist.tracks.filter((_, i) => checkedTracks[i]);
          const tracksToDownloadId = tracksToDownload.map(track => track.id)

          console.log(tracksToDownloadId);
        }
      }
      >Download selected
      </button>
      <table className='Table'>
        <tr>
          <th></th>
          <th>Logo</th>
          <th>Artist</th>
          <th>Track Name</th>
          <th>Album</th>
          <th>Alailiable/Status</th>
        </tr>
        {
          playlist.tracks.map((track, index) =>
            (
              <tr>
                <td><input type="checkbox" checked={checkedTracks[index]} onChange={
                  e => {
                    const clonedCheckedTracks = [...checkedTracks];
                    clonedCheckedTracks[index] = !clonedCheckedTracks[index];
                    updateCheckedTracks(clonedCheckedTracks);
                  }
                }/></td>
                <td><img src={track.album_image}
                height="30" px/>
                </td>
                <td>{track.artists}</td>
                <td>{track.title}</td>
                <td>{track.album_title}</td>
                <td>YouTubeAvailiable?</td>
              </tr>
            )
          )
        }
      </table>
    </>
  )
}
