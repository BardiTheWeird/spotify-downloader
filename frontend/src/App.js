import logo from './logo.svg';
import './App.css';
import React from 'react';
const {ipcRenderer} = window.require('electron');

export const isDarkInitialValue = localStorage.getItem("DarkMode") === "true";

const BaseUrlContex = React.createContext();

export function App() {
  const [isDark, updateisDark] = 
    React.useState(isDarkInitialValue);
  React.useEffect(() => {
    localStorage.setItem("DarkMode", isDark.toString())
  }, [isDark]);

  const [baseUrl, updateBaseUrl] = React.useState();
  React.useEffect(() => {
    (async () => {
      while (true) {
        const backendStatus = await ipcRenderer.invoke('backendStatus');
        if (backendStatus) {
          if (backendStatus.running) {
            updateBaseUrl(backendStatus.address);
          }
          else {
            updateBaseUrl(null);
          }
          break;
        }
        await new Promise(r => setTimeout(r, 500));
      }
    })();
  }, []);

  function LightDark() {
    let returnVal;
    if (isDark === true) {
      returnVal = "Dark";
    }
    else {
      returnVal = "Light";
    }
    return returnVal;
  }
  return (
    <div className={`App ${LightDark()}`}>
      {
        baseUrl === undefined &&
          <div>Backend is starting...</div>
        || baseUrl === null &&
          <div>Backend could not be started</div>
        || <BaseUrlContex.Provider value={baseUrl}>
          <div  className='App-header-info'>Light/Dark</div>
          <label className="switch">
            <input type="checkbox" onChange={e => updateisDark(!isDark)} checked={!isDark}></input>
            <span className="slider round"></span>
          </label>
          <header className="App-header">
            <div>Enter Spotify Playlist Or Album URL:</div>
          </header>
          <header>
          <div className='App-header-info'>If directory is not defined PL will be loaded into "userDownloads" folder</div>
          </header>
          <InputBar />
        </BaseUrlContex.Provider>
      }
    </div>
  );
}

export function InputBar() {
  const baseUrl = React.useContext(BaseUrlContex);
  const [formData, updateFormData] = React.useState();
  const [playlist, updatePlaylist] = React.useState();
  const [downloadPath, updateDownloadPath] = React.useState('');

  React.useEffect(() => {
    ipcRenderer.on('returnDirectory', (e, path) => {
      updateDownloadPath(path[0]);
    });
  }, []);

  const submitPlaylistLink = async (e) => {
    e.preventDefault();
    let response = await fetch(`${baseUrl}/spotify/playlist?link=${formData}`);
    let playlist = await response.json();
    updatePlaylist(playlist);
  }

  return (
    <>
      <div className="Bar">
        <div>
          <div className="SearchBar">
            <form onSubmit={submitPlaylistLink} className="inputForm">
              <input placeholder='Spotify Link (https://open.spotify.com/playlist/etc...)' type="text" name='PL-URL' required className="inputForm" onChange={
                e => updateFormData(e.target.value.trim())}/>
              <input type="submit" className="uselessButton" value="Submit"/>
            </form>
          </div>
          <div className="SearchBar">
            <form onSubmit={e => e.preventDefault()} className="inputForm">
              <input placeholder='Insert Download Directory' type="text" name='DL-path' required className="inputForm" onChange={
                e => updateDownloadPath(e.target.value)}
                  value={downloadPath}
                />
              <button className="uselessButton" onClick={e => {
                e.preventDefault();
                ipcRenderer.send('openDirectory');
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

export default App;

export function PlaylistTable({playlist, downloadPath}) {
  const baseUrl = React.useContext(BaseUrlContex);
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
      let url = `${baseUrl}/download/start`;
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
        }),
      });
      if (downloadResponse.status !== 204) {
        downloadCounter.current--;
        if (downloadCounter.current == 0) {
          isDownloading.current = false;
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
          track.status = "Bad request";
          console.log(await downloadResponse.json());
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
        const getStatusResponse = await fetch(`${baseUrl}/download/status?folder=${downloadFolder}&filename=${track.artists} - ${track.title}`);
        const downloadEntry = await getStatusResponse.json();

        switch (downloadEntry.status) {
          case 0: 
            const percentage = Math.floor(downloadEntry.downloaded_bytes/downloadEntry.total_bytes*100);
            copiedTracks[trackIndex].status = `Downloading ${percentage || 0}%`;
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
            setForceUpdate(Date.now());
          }
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
    let url = `${baseUrl}/download/cancel?folder=${downloadFolder}&filename=${track.artists} - ${track.title}`;
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
              <input type="checkbox" className="checkmark" onChange={
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
                <td><input type="checkbox" className="checkmark" checked={track.checked} onChange={
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
