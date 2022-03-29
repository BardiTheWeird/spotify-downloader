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

ReactDOM.render(
  <React.StrictMode>
    <BrowserRouter>

      <BaseUrlProvider>
      <UserProvider>
      <FaqStatusProvider>
      <OAuthUrlProvider>

        <App />

      </OAuthUrlProvider>
      </FaqStatusProvider>
      </UserProvider>
      </BaseUrlProvider>

    </BrowserRouter>
  </React.StrictMode>,
  document.getElementById('root')
);
