const { app, BrowserWindow, ipcMain, dialog } = require('electron');
const {spawn, exec} = require('child_process');
const kill = require('tree-kill');
const isDev = require('electron-is-dev');

const child = spawn('../backend/build/main.exe', 
    ['-settings', '../backend/settings.json']
);
child.stdout.on('data', data => {
    console.log('BACKEND:', data.toString());
});
child.stderr.on('data', data => {
    console.log('BACKEND:', data.toString());
});

ipcMain.on('openDirectory', (e, a) => {
    const path = dialog.showOpenDialogSync({
        properties: ["openDirectory"]
    });
    e.sender.send('returnDirectory', path);
})

function createWindow() {
  // Create the browser window.
  const win = new BrowserWindow({
    width: 800,
    height: 600,
    webPreferences: {
      nodeIntegration: true,
      enableRemoteModule: true,
      contextIsolation: false,
    },
  });

  // and load the index.html of the app.
  // win.loadFile("index.html");
  win.loadURL(
    isDev
      ? 'http://localhost:3000'
      : `file://${path.join(__dirname, '../build/index.html')}`
  );
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
    console.log('child.pid', child.pid);
    kill(child.pid);
    app.quit();
    // spawn("powershell", ['taskkill', '/pid', child.pid, '/F', '/T']);
  }
});

app.on('activate', () => {
  if (BrowserWindow.getAllWindows().length === 0) {
    createWindow();
  }
});
