import axios from 'axios';

const base = axios.create({
  baseURL: process.env.REACT_APP_API_BASE_URL,
});

const fetchInbox = (address) => {
  return base.get(`/inboxes/${address}`);
}

const fetchMessage = (id) => {
  return base.get(`/messages/${id}`);
}

const api = {
  fetchInbox,
  fetchMessage,
}

export default api;