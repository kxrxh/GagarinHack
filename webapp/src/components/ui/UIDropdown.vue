<script setup>
import DownArrowIcon from '@/assets/icons/DownArrowIcon.vue';
</script>
<template>
	<div class="w-full p-0 m-0">
		<div :class="['relative', showSelector ? 'z-20' : '']">
			<button
				class="flex-shrink-0 inline-flex items-center py-2.5 px-4 text-sm font-medium text-center text-gray-500 bg-gray-100 border border-gray-300 rounded-lg hover:bg-gray-200 focus:ring-4 focus:outline-none focus:ring-gray-100 dark:bg-gray-700 dark:hover:bg-gray-600 dark:focus:ring-gray-700 dark:text-white dark:border-gray-600 w-full disabled:opacity-50 disabled:cursor-not-allowed transition-opacity"
				type="button"
				@click="showSelector = !showSelector"
				ref="button"
				:disabled="disabled"
			>
				<slot v-if="!options || options.length <= 0 || !modelValue || modelValue.length === 0 || modelValue * 1 === 0 || !options[modelValue]"></slot>
				<span>
					{{ options[modelValue] }}
				</span>
				<span class="ml-auto">
					<DownArrowIcon />
				</span>
			</button>
		</div>
		<div
			v-if="showSelector"
			class="z-10 bg-white rounded-lg shadow dark:bg-gray-700 absolute animate-fade-down animate-once animate-duration-200 animate-ease-in"
			:style="`width: ${$refs.button.offsetWidth}px;`"
		>
			<div class="h-3"></div>
			<ul v-if="optionsFiltered.length > 0" class="max-h-48 px-3 pb-3 overflow-y-auto text-sm text-gray-700 dark:text-gray-200">
				<li v-for="key in optionsFiltered" :value="key" :key="key" @click="select(key)" class="cursor-pointer">
					<div class="flex items-center ps-2 rounded hover:bg-gray-100 dark:hover:bg-gray-600">
						<span class="py-2 ms-2 text-sm font-medium text-gray-900 rounded dark:text-gray-300">
							{{ options[key] }}
						</span>
					</div>
				</li>
			</ul>
		</div>
	</div>
</template>
<script>
export default {
	name: 'UIDropdownWithSearch',
	props: {
		modelValue: null,
		options: {
			type: Object,
			default: () => ({})
		},
		disabled: {
			type: Boolean,
			default: false
		}
	},
	computed: {
		/**
		 * Filters the options object based on the search string.
		 */
		optionsFiltered() {
			return Object.keys(this.options).filter(key => this.options[key].toLowerCase().includes(this.search.toLowerCase()));
		}
	},
	emits: ['update:modelValue', 'changed'],
	data() {
		return {
			showSelector: false,
			search: ''
		};
	},
	methods: {
		select(key) {
			this.$emit('update:modelValue', key);
			if (this.modelValue != key) this.$emit('changed');
			this.showSelector = false;
		},
		/**
		 * Handles the click event outside of the component.
		 */
		handleClick() {
			const container = this.$el;
			if (!container.contains(event.target)) {
				this.showSelector = false;
			}
		}
	},
	mounted() {
		document.addEventListener('click', this.handleClick);
	},
	unmounted() {
		document.removeEventListener('click', this.handleClick);
	}
};
</script>
