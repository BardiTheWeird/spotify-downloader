// import logo from './logo.svg';
import './App.css';
import React from 'react';
import {
  Routes,
  Route,
} from "react-router-dom";

import {BaseUrlContex, IsLoggedInContext} from './contexts'
import { IsLoggedIn } from './components/IsLoggedIn';
import { AuthCallback } from './components/AuthCallback';
import { InputBar } from './components/InputBar';
import { Faq } from './components/Faq';
import { LightDarkToggle } from './components/LightDarkToggle';

const {ipcRenderer} = window.require('electron');

export async function BackendPull() {
  const [isResourses, updateIsResourses] = React.useState();
  updateIsResourses(await fetch())
}


export const isDarkInitialValue = localStorage.getItem("DarkMode") === "true";

export function App() {
  const [isUserLogged, updateIsUserLogged] = React.useState();
  const [isDark, updateisDark] = React.useState(isDarkInitialValue);
  React.useEffect(() => {
    localStorage.setItem("DarkMode", isDark.toString())
  }, [isDark]);

  const [baseUrl, updateBaseUrl] = React.useState();
  React.useEffect(() => {
    (async () => {
      while (true) {
        const backendStatus = await ipcRenderer.invoke('backendStatus');
        if (backendStatus) {
          if (backendStatus.running) {
            updateBaseUrl(backendStatus.address);
          }
          else {
            updateBaseUrl(null);
          }
          break;
        }
        await new Promise(r => setTimeout(r, 500));
      }
    })();
  }, []);

  const [faqStatus, updateFAQStatus] = React.useState();

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
    <BaseUrlContex.Provider value={baseUrl}>
    <IsLoggedInContext.Provider value={[isUserLogged, updateIsUserLogged]}>

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
            <Faq faqStatus={faqStatus} updateFAQStatus={updateFAQStatus} />
            
            <header className="App-header">
              <IsLoggedIn/>
              <div>Enter Spotify Playlist Or Album URL:</div>
            </header>
            <InputBar />
            </>
        }/>
    </Routes>
    </div>

    </IsLoggedInContext.Provider>
    </BaseUrlContex.Provider>
  );
}

export default App;
