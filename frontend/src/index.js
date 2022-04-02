import React from 'react';
import ReactDOM from 'react-dom';
import './index.css';
import { App } from './App';
import { BrowserRouter } from 'react-router-dom';
import './css/all.min.css'
import { BaseUrlProvider } from './services/BaseUrlService';
import { UserProvider } from './services/UserService';
import { FaqStatusProvider } from './services/FaqService';
import { OAuthUrlProvider } from './services/OAuthUrlService';
import { FeatureFoundProvider } from './services/FeaturesFoundService';
import { PlayerProvider } from './services/PlayerService'

ReactDOM.render(
  <React.StrictMode>
    <BrowserRouter>

      <BaseUrlProvider>
      <UserProvider>
      <FaqStatusProvider>
      <OAuthUrlProvider>
      <FeatureFoundProvider>
      <PlayerProvider>

        <App />

      </PlayerProvider>
      </FeatureFoundProvider>
      </OAuthUrlProvider>
      </FaqStatusProvider>
      </UserProvider>
      </BaseUrlProvider>

    </BrowserRouter>
  </React.StrictMode>,
  document.getElementById('root')
);
