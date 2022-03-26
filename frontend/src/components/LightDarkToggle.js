import React from "react";

export function LightDarkToggle({isDark, updateisDark, LightDark}) {
    return <div className="userright">
        <label className="switch">
        <input type="checkbox" onChange={() => updateisDark(!isDark)} checked={!isDark}>
        </input>
        <span className="slider round"></span>
        </label>
        <div className='App-header-info symbolTranslate'>
        { 
            LightDark() == "Light" &&
                <i className="fa-solid fa-moon"></i>
            || <i className="fa-solid fa-sun"></i>
        }
        </div>
    </div>
}