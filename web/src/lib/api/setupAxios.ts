import Axios, { AxiosError } from 'axios';

export const request = Axios.create({
  baseURL: process.env.REACT_APP_API_URL,
  withCredentials: true,
});

request.interceptors.response.use(async response => {
  return response;
}, (error: AxiosError) => {
  const { status } = error.response!;
  switch (status) {
    case 400:
      break;
    case 403:
      localStorage.removeItem("user-storage");
      window.location.replace('/login');
      break;
    case 404:
      break;
    case 500:
      break;
  }
  return Promise.reject(error);
})
