import { describe, expect, it, vi } from 'vitest';
import TokenOnceReveal from '@/modules/token/components/TokenOnceReveal.vue';
import { expectComponentSnapshot, mountComponent } from '../../../helpers/mount';

describe('TokenOnceReveal (component)', () => {
  it('renders the token and copy button', async () => {
    const { wrapper } = await mountComponent(TokenOnceReveal, {
      props: { token: 'ev_secretvalue' },
    });

    expect(wrapper.get('[data-testid="api-token-value"]').element).toMatchObject({
      value: 'ev_secretvalue',
    });
    expect(wrapper.text()).toContain('Copiar');
    expectComponentSnapshot(wrapper);
  });

  it('copies the current token without changing it', async () => {
    const writeText = vi.fn().mockResolvedValue(undefined);
    vi.stubGlobal('navigator', { clipboard: { writeText } });

    const { wrapper } = await mountComponent(TokenOnceReveal, {
      props: { token: 'ev_secretvalue' },
    });

    await wrapper.get('[data-testid="copy-token"]').trigger('click');

    expect(writeText).toHaveBeenCalledWith('ev_secretvalue');
    expect(wrapper.get('[data-testid="api-token-value"]').element).toMatchObject({
      value: 'ev_secretvalue',
    });
    expect(wrapper.text()).toContain('Copiada');
    expectComponentSnapshot(wrapper);

    vi.unstubAllGlobals();
  });
});
