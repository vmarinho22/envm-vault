import { flushPromises } from '@vue/test-utils';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { environmentService } from '@/modules/environment/services/EnvironmentService';
import { secretService } from '@/modules/secret/services/SecretService';
import SecretDetailView from '@/modules/secret/views/SecretDetailView.vue';
import { expectComponentSnapshot, mountComponent } from '../../../helpers/mount';

vi.mock('@/modules/environment/services/EnvironmentService', () => ({
  environmentService: {
    list: vi.fn(),
  },
}));

vi.mock('@/modules/secret/services/SecretService', () => ({
  secretService: {
    getById: vi.fn(),
    remove: vi.fn(),
  },
}));

vi.mock('@/shared/confirm/useConfirmDialog', () => ({
  useConfirmDialog: () => ({
    confirm: vi.fn().mockResolvedValue(true),
  }),
}));

const getByIdMock = vi.mocked(secretService.getById);
const removeMock = vi.mocked(secretService.remove);
const envListMock = vi.mocked(environmentService.list);

const session = {
  token: 'jwt-token',
  user: { id: '1', email: 'admin@example.com', role: 'super_admin' as const },
};

const detail = {
  id: 's1',
  name: 'DATABASE_URL',
  project_count: 1,
  values: {
    e1: { has_value: true, updated_at: '2026-01-02T00:00:00Z' },
  },
  projects: [{ id: 'p1', title: 'Payment', slug: 'payment' }],
  created_at: '2026-01-01T00:00:00Z',
  updated_at: '2026-01-02T00:00:00Z',
};

describe('SecretDetailView', () => {
  beforeEach(() => {
    getByIdMock.mockReset();
    removeMock.mockReset();
    envListMock.mockReset();
    getByIdMock.mockResolvedValue(detail);
    envListMock.mockResolvedValue([{ id: 'e1', title: 'dev', created_at: '', updated_at: '' }]);
    removeMock.mockResolvedValue(undefined);
  });

  it('renders secret detail with linked projects and environment badges', async () => {
    const { wrapper } = await mountComponent(SecretDetailView, {
      session,
      routes: [
        { path: '/secrets', name: 'secrets', component: { template: '<div />' } },
        { path: '/secrets/:id', name: 'secret-detail', component: SecretDetailView },
        { path: '/project/:slug', name: 'project-detail', component: { template: '<div />' } },
      ],
      initialRoute: '/secrets/s1',
    });
    await flushPromises();

    expect(wrapper.text()).toContain('DATABASE_URL');
    expect(wrapper.text()).toContain('Payment');
    expect(wrapper.text()).toContain('dev');
    expect(wrapper.get('a[href="/project/payment"]').text()).toContain('Payment');
    expectComponentSnapshot(wrapper);
  });

  it('detaches project when unlink is confirmed', async () => {
    getByIdMock
      .mockResolvedValueOnce(detail)
      .mockResolvedValueOnce({ ...detail, project_count: 0, projects: [] });

    const { wrapper } = await mountComponent(SecretDetailView, {
      session,
      routes: [
        { path: '/secrets', name: 'secrets', component: { template: '<div />' } },
        { path: '/secrets/:id', name: 'secret-detail', component: SecretDetailView },
        { path: '/project/:slug', name: 'project-detail', component: { template: '<div />' } },
      ],
      initialRoute: '/secrets/s1',
    });
    await flushPromises();

    const unlinkBtn = wrapper.findAll('button').find((b) => b.text().includes('Desvincular'));
    expect(unlinkBtn).toBeTruthy();
    await unlinkBtn!.trigger('click');
    await flushPromises();

    expect(removeMock).toHaveBeenCalledWith('p1', 's1');
    expect(wrapper.text()).toContain('Nenhum projeto vinculado');
    expectComponentSnapshot(wrapper);
  });
});
