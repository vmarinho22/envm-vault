import { flushPromises, type VueWrapper } from '@vue/test-utils';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { authService } from '@/modules/auth/services/AuthService';
import LoginView from '@/modules/auth/views/LoginView.vue';
import HomeView from '@/modules/home/views/HomeView.vue';
import { tokenService } from '@/modules/token/services/TokenService';
import { ApiError } from '@/shared/http/ApiError';
import { expectComponentSnapshot, mountComponent } from '../../../helpers/mount';

vi.mock('@/modules/auth/services/AuthService', () => ({
  authService: {
    requestLogin: vi.fn(),
    verifyLogin: vi.fn(),
  },
}));

vi.mock('@/modules/token/services/TokenService', () => ({
  tokenService: {
    getMine: vi.fn(),
    create: vi.fn(),
    revokeOwn: vi.fn(),
    revokeForUser: vi.fn(),
  },
}));

const requestLoginMock = vi.mocked(authService.requestLogin);
const verifyLoginMock = vi.mocked(authService.verifyLogin);
const getMineMock = vi.mocked(tokenService.getMine);

async function fillOtpCode(wrapper: VueWrapper, code: string): Promise<void> {
  const digitInputs = wrapper.findAll('input[aria-label^="Dígito"]');
  expect(digitInputs).toHaveLength(8);

  const first = digitInputs[0];
  expect(first).toBeDefined();
  await first!.trigger('paste', {
    clipboardData: {
      getData: () => code,
    },
  });
}

describe('LoginView (component)', () => {
  beforeEach(() => {
    requestLoginMock.mockReset();
    verifyLoginMock.mockReset();
    getMineMock.mockReset();
    getMineMock.mockResolvedValue({ active: true, prefix: 'ev_a1b2c3d4' });
  });

  it('renders email step with brand and send CTA', async () => {
    const { wrapper } = await mountComponent(LoginView, {
      routes: [
        { path: '/login', name: 'login', component: LoginView },
        { path: '/', name: 'home', component: HomeView },
        { path: '/setup-token', name: 'setup-token', component: { template: '<div />' } },
      ],
      initialRoute: '/login',
    });

    expect(wrapper.text()).toContain('envm-vault');
    expect(wrapper.find('input[name="email"]').exists()).toBe(true);
    expect(wrapper.text()).toContain('Enviar código');
    expectComponentSnapshot(wrapper);
  });

  it('moves from email step to code step after requestLogin', async () => {
    requestLoginMock.mockResolvedValue({ message: 'verification code sent' });
    const { wrapper } = await mountComponent(LoginView, {
      routes: [
        { path: '/login', name: 'login', component: LoginView },
        { path: '/', name: 'home', component: HomeView },
        { path: '/setup-token', name: 'setup-token', component: { template: '<div />' } },
      ],
      initialRoute: '/login',
    });

    await wrapper.get('input[name="email"]').setValue('admin@example.com');
    await wrapper.get('form').trigger('submit');
    await flushPromises();

    expect(requestLoginMock).toHaveBeenCalledWith('admin@example.com');
    expect(wrapper.findAll('input[aria-label^="Dígito"]')).toHaveLength(8);
    expect(wrapper.text()).toContain('Entrar');
    expectComponentSnapshot(wrapper);
  });

  it('shows API error on the email step', async () => {
    requestLoginMock.mockRejectedValue(new ApiError(404, 'user not found'));
    const { wrapper } = await mountComponent(LoginView, {
      routes: [
        { path: '/login', name: 'login', component: LoginView },
        { path: '/', name: 'home', component: HomeView },
        { path: '/setup-token', name: 'setup-token', component: { template: '<div />' } },
      ],
      initialRoute: '/login',
    });

    await wrapper.get('input[name="email"]').setValue('missing@example.com');
    await wrapper.get('form').trigger('submit');
    await flushPromises();

    expect(wrapper.get('[role="alert"]').text()).toBe('user not found');
    expectComponentSnapshot(wrapper);
  });

  it('calls verifyLogin and stores the session on success', async () => {
    requestLoginMock.mockResolvedValue({ message: 'verification code sent' });
    verifyLoginMock.mockResolvedValue({
      token: 'jwt-token',
      user: { id: '1', email: 'admin@example.com', role: 'super_admin' },
    });
    const { wrapper } = await mountComponent(LoginView, {
      routes: [
        { path: '/login', name: 'login', component: LoginView },
        { path: '/', name: 'home', component: HomeView },
        { path: '/setup-token', name: 'setup-token', component: { template: '<div />' } },
      ],
      initialRoute: '/login',
    });

    await wrapper.get('input[name="email"]').setValue('admin@example.com');
    await wrapper.get('form').trigger('submit');
    await flushPromises();

    await fillOtpCode(wrapper, '12345678');
    await wrapper.get('form').trigger('submit');
    await flushPromises();

    expect(verifyLoginMock).toHaveBeenCalledWith('admin@example.com', '12345678');
    expect(getMineMock).toHaveBeenCalled();
    expect(sessionStorage.getItem('envm-vault.auth')).toContain('jwt-token');
    expectComponentSnapshot(wrapper);
  });
});
