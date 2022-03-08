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

  const handleSubmit = async (e) => {
    e.preventDefault();
    let response = await fetch("http://localhost:8080/api/v1/spotify/playlist?link=" + formData);
    let playlist = await response.json();
    updatePlaylist(playlist);
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
  const [tracks, updateTracks] = React.useState(playlist.tracks.map(track => {
    return {...track,
      checked: true,
      status: "NA"
    };
  }));

  return (
    <>
      <button className='DonwloadButton' onClick={() => {
          const tracksToDownload = tracks.filter(track => !track.checked);
          const tracksToDownloadId = tracksToDownload.map(track => track.id)
          console.log(tracksToDownload)
        }
      }
      >Download selected
      </button>
      <table className='Table'>
        <tr>
          <th></th>
          <th></th>
          <th>Artist</th>
          <th>Track Name</th>
          <th>Album</th>
          <th>Status</th>
        </tr>
        {
          tracks.map((track, index) =>
            (
              <tr>
                <td><input type="checkbox" checked={track.checked} onChange={
                  e => {
                    const clonedTracks = [...tracks];
                    clonedTracks[index].checked = !clonedTracks[index].checked;
                    updateTracks(clonedTracks);
                  }
                }/></td>
                <td><img src={track.album_image}
                height="30" px/>
                </td>
                <td>{track.artists}</td>
                <td>{track.title}</td>
                <td>{track.album_title}</td>
                <td>{track.status}</td>
              </tr>
            )
          )
        }
      </table>
    </>
  )
}
