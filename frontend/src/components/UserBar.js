import React from "react";
import { OAuthUrlContext } from "../services/OAuthUrlService";

import {UserContext} from '../services/UserService'
import {authorizedFetch} from '../utilities'

const { ipcRenderer } = window.require('electron');


export function UserBar() {
    const [user, updateUser] = React.useContext(UserContext);
    const oauthUrl = React.useContext(OAuthUrlContext);

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
        return <>{
            oauthUrl && <a href={oauthUrl} className="Login">Log In</a>
        }</>
    }
    else {
        return <>
        <button className="userleft">
            <img src={userImage} className='userImage'/>
            <span>{user.display_name}</span><i className="fa-solid fa-caret-down arrowdown"></i>
            <button className="logout" onClick={Logout}>
            Log Out
            </button>
        </button>
        </>
    }
}

