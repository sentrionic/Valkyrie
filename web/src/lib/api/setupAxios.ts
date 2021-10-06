import Axios from 'axios';

export const request = Axios.create({
  baseURL: `${process.env.REACT_APP_API!}/api`,
  withCredentials: true,
});
