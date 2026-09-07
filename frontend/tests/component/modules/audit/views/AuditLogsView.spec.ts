import { flushPromises } from '@vue/test-utils';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { auditService } from '@/modules/audit/services/AuditService';
import AuditLogsView from '@/modules/audit/views/AuditLogsView.vue';
import { projectService } from '@/modules/project/services/ProjectService';
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

vi.mock('@/modules/project/services/ProjectService', () => ({
  projectService: {
    list: vi.fn(),
  },
}));

const listLogsMock = vi.mocked(auditService.list);
const listUsersMock = vi.mocked(userService.list);
const listProjectsMock = vi.mocked(projectService.list);

const session = {
  token: 'jwt-token',
  user: { id: '1', email: 'admin@example.com', role: 'super_admin' as const },
};

describe('AuditLogsView (component)', () => {
  beforeEach(() => {
    listLogsMock.mockReset();
    listUsersMock.mockReset();
    listProjectsMock.mockReset();
    listUsersMock.mockResolvedValue([
      {
        id: '1',
        email: 'admin@example.com',
        role: 'super_admin',
        token: null,
      },
    ]);
    listProjectsMock.mockResolvedValue({
      items: [
        {
          id: 'p1',
          title: 'Payment',
          slug: 'payment',
          teams: [],
          environments: [],
          tags: [],
          created_at: '2026-01-01T00:00:00Z',
          updated_at: '2026-01-01T00:00:00Z',
        },
      ],
      page: 1,
      page_size: 12,
      total: 1,
    });
    listLogsMock.mockResolvedValue({
      items: [
        {
          id: 'log-1',
          action: 'project.created',
          actor_user_id: '1',
          actor_email: 'admin@example.com',
          actor_role: 'super_admin',
          target_user_id: null,
          target_email: null,
          project_id: 'p1',
          project_title: 'Payment',
          project_slug: 'payment',
          secret_id: null,
          secret_name: '',
          auth_channel: 'logged_in_user',
          ip: '127.0.0.1',
          user_agent: 'vitest',
          origin: 'http://localhost:3100',
          occurred_at: '2026-09-06T15:00:00Z',
          created_at: '2026-09-06T15:00:00Z',
        },
        {
          id: 'log-2',
          action: 'secret_value.set',
          actor_user_id: '1',
          actor_email: 'admin@example.com',
          actor_role: 'super_admin',
          target_user_id: null,
          target_email: null,
          project_id: 'p1',
          project_title: 'Payment',
          project_slug: 'payment',
          secret_id: 's1',
          secret_name: 'DATABASE_URL',
          auth_channel: 'logged_in_user',
          ip: '127.0.0.1',
          user_agent: 'vitest',
          origin: 'http://localhost:3100',
          occurred_at: '2026-09-06T15:01:00Z',
          created_at: '2026-09-06T15:01:00Z',
        },
      ],
      page: 1,
      page_size: 25,
      total: 2,
    });
  });

  it('renders audit logs table with project column and filter', async () => {
    const { wrapper } = await mountComponent(AuditLogsView, { session });
    await flushPromises();

    expect(wrapper.get('h1').text()).toBe('Logs de auditoria');
    expect(wrapper.text()).toContain('Projeto criado');
    expect(wrapper.text()).toContain('Payment');
    expect(wrapper.text()).toContain('Secret');
    expect(wrapper.text()).toContain('DATABASE_URL');
    expect(wrapper.text()).toContain('admin@example.com');
    expect(wrapper.text()).toContain('Projeto');
    expect(listUsersMock).toHaveBeenCalledTimes(1);
    expect(listProjectsMock).toHaveBeenCalled();
    expect(listLogsMock).toHaveBeenCalled();
    expectComponentSnapshot(wrapper);
  });
});
