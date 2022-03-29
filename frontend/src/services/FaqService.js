import React from "react";

export const FaqStatusContext = React.createContext();

export function FaqStatusProvider(props) {
    const [faqStatus, updateFAQStatus] = React.useState();
    
    return (
        <FaqStatusContext.Provider value={[faqStatus, updateFAQStatus]}>
            {props.children}
        </FaqStatusContext.Provider>
    )
}
