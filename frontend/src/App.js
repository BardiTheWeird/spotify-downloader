import logo from './logo.svg';
import './App.css';
import React from 'react';
import {
  Routes,
  Route,
  useNavigate,
} from "react-router-dom";
const {ipcRenderer} = window.require('electron');

export const isDarkInitialValue = localStorage.getItem("DarkMode") === "true";

const BaseUrlContex = React.createContext();
const IsLoggedInContext = React.createContext();

async function authorizedFetch(input, init) {
  async function innerFunction() {
    let accessToken = localStorage.getItem('access token');
    if (!accessToken) {
      return null;
    }
    init = init || {};
    init.headers = {
      ...init.headers,
      'Authorization': `Bearer ${accessToken}`,
    }

    let response = await fetch(input, init);
    if (response.status === 401) {
      // refresh access token
      const refreshToken = localStorage.getItem('refresh token');
      if (!refreshToken) {
        return null;
      }

      const refreshResponse = await fetch(`https://accounts.spotify.com/api/token?grant_type=refresh_token&refresh_token=${refreshToken}&client_id=63d55a793f9c4a9e8d5aacba30069a23`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/x-www-form-urlencoded'
        }
      });
      if (refreshResponse.status !== 200) {
        return null;
      }
      const refreshResponseBody = refreshResponse.json();
      accessToken = refreshResponseBody.access_token;
      localStorage.setItem('access token', accessToken);

      response = await fetch(input, init);
    }
    return response;
  }

  const result = await innerFunction();
  // clean up token if they're invalid
  if (result === null) {
    localStorage.setItem('access token', '');
    localStorage.setItem('refresh token', '');
  }
  return result;
}

export function IsLoggedIn() {
  const [isUserLogged, updateIsUserLogged] = React.useContext(IsLoggedInContext);
  const [user, updateUser] = React.useState();
  const [code_challenge, updateCode_challenge] = React.useState();
  const [appUrl, _updateAppUrl] = React.useState();

  // returns userObj or null if not logged in
  async function getUser() {
    const userInfoResponse = await authorizedFetch('https://api.spotify.com/v1/me', {
      headers: { 
        'Accept': 'application/json', 
        'Content-Type': 'application/json',
      }
    });
    if (!userInfoResponse || userInfoResponse.status !== 200) {
      return null
    }
    const userInfo = await userInfoResponse.json();

    return {
      display_name: userInfo.display_name,
      image: userInfo.images.length === 0 ? null : userInfo.images[0].url
    };
  }

  async function updateCodeChallenge() {
    const code_verifier = generateRandomString(64);
      updateCode_challenge(await generateCodeChallenge(code_verifier));
      localStorage.setItem('code_verifier', code_verifier);
  }

  async function updateAppUrl() {
    _updateAppUrl(await ipcRenderer.invoke('appUrl'));
  }

  React.useEffect(async () => {
    const user = await getUser();
    updateUser(user);
    if (!user) {
      updateCodeChallenge();
      updateAppUrl();
      updateIsUserLogged(false)
    }
    else {updateIsUserLogged(true)}
  },[]);

  function Logout() {
    localStorage.setItem('access token', '');
    localStorage.setItem('refresh token', '');
    updateUser(null);
    updateCodeChallenge();
    updateAppUrl();
    updateIsUserLogged(false);
  }
  
  if (!user) {
    if (!appUrl) {
      return <></>;
    }

    return <>{
      code_challenge && <a href={`https://accounts.spotify.com/authorize?response_type=code&client_id=63d55a793f9c4a9e8d5aacba30069a23&redirect_uri=${appUrl}/callback&code_challenge_method=S256&code_challenge=${code_challenge}`} className="Login">Log In</a>
    }</>
  }
  else {
    return <>
      <button className="userleft">
        <img src={user.image} className='userImage'/>
        <span>{user.display_name}</span><i className="fa-solid fa-caret-down arrowdown"></i>
          <button className="logout" onClick={Logout}>
          Log Out
          </button>
      </button>
    </>
  }
}

function generateRandomString(length) {
  let text = '';
  const possible =
    'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';

  for (let i = 0; i < length; i++) {
    text += possible.charAt(Math.floor(Math.random() * possible.length));
  }
  return text;
}

async function generateCodeChallenge(codeVerifier) {
  const digest = await crypto.subtle.digest(
    'SHA-256',
    new TextEncoder().encode(codeVerifier),
  );

  return btoa(String.fromCharCode(...new Uint8Array(digest)))
    .replace(/=/g, '')
    .replace(/\+/g, '-')
    .replace(/\//g, '_');
}

export function App() {
  const [isUserLogged, updateIsUserLogged] = React.useState();
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
    <Routes>
      <Route path="/callback" element={<AuthCallback />}/>
      <Route path="/" element=
        {
          baseUrl === undefined &&
            <div>Backend is starting...</div>
          || baseUrl === null &&
            <div>Backend could not be started</div>
          || <BaseUrlContex.Provider value={baseUrl}>
          <IsLoggedInContext.Provider value={[isUserLogged, updateIsUserLogged]}>
            <div className="userright">
              <label className="switch">
                <input type="checkbox" onChange={e => updateisDark(!isDark)} checked={!isDark}>
                </input>
                <span className="slider round"></span>
              </label>
              <div className='App-header-info symbolTranslate'>
                { 
                  LightDark() == "Light" &&
                     <i className="fa-solid fa-moon"></i>
                  || <i className="fa-solid fa-sun"></i>
                }
              </div>
            </div>
            
            <header className="App-header">
              <IsLoggedIn/>
              <div>Enter Spotify Playlist Or Album URL:</div>
            </header>
            <InputBar />
          </IsLoggedInContext.Provider>
          </BaseUrlContex.Provider>
        }/>
    </Routes>
    </div>
  );
}

export function AuthCallback() {
  const navigate = useNavigate();
  React.useEffect(async () => {
    const code_verifier = localStorage.getItem('code_verifier');
    const url = new URL(document.URL)
    const authorizationCode = url.searchParams.get('code');

    const appUrl = await ipcRenderer.invoke('appUrl');

    const response = await fetch(`https://accounts.spotify.com/api/token?grant_type=authorization_code&code=${authorizationCode}&redirect_uri=${appUrl}/callback&client_id=63d55a793f9c4a9e8d5aacba30069a23&code_verifier=${code_verifier}`, {
      method: "POST",
      headers: { 'Content-Type': 'application/x-www-form-urlencoded' }
    });
    const responseBody = await response.json();

    localStorage.setItem('access token', responseBody.access_token || '');
    localStorage.setItem('refresh token', responseBody.refresh_token || '');

    navigate('/');
  }, [])

  return <></>;
}

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
            if (!downloadPath) {
              return alert("Please choose directory using the Browse button")
            }
            else {
              SwitchIsDownloading();
              DownloadSelected();
            }
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
      <thead>
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
        </thead>
      <tbody>
        {
          tracks.map((track, index) =>
            (
              <tr key={track.id}>
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
                height="30px"/>
                </td>
                <td>{track.artists.join(', ')}</td>
                <td>{track.title}</td>
                <td>{track.album_title}</td>
                <td>{track.status}</td>
              </tr>
            )
          )
        }
        </tbody>
      </table>
    </>
  )
}
