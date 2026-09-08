import { flushPromises } from '@vue/test-utils';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { environmentService } from '@/modules/environment/services/EnvironmentService';
import EnvironmentsView from '@/modules/environment/views/EnvironmentsView.vue';
import { expectComponentSnapshot, mountComponent } from '../../../helpers/mount';

vi.mock('@/modules/environment/services/EnvironmentService', () => ({
  environmentService: {
    list: vi.fn(),
    create: vi.fn(),
    update: vi.fn(),
    remove: vi.fn(),
  },
}));

const listMock = vi.mocked(environmentService.list);

const adminSession = {
  token: 'jwt-token',
  user: { id: '1', email: 'admin@example.com', role: 'super_admin' as const },
};

const userSession = {
  token: 'jwt-token',
  user: { id: '2', email: 'user@example.com', role: 'user' as const },
};

const sampleEnvironments = [
  {
    id: 'e1',
    title: 'dev',
    created_at: '2026-01-01T00:00:00Z',
    updated_at: '2026-01-02T00:00:00Z',
  },
  {
    id: 'e2',
    title: 'production',
    created_at: '2026-01-01T00:00:00Z',
    updated_at: '2026-01-02T00:00:00Z',
  },
];

describe('EnvironmentsView (component)', () => {
  beforeEach(() => {
    listMock.mockReset();
    listMock.mockResolvedValue(sampleEnvironments);
  });

  it('renders environments and admin actions', async () => {
    const { wrapper } = await mountComponent(EnvironmentsView, { session: adminSession });
    await flushPromises();

    expect(wrapper.get('h1').text()).toBe('Ambientes');
    expect(wrapper.text()).toContain('dev');
    expect(wrapper.text()).toContain('production');
    expect(wrapper.text()).toContain('Criar ambiente');
    const editButtons = wrapper.findAll('button[aria-label="Editar"]');
    const deleteButtons = wrapper.findAll('button[aria-label="Excluir"]');
    expect(editButtons).toHaveLength(2);
    expect(deleteButtons).toHaveLength(2);
    expect(editButtons[1]?.attributes('disabled')).toBeDefined();
    expect(deleteButtons[1]?.attributes('disabled')).toBeDefined();
    expectComponentSnapshot(wrapper);
  });

  it('hides mutate actions for standard user', async () => {
    const { wrapper } = await mountComponent(EnvironmentsView, { session: userSession });
    await flushPromises();

    expect(wrapper.text()).toContain('dev');
    expect(wrapper.text()).not.toContain('Criar ambiente');
    expect(wrapper.find('button[aria-label="Editar"]').exists()).toBe(false);
    expect(wrapper.find('button[aria-label="Excluir"]').exists()).toBe(false);
    expectComponentSnapshot(wrapper);
  });

  it('searches by title via API', async () => {
    listMock.mockResolvedValueOnce(sampleEnvironments).mockResolvedValueOnce([]);

    const { wrapper } = await mountComponent(EnvironmentsView, { session: userSession });
    await flushPromises();

    await wrapper.get('input[type="search"]').setValue('stag');
    await wrapper.get('form').trigger('submit');
    await flushPromises();

    expect(listMock).toHaveBeenLastCalledWith({ title: 'stag' });
    expect(wrapper.text()).toContain('Nenhum ambiente encontrado para essa busca.');
    expectComponentSnapshot(wrapper);
  });
});
