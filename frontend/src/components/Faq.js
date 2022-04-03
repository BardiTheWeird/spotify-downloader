import React from "react";

import { FaqStatusContext } from '../services/FaqService';
import { UserContext } from "../services/UserService";
import { useYtdlFound, useFfmpegFound } from "../services/FeaturesFoundService";

export function Faq() {
    const [user] = React.useContext(UserContext);
    const [faqStatus, updateFAQStatus] = React.useContext(FaqStatusContext);

    const [ytdlFound] = useYtdlFound();
    const [ffmpegFound] = useFfmpegFound();

    let FaqButtonStyle = "";
    if (ytdlFound && ffmpegFound) {
        FaqButtonStyle = "FAQ"
    }
    else {
        FaqButtonStyle = "FAQpulse"
    }

    let userName = "User";
    if (user) {
        userName = user.display_name;
    }

    return <>
        <div className={FaqButtonStyle} onClick={() => updateFAQStatus(true)}>
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
                
                    <div className='infotext'>
                        Please notice: for a search of a private playlist e.g. Liked Tracks, please login using a dropping button in upper-left corner. You can log out any time you want using the dropping button under the profile name.
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
                
                <button onClick={() => updateFAQStatus(false)} className='uselessButton'>Got It</button>
            </div>
        </div>
        }
    </>
}