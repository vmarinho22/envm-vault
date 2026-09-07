import { flushPromises } from '@vue/test-utils';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { projectService } from '@/modules/project/services/ProjectService';
import { secretService } from '@/modules/secret/services/SecretService';
import ProjectDetailView from '@/modules/secret/views/ProjectDetailView.vue';
import { expectComponentSnapshot, mountComponent } from '../../../helpers/mount';

vi.mock('@/modules/project/services/ProjectService', () => ({
  projectService: {
    getBySlug: vi.fn(),
  },
}));

vi.mock('@/modules/secret/services/SecretService', () => ({
  secretService: {
    list: vi.fn(),
    create: vi.fn(),
    updateName: vi.fn(),
    remove: vi.fn(),
    setValue: vi.fn(),
    revealValue: vi.fn(),
    clearValue: vi.fn(),
    exportEnvironment: vi.fn(),
  },
  formatDotEnv: vi.fn(() => ''),
}));

const getBySlugMock = vi.mocked(projectService.getBySlug);
const listMock = vi.mocked(secretService.list);

const session = {
  token: 'jwt-token',
  user: { id: '1', email: 'admin@example.com', role: 'super_admin' as const },
};

const sampleProject = {
  id: 'p1',
  title: 'Payment',
  slug: 'payment',
  teams: [{ id: 't1', name: 'Platform', slug: 'platform' }],
  environments: [
    { id: 'e1', title: 'dev' },
    { id: 'e2', title: 'production' },
  ],
  tags: [],
  created_at: '2026-01-01T00:00:00Z',
  updated_at: '2026-01-01T00:00:00Z',
};

describe('ProjectDetailView', () => {
  beforeEach(() => {
    getBySlugMock.mockReset();
    listMock.mockReset();
    getBySlugMock.mockResolvedValue(sampleProject);
    listMock.mockResolvedValue({
      items: [
        {
          id: 's1',
          name: 'DATABASE_URL',
          values: {
            e1: { has_value: true, updated_at: '2026-01-02T00:00:00Z' },
          },
          created_at: '2026-01-01T00:00:00Z',
          updated_at: '2026-01-02T00:00:00Z',
        },
      ],
    });
  });

  it('renders secret list with environments and actions', async () => {
    const { wrapper } = await mountComponent(ProjectDetailView, {
      session,
      routes: [
        { path: '/', name: 'home', component: { template: '<div />' } },
        { path: '/project/:slug', name: 'project-detail', component: ProjectDetailView },
      ],
      initialRoute: '/project/payment',
    });

    await flushPromises();

    expect(wrapper.text()).toContain('Payment');
    expect(wrapper.text()).toContain('DATABASE_URL');
    expect(wrapper.text()).toContain('dev');
    expect(wrapper.text()).not.toContain('Exportar');
    expect(wrapper.get('input[name="secret-environment-search"]')).toBeTruthy();
    expect(wrapper.get('button[aria-label="Editar"]')).toBeTruthy();
    expect(wrapper.get('button[aria-label="Excluir"]')).toBeTruthy();
    expectComponentSnapshot(wrapper);
  });

  it('filters secrets by environment title', async () => {
    listMock.mockResolvedValue({
      items: [
        {
          id: 's1',
          name: 'DATABASE_URL',
          values: {
            e1: { has_value: true, updated_at: '2026-01-02T00:00:00Z' },
          },
          created_at: '2026-01-01T00:00:00Z',
          updated_at: '2026-01-02T00:00:00Z',
        },
        {
          id: 's2',
          name: 'API_KEY',
          values: {
            e2: { has_value: true, updated_at: '2026-01-02T00:00:00Z' },
          },
          created_at: '2026-01-01T00:00:00Z',
          updated_at: '2026-01-02T00:00:00Z',
        },
      ],
    });

    const { wrapper } = await mountComponent(ProjectDetailView, {
      session,
      routes: [
        { path: '/', name: 'home', component: { template: '<div />' } },
        { path: '/project/:slug', name: 'project-detail', component: ProjectDetailView },
      ],
      initialRoute: '/project/payment',
    });
    await flushPromises();

    await wrapper.get('input[name="secret-environment-search"]').setValue('prod');
    await flushPromises();

    const rows = wrapper.findAll('tbody tr');
    expect(rows).toHaveLength(1);
    expect(rows[0]?.text()).toContain('API_KEY');
    expect(rows[0]?.text()).toContain('production');
    expectComponentSnapshot(wrapper);
  });

  it('filters secrets by name substring', async () => {
    listMock.mockResolvedValue({
      items: [
        {
          id: 's1',
          name: 'TEST_02',
          values: {
            e1: { has_value: true, updated_at: '2026-01-02T00:00:00Z' },
          },
          created_at: '2026-01-01T00:00:00Z',
          updated_at: '2026-01-02T00:00:00Z',
        },
        {
          id: 's2',
          name: 'API_KEY',
          values: {
            e2: { has_value: true, updated_at: '2026-01-02T00:00:00Z' },
          },
          created_at: '2026-01-01T00:00:00Z',
          updated_at: '2026-01-02T00:00:00Z',
        },
      ],
    });

    const { wrapper } = await mountComponent(ProjectDetailView, {
      session,
      routes: [
        { path: '/', name: 'home', component: { template: '<div />' } },
        { path: '/project/:slug', name: 'project-detail', component: ProjectDetailView },
      ],
      initialRoute: '/project/payment',
    });
    await flushPromises();

    await wrapper.get('input[name="secret-environment-search"]').setValue('02');
    await flushPromises();

    const rows = wrapper.findAll('tbody tr');
    expect(rows).toHaveLength(1);
    expect(rows[0]?.text()).toContain('TEST_02');
    expectComponentSnapshot(wrapper);
  });
});
