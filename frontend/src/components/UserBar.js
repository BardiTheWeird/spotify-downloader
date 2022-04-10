import React from "react";
import { useClientId, useOAuthUrl } from "../services/OAuthUrlService";
import { useGetFavourites } from "../services/PlaylistService";

import { UserContext } from '../services/UserService'
import { authorizedFetch } from '../utilities'

export function UserBar() {
    const oauthUrl = useOAuthUrl();
    const [clientId, updateClientId] = useClientId();
    const [user, updateUser] = React.useContext(UserContext);
    const [loginStatus, updateLoginStatus] = React.useState();
    const getFavourites = useGetFavourites();
    const [loginPopup, updateLoginPopup] = React.useState(null);
    const [heart, updateHeart] = React.useState("regular");

    // returns userObj or null if not logged in
    async function getUser() {
        const userInfoResponse = await authorizedFetch('https://api.spotify.com/v1/me', {
            headers: { 
                'Accept': 'application/json', 
                'Content-Type': 'application/json',
            }
        });
        if (!userInfoResponse || userInfoResponse.status !== 200) {
        return null
        }
        const userInfo = await userInfoResponse.json();

        return {
            display_name: userInfo.display_name,
            image: userInfo.images.length === 0 ? null : userInfo.images[0].url
        };
    }

    React.useEffect(() => {
        window.self.onLoginSuccess = () => {
            updateLoginPopup(null);
            (async () => updateUser(await getUser()))();
        }

        (async () => {
            updateUser(await getUser());
        })();
    }, []);

    function Logout() {
        localStorage.setItem('access token', '');
        updateLoginStatus(false);
        updateUser(null);
    }

    let userImage = "./icon.ico";
    if (user && user.image) {
        userImage = user.image;
    }
    
return <>{
!user && <>
    <div className="userleft Login">
        <div>Logged Out</div>
        <i className="fa-solid fa-caret-down PublicDataArrow"></i>            
        { oauthUrl && 
            <div className="logout privatePlaylist" 
                onClick={() => updateLoginStatus(true)}>
                Log In
            </div> 
        }
    </div>
            

    { loginStatus && 
        <div className='infoBack'>
            <div className='infobox'>
                <h3>Log In FAQ</h3>
                <p>
                    <span>Dear User,</span><br/>
                    <span>welcome to Spotify Downloader</span>
                </p>                
                <div className='infotext'>
                    To login properly please login as a developer using a</div>
                <div className='infotext'>
                <i className="fa-solid fa-sign-in-alt signInArrow"></i>
                <a href="https://developer.spotify.com/dashboard/" className='link infotext' target="_blank" rel="noreferrer">
                    Developer website 
                </a>
                </div>
                        
                <div className='infotext'>
                    Click "Create New App" and give it any name you see fit. </div>
                <div className='infotext'>
                    Then click "Edit settings" in the top right corner.
                </div>
                <div className='infotext'>
                    Into the "Redirect URIs" field insert " app://-/callback " and click add, then "Save".
                </div>
                <div className='infotext'>
                    Then in the top left corner copy Client ID and insert into a field below. Then click "Log In" button and login into the Spotify using it's interface.
                </div>
                <input type="text" className='ClientIDField inputForm' 
                    value={clientId} 
                    onChange={e => updateClientId(e.target.value.trim())}
                    placeholder='Client ID'>
                </input>
                <button className='uselessButton' 
                    onClick={() => {
                        if (clientId.length !== 32) {
                            alert('ClientID is not 32 characters');
                            return;
                        }
                        if (loginPopup && !loginPopup.closed) {
                            loginPopup.location.href = oauthUrl;
                            loginPopup.focus();
                            return;
                        }
                        updateLoginPopup(window.self.open(oauthUrl, 'sharer'));
                    }}>
                    Log In
                </button>
                <button className='uselessButton'
                    onClick={() => updateLoginStatus(false)}>
                    Nah...
                </button>
            </div>
        </div>
    }
</>
|| 
    <div className="userContainer">
        <button className="userleft">
            <img src={userImage} className='userImage'/>
            <span>{user.display_name}</span>
            <i className="fa-solid fa-caret-down arrowdown"></i>
            <button className="logout" onClick={Logout}>
                Log Out
            </button>
        </button>
        <button className="likedSongs" 
            onClick={getFavourites} 
            onMouseEnter={() => updateHeart("solid")}
            onMouseLeave={() => updateHeart("regular")}>
            
            <i className={`fa-${heart} fa-heart heartAction`}></i>
            <span className="LikedSongsExpansion">Get All Liked</span>
        </button>
        <span className='dragBarLeft'></span>
        <span className='dragBarRight'></span>
    </div>
}</>
}

