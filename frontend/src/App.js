import logo from './logo.svg';
import './App.css';
import React from 'react';

export function App() {
  return (
    <div className="App">
      <header className="App-header">
        <div>Please make sure you installed</div>
        <a href='https://youtube-dl.org' className="App-link">"YOUTUBE downloader"</a>
        <div>Prior to beginning the search</div>
        <div>Enter Spotify Playlist URL:</div>
      </header>
    </div>
    
  );
}

export function InputBar() {
  const [formData, updateFormData] = React.useState();

  const handleChange = (e) => {
    updateFormData(e.target.value.trim());
  }

  const handleSubmit = (e) => {
    e.preventDefault();
    console.log(formData);
  }
  return (
    <div className="Bar">
      <div>
        <div className="SearchBar">
          <form onSubmit={handleSubmit} class="inputForm">
            <input type="text" name='PL-URL' required class="inputForm" onChange={handleChange}/>
            <input type="submit"/>
          </form>
        </div>
      </div>
    </div>
    
  );
}

export default App;
