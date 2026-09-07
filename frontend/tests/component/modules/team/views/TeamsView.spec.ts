import { flushPromises } from '@vue/test-utils';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { teamService } from '@/modules/team/services/TeamService';
import TeamsView from '@/modules/team/views/TeamsView.vue';
import { userService } from '@/modules/user/services/UserService';
import { expectComponentSnapshot, mountComponent } from '../../../helpers/mount';

vi.mock('@/modules/team/services/TeamService', () => ({
  teamService: {
    list: vi.fn(),
    create: vi.fn(),
    update: vi.fn(),
    remove: vi.fn(),
  },
}));

vi.mock('@/modules/user/services/UserService', () => ({
  userService: {
    register: vi.fn(),
    list: vi.fn(),
  },
}));

const listMock = vi.mocked(teamService.list);
const userListMock = vi.mocked(userService.list);

const adminSession = {
  token: 'jwt-token',
  user: { id: '1', email: 'admin@example.com', role: 'super_admin' as const },
};

const userSession = {
  token: 'jwt-token',
  user: { id: '2', email: 'user@example.com', role: 'user' as const },
};

const sampleTeams = [
  {
    id: 't1',
    name: 'Platform',
    slug: 'platform',
    accountable: { id: '1', email: 'admin@example.com' },
    created_at: '2026-01-01T00:00:00Z',
    updated_at: '2026-01-02T00:00:00Z',
  },
];

describe('TeamsView (component)', () => {
  beforeEach(() => {
    listMock.mockReset();
    userListMock.mockReset();
    listMock.mockResolvedValue(sampleTeams);
    userListMock.mockResolvedValue([
      {
        id: '1',
        email: 'admin@example.com',
        role: 'super_admin',
        token: null,
      },
    ]);
  });

  it('renders teams and admin actions', async () => {
    const { wrapper } = await mountComponent(TeamsView, { session: adminSession });
    await flushPromises();

    expect(wrapper.get('h1').text()).toBe('Times');
    expect(wrapper.text()).toContain('Platform');
    expect(wrapper.text()).toContain('platform');
    expect(wrapper.text()).toContain('admin@example.com');
    expect(wrapper.text()).toContain('Criar time');
    expect(wrapper.get('button[aria-label="Editar"]')).toBeTruthy();
    expect(wrapper.get('button[aria-label="Excluir"]')).toBeTruthy();
    expect(userListMock).toHaveBeenCalled();
    expectComponentSnapshot(wrapper);
  });

  it('hides mutate actions for standard user', async () => {
    const { wrapper } = await mountComponent(TeamsView, { session: userSession });
    await flushPromises();

    expect(wrapper.text()).toContain('Platform');
    expect(wrapper.text()).not.toContain('Criar time');
    expect(wrapper.find('button[aria-label="Editar"]').exists()).toBe(false);
    expect(wrapper.find('button[aria-label="Excluir"]').exists()).toBe(false);
    expect(userListMock).not.toHaveBeenCalled();
    expectComponentSnapshot(wrapper);
  });

  it('searches by name via API', async () => {
    listMock.mockResolvedValueOnce(sampleTeams).mockResolvedValueOnce([]);

    const { wrapper } = await mountComponent(TeamsView, { session: userSession });
    await flushPromises();

    await wrapper.get('input[type="search"]').setValue('plat');
    await wrapper.get('form').trigger('submit');
    await flushPromises();

    expect(listMock).toHaveBeenLastCalledWith({ name: 'plat' });
    expect(wrapper.text()).toContain('Nenhum time encontrado para essa busca.');
    expectComponentSnapshot(wrapper);
  });
});
