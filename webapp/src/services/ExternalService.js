import { NetworkService } from './NetworkService';

export const ExternalService = {
	/**
	 * Creates request options for a specific page with detailed information.
	 *
	 * @param {Page} page - The page object.
	 * @param {boolean} isMale - Indicates whether the person is male.
	 * @param {string} name - The name of the person.
	 * @param {string} epitaph - The epitaph text.
	 * @param {string} author_epitaph - The author of the epitaph.
	 * @param {string} firstName - The first name of the person.
	 * @param {string} lastName - The last name of the person.
	 * @param {string} surname - The surname of the person.
	 * @param {string} placeOfBirth - The place of birth.
	 * @param {string} placeOfDeath - The place of death.
	 * @param {string} children - The children of the person.
	 * @param {string} spouse - The spouse of the person.
	 * @param {string} citizenship - The citizenship of the person.
	 * @param {string} education - The education of the person.
	 * @param {string} occupation - The occupation of the person.
	 * @param {string} awards - The awards received by the person.
	 * @param {string} bio_title1 - The title of biography section 1.
	 * @param {string} bio1_text - The text for biography section 1.
	 * @param {string} bio_title2 - The title of biography section 2.
	 * @param {string} bio2_text - The text for biography section 2.
	 * @param {string} bio_title3 - The title of biography section 3.
	 * @param {string} bio3_text - The text for biography section 3.
	 * @param {string} bioend_text - The concluding text of the biography.
	 * @param {Date} birth - The birth date of the person.
	 * @param {Date} death - The death date of the person.
	 * @return {Object} The request options object with detailed data.
	 */
	createRequestOptions(
		cookies,
		page,
		isMale,
		name,
		epitaph,
		author_epitaph,
		firstName,
		lastName,
		surname,
		placeOfBirth,
		placeOfDeath,
		children,
		spouse,
		citizenship,
		education,
		occupation,
		awards,
		bio_title1,
		bio1_text,
		bio_title2,
		bio2_text,
		bio_title3,
		bio3_text,
		bioend_text,
		birth,
		death
	) {
		function formatDateToString(date) {
			const year = date.getFullYear();
			const month = ('0' + (date.getMonth() + 1)).slice(-2);
			const day = ('0' + date.getDate()).slice(-2);
			const hours = ('0' + date.getHours()).slice(-2);
			const minutes = ('0' + date.getMinutes()).slice(-2);
			const seconds = ('0' + date.getSeconds()).slice(-2);
			const formattedDate = `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
			return formattedDate;
		}
		function formatDate(date) {
			let day = String(date.getDate()).padStart(2, '0');
			let month = String(date.getMonth() + 1).padStart(2, '0');
			let year = date.getFullYear();
			return { day: day, month: month, year: year };
		}
		return {
			method: 'PUT',
			url: 'https://mc.dev.rand.agency/api/page/54553748',
			headers: {
				'User-Agent': 'Thunder Client (https://www.thunderclient.com)',
				Accept: 'application/json',
				'Content-Type': 'application/json;charset=UTF-8',
				Authorization: 'Bearer ' + cookies.get('token')
			},
			data: {
				id: page.id,
				name: name,
				surname: surname,
				patronym: null,
				birthday_at: formatDateToString(birth),
				died_at: formatDateToString(death),
				start: formatDate(birth),
				end: formatDate(death),
				epitaph: epitaph,
				author_epitaph: author_epitaph,
				page_type_id: 1,
				slug: page.slug,
				filled_fields: ['biography_1', 'biography_2', 'biography_3', 'end_of_biography', 'epitaph'],
				lastName: lastName,
				firstName: firstName,
				link: page.link,
				full_name: lastName + ' ' + firstName + ' ' + surname,
				page_type_name: 'pageType.full',
				biographies: [
					{
						id: null,
						title: bio_title1,
						description: bio1_text,
						page_id: page.id,
						checked: true
					},
					{
						id: null,
						title: bio_title2,
						description: bio2_text,
						page_id: page.id,
						checked: true
					},
					{
						id: null,
						title: bio_title3,
						description: bio3_text,
						page_id: page.id,
						checked: true
					},
					{
						id: null,
						title: 'Заключение',
						description: bioend_text,
						page_id: page.id,
						checked: true
					}
				],
				page_information: [
					{ page_id: page.id, title: 'pageInformation.placeOfBirth', description: placeOfBirth },
					{ page_id: page.id, title: 'pageInformation.placeOfDeath', description: placeOfDeath },
					{ page_id: page.id, title: 'pageInformation.children', description: children },
					{ page_id: page.id, title: isMale ? 'pageInformation.wife' : 'pageInformation.husband', description: spouse },
					{ page_id: page.id, title: 'pageInformation.citizenship', description: citizenship },
					{ page_id: page.id, title: 'pageInformation.education', description: education },
					{ page_id: page.id, title: 'pageInformation.occupation', description: occupation },
					{ page_id: page.id, title: 'pageInformation.awards', description: awards }
				]
			}
		};
	},
	/**
	 * Sends a page using the specified request options.
	 *
	 * @param {Object} requestOptions - The request options, including the method, URL, and data.
	 * @param {function} success - The callback function to be executed if the request is successful.
	 * @param {function} fail - The callback function to be executed if the request fails.
	 * @return {Promise} The response from the request.
	 */
	sendPage(requestOptions, success, fail) {
		return NetworkService.RawRequest(
			requestOptions,
			response => {
				success(response.data);
			},
			fail
		);
	},
	/**
	 * A description of the entire function.
	 *
	 * @param {function} success - The callback function to be executed if the request is successful.
	 * @param {function} fail - The callback function to be executed if the request fails.
	 * @param {Object} cookies - The cookies object containing the token for authorization.
	 * @return {Promise} The response from the request.
	 */
	getPages(success, fail, cookies) {
		return NetworkService.AuthRequest(
			'GET',
			`cabinet/individual-pages`,
			null,
			cookies,
			response => {
				success(response.data);
			},
			fail
		);
	},
	/**
	 * Sends a report with page ID, user information, text, and relation role.
	 *
	 * @param {type} page_id - The ID of the page.
	 * @param {type} fio - The user's full name.
	 * @param {type} email - The user's email address.
	 * @param {type} text - The report text.
	 * @param {type} relation_role - The role of the user in relation to the reported content.
	 * @param {function} success - The callback function to be called on successful report submission.
	 * @param {function} fail - The callback function to be called on failed report submission.
	 * @return {type} description of return value
	 */
	sendReport(page_id, fio, email, text, relation_role, success, fail) {
		return NetworkService.DirectRequest(
			'POST',
			`comment`,
			{ page_id, fio, email, text, relation_role, checked: true, hasEmail: email.length > 0 },
			response => {
				success(response.data);
			},
			fail
		);
	}
};
