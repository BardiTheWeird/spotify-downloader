export async function authorizedFetch(input, init) {
    async function innerFunction() {
      let accessToken = localStorage.getItem('access token');
      if (!accessToken) {
        return null;
      }
      init = init || {};
      init.headers = {
        ...init.headers,
        'Authorization': `Bearer ${accessToken}`,
      }
  
      let response = await fetch(input, init);
      if (response.status === 401) {
        // refresh access token
        const refreshToken = localStorage.getItem('refresh token');
        if (!refreshToken) {
          return null;
        }
  
        const refreshResponse = await fetch(`https://accounts.spotify.com/api/token?grant_type=refresh_token&refresh_token=${refreshToken}&client_id=63d55a793f9c4a9e8d5aacba30069a23`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/x-www-form-urlencoded'
          }
        });
        if (refreshResponse.status !== 200) {
          return null;
        }
        const refreshResponseBody = refreshResponse.json();
        accessToken = refreshResponseBody.access_token;
        localStorage.setItem('access token', accessToken);
  
        response = await fetch(input, init);
      }
      return response;
    }
  
    const result = await innerFunction();
    // clean up token if they're invalid
    if (result === null) {
      localStorage.setItem('access token', '');
      localStorage.setItem('refresh token', '');
    }
    return result;
  }
