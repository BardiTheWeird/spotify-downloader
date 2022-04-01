import React from "react";

const FeaturesFoundContext = React.createContext([]);

export function useYtdlFound() {
    return React.useContext(FeaturesFoundContext)[0];
}

export function useFfmpegFound() {
    return React.useContext(FeaturesFoundContext)[1];
}

export function FeatureFoundProvider(props) {
    const ytdlFoundState = React.useState();
    const ffmpegFoundState = React.useState();

    return (
        <FeaturesFoundContext.Provider value={[ytdlFoundState, ffmpegFoundState]}>
            {props.children}
        </FeaturesFoundContext.Provider>
    )
}
