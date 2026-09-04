import { flushPromises } from '@vue/test-utils';
import { describe, expect, it } from 'vitest';
import DefaultLayout from '@/app/layouts/DefaultLayout.vue';
import { useAuthStore } from '@/modules/auth/stores/auth.store';
import LoginView from '@/modules/auth/views/LoginView.vue';
import HomeView from '@/modules/home/views/HomeView.vue';
import { expectComponentSnapshot, mountComponent } from '../../helpers/mount';

describe('DefaultLayout (component)', () => {
  const session = {
    token: 'jwt-token',
    user: {
      id: '1',
      email: 'admin@example.com',
      role: 'super_admin' as const,
    },
  };

  const routes = [
    {
      path: '/',
      component: DefaultLayout,
      children: [{ path: '', name: 'home', component: HomeView }],
    },
    { path: '/login', name: 'login', component: LoginView },
  ];

  it('renders sidebar with Início and top bar with user email', async () => {
    const { wrapper } = await mountComponent(DefaultLayout, {
      routes,
      initialRoute: '/',
      session,
    });

    expect(wrapper.get('aside[aria-label="Menu principal"]').text()).toContain('Início');
    expect(wrapper.get('header').text()).toContain('admin@example.com');
    expect(wrapper.get('header').text()).toContain('Sair');
    expectComponentSnapshot(wrapper);
  });

  it('logs out and navigates to login when Sair is clicked', async () => {
    const { wrapper, router, pinia } = await mountComponent(DefaultLayout, {
      routes,
      initialRoute: '/',
      session,
    });

    await wrapper.get('header button').trigger('click');
    await flushPromises();

    expect(useAuthStore(pinia).isAuthenticated).toBe(false);
    expect(router.currentRoute.value.name).toBe('login');
    expect(sessionStorage.getItem('envm-vault.auth')).toBeNull();
    expectComponentSnapshot(wrapper);
  });
});
