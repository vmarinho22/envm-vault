import { createRouter, createWebHistory } from 'vue-router';
import DefaultLayout from '@/app/layouts/DefaultLayout.vue';
import { useAuthStore } from '@/modules/auth/stores/auth.store';
import LoginView from '@/modules/auth/views/LoginView.vue';
import HomeView from '@/modules/home/views/HomeView.vue';

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: LoginView,
    },
    {
      path: '/',
      component: DefaultLayout,
      meta: { requiresAuth: true },
      children: [
        {
          path: '',
          name: 'home',
          component: HomeView,
        },
      ],
    },
  ],
});

router.beforeEach((to) => {
  const auth = useAuthStore();

  if (to.meta.requiresAuth && !auth.isAuthenticated) {
    return { name: 'login' };
  }

  if (to.name === 'login' && auth.isAuthenticated) {
    return { name: 'home' };
  }

  return true;
});

export default router;
