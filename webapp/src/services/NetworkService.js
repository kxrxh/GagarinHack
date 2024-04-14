import { configuration } from '@/assets/configuration';

import axios from 'axios';

export const NetworkService = {
    /**
     * Executes a HTTP request using Axios.
     *
     * @param {string} method - The HTTP method to use for the request.
     * @param {string} url - The URL to send the request to.
     * @param {Object} data - The data to send with the request.
     * @param {function} success - The callback function to be executed if the request is successful.
     * @param {function} fail - The callback function to be executed if the request fails.
     */
    ClassicRequest(method, url, data, success, fail, contentType = 'application/json') {
        NetworkService.RawRequest({
            method: method,
            url: `${configuration.serverUrl}${url}`,
            data: data,
            headers: { 'Content-Type': contentType }
        }, success, fail);
    },
    AuthRequest(method, url, data, cookies, success, fail) {
        NetworkService.RawRequest({
            method: method,
            url: `${configuration.directUrl}${url}`,
            data: data,
            headers: {'Authorization': `Bearer ${cookies.get("token")}`}
        }, success, fail);
    },
    DirectRequest(method, url, data, success, fail, contentType = 'application/json') {
        NetworkService.RawRequest({
            method: method,
            url: `${configuration.directUrl}${url}`,
            data: data,
            headers: { 'Content-Type': contentType }
        }, success, fail);
    },
    RawRequest(options, success, fail) {
        (async () => {
            let response = await axios(options).catch(fail);
            if (response) success(response);
        })();
    }
};
