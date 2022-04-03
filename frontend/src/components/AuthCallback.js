import React from "react";
import { useNavigate } from "react-router-dom";
import { useClientId } from "../services/OAuthUrlService";
const { ipcRenderer } = window.require('electron');

export function AuthCallback() {
    const navigate = useNavigate();
    const [clientId] = useClientId();

    React.useEffect(async () => {
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
      const responseBody = await response.json();
  
      localStorage.setItem('access token', responseBody.access_token || '');
  
      navigate('/');
    }, [clientId])
  
    return <></>;
  }