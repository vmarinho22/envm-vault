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
    search: vi.fn(),
    list: vi.fn(),
    create: vi.fn(),
    attach: vi.fn(),
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
const searchMock = vi.mocked(secretService.search);
const attachMock = vi.mocked(secretService.attach);

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

function secretItem(
  overrides: Partial<{
    id: string;
    name: string;
    project_count: number;
    values: Record<string, { has_value: boolean; updated_at: string }>;
  }> = {},
) {
  return {
    id: 's1',
    name: 'DATABASE_URL',
    project_count: 1,
    values: {
      e1: { has_value: true, updated_at: '2026-01-02T00:00:00Z' },
    },
    created_at: '2026-01-01T00:00:00Z',
    updated_at: '2026-01-02T00:00:00Z',
    ...overrides,
  };
}

const detailRoutes = [
  { path: '/', name: 'home', component: { template: '<div />' } },
  { path: '/project/:slug', name: 'project-detail', component: ProjectDetailView },
];

describe('ProjectDetailView', () => {
  beforeEach(() => {
    getBySlugMock.mockReset();
    listMock.mockReset();
    searchMock.mockReset();
    attachMock.mockReset();
    getBySlugMock.mockResolvedValue(sampleProject);
    listMock.mockResolvedValue({ items: [secretItem()] });
    searchMock.mockResolvedValue({ items: [] });
  });

  it('renders secret list with environments and actions', async () => {
    const { wrapper } = await mountComponent(ProjectDetailView, {
      session,
      routes: detailRoutes,
      initialRoute: '/project/payment',
    });

    await flushPromises();

    expect(wrapper.text()).toContain('Payment');
    expect(wrapper.text()).toContain('DATABASE_URL');
    expect(wrapper.text()).toContain('dev');
    expect(wrapper.text()).not.toContain('Exportar');
    expect(wrapper.get('input[name="secret-environment-search"]')).toBeTruthy();
    expect(wrapper.get('button[aria-label="Editar"]')).toBeTruthy();
    expect(wrapper.get('button[aria-label="Remover do projeto"]')).toBeTruthy();
    expectComponentSnapshot(wrapper);
  });

  it('filters secrets by environment title', async () => {
    listMock.mockResolvedValue({
      items: [
        secretItem(),
        secretItem({
          id: 's2',
          name: 'API_KEY',
          values: {
            e2: { has_value: true, updated_at: '2026-01-02T00:00:00Z' },
          },
        }),
      ],
    });

    const { wrapper } = await mountComponent(ProjectDetailView, {
      session,
      routes: detailRoutes,
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
        secretItem({ id: 's1', name: 'TEST_02' }),
        secretItem({
          id: 's2',
          name: 'API_KEY',
          values: {
            e2: { has_value: true, updated_at: '2026-01-02T00:00:00Z' },
          },
        }),
      ],
    });

    const { wrapper } = await mountComponent(ProjectDetailView, {
      session,
      routes: detailRoutes,
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

  it('opens create modal when Buscar finds no catalog hits', async () => {
    searchMock.mockResolvedValue({ items: [] });

    const { wrapper } = await mountComponent(ProjectDetailView, {
      session,
      routes: detailRoutes,
      initialRoute: '/project/payment',
    });
    await flushPromises();

    await wrapper.get('input[name="secret-environment-search"]').setValue('NEW_VAR');
    await wrapper.get('form').trigger('submit');
    await flushPromises();

    expect(searchMock).toHaveBeenCalledWith({
      title: 'NEW_VAR',
      excludeProjectId: 'p1',
      limit: 20,
    });
    expect((wrapper.get('input[name="secret-name"]').element as HTMLInputElement).value).toBe(
      'NEW_VAR',
    );
    expectComponentSnapshot(wrapper);
  });

  it('attaches catalog hit when user clicks Anexar', async () => {
    searchMock.mockResolvedValue({
      items: [
        {
          id: 's9',
          name: 'SHARED_KEY',
          project_count: 2,
          created_at: '2026-01-01T00:00:00Z',
          updated_at: '2026-01-01T00:00:00Z',
        },
      ],
    });
    attachMock.mockResolvedValue(
      secretItem({ id: 's9', name: 'SHARED_KEY', project_count: 3, values: {} }),
    );
    listMock
      .mockResolvedValueOnce({ items: [secretItem()] })
      .mockResolvedValueOnce({
        items: [
          secretItem(),
          secretItem({ id: 's9', name: 'SHARED_KEY', project_count: 3, values: {} }),
        ],
      });

    const { wrapper } = await mountComponent(ProjectDetailView, {
      session,
      routes: detailRoutes,
      initialRoute: '/project/payment',
    });
    await flushPromises();

    await wrapper.get('input[name="secret-environment-search"]').setValue('SHARED');
    await wrapper.get('form').trigger('submit');
    await flushPromises();

    expect(wrapper.text()).toContain('No catálogo do sistema');
    expect(wrapper.text()).toContain('SHARED_KEY');

    const attachBtn = wrapper.findAll('button').find((b) => b.text().includes('Anexar'));
    expect(attachBtn).toBeTruthy();
    await attachBtn!.trigger('click');
    await flushPromises();

    expect(attachMock).toHaveBeenCalledWith('p1', 's9');
    expect(wrapper.text()).toContain('SHARED_KEY');
    expectComponentSnapshot(wrapper);
  });
});
