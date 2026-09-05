import { flushPromises } from '@vue/test-utils';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { tokenService } from '@/modules/token/services/TokenService';
import TokenSetupView from '@/modules/token/views/TokenSetupView.vue';
import { ApiError } from '@/shared/http/ApiError';
import { expectComponentSnapshot, mountComponent } from '../../../helpers/mount';

vi.mock('@/modules/token/services/TokenService', () => ({
  tokenService: {
    create: vi.fn(),
    getMine: vi.fn(),
    revokeOwn: vi.fn(),
    revokeForUser: vi.fn(),
  },
}));

const createMock = vi.mocked(tokenService.create);

const session = {
  token: 'jwt-token',
  user: { id: '1', email: 'admin@example.com', role: 'super_admin' as const },
};

describe('TokenSetupView (component)', () => {
  beforeEach(() => {
    createMock.mockReset();
  });

  it('renders generate CTA', async () => {
    const { wrapper } = await mountComponent(TokenSetupView, {
      routes: [
        { path: '/setup-token', name: 'setup-token', component: TokenSetupView },
        { path: '/', name: 'home', component: { template: '<div />' } },
        { path: '/login', name: 'login', component: { template: '<div />' } },
      ],
      initialRoute: '/setup-token',
      session,
    });

    expect(wrapper.get('h1').text()).toContain('Gere sua chave de API');
    expect(wrapper.text()).toContain('Gerar chave');
    expectComponentSnapshot(wrapper);
  });

  it('shows the plaintext token after generate', async () => {
    createMock.mockResolvedValue({
      token: 'ev_testtokenvalue',
      prefix: 'ev_testtok',
      created_at: '2026-01-01T00:00:00Z',
    });

    const { wrapper } = await mountComponent(TokenSetupView, {
      routes: [
        { path: '/setup-token', name: 'setup-token', component: TokenSetupView },
        { path: '/', name: 'home', component: { template: '<div />' } },
      ],
      initialRoute: '/setup-token',
      session,
    });

    await wrapper.get('button.btn-primary').trigger('click');
    await flushPromises();

    expect(wrapper.get('[data-testid="api-token-value"]').element).toMatchObject({
      value: 'ev_testtokenvalue',
    });
    expect(wrapper.text()).toContain('Continuar');
    expectComponentSnapshot(wrapper);
  });

  it('does not generate a new token when Copiar is clicked', async () => {
    createMock.mockResolvedValue({
      token: 'ev_testtokenvalue',
      prefix: 'ev_testtok',
      created_at: '2026-01-01T00:00:00Z',
    });
    const writeText = vi.fn().mockResolvedValue(undefined);
    vi.stubGlobal('navigator', { clipboard: { writeText } });

    const { wrapper } = await mountComponent(TokenSetupView, {
      routes: [
        { path: '/setup-token', name: 'setup-token', component: TokenSetupView },
        { path: '/', name: 'home', component: { template: '<div />' } },
      ],
      initialRoute: '/setup-token',
      session,
    });

    await wrapper.get('button.btn-primary').trigger('click');
    await flushPromises();

    await wrapper.get('[data-testid="copy-token"]').trigger('click');
    await flushPromises();

    expect(createMock).toHaveBeenCalledTimes(1);
    expect(writeText).toHaveBeenCalledWith('ev_testtokenvalue');
    expect(wrapper.get('[data-testid="api-token-value"]').element).toMatchObject({
      value: 'ev_testtokenvalue',
    });
    expectComponentSnapshot(wrapper);

    vi.unstubAllGlobals();
  });

  it('shows API error when generate fails', async () => {
    createMock.mockRejectedValue(new ApiError(500, 'internal server error'));

    const { wrapper } = await mountComponent(TokenSetupView, {
      routes: [{ path: '/setup-token', name: 'setup-token', component: TokenSetupView }],
      initialRoute: '/setup-token',
      session,
    });

    await wrapper.get('button.btn-primary').trigger('click');
    await flushPromises();

    expect(wrapper.get('[role="alert"]').text()).toBe('internal server error');
    expectComponentSnapshot(wrapper);
  });
});
