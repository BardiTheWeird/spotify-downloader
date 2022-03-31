import React from "react";

const {ipcRenderer} = window.require('electron');

const BaseUrlContex = React.createContext();

export function useBaseUrl() {
    return React.useContext(BaseUrlContex);
}

export function BaseUrlProvider(props) {
    const [baseUrl, updateBaseUrl] = React.useState();
  React.useEffect(() => {
    (async () => {
      while (true) {
        const backendStatus = await ipcRenderer.invoke('backendStatus');
        if (backendStatus) {
          if (backendStatus.running) {
            updateBaseUrl(backendStatus.address);
          }
          else {
            updateBaseUrl(null);
          }
          break;
        }
        await new Promise(r => setTimeout(r, 500));
      }
    })();
  }, []);

  console.log('backend base url:', baseUrl);

  return (
    <BaseUrlContex.Provider value={baseUrl}>
        {props.children}
    </BaseUrlContex.Provider>
  );
}
