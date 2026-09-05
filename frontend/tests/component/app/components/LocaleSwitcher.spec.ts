import { describe, expect, it } from 'vitest';
import LocaleSwitcher from '@/app/components/LocaleSwitcher.vue';
import { expectComponentSnapshot, mountComponent } from '../../helpers/mount';

describe('LocaleSwitcher (component)', () => {
  it('renders supported locales and keeps pt-BR selected by default', async () => {
    const { wrapper } = await mountComponent(LocaleSwitcher);

    const select = wrapper.get('[data-testid="locale-switcher"]');
    expect(select.element).toBeInstanceOf(HTMLSelectElement);
    expect((select.element as HTMLSelectElement).value).toBe('pt-BR');
    expect(select.text()).toContain('Português');
    expect(select.text()).toContain('English');
    expectComponentSnapshot(wrapper);
  });

  it('updates labels when locale changes', async () => {
    const { wrapper } = await mountComponent(LocaleSwitcher);

    const select = wrapper.get('[data-testid="locale-switcher"]');
    await select.setValue('en-US');

    expect(select.attributes('aria-label')).toBe('Language');
    expectComponentSnapshot(wrapper);
  });
});
