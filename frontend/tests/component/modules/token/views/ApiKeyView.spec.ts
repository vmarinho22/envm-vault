import { flushPromises } from '@vue/test-utils';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { tokenService } from '@/modules/token/services/TokenService';
import ApiKeyView from '@/modules/token/views/ApiKeyView.vue';
import { expectComponentSnapshot, mountComponent } from '../../../helpers/mount';

vi.mock('@/modules/token/services/TokenService', () => ({
  tokenService: {
    create: vi.fn(),
    getMine: vi.fn(),
    revokeOwn: vi.fn(),
    revokeForUser: vi.fn(),
  },
}));

const getMineMock = vi.mocked(tokenService.getMine);
const createMock = vi.mocked(tokenService.create);
const revokeOwnMock = vi.mocked(tokenService.revokeOwn);

const session = {
  token: 'jwt-token',
  user: { id: '1', email: 'admin@example.com', role: 'super_admin' as const },
};

describe('ApiKeyView (component)', () => {
  beforeEach(() => {
    getMineMock.mockReset();
    createMock.mockReset();
    revokeOwnMock.mockReset();
    getMineMock.mockResolvedValue({
      active: true,
      prefix: 'ev_a1b2c3d4',
      created_at: '2026-01-01T00:00:00.000Z',
    });
  });

  it('renders prefix and actions after loading', async () => {
    const { wrapper } = await mountComponent(ApiKeyView, { session });
    await flushPromises();

    expect(wrapper.get('h1').text()).toBe('Chave de API');
    expect(wrapper.text()).toContain('ev_a1b2c3d4');
    expect(wrapper.text()).toContain('Rotacionar');
    expect(wrapper.text()).toContain('Revogar');
    expectComponentSnapshot(wrapper);
  });

  it('shows the new token after rotate is confirmed', async () => {
    vi.stubGlobal('confirm', vi.fn(() => true));
    createMock.mockResolvedValue({
      token: 'ev_rotatedtoken',
      prefix: 'ev_rotated',
      created_at: '2026-02-01T00:00:00Z',
    });

    const { wrapper } = await mountComponent(ApiKeyView, { session });
    await flushPromises();

    await wrapper.get('[data-testid="rotate-token"]').trigger('click');
    await flushPromises();

    expect(createMock).toHaveBeenCalledTimes(1);
    expect(wrapper.get('[data-testid="api-token-value"]').element).toMatchObject({
      value: 'ev_rotatedtoken',
    });
    expect(wrapper.find('[data-testid="rotate-token"]').exists()).toBe(false);

    const writeText = vi.fn().mockResolvedValue(undefined);
    vi.stubGlobal('navigator', { clipboard: { writeText } });
    await wrapper.get('[data-testid="copy-token"]').trigger('click');
    await flushPromises();

    expect(createMock).toHaveBeenCalledTimes(1);
    expect(writeText).toHaveBeenCalledWith('ev_rotatedtoken');
    expectComponentSnapshot(wrapper);

    vi.unstubAllGlobals();
  });

  it('hides rotate while the new key is visible and restores it after dismiss', async () => {
    vi.stubGlobal('confirm', vi.fn(() => true));
    createMock.mockResolvedValue({
      token: 'ev_rotatedtoken',
      prefix: 'ev_rotated',
      created_at: '2026-02-01T00:00:00Z',
    });

    const { wrapper } = await mountComponent(ApiKeyView, { session });
    await flushPromises();

    await wrapper.get('[data-testid="rotate-token"]').trigger('click');
    await flushPromises();

    expect(wrapper.find('[data-testid="rotate-token"]').exists()).toBe(false);

    await wrapper.get('button.btn-primary').trigger('click');
    await flushPromises();

    expect(wrapper.find('[data-testid="rotate-token"]').exists()).toBe(true);
    expect(wrapper.find('[data-testid="api-token-value"]').exists()).toBe(false);
    expectComponentSnapshot(wrapper);

    vi.unstubAllGlobals();
  });

  it('does not rotate when confirmation is cancelled', async () => {
    vi.stubGlobal('confirm', vi.fn(() => false));

    const { wrapper } = await mountComponent(ApiKeyView, { session });
    await flushPromises();

    await wrapper.get('[data-testid="rotate-token"]').trigger('click');
    await flushPromises();

    expect(createMock).not.toHaveBeenCalled();
    expect(wrapper.find('[data-testid="api-token-value"]').exists()).toBe(false);
    expectComponentSnapshot(wrapper);

    vi.unstubAllGlobals();
  });

  it('navigates to setup-token after revoke', async () => {
    revokeOwnMock.mockResolvedValue(undefined);

    const { wrapper, router } = await mountComponent(ApiKeyView, {
      session,
      routes: [
        { path: '/api-key', name: 'api-key', component: ApiKeyView },
        { path: '/setup-token', name: 'setup-token', component: { template: '<div />' } },
      ],
      initialRoute: '/api-key',
    });
    await flushPromises();

    await wrapper.get('button.btn-error').trigger('click');
    await flushPromises();

    expect(revokeOwnMock).toHaveBeenCalled();
    expect(router.currentRoute.value.name).toBe('setup-token');
    expectComponentSnapshot(wrapper);
  });
});
