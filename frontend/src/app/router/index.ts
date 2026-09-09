import { createRouter, createWebHistory } from 'vue-router';
import DefaultLayout from '@/app/layouts/DefaultLayout.vue';
import { useAuthStore } from '@/modules/auth/stores/auth.store';
import LoginView from '@/modules/auth/views/LoginView.vue';
import ApiKeyView from '@/modules/token/views/ApiKeyView.vue';
import TokenSetupView from '@/modules/token/views/TokenSetupView.vue';
import UsersView from '@/modules/user/views/UsersView.vue';
import AuditLogsView from '@/modules/audit/views/AuditLogsView.vue';
import TeamsView from '@/modules/team/views/TeamsView.vue';
import TagsView from '@/modules/tag/views/TagsView.vue';
import EnvironmentsView from '@/modules/environment/views/EnvironmentsView.vue';
import ProjectsView from '@/modules/project/views/ProjectsView.vue';
import ProjectDetailView from '@/modules/secret/views/ProjectDetailView.vue';
import SecretsView from '@/modules/secret/views/SecretsView.vue';
import SecretDetailView from '@/modules/secret/views/SecretDetailView.vue';

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: LoginView,
    },
    {
      path: '/setup-token',
      name: 'setup-token',
      component: TokenSetupView,
      meta: { requiresAuth: true },
    },
    {
      path: '/',
      component: DefaultLayout,
      meta: { requiresAuth: true, requiresApiToken: true },
      children: [
        {
          path: '',
          name: 'home',
          component: ProjectsView,
        },
        {
          path: 'project/:slug',
          name: 'project-detail',
          component: ProjectDetailView,
        },
        {
          path: 'api-key',
          name: 'api-key',
          component: ApiKeyView,
        },
        {
          path: 'teams',
          name: 'teams',
          component: TeamsView,
        },
        {
          path: 'tags',
          name: 'tags',
          component: TagsView,
        },
        {
          path: 'environments',
          name: 'environments',
          component: EnvironmentsView,
        },
        {
          path: 'secrets',
          name: 'secrets',
          component: SecretsView,
        },
        {
          path: 'secrets/:id',
          name: 'secret-detail',
          component: SecretDetailView,
        },
        {
          path: 'users',
          name: 'users',
          component: UsersView,
          meta: { requiresAdmin: true },
        },
        {
          path: 'audit-logs',
          name: 'audit-logs',
          component: AuditLogsView,
          meta: { requiresAdmin: true },
        },
      ],
    },
  ],
});

router.beforeEach(async (to) => {
  const auth = useAuthStore();

  if (to.meta.requiresAuth && !auth.isAuthenticated) {
    return { name: 'login' };
  }

  if (auth.isAuthenticated) {
    await auth.ensureHasApiToken();

    if (!auth.hasApiToken && to.name !== 'setup-token') return { name: 'setup-token' };
    if (auth.hasApiToken && to.name === 'setup-token') return { name: 'home' };
  }

  if (to.name === 'login' && auth.isAuthenticated) {
    const name = auth.hasApiToken ? 'home' : 'setup-token';
    return { name };
  }

  if (to.meta.requiresAdmin && !auth.isAdmin) return { name: 'home' };

  return true;
});

export default router;
