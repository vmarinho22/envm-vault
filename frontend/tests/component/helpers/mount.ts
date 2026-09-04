import { mount, type ComponentMountingOptions, type VueWrapper } from '@vue/test-utils';
import { createPinia, setActivePinia, type Pinia } from 'pinia';
import { expect } from 'vitest';
import { createMemoryHistory, createRouter, type RouteRecordRaw, type Router } from 'vue-router';
import type { Component } from 'vue';
import { useAuthStore } from '@/modules/auth/stores/auth.store';
import type { User } from '@/modules/user/types';

export interface MountComponentOptions {
  routes?: RouteRecordRaw[];
  initialRoute?: string;
  session?: { token: string; user: User };
  props?: Record<string, unknown>;
}

/** Assert the rendered HTML against a Vitest snapshot. Required in every component test. */
export function expectComponentSnapshot(wrapper: VueWrapper): void {
  expect(wrapper.html()).toMatchSnapshot();
}

export async function mountComponent(
  component: Component,
  options: MountComponentOptions = {},
): Promise<{ wrapper: VueWrapper; router: Router; pinia: Pinia }> {
  const pinia = createPinia();
  setActivePinia(pinia);
  sessionStorage.clear();

  if (options.session) {
    useAuthStore(pinia).setSession(options.session);
  }

  const routes = options.routes ?? [
    { path: '/', name: 'home', component: { template: '<div />' } },
  ];
  const router = createRouter({
    history: createMemoryHistory(),
    routes,
  });

  await router.push(options.initialRoute ?? routes[0]?.path ?? '/');
  await router.isReady();

  const mountingOptions: ComponentMountingOptions<Component> = {
    props: options.props,
    global: {
      plugins: [pinia, router],
    },
  };

  const wrapper = mount(component, mountingOptions);

  return { wrapper, router, pinia };
}
