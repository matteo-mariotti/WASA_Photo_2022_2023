import axios from "axios";

const backend = axios.create({
  baseURL: __API_URL__,
  timeout: 1000 * 5,
});

backend.interceptors.request.use(config => {
  config.headers["Authorization"] = "Bearer " + sessionStorage.getItem("token");
  return config;
});

export default backend;
