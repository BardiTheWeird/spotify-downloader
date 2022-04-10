import React from "react";

export function LightDarkToggle({isDark, updateisDark, LightDark}) {
    return <div className="userright" onClick={() => updateisDark(!isDark)}>
        <div className=''>
        { 
            LightDark() == "Light" &&
                <i className="fa-regular fa-moon symbTransl"></i>
            || <i className="fa-solid fa-sun symbTransl"></i>
        }
        </div>
    </div>
}