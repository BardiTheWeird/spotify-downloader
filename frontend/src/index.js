import React from 'react';
import ReactDOM from 'react-dom';
import './index.css';
import {App, InputBar, playlistTable} from './App';

ReactDOM.render(
  <React.StrictMode>
    <App />
    <InputBar />
  </React.StrictMode>,
  document.getElementById('root')
);