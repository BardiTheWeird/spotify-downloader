import React from "react";
import { useNavigate } from "react-router-dom";
const {ipcRenderer} = window.require('electron');

export function AuthCallback() {
    const navigate = useNavigate();
    React.useEffect(async () => {
      const code_verifier = localStorage.getItem('code_verifier');
      const url = new URL(document.URL)
      const authorizationCode = url.searchParams.get('code');
  
      const appUrl = await ipcRenderer.invoke('appUrl');
  
      const response = await fetch(`https://accounts.spotify.com/api/token?grant_type=authorization_code&code=${authorizationCode}&redirect_uri=${appUrl}/callback&client_id=63d55a793f9c4a9e8d5aacba30069a23&code_verifier=${code_verifier}`, {
        method: "POST",
        headers: { 'Content-Type': 'application/x-www-form-urlencoded' }
      });
      const responseBody = await response.json();
  
      localStorage.setItem('access token', responseBody.access_token || '');
      localStorage.setItem('refresh token', responseBody.refresh_token || '');
  
      navigate('/');
    }, [])
  
    return <></>;
  }