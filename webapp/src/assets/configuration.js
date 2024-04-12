export const configuration = {
  serverUrl: import.meta.env.MODE === 'development' ? `http://localhost/api/` : '/api/'
}
