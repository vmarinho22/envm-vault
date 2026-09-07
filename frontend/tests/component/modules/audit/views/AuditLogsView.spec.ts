import { flushPromises } from '@vue/test-utils';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { auditService } from '@/modules/audit/services/AuditService';
import AuditLogsView from '@/modules/audit/views/AuditLogsView.vue';
import { userService } from '@/modules/user/services/UserService';
import { expectComponentSnapshot, mountComponent } from '../../../helpers/mount';

vi.mock('@/modules/audit/services/AuditService', () => ({
  auditService: {
    list: vi.fn(),
  },
}));

vi.mock('@/modules/user/services/UserService', () => ({
  userService: {
    register: vi.fn(),
    list: vi.fn(),
  },
}));

const listLogsMock = vi.mocked(auditService.list);
const listUsersMock = vi.mocked(userService.list);

const session = {
  token: 'jwt-token',
  user: { id: '1', email: 'admin@example.com', role: 'super_admin' as const },
};

describe('AuditLogsView (component)', () => {
  beforeEach(() => {
    listLogsMock.mockReset();
    listUsersMock.mockReset();
    listUsersMock.mockResolvedValue([
      {
        id: '1',
        email: 'admin@example.com',
        role: 'super_admin',
        token: null,
      },
    ]);
    listLogsMock.mockResolvedValue({
      items: [
        {
          id: 'log-1',
          action: 'login.succeeded',
          actor_user_id: '1',
          actor_email: 'admin@example.com',
          actor_role: 'super_admin',
          target_user_id: null,
          target_email: null,
          auth_channel: 'login_attempt',
          ip: '127.0.0.1',
          user_agent: 'vitest',
          origin: 'http://localhost:3100',
          occurred_at: '2026-09-06T15:00:00Z',
          created_at: '2026-09-06T15:00:00Z',
        },
      ],
      page: 1,
      page_size: 25,
      total: 1,
    });
  });

  it('renders audit logs table with action labels', async () => {
    const { wrapper } = await mountComponent(AuditLogsView, { session });
    await flushPromises();

    expect(wrapper.get('h1').text()).toBe('Logs de auditoria');
    expect(wrapper.text()).toContain('Login bem-sucedido');
    expect(wrapper.text()).toContain('admin@example.com');
    expect(wrapper.text()).toContain('Web');
    expect(listUsersMock).toHaveBeenCalledTimes(1);
    expect(listLogsMock).toHaveBeenCalled();
    expectComponentSnapshot(wrapper);
  });
});
