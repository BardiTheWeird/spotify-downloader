import React from "react";

const { ipcRenderer } = window.require('electron');

const OAuthUrlContext = React.createContext();

export function useOAuthUrl() {
  return React.useContext(OAuthUrlContext)[0];
}

export function useClientId() {
  return React.useContext(OAuthUrlContext)[1];
}

export function OAuthUrlProvider(props) {
    const [appUrl, updateAppUrl] = React.useState();
    const [codeChallenge, updateCodeChallenge] = React.useState();
    const [clientId, updateClientId] = React.useState("");

    React.useEffect(async () => {
      // read clientId from localStorage
      updateClientId(localStorage.getItem('clientId') || '');

      // update code challenge
      const codeVerifier = generateRandomString(64);
      updateCodeChallenge(await generateCodeChallenge(codeVerifier));
      localStorage.setItem('code_verifier', codeVerifier);

      updateAppUrl(await ipcRenderer.invoke('appUrl'));
    }, []);

    React.useEffect(() => {
      localStorage.setItem('clientId', clientId);
    }, [clientId]);

    return <OAuthUrlContext.Provider value={[
      `https://accounts.spotify.com/authorize?response_type=code&client_id=${clientId}&redirect_uri=${appUrl}/callback&code_challenge_method=S256&code_challenge=${codeChallenge}&scope=user-read-private,playlist-read-private,user-library-read`, 
      [clientId, updateClientId]]}>
        {props.children}
    </OAuthUrlContext.Provider>
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