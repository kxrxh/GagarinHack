<script setup>
import { RouterView } from 'vue-router';
import NotificationModal from '@/components/modal/NotificationModal.vue';
</script>

<template>
	<RouterView />
	<notifications class="p-3" position="top left">
		<template #body="props">
			<NotificationModal :data="props" />
		</template>
	</notifications>
</template>

<style>
*::-webkit-scrollbar-track {
	-webkit-box-shadow: inset 0 0 6px rgba(0, 0, 0, 0.3);
	background-color: transparent;
	border-radius: 10px;
}

*::-webkit-scrollbar {
	width: 10px;
	background-color: transparent;
}

*::-webkit-scrollbar-thumb {
	border-radius: 10px;
	background-image: -webkit-gradient(linear, left bottom, left top, color-stop(0.44, rgb(122, 153, 217)), color-stop(0.72, rgb(73, 125, 189)), color-stop(0.86, rgb(28, 58, 148)));
}
</style>
