import { NetworkService } from './NetworkService';

export const MemoryService = {
	/**
	 * Sends a request to get questions based on the specified model, name, sex, and birth date.
	 *
	 * @param {string} model - The model to use for completion.
	 * @param {string} name - The name for the human info.
	 * @param {string} sex - The sex for the human info.
	 * @param {string} birth_date - The birth date for the human info.
	 * @param {function} success - The success callback function.
	 * @param {function} fail - The fail callback function.
	 */
	getQuestions: (model, name, sex, birth_date, success, fail) => {
		NetworkService.ClassicRequest(
			'POST',
			`v1/completion/${model}/questions`,
			{ human_info: { name, sex, birth_date } },
			response => {
				success(response.data);
			},
			fail
		);
	},
	/**
	 * Sends a request to get biography questions based on the specified model, type of story, name, sex, and birth date.
	 *
	 * @param {string} model - The model to use for completion.
	 * @param {string} type_of_story - The type of story for the biography.
	 * @param {string} name - The name for the human info.
	 * @param {string} sex - The sex for the human info.
	 * @param {string} birth_date - The birth date for the human info.
	 * @param {function} success - The success callback function.
	 * @param {function} fail - The fail callback function.
	 */
	getBiographyQuestions: (model, type_of_story, name, sex, birth_date, success, fail) => {
		NetworkService.ClassicRequest(
			'POST',
			`v1/completion/${model}/questions/biography`,
			{ type_of_story, human_info: { name, sex, birth_date } },
			response => {
				success(response.data);
			},
			fail
		);
	},
	/**
	 * Sends a request to get an epitaph based on the specified model, name, sex, birth date, death date, questions.
	 *
	 * @param {string} model - The model to use for completion.
	 * @param {string} name - The name for the human info.
	 * @param {string} sex - The sex for the human info.
	 * @param {string} birth_date - The birth date for the human info.
	 * @param {string} death_date - The death date for the human info.
	 * @param {Array} questions - The array of questions for the epitaph.
	 * @param {function} success - The success callback function.
	 * @param {function} fail - The fail callback function.
	 * @return {type} description of return value
	 */
	getEpitaph: (model, name, sex, birth_date, death_date, questions, success, fail) => {
		NetworkService.ClassicRequest(
			'POST',
			`v1/completion/${model}/epitaph`,
			{ human_info: { name, sex, birth_date, death_date, questions } },
			response => {
				success(response.data);
			},
			fail
		);
	},
	/**
	 * Sends a request to get a biography based on the specified model, type of story, name, sex, birth date, death date, questions, and previous information.
	 *
	 * @param {string} model - The model to use for completion.
	 * @param {string} type_of_story - The type of story for the biography.
	 * @param {string} name - The name for the human info.
	 * @param {string} sex - The sex for the human info.
	 * @param {string} birth_date - The birth date for the human info.
	 * @param {string} death_date - The death date for the human info.
	 * @param {Array} questions - The array of questions for the biography.
	 * @param {Object} previous - The previous information related to the biography.
	 * @param {function} success - The success callback function.
	 * @param {function} fail - The fail callback function.
	 */
	getBiography: (model, type_of_story, name, sex, birth_date, death_date, questions, previous, success, fail) => {
		NetworkService.ClassicRequest(
			'POST',
			`v1/completion/${model}/biography`,
			{ type_of_story, human_info: { name, sex, birth_date, death_date, questions }, previous },
			response => {
				success(response.data);
			},
			fail
		);
	},
	/**
	 * Sends a request to get a shortened biography based on the specified model, biography parts, name, sex, birth date, death date.
	 *
	 * @param {string} model - The model to use for completion.
	 * @param {Array} biography - The parts of the biography.
	 * @param {string} name - The name for the human info.
	 * @param {string} sex - The sex for the human info.
	 * @param {string} birth_date - The birth date for the human info.
	 * @param {string} death_date - The death date for the human info.
	 * @param {function} success - The success callback function.
	 * @param {function} fail - The fail callback function.
	 */
	getBiographyShort: (model, biography, name, sex, birth_date, death_date, success, fail) => {
		NetworkService.ClassicRequest(
			'POST',
			`v1/completion/${model}/biography/short`,
			{ part_young: biography[0].text, part_middle: biography[1].text, part_old: biography[2].text, human_info: { name, sex, birth_date, death_date } },
			response => {
				success(response.data);
			},
			fail
		);
	},
	/**
	 * Sends a request to get GPT based on the specified model and request message.
	 *
	 * @param {string} model - The model to use for completion.
	 * @param {string} request_message - The message used for the request.
	 * @param {function} success - The success callback function.
	 * @param {function} fail - The fail callback function.
	 * @return {type} description of return value
	 */
	getGPT: (model, request_message, success, fail) => {
		NetworkService.ClassicRequest(
			'POST',
			`v1/completion/${model}`,
			{ request_message },
			response => {
				success(response.data);
			},
			fail
		);
	}
};
