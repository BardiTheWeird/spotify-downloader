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

export const isDarkInitialValue = localStorage.getItem("DarkMode") === "true";

export function App() {
  const baseUrl = useBaseUrl();

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
            
            <header className="App-header">
              <UserBar/>
              <div>Enter Spotify Playlist Or Album URL:</div>
            </header>
            <InputBar />
            </>
        }/>
    </Routes>
    </div>
  );
}

export default App;
