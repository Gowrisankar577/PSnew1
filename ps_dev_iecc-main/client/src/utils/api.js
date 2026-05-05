import axios from "axios";
import { appBase } from "./settings";

const apiUrl = process.env.REACT_APP_API;


const apiPostRequest = async (url, data, isMultipart = false) => {
  try {
    const res = await axios.post(`${apiUrl}${url}`, data, {
      withCredentials: true,
      headers: {
        ...(isMultipart ? {} : { "Content-Type": "application/json" }),
      },
    });
    if (res.status === 200) {
      return { success: true, data: res.data };
    } else {
      if (res.status === 401) {
        window.location = appBase + "/auth/login";
      }
      return { success: false, error: res.data };
    }
  } catch (error) {
    if (error.status === 401) {
      window.location = appBase + "/auth/login";
    }
    return { success: false, error: error.response.data.error };
  }
};

const apiGetRequest = async (url) => {
  try {
    const res = await axios.get(`${apiUrl}${url}`, { withCredentials: true });
    if (res.status === 200) {
      return { success: true, data: res.data };
    } else {
      if (res.status === 401) {
        window.location = appBase + "/auth/login";
      }
      return { success: false, error: res.data };
    }
  } catch (error) {
    if (error.status === 401) {
      window.location = appBase + "/auth/login";
    }
    console.log(error);
    return { success: false, error: error.response.data.error };
  }
};

const apiDeleteRequest = async (url, data = null) => {
  try {
    const config = {
      withCredentials: true,
    };
    
    // Add data to config if provided
    if (data) {
      config.data = data;
    }
    
    const res = await axios.delete(`${apiUrl}${url}`, config);
    if (res.status === 200) {
      return { success: true, data: res.data };
    } else {
      if (res.status === 401) {
        window.location = appBase + "/auth/login";
      }
      return { success: false, error: res.data };
    }
  } catch (error) {
    if (error.status === 401) {
      window.location = appBase + "/auth/login";
    }
    return { success: false, error: "Something went wrong" };
  }
};

const apiPutRequest = async (url, data) => {
  try {
    const res = await axios.put(`${apiUrl}${url}`, data, {
      withCredentials: true,
    });
    if (res.status === 200) {
      return { success: true, data: res.data };
    } else {
      if (res.status === 401) {
        window.location = appBase + "/auth/login";
      }
      return { success: false, error: res.data };
    }
  } catch (error) {
    if (error.status === 401) {
      window.location = appBase + "/auth/login";
    }
    return { success: false, error: error.response.data.error };
  }
};
export {
  apiPostRequest,
  apiGetRequest,
  apiPutRequest,
  apiUrl,
  apiDeleteRequest,
};
