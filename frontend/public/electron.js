const { app, BrowserWindow, ipcMain, dialog } = require('electron');
const { spawn, spawnSync } = require('child_process');
const path = require('path');
const kill = require('tree-kill');
const isDev = require('electron-is-dev');
const os = require('os');
const fs = require('fs');

const fetch = require('make-fetch-happen');

const isWin = os.platform() === 'win32';
const excutableExtension = isWin && '.exe' || '';

const resourcesDir = process.resourcesPath;
// The directory for storing your app's configuration files, which by default it is the appData directory appended with your app's name
const userData = app.getPath('userData');
const userSettingsPath = path.join(userData, 'settings.json');

const serve = require('electron-serve');
const serveDirectory = serve({directory: path.join(resourcesDir, 'front')});

const appUrl = isDev
  ? 'http://localhost:3000'
  : 'app://-';
const backendExecutablePath = isDev
  ? '../backend/build/backend' + excutableExtension
  : path.join(resourcesDir, 'backend' + excutableExtension);

let backendStatus;
function getBaseUrl() {
  if (backendStatus && backendStatus.running) {
    return backendStatus.address;
  }
  return null;
}

if (isDev) {
  console.log('building backend...')
  const buildBackend = spawnSync("go",
    ['build', '-o', './build/backend' + excutableExtension],
    {
      cwd: '../backend'
    });

  if (buildBackend.error) {
    console.log('error building backend:', buildBackend.error);
  }
  else {
    if (buildBackend.status !== 0) {
      console.log('status:', buildBackend.status);
      console.log('output:', buildBackend.output.join('\n'));
    }
  }
}

const backend = spawn(backendExecutablePath, [
  "--settings-path", userSettingsPath
]);

backend.on('error', err => {
  backendStatus = {
    running: false
  }
  console.log('BACKEND error:', err);
})

const portRegex = /listening on port (\d+)/
backend.stdout.on('data', data => {
  const stringData = data.toString();
  const match = stringData.match(portRegex);
  if (match && match.length > 1) {
    const port = match[1];
    backendStatus = {
      running: true,
      address: `http://localhost:${port}/api/v1`
    }
    console.log('backend status:', backendStatus);
  }
  console.log('BACKEND:', stringData);
});
backend.stderr.on('data', data => {
  console.log('BACKEND:', data.toString());
});

ipcMain.handle('configureFeaturePath', async (_, featureName) => {
  const path = dialog.showOpenDialogSync()[0];
  console.log(`${featureName} path is:`, path);
  
  const response = await fetch(`${getBaseUrl()}/configure/${featureName}?path=${path}`, {
    method: 'POST'
  });

  if (response.status == 204) {
    return true;
  }
  else if (response.status == 400 || response.status == 404) {
    return false;
  }
});
ipcMain.handle('openDirectory', () => {
  const path = dialog.showOpenDialogSync({
    properties: ["openDirectory"]
  });
  return path;
});
ipcMain.handle('backendStatus', () => backendStatus);
ipcMain.handle('appUrl', () => appUrl);

function createWindow() {
  // Create the browser window.
  const win = new BrowserWindow({
    width: 800,
    height: 1080,
    webPreferences: {
      nodeIntegration: true,
      enableRemoteModule: true,
      contextIsolation: false,
    },
  });

  // and load the index.html of the app.
  if (isDev) {
    win.loadURL(appUrl);
  }
  else {
    serveDirectory(win);
  }
  
  // Open the DevTools.
  if (isDev) {
    win.webContents.openDevTools();
  }
}

// This method will be called when Electron has finished
// initialization and is ready to create browser windows.
// Some APIs can only be used after this event occurs.
app.whenReady().then(createWindow);

// Quit when all windows are closed, except on macOS. There, it's common
// for applications and their menu bar to stay active until the user quits
// explicitly with Cmd + Q.
app.on('window-all-closed', () => {
  if (process.platform !== 'darwin') {
    try {
      kill(backend.pid);
    }
    catch (error) {
      console.log("error killing a backend process:", error);
    }
    app.quit();
  }
});

app.on('activate', () => {
  if (BrowserWindow.getAllWindows().length === 0) {
    createWindow();
  }
});
