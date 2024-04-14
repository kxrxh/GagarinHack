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
		NetworkService.RawRequest(
			{
				method: method,
				url: `${configuration.serverUrl}${url}`,
				data: data,
				headers: { 'Content-Type': contentType }
			},
			success,
			fail
		);
	},
	/**
	 * Executes an authenticated HTTP request using RawRequest.
	 *
	 * @param {string} method - The HTTP method to use for the request.
	 * @param {string} url - The URL to send the request to.
	 * @param {Object} data - The data to send with the request.
	 * @param {Object} cookies - The cookies object containing the token for authorization.
	 * @param {function} success - The callback function to be executed if the request is successful.
	 * @param {function} fail - The callback function to be executed if the request fails.
	 */
	AuthRequest(method, url, data, cookies, success, fail) {
		NetworkService.RawRequest(
			{
				method: method,
				url: `${configuration.directUrl}${url}`,
				data: data,
				headers: { Authorization: `Bearer ${cookies.get('token')}` }
			},
			success,
			fail
		);
	},
	/**
	 * Executes a direct HTTP request using RawRequest.
	 *
	 * @param {string} method - The HTTP method to use for the request.
	 * @param {string} url - The URL to send the request to.
	 * @param {Object} data - The data to send with the request.
	 * @param {function} success - The callback function to be executed if the request is successful.
	 * @param {function} fail - The callback function to be executed if the request fails.
	 * @param {string} [contentType='application/json'] - The content type of the request.
	 */
	DirectRequest(method, url, data, success, fail, contentType = 'application/json') {
		NetworkService.RawRequest(
			{
				method: method,
				url: `${configuration.directUrl}${url}`,
				data: data,
				headers: { 'Content-Type': contentType }
			},
			success,
			fail
		);
	},
	/**
	 * Executes a raw HTTP request using Axios.
	 *
	 * @param {Object} options - The request options, including the method, URL, and data.
	 * @param {function} success - The callback function to be executed if the request is successful.
	 * @param {function} fail - The callback function to be executed if the request fails.
	 * @return {Promise} The response from the request.
	 */
	RawRequest(options, success, fail) {
		(async () => {
			let response = await axios(options).catch(fail);
			if (response) success(response);
		})();
	}
};
