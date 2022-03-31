import React from "react";

import { FaqStatusContext } from '../services/FaqService';
import { UserContext } from "../services/UserService";

export function Faq() {
    const [user] = React.useContext(UserContext);
    const [faqStatus, updateFAQStatus] = React.useContext(FaqStatusContext);

    let userName = "User";
    if (user) {
        userName = user.display_name;
    }
    const system = navigator.platform;

    let ffmperUrl = "";

    if (system == "Win32") {
        ffmperUrl = <a href='https://github.com/BtbN/FFmpeg-Builds/releases/download/latest/ffmpeg-n5.0-latest-win64-gpl-5.0.zip' className='link'>FFMPEG</a>
    }
    else {
        ffmperUrl = <>
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
                <div>Dear {userName},</div>
                <div>welcome to Spotify Downloader</div>
                </p>                
                <body className='infoBody'>
                <div className='infotext'>Please notice: for application to work properly you need to install:</div>
                <div className='infotext'>
                <i className="fa-solid fa-download infotext"></i>
                <a href='https://youtube-dl.org/' className='link'>youtube-dl</a>
                </div>
                <div className='infotext'>
                <i className="fa-solid fa-download infotext"></i>
                {ffmperUrl}
                </div>
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
                </body>
                <button onClick={() => updateFAQStatus(false)} className='uselessButton'>Goi It</button>
            </div>
        </div>
        }
    </>
}