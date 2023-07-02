import axios from "axios";

import env from "../../env.json";
import Utils from "../share/Utils";

const Result = {
    SUCCESS: 1,
    FAIL: 2
}

const baseApi = axios.create({
    baseURL: `${env.API_BASE_URL}/api/v1`,
    headers: {
        "content-type": "application/json",
    },
});

baseApi.interceptors.request.use(async (config) => {
    // debug call api
    console.log(
        `[Axios send request]\n\t- [URL] ${config.baseURL}${config.url},\n\t- [Method]: ${config.method.toUpperCase()},\n\t- [Payload]: ${config.method === "post" || config.method === "put"
            ? JSON.stringify(config.data)
            : JSON.stringify(config.params)
        }`
    );
    return config;
});

baseApi.interceptors.response.use(
    (response) => {
        console.log(JSON.stringify(response));
        return {
            result: Result.SUCCESS,
            data: response.data,
            headers: response.headers
        }
    },
    async (error) => {
        // debugger
        // console.log(JSON.stringify(error));
        // console.log("==========================");
        // console.log(error.response);
        // console.log("==========================");
        if (!error.response) {
            Utils.showCheckNetwork();
            return {
                result: Result.FAIL,
                data: error.response.data
            }
        }
        if (error.response) {
            return {
                result: Result.FAIL,
                data: error.response.data
            }
        }
    }
);

export default baseApi
export { Result }