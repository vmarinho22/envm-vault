import { flushPromises } from '@vue/test-utils';
import { describe, expect, it } from 'vitest';
import DefaultLayout from '@/app/layouts/DefaultLayout.vue';
import { useAuthStore } from '@/modules/auth/stores/auth.store';
import LoginView from '@/modules/auth/views/LoginView.vue';
import HomeView from '@/modules/home/views/HomeView.vue';
import { expectComponentSnapshot, mountComponent } from '../../helpers/mount';

describe('DefaultLayout (component)', () => {
  const adminSession = {
    token: 'jwt-token',
    user: {
      id: '1',
      email: 'admin@example.com',
      role: 'super_admin' as const,
    },
  };

  const userSession = {
    token: 'jwt-token',
    user: {
      id: '2',
      email: 'user@example.com',
      role: 'user' as const,
    },
  };

  const routes = [
    {
      path: '/',
      component: DefaultLayout,
      children: [
        { path: '', name: 'home', component: HomeView },
        { path: 'api-key', name: 'api-key', component: { template: '<div />' } },
        { path: 'teams', name: 'teams', component: { template: '<div />' } },
        { path: 'users', name: 'users', component: { template: '<div />' } },
        { path: 'audit-logs', name: 'audit-logs', component: { template: '<div />' } },
      ],
    },
    { path: '/login', name: 'login', component: LoginView },
  ];

  it('renders sidebar with Início, Chave de API, Times, Usuários and Auditoria for admin', async () => {
    const { wrapper } = await mountComponent(DefaultLayout, {
      routes,
      initialRoute: '/',
      session: adminSession,
    });

    const nav = wrapper.get('aside[aria-label="Menu principal"]');
    expect(nav.text()).toContain('Início');
    expect(nav.text()).toContain('Chave de API');
    expect(nav.text()).toContain('Times');
    expect(nav.text()).toContain('Usuários');
    expect(nav.text()).toContain('Auditoria');
    expect(wrapper.get('header').text()).toContain('admin@example.com');
    expectComponentSnapshot(wrapper);
  });

  it('highlights only Usuários when that route is active', async () => {
    const { wrapper } = await mountComponent(DefaultLayout, {
      routes,
      initialRoute: '/users',
      session: adminSession,
    });

    const links = wrapper.findAll('aside[aria-label="Menu principal"] nav a');
    const homeLink = links.find((link) => link.text().includes('Início'));
    const usersLink = links.find((link) => link.text().includes('Usuários'));

    expect(homeLink?.classes().includes('!bg-base-300')).toBe(false);
    expect(usersLink?.classes().includes('!bg-base-300')).toBe(true);
    expectComponentSnapshot(wrapper);
  });

  it('hides Usuários for a standard user but shows Times', async () => {
    const { wrapper } = await mountComponent(DefaultLayout, {
      routes,
      initialRoute: '/',
      session: userSession,
    });

    const nav = wrapper.get('aside[aria-label="Menu principal"]');
    expect(nav.text()).toContain('Chave de API');
    expect(nav.text()).toContain('Times');
    expect(nav.text()).not.toContain('Usuários');
    expect(nav.text()).not.toContain('Auditoria');
    expectComponentSnapshot(wrapper);
  });

  it('logs out and navigates to login when Sair is clicked', async () => {
    const { wrapper, router, pinia } = await mountComponent(DefaultLayout, {
      routes,
      initialRoute: '/',
      session: adminSession,
    });

    await wrapper.get('[data-testid="logout"]').trigger('click');
    await flushPromises();

    expect(useAuthStore(pinia).isAuthenticated).toBe(false);
    expect(router.currentRoute.value.name).toBe('login');
    expect(sessionStorage.getItem('envm-vault.auth')).toBeNull();
    expectComponentSnapshot(wrapper);
  });

  it('switches UI language when locale is changed', async () => {
    const { wrapper } = await mountComponent(DefaultLayout, {
      routes,
      initialRoute: '/',
      session: adminSession,
    });

    const select = wrapper.get('[data-testid="locale-switcher"]');
    await select.setValue('en-US');
    await flushPromises();

    expect(wrapper.get('header').text()).toContain('Log out');
    expect(wrapper.get('aside[aria-label="Main menu"]').text()).toContain('Home');
    expectComponentSnapshot(wrapper);
  });
});
