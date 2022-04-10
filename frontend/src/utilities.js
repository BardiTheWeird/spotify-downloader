export async function authorizedFetch(input, init) {
  const accessToken = localStorage.getItem('access token');
  if (accessToken) {
    init = init || {};
    init.headers = {
      ...init.headers,
      'Authorization': `Bearer ${accessToken}`,
    }
  }

  let response = await fetch(input, init);
  if (response.status >= 300) {
    localStorage.setItem('access token', '');
    return null;
  }
  return response;
}
