<template>
	<button
		type="button"
		:class="`items-center justify-center ${classExtension.includes('text') ? '' : 'text-white'} bg-${color}-600 hover:bg-${color}-700 focus:ring-4 focus:ring-${color}-300 font-medium rounded-lg text-sm px-4 py-2 dark:bg-${color}-600 dark:hover:bg-${color}-500 focus:outline-none dark:focus:ring-${color}-700 disabled:opacity-25 transition-opacity ${classExtension}`"
		:disabled="disabled"
	>
		<slot></slot>
	</button>
</template>

<script>
export default {
	name: 'UIButton',
	props: {
		color: {
			type: String,
			default: 'primary'
		},
		classExtension: {
			type: String,
			default: ''
		},
		disabled: {
			type: null,
			default: false
		}
	}
};
</script>
