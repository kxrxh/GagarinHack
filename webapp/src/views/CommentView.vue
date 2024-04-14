<script setup>
import UILabel from '@/components/ui/UILabel.vue';
import UILabeledInput from '@/components/ui/UILabeledInput.vue';
import UIButton from '@/components/ui/UIButton.vue';
import UIDropdown from '../components/ui/UIDropdown.vue';
</script>

<template>
	<section class="bg-gray-50 dark:bg-gray-900 h-screen fixed overflow-auto w-screen">
		<div class="flex flex-col items-center justify-center px-6 py-8 mx-auto md:h-screen lg:py-0">
			<div class="w-full bg-white rounded-lg shadow dark:border md:mt-0 sm:max-w-md xl:p-0 dark:bg-gray-800 dark:border-gray-700">
				<div class="p-6 space-y-4 md:space-y-6 sm:p-8">
					<h1 class="text-xl font-bold leading-tight tracking-tight text-gray-900 md:text-2xl dark:text-white text-center">Добавление отзыва</h1>
					<div class="space-y-4 md:space-y-6">
						<div>
							<UILabel> Выберите страницу </UILabel>
							<UIDropdown v-model="selectedPage" :options="wrappedPages"> Страница не выбрана </UIDropdown>
						</div>
						<UILabeledInput v-model="fio" type="text" property="fio" placeholder="Иван Иванович"> Ваше имя и фамилия </UILabeledInput>
						<UILabeledInput v-model="email" type="email" property="email" placeholder="example@email.com"> Ваша почта </UILabeledInput>
						<UILabeledInput v-model="relative" type="text" property="relative" placeholder="Близкий друг"> Кем вы приходитесь? </UILabeledInput>
						<UILabeledInput v-model="words" :textarea="true" :disabled="generating" type="text" property="words" placeholder="Ты останешься навсегда в моем сердце">
							Ваши слова
							<UIButton @click="generate" :disabled="generating || !wrappedPages?.hasOwnProperty(selectedPage)" classExtension="w-full px-5 py-2.5"> Сгенерировать </UIButton>
						</UILabeledInput>
						<UIButton @click="add" :disabled="loading || !isInputsSet" classExtension="w-full px-5 py-2.5"> Добавить отзыв </UIButton>
					</div>
				</div>
			</div>
		</div>
	</section>
</template>

<script>
import { useWebApp } from 'vue-tg';
import { MemoryService } from '@/services/MemoryService';
import { ExternalService } from '@/services/ExternalService';
export default {
	name: 'CommentView',
	data() {
		return {
			pages: [],
			wrappedPages: {},
			selectedPage: null,

			fio: '',
			email: '',
			relative: '',
			words: '',

			loading: false,
			generating: false
		};
	},
	computed: {
		isInputsSet() {
			return this.fio.length > 0 && this.relative.length > 0 && this.words.length > 0 && this.wrappedPages.hasOwnProperty(this.selectedPage);
		}
	},
	methods: {
		add() {
			if (!this.wrappedPages.hasOwnProperty(this.selectedPage)) {
				this.$notify({ text: 'Выберите страницу', type: 'error' });
				return;
			}
			if (this.fio.length < 1) {
				this.$notify({ text: 'Введите ФИО', type: 'error' });
				return;
			}
			if (this.relative.length < 1) {
				this.$notify({ text: 'Введите кем вы являетесь', type: 'error' });
				return;
			}
			if (this.words.length < 1) {
				this.$notify({ text: 'Введите ваши слова', type: 'error' });
				return;
			}
			ExternalService.sendReport(
				this.selectedPage,
				this.fio,
				this.email,
				this.words,
				this.relative,
				data => {
					this.$notify({ text: 'Отзыв добавлен', type: 'success' });
					useWebApp().sendData(
						JSON.stringify({
							type: 'success_comment'
						})
					);
					//useWebApp().close();
				},
				err => {
					this.$notify({ text: 'Не удалось добавить отзыв, попробуйте позже.', type: 'error' });
				}
			);
		},
		/**
		 * Generate the final words for a deceased person based on user input and selected page.
		 */
		generate() {
			this.words = '';
			this.generating = true;
			let page = null;
			for (let i of this.pages) {
				if (i.id == this.selectedPage) {
					page = i;
					break;
				}
			}
			MemoryService.getGPT(
				'gigachat',
				`Сгенерируй последнее прощание от ${this.relative} ${this.fio}, для умершего человека по имени ${page.name}`,
				data => {
					this.words = data.response;
					this.generating = false;
				},
				err => {
					this.$notify({ text: 'Не удалось сгенерировать эпитафию с использованием Gigachat, попробуйте позже.', type: 'error' });
					this.generating = false;
				}
			);
		}
	},
	/**
	 * Fetches pages if the user is authenticated and handles the result and error cases.
	 */
	mounted() {
		if (!this.$cookies.get('token')) {
			sessionStorage.setItem('fallback', 'comment');
			this.$router.push({ path: 'auth' });
		} else {
			ExternalService.getPages(
				data => {
					this.pages = data;
					this.wrappedPages = data.reduce((acc, obj) => {
						if (obj.page_type_name === 'pageType.full') {
							acc[obj.id] = obj.name;
						}
						return acc;
					}, {});
					console.log(this.wrappedPages);
				},
				err => {
					this.$notify({ text: 'Не удалось загрузить страницы, попробуйте позже.', type: 'error' });
					console.log(err);
					this.$cookies.remove('token');
					sessionStorage.setItem('fallback', 'comment');
					this.$router.push({ path: 'auth' });
				},
				this.$cookies
			);
		}
	}
};
</script>
