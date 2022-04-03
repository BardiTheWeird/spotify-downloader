import React from "react";
import { OAuthUrlContext } from "../services/OAuthUrlService";

import {UserContext} from '../services/UserService'
import {authorizedFetch} from '../utilities'

const { ipcRenderer } = window.require('electron');


export function UserBar() {
    const [user, updateUser] = React.useContext(UserContext);
    const oauthUrl = React.useContext(OAuthUrlContext);
    const [loginStatus, updateLoginStatus] = React.useState();

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

    React.useEffect(async () => {
        updateUser(await getUser());
    },[]);

    function Logout() {
        localStorage.setItem('access token', '');
        localStorage.setItem('refresh token', '');
        updateUser(null);
    }

    let userImage = "./icon.ico";
    if (user && user.image) {
        userImage = user.image;
    }

    if (!user) {
        return <>
            <div className="userleft">
                <span className="login">Logged Out</span>
                <i className="fa-solid fa-caret-down PublicDataArrow"></i>            
            {
                oauthUrl && <span className="logout privatePlaylist" onClick={() => {
                    updateLoginStatus(true);
                }}>Log In</span>
            }
            </div>
            

{ loginStatus && 
    <div className='infoBack'>
        <div className='infobox'>
            <h3>Log In FAQ</h3>
            <p>
                <span>Dear User,</span><br/>
                <span>welcome to Spotify Downloader</span>
            </p>                
            <div className='infotext'>
                To login properly please login as a developer using a</div>
            <div className='infotext'>
            <i className="fa-solid fa-sign-in-alt signInArrow"></i>
            <a href="https://developer.spotify.com" className='link infotext'>
                Developer website 
            </a>
            </div>
                       
            <div className='infotext'>
                Click "Create New App" and give it any name you see fit. </div>
            <div className='infotext'>
                Then click "Edit settings" in the top right corner.
            </div>
            <div className='infotext'>
                Into the "Redirect URIs" field insert " app://-/callback " and click add, then "Save".
            </div>
            <div className='infotext'>
                Then in the top left corner copy Client ID and insert into a field below. Then click "Log In" button and login into the Spotify using it's interface.
            </div>
            <input type="text" className='ClientIDField inputForm'></input>
            <a href={oauthUrl} className='loginButton uselessButton' >Log In</a>
            <button onClick={() => updateLoginStatus(false)} className='uselessButton'>Nah...</button>
        </div>
    </div>
    }
    </>      
    }
    else {
        return <>
        <button className="userleft">
            <img src={userImage} className='userImage'/>
            <span>{user.display_name}</span><i className="fa-solid fa-caret-down"></i>
            <button className="logout" onClick={Logout}>
            Log Out
            </button>
        </button>
        </>
    }
}

