import { createRouter, createWebHistory } from 'vue-router';

const router = createRouter({
	history: createWebHistory(import.meta.env.BASE_URL),
	routes: [
		{
			path: '/',
			name: 'main',
			component: () => import('../views/MainView.vue')
		},
		{
			path: '/auth',
			name: 'auth',
			component: () => import('../views/AuthView.vue')
		},
		{
			path: '/comment',
			name: 'comment',
			component: () => import('../views/CommentView.vue')
		}
	]
});

export default router;
