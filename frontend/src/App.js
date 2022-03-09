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
      <header>
      <div className='App-header-info'>If directory is not defined PL will be loaded into "userDownloads" folder</div>
      </header>
    </div>
    
  );
}

export function InputBar() {
  const [formData, updateFormData] = React.useState();
  const [playlist, updatePlaylist] = React.useState();
  const [downloadPath, updateDownloadPath] = React.useState();

  const submitPlaylistLink = async (e) => {
    e.preventDefault();
    let response = await fetch("http://localhost:8080/api/v1/spotify/playlist?link=" + formData);
    let playlist = await response.json();
    updatePlaylist(playlist);
    console.log(response)
  }
  
  return (
    <>
      <div className="Bar">
        <div>
          <div className="SearchBar">
            <form onSubmit={submitPlaylistLink} class="inputForm">
              <input placeholder='Spotify Playlist Link' type="text" name='PL-URL' required class="inputForm" onChange={
                e => updateFormData(e.target.value.trim())}/>
              <input type="submit"/>
            </form>
          </div>
          <div className="SearchBar">
            <form class="inputForm">
              <input placeholder='Download Directory' type="text" name='DL-path' required class="inputForm" onChange={
                e => updateDownloadPath(e.target.value)}/>
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

export default App;

export function PlaylistTable({playlist, downloadPath}) {
  const [tracks, updateTracks] = React.useState(playlist.tracks.map(track => {
    return {...track,
      checked: true,
      status: "NA"
    };
  }));
  // WTF IS THIS AND WHY DOES IT WORK
  React.useEffect(() => {updateTracks(playlist.tracks.map(track => {
    return {...track,
      checked: true,
      status: "NA"
    };
  }))}, [playlist]);

  return (
    <>
      <div className='inline-buttons'>
        <button className='DownloadButton' onClick={() => {
            tracks.forEach(async (track, index) => {
              if (!track.checked) {
                return;
              }
              let url = 'http://localhost:8080/api/v1/download/start';
              let downloadFolder = downloadPath;
              if (!downloadFolder) {
                downloadFolder = "./userDownloads/"
              }
              let downloadResponse = await fetch(url, {
                method: 'POST',
                body: JSON.stringify({
                  id: track.id,
                  folder: downloadPath,
                  filename: `${track.artists} - ${track.title}`,
                  title: track.title,
                  artist: track.artists.join(' '),
                  album: track.album_title,
                  image: track.album_image
                })
              });
              console.log(downloadResponse);
            })
          }
        }
        >Download selected
        </button>
        <button className='CancelDownloadButton' disabled={true}>Cancel Download</button>
      </div>
      
      <table className='Table'>
        <tr>
          <th>
              <input type="checkbox" class="checkmark" onChange={
                  e => {
                    const isChecked = e.target.checked;
                    
                    const copiedTracks = [...tracks];
                    copiedTracks.forEach(
                      track => track.checked = isChecked
                    )
                    updateTracks(copiedTracks);
                  }
                }
              />
          </th>
          <th>All</th>
          <th>Artist</th>
          <th>Track Name</th>
          <th>Album</th>
          <th>Status</th>
        </tr>
        {
          tracks.map((track, index) =>
            (
              <tr>
                <td><input type="checkbox" class="checkmark" checked={track.checked} onChange={
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
