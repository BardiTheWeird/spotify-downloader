import React from "react";
import { useBaseUrl } from "../services/BaseUrlService";
import { usePlayPause } from "../services/PlayerService";

const illegalFilenameChars = ['<', '>', ':', '"', '\\', '/', '|', '?', '*'];

export function PlaylistTable({playlist, downloadPath}) {
    const baseUrl = useBaseUrl();
    // a workaround for forcing a rerender
    const [, setForceUpdate] = React.useState(Date.now());
    const isDownloading = React.useRef(false);
    const [tracks, updateTracks] = React.useState([]);
    const isCancelling = React.useRef(false);

    const downloadCounter = React.useRef(0);
    function incrementDownloadCounter() {
      downloadCounter.current++;
    }
    function decrementDownloadCounter() {
      downloadCounter.current--;
      if (downloadCounter.current == 0) {
        isDownloading.current = false;
        isCancelling.current = false
        setForceUpdate(Date.now());
      }
    }

    React.useEffect(() => {
      if (playlist == 'updating') {
        return;
      }

      updateTracks(playlist.map(track => {
        return {...track,
          isPlaying: false,
          checked: true,
          status: "NA"
        };
      }));
    }, [playlist]);

    const playPause = usePlayPause();
  
    function DownloadSelected() {
      // set pending status
      tracks.forEach(track => { 
        if (track.checked) {
          track.status = 'Pending'
          incrementDownloadCounter();
        }
      });
      updateTracks([...tracks]);
  
      // start download
      const pushTrackDownload = createThrottledFunction(async (track, index) => {
        if (!track.checked) {
          return;
        }
        if (isCancelling.current) {
          track.status = 'Cancelled';
          updateTracks([...tracks]);

          decrementDownloadCounter();
          return;
        }

        let url = `${baseUrl}/download/start`;
        let downloadFolder = downloadPath;
        if (!downloadFolder) {
          downloadFolder = "./userDownloads/"
        }
  
        let filename = `${track.artists.join(', ')} - ${track.title}`;
        filename = filename.split('').map(char => {
        if (illegalFilenameChars.includes(char)) {
            return "_"
        }
        else {
            return char
        }
        }).join('');
  
        let downloadResponse = await fetch(url, {
          method: 'POST',
          body: JSON.stringify({
            id: track.id,
            folder: downloadPath,
            filename: filename,
            title: track.title,
            artist: track.artists.join(', '),
            album: track.album_title,
            image: track.album_image
          }),
        });
        if (downloadResponse.status !== 204) {
          decrementDownloadCounter();
        }
        let updatePromise = null;
        switch (downloadResponse.status) {
          case 204:
            updatePromise = UpdateDownloadStatus(index);
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
          case 408:
          case 500:
            track.status = "Download Error";
            break;
          default:
            track.status = "Unexpected status";
            break;
        }
        
        updateTracks([...tracks]);

        await updatePromise;
      }, 10);

      tracks.forEach((track, index) => pushTrackDownload(track, index));
    }
  
    async function UpdateDownloadStatus(trackIndex) {
        while (true) {
          const copiedTracks = [...tracks];
          const track = copiedTracks[trackIndex];
  
          const getStatusResponse = await fetch(`${baseUrl}/download/status?id=${track.id}`);
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
            default:
              copiedTracks[trackIndex].status = 'Unexpected Status'
              break;
          }
          updateTracks(copiedTracks);
          await new Promise(r => setTimeout(r, 1000));
          if (downloadEntry.status >= 2) {
            decrementDownloadCounter();
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
      let url = `${baseUrl}/download/cancel?id=${track.id}`;
      let downloadResponse = await fetch(url, {
        method: 'POST',
      });
    }
    
    return (
      <>
      {playlist == 'updating' && <div><img src={"./icon.ico"} className='Spinney-vinney'/></div> ||
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
              isCancelling.current = true;
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
                      const copiedTracks = [...tracks];
                      copiedTracks[index].checked = !copiedTracks[index].checked;
                      updateTracks(copiedTracks);
                    }
                  }
                  disabled={isDownloading.current}
                  /></td>
                  <td onMouseEnter={(e) => {e.target.style = "Preview"}} onMouseLeave={(e) => {e.target.style = "PreviewNone"}}
                  onClick={() => {
                    const pausedTrackIndex = playPause(index);
                    tracks[index].isPlaying = !tracks[index].isPlaying;
                    if (pausedTrackIndex !== index) {
                      tracks[pausedTrackIndex].isPlaying = !tracks[pausedTrackIndex].isPlaying;
                    }
                    updateTracks([...tracks]);
                  }}>{
                    tracks[index].isPlaying == false &&
                    <i className="fa-solid fa-play Preview"></i> ||
                    <i className="fa-solid fa-pause PreviewPause"></i>
                  }
                  <img src={track.album_image}
                  className="album_image"/>
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
      }</>
    )
  }

function createThrottledFunction(f, maxParallel) {
  const queuedWork = [];
  let capacity = maxParallel;
  
  async function doWork(...args) {
      capacity--;
      while (true) {
          await f(...args);
          if (!queuedWork.length) {
              break;
          }
          args = queuedWork.shift();
      }
      capacity++;
  }
  
  return (...args) => {
      if (capacity > 0) {
          doWork(...args);
      }
      else {
          queuedWork.push(args);
      }
  }
}