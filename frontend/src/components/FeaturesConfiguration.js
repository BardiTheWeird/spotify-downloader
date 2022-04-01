import React from "react";

import { useYtdlFound, useFfmpegFound } from "../services/FeaturesFoundService";
import { useBaseUrl } from "../services/BaseUrlService";
import { FaqStatusContext } from "../services/FaqService";

const { ipcRenderer } = window.require('electron');

export function FeatureConfiguration() {
    const [faqStatus, updateFAQStatus] = React.useContext(FaqStatusContext)
    const [ytdlFound, updateYtdlFound] = useYtdlFound();
    const [ffmpegFound, updateFfmpegFound] = useFfmpegFound();
    const baseUrl = useBaseUrl();

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

    return (
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
                    </>
                }
        </div>
    )
}