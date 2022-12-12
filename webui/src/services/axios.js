import axios from "axios";

const instance = axios.create({
  baseURL: __API_URL__,
  timeout: 1000 * 5,
});

/* instance.interceptors.request.use((config) => {
  config.headers["Authorization"] =
    "Bearer" + localStorage.getItem("access_token");
  return config;
}); */

/* 
instance.interceptors.request.use((request) => {
  console.log("Starting Request", JSON.stringify(request, null, 2));
  return request;
});  */

export default instance;
