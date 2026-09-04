import { describe, expect, it } from 'vitest';
import HomeView from '@/modules/home/views/HomeView.vue';
import { expectComponentSnapshot, mountComponent } from '../../../helpers/mount';

describe('HomeView (component)', () => {
  it('renders title and welcome message', async () => {
    const { wrapper } = await mountComponent(HomeView);

    expect(wrapper.get('h1').text()).toBe('Início');
    expect(wrapper.text()).toContain('Bem-vindo ao envm-vault.');
    expectComponentSnapshot(wrapper);
  });
});
