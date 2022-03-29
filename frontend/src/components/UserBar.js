import React from "react";

import {UserContext} from '../services/UserService'
import {authorizedFetch} from '../utilities'

const { ipcRenderer } = window.require('electron');


export function IsLoggedIn() {
    const [user, updateUser] = React.useContext(UserContext);
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
        }
    },[]);

    function Logout() {
        localStorage.setItem('access token', '');
        localStorage.setItem('refresh token', '');
        updateUser(null);
        updateCodeChallenge();
        updateAppUrl();
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
