const { app, BrowserWindow, ipcMain, dialog } = require('electron');
const {spawn} = require('child_process');
const path = require('path');
const kill = require('tree-kill');
const isDev = require('electron-is-dev');

const resourcesDir = process.resourcesPath;

const serve = require('electron-serve');
const serveDirectory = serve({directory: path.join(resourcesDir, 'front')});

const appUrl = isDev
  ? 'http://localhost:3000'
  : 'app://-';
const backendExecutablePath = isDev
  ? '../backend/build/backend.exe'
  : path.join(resourcesDir, 'backend.exe');

let backendStatus;
const backend = spawn(backendExecutablePath);

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
