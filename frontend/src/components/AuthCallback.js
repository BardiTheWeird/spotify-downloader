import React from "react";
import { useClientId } from "../services/OAuthUrlService";
const { ipcRenderer } = (window.require && window.require('electron')) || (window.opener && window.opener.require('electron'));

export function AuthCallback() {
  const [clientId] = useClientId();
  React.useEffect(() => {
    (async () => {
      console.log('clientId:', clientId);
      if (!clientId) {
        return;
      }

      const code_verifier = localStorage.getItem('code_verifier');
      const url = new URL(document.URL)
      const authorizationCode = url.searchParams.get('code');
  
      const appUrl = await ipcRenderer.invoke('appUrl');
  
      const response = await fetch(`https://accounts.spotify.com/api/token?grant_type=authorization_code&code=${authorizationCode}&redirect_uri=${appUrl}/callback&client_id=${clientId}&code_verifier=${code_verifier}`, {
        method: "POST",
        headers: { 'Content-Type': 'application/x-www-form-urlencoded' }
      });
      console.log('response:', response);
      const responseBody = await response.json();
      console.log('responseBody:', responseBody);
  
      localStorage.setItem('access token', responseBody.access_token || '');
  
      window.opener.onLoginSuccess(responseBody.access_token || '');
      window.close();
    })();
  }, [clientId]);
  
  return <></>;
}