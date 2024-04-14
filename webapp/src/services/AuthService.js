import { NetworkService } from './NetworkService';

export const AuthService = {
	/**
	 * Logs in the user with the provided credentials.
	 *
	 * @param {string} email - The email address of the user.
	 * @param {string} password - The password of the user.
	 * @param {function} success - The callback function to be called on successful login.
	 * @param {function} fail - The callback function to be called on failed login.
	 */
	login(email, password, success, fail) {
		NetworkService.DirectRequest(
			'POST',
			'v1/get-access-token',
			{ email, password, device: 'bot-v0.0.1' },
			response => {
				if (response.data) success(response.data);
				else fail();
			},
			fail
		);
	}
};
