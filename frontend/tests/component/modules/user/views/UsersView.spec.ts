import { flushPromises } from '@vue/test-utils';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { tokenService } from '@/modules/token/services/TokenService';
import { userService } from '@/modules/user/services/UserService';
import UsersView from '@/modules/user/views/UsersView.vue';
import { expectComponentSnapshot, mountComponent } from '../../../helpers/mount';

vi.mock('@/modules/user/services/UserService', () => ({
  userService: {
    register: vi.fn(),
    list: vi.fn(),
  },
}));

vi.mock('@/modules/token/services/TokenService', () => ({
  tokenService: {
    create: vi.fn(),
    getMine: vi.fn(),
    revokeOwn: vi.fn(),
    revokeForUser: vi.fn(),
  },
}));

const listMock = vi.mocked(userService.list);
const revokeForUserMock = vi.mocked(tokenService.revokeForUser);

const session = {
  token: 'jwt-token',
  user: { id: '1', email: 'admin@example.com', role: 'super_admin' as const },
};

describe('UsersView (component)', () => {
  beforeEach(() => {
    listMock.mockReset();
    revokeForUserMock.mockReset();
    listMock.mockResolvedValue([
      {
        id: '1',
        email: 'admin@example.com',
        role: 'super_admin',
        token: { prefix: 'ev_adminxxx', created_at: '2026-01-01T00:00:00Z' },
      },
      {
        id: '2',
        email: 'user@example.com',
        role: 'user',
        token: null,
      },
    ]);
  });

  it('renders users and token status', async () => {
    const { wrapper } = await mountComponent(UsersView, { session });
    await flushPromises();

    expect(wrapper.get('h1').text()).toBe('Usuários');
    expect(wrapper.text()).toContain('admin@example.com');
    expect(wrapper.text()).toContain('ev_adminxxx');
    expect(wrapper.text()).toContain('Sem chave');
    expectComponentSnapshot(wrapper);
  });

  it('revokes an active token and reloads the list', async () => {
    revokeForUserMock.mockResolvedValue(undefined);
    listMock
      .mockResolvedValueOnce([
        {
          id: '2',
          email: 'user@example.com',
          role: 'user',
          token: { prefix: 'ev_userxxxx', created_at: '2026-01-01T00:00:00Z' },
        },
      ])
      .mockResolvedValueOnce([
        {
          id: '2',
          email: 'user@example.com',
          role: 'user',
          token: {
            prefix: 'ev_userxxxx',
            created_at: '2026-01-01T00:00:00Z',
            revoked_at: '2026-01-02T00:00:00Z',
          },
        },
      ]);

    const { wrapper } = await mountComponent(UsersView, { session });
    await flushPromises();

    await wrapper.get('button.btn-error').trigger('click');
    await flushPromises();

    expect(revokeForUserMock).toHaveBeenCalledWith('2');
    expect(wrapper.text()).toContain('Revogada');
    expectComponentSnapshot(wrapper);
  });
});
