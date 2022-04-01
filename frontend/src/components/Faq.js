import React from "react";
import { useBaseUrl } from "../services/BaseUrlService";

import { FaqStatusContext } from '../services/FaqService';
import { useYtdlFound, useFfmpegFound } from "../services/FeaturesFoundService";
import { UserContext } from "../services/UserService";

const { ipcRenderer } = window.require('electron');

export function Faq() {
    const [user] = React.useContext(UserContext);
    const [faqStatus, updateFAQStatus] = React.useContext(FaqStatusContext);
    const baseUrl = useBaseUrl();

    const [ytdlFound, updateYtdlFound] = useYtdlFound();
    const [ffmpegFound, updateFfmpegFound] = useFfmpegFound();

    async function featuresStatus() {
        let featuresStatus = await fetch(`${baseUrl}/features`, {
            method: 'GET',
        })
        featuresStatus = await featuresStatus.json();

        updateYtdlFound(featuresStatus.youtube_dl);
        updateFfmpegFound(featuresStatus.ffmpeg);
    }

    React.useEffect(() => {
        featuresStatus();
    }, [baseUrl, faqStatus]);
    
    let userName = "User";
    if (user) {
        userName = user.display_name;
    }
    const system = navigator.platform;

    let ffmpegUrl = "";

    if (system == "Win32") {
        ffmpegUrl = <a href='https://github.com/BtbN/FFmpeg-Builds/releases/download/latest/ffmpeg-n5.0-latest-win64-gpl-5.0.zip' className='link'>FFMPEG</a>
    }
    else if (system == "darwin"){
        ffmpegUrl = <a href='https://evermeet.cx/ffmpeg/' className='link'>FFMPEG</a>
    }
    else {
        ffmpegUrl = <>
            Arch: pacman -S ffmpeg <br/>
            Debian: apt-get install ffmpeg
        </>; 
    }

    return <>
        <div className='FAQ' onClick={() => updateFAQStatus(true)}>
            <img src={"./icon.ico"} width="40px" height="40px"></img>
        </div>

        { faqStatus && 
        <div className='infoBack'>
            <div className='infobox'>
                <h3>FAQ</h3>
                <p>
                    <span>Dear {userName},</span><br/>
                    <span>welcome to Spotify Downloader</span>
                </p>                
                <div className='infoBody'>
                    {!(ffmpegFound && ytdlFound) && <>
                        <div className='infotext'>Please notice: for application to work properly you need to install applications using the links below and provide the route to executables after the downloading using respective buttons:</div>
                        <div className='infotext'>
                        <i className="fa-solid fa-download infotext"></i>
                        <a href='https://youtube-dl.org/' className='link'>youtube-dl</a>
                        {!ytdlFound &&
                            <button onClick={async () => {
                                const pathOk = await ipcRenderer.invoke('configureFeaturePath', 'youtube-dl');
                                updateYtdlFound(pathOk)
                            }}
                            className="uselessButton">
                                Path to
                            </button>
                        }
                        </div>
                        <div className='infotext'>
                            <i className="fa-solid fa-download infotext"></i>
                            {ffmpegUrl}
                            {!ffmpegFound &&
                                <button onClick={async () => {
                                        const ffmpegOk = await ipcRenderer.invoke('configureFeaturePath', 'ffmpeg');
                                        updateFfmpegFound(ffmpegOk);
                                    }}
                                    className="uselessButton">
                                    Path to
                                </button>
                            }
                        </div>
                    </>}
                    <div className='infotext'>
                        Before searching of a playlist, please login using a button in upper-left corner. You can log out any time you want using the dropping button under the profile name.
                    </div>
                    <div className='infotext'>
                        Insert a copied link to the Spotify playlist or album into the upper submission field and click Submit button. If the link is incorrect you'll receive a message from the application.
                    </div>
                    <div className='infotext'>
                        Before the download, either insert a directory into the second submission field or use the Browse button to select a desired folder.
                    </div>
                    <div className='infotext'>
                        Select tracks you want to download using checkboxes on the left; begin the download by clicking the Download Selected button. While in process it can be cancelled by the respective button. Download status will be displayed in the Status column.
                    </div>
                </div>
                <button onClick={() => updateFAQStatus(false)} className='uselessButton'>Got It</button>
            </div>
        </div>
        }
    </>
}