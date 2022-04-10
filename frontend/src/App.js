// import logo from './logo.svg';
import './App.css';
import React from 'react';
import {
  Routes,
  Route,
} from "react-router-dom";

import { UserBar } from './components/UserBar';
import { AuthCallback } from './components/AuthCallback';
import { InputBar } from './components/InputBar';
import { Faq } from './components/Faq';
import { LightDarkToggle } from './components/LightDarkToggle';
import { useBaseUrl } from './services/BaseUrlService';
import { FaqStatusContext } from './services/FaqService';
import { useFfmpegFound, useYtdlFound } from './services/FeaturesFoundService';
import { FeatureConfiguration } from './components/FeaturesConfiguration';

export const isDarkInitialValue = localStorage.getItem("DarkMode") === "true";

const { ipcRenderer } = (window.require && window.require('electron')) || (window.opener && window.opener.require('electron'));

export function App() {
  const baseUrl = useBaseUrl();

  const [ytdlFound] = useYtdlFound();
  const [ffmpegFound] = useFfmpegFound();

  const [isDark, updateisDark] = React.useState(isDarkInitialValue);
  React.useEffect(() => {
    localStorage.setItem("DarkMode", isDark.toString())
  }, [isDark]);

  function LightDark() {
    let returnVal;
    if (isDark === true) {
      returnVal = "Dark";
    }
    else {
      returnVal = "Light";
    }
    return returnVal;
  }

  return (
    <div className={`App ${LightDark()}`}>
    <Routes>
      <Route path="/callback" element={<AuthCallback />}/>
      <Route path="/" element=
        {
          baseUrl === undefined &&
            <div>Backend is starting...</div>
          || baseUrl === null &&
            <div>Backend could not be started</div>
          || <>
            
            <LightDarkToggle isDark={isDark} updateisDark={updateisDark} LightDark={LightDark}/>
            <Faq />
            <UserBar/>
            {!(ffmpegFound && ytdlFound) &&
                <FeatureConfiguration />
              || <>
                <header className="App-header">
                  <div>Enter Spotify Playlist Or Album URL:</div>
                </header>
                <InputBar />
            </>
            }
            
            </>
        }/>
    </Routes>
    <div className='controls'>
      <button className='controlButton' onClick={() => ipcRenderer.invoke('winMinimize')}><i className='fa-regular fa-window-minimize controlSymb'></i></button>
      <button className='controlButton' onClick={() => ipcRenderer.invoke('winMaximize')}><i className='fa-regular fa-window-maximize controlSymb'></i></button>
      <button className='controlButton' onClick={() => ipcRenderer.invoke('winClose')}><i className='fa-regular fa-window-close controlSymb'></i></button>
    </div>
    
    </div>
  );
}

export default App;
