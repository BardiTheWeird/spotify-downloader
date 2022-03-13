import logo from './logo.svg';
import './App.css';
import React from 'react';


export function App() {
  const [isLight, updateisLight] = React.useState(true);
  function LightDark() {
    
    if (!isLight) {
      return "Light"
    }
    else {
      return "Dark"
    }
  }
  return (
    <div className={`App ${LightDark()}`}>
      <div  className='App-header-info'>Light/Dark</div>
      <label class="switch">
        <input type="checkbox" onChange={e => updateisLight(!isLight)}></input>
        <span class="slider round"></span>
      </label>
      <header className="App-header">
        <div>Enter Spotify Playlist Or Album URL:</div>
      </header>
      <header>
      <div className='App-header-info'>If directory is not defined PL will be loaded into "userDownloads" folder</div>
      </header>
      <InputBar />
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
              <input type="submit" value="Submit"/>
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
  // a workaround for forcing a rerender
  const [, setForceUpdate] = React.useState(Date.now());
  const isDownloading = React.useRef(false);
  const downloadCounter = React.useRef(0);
  const [tracks, updateTracks] = React.useState(playlist.tracks.map(track => {
    return {...track,
      checked: true,
      status: "NA"
    };
  }));

  React.useEffect(() => {updateTracks(playlist.tracks.map(track => {
    return {...track,
      checked: true,
      status: "NA"
    };
  }))}, [playlist]);

  function DownloadSelected() {
    // set pending status
    tracks.forEach(track => { 
      if (track.checked) {
        track.status = 'Pending'
        downloadCounter.current++;
      }
    })
    updateTracks(tracks);

    // start download
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
      if (downloadResponse.status !== 204) {
        downloadCounter.current--;
        if (downloadCounter.current == 0) {
          isDownloading.current = false;
          console.log("finished downloading");
          setForceUpdate(Date.now());
        }
      }
      switch (downloadResponse.status) {
        case 204:
          UpdateDownloadStatus(index);
          if (!isDownloading.current) {
            CancelDownload(index);
          }
          else {
            track.status = 'Downloading';
          }
          break;
        case 400:
          track.status = "Not Available";
          break;
        case 403:
          track.status = "Invalid path";
          break;
        case 404:
          track.status = "No YouTube Link";
          break;
        case 500:
          track.status = "Download Error";
          break;           
      }
      
      const copiedTracks = [...tracks];
      updateTracks(copiedTracks);
    })
  }

  async function UpdateDownloadStatus(trackIndex) {
      while (true) {
        const copiedTracks = [...tracks];
        const track = copiedTracks[trackIndex];

        let downloadFolder = downloadPath;
        if (!downloadFolder) {
          downloadFolder = "./userDownloads/"
        }
        const getStatusResponse = await fetch(`http://localhost:8080/api/v1/download/status?folder=${downloadFolder}&filename=${track.artists} - ${track.title}`);
        const downloadEntry = await getStatusResponse.json();

        switch (downloadEntry.status) {
          case 0: 
            const percentage = Math.floor(downloadEntry.downloaded_bytes/downloadEntry.total_bytes*100);
            copiedTracks[trackIndex].status = `Downloading ${percentage}%`;
            break;
          case 1: 
            copiedTracks[trackIndex].status = 'Converting'
            break;
          case 2: 
            copiedTracks[trackIndex].status = 'Completed'
            break;
          case 3: 
            copiedTracks[trackIndex].status = 'Convertation Failed'
            break;
          case 4: 
            copiedTracks[trackIndex].status = 'Failed'
            break;
          case 5: 
            copiedTracks[trackIndex].status = 'Cancelled'
            break;
        }
        updateTracks(copiedTracks);
        // if (SwitchButtonAfterDownload() === downloadEntry.length) {
        //   isDownloading.current = !isDownloading.current
        // }
        // sleep 1s
        await new Promise(r => setTimeout(r, 1000));
        if (downloadEntry.status >= 2) {
          downloadCounter.current--;
          if (downloadCounter.current == 0) {
            isDownloading.current = false;
            console.log("finished downloading");
            setForceUpdate(Date.now());
          }
          console.log(downloadCounter.current)
          break;
        }
      }
  }

  function SwitchIsDownloading () {
    isDownloading.current = !isDownloading.current;
  }

  async function CancelDownload(trackIndex) {
    const copiedTracks = [...tracks];
    const track = copiedTracks[trackIndex];
    if (!track.checked) {
      return;
    }
    let downloadFolder = downloadPath;
    if (!downloadFolder) {
      downloadFolder = "./userDownloads/"
    }
    let url = `http://localhost:8080/api/v1/download/cancel?folder=${downloadFolder}&filename=${track.artists} - ${track.title}`;
    let downloadResponse = await fetch(url, {
      method: 'POST',
    })
  }

  return (
    <>
      <div className='inline-buttons'>
        <button className='DownloadButton' onClick={() => {
            SwitchIsDownloading();
            DownloadSelected();
          }
        }
        disabled={isDownloading.current}
        >Download selected
        </button>
        <button className='CancelDownloadButton' onClick={() => {
            SwitchIsDownloading();
            tracks.forEach((track, id) => CancelDownload(id))
          }
        }
          disabled={!isDownloading.current}>Cancel Download</button>
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
              disabled={isDownloading.current}
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
                }
                disabled={isDownloading.current}
                /></td>
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
