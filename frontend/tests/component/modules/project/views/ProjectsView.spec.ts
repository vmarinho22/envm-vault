import { flushPromises } from '@vue/test-utils';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { environmentService } from '@/modules/environment/services/EnvironmentService';
import { projectService } from '@/modules/project/services/ProjectService';
import ProjectsView from '@/modules/project/views/ProjectsView.vue';
import { tagService } from '@/modules/tag/services/TagService';
import { teamService } from '@/modules/team/services/TeamService';
import { userService } from '@/modules/user/services/UserService';
import { expectComponentSnapshot, mountComponent } from '../../../helpers/mount';

vi.mock('@/modules/project/services/ProjectService', () => ({
  projectService: {
    list: vi.fn(),
    getBySlug: vi.fn(),
    create: vi.fn(),
    update: vi.fn(),
    remove: vi.fn(),
  },
}));

vi.mock('@/modules/team/services/TeamService', () => ({
  teamService: {
    list: vi.fn(),
    create: vi.fn(),
  },
}));

vi.mock('@/modules/environment/services/EnvironmentService', () => ({
  environmentService: {
    list: vi.fn(),
    create: vi.fn(),
  },
}));

vi.mock('@/modules/tag/services/TagService', () => ({
  tagService: {
    list: vi.fn(),
    create: vi.fn(),
  },
}));

vi.mock('@/modules/user/services/UserService', () => ({
  userService: {
    list: vi.fn(),
  },
}));

vi.mock('@/modules/tag/types', async (importOriginal) => {
  const actual = await importOriginal<typeof import('@/modules/tag/types')>();
  return {
    ...actual,
    randomTagColor: () => '#6366F1',
  };
});

const listMock = vi.mocked(projectService.list);
const teamListMock = vi.mocked(teamService.list);
const teamCreateMock = vi.mocked(teamService.create);
const envListMock = vi.mocked(environmentService.list);
const envCreateMock = vi.mocked(environmentService.create);
const tagListMock = vi.mocked(tagService.list);
const tagCreateMock = vi.mocked(tagService.create);
const userListMock = vi.mocked(userService.list);

const adminSession = {
  token: 'jwt-token',
  user: { id: '1', email: 'admin@example.com', role: 'super_admin' as const },
};

const userSession = {
  token: 'jwt-token',
  user: { id: '2', email: 'user@example.com', role: 'user' as const },
};

const projectRoutes = [
  { path: '/', name: 'home', component: ProjectsView },
  { path: '/project/:slug', name: 'project-detail', component: { template: '<div />' } },
];

const sampleProjects = {
  items: [
    {
      id: 'p1',
      title: 'Payment',
      slug: 'payment',
      teams: [{ id: 't1', name: 'Platform', slug: 'platform' }],
      environments: [{ id: 'e1', title: 'dev' }],
      tags: [{ id: 'g1', title: 'backend', color: '#000000' }],
      created_at: '2026-01-01T00:00:00Z',
      updated_at: '2026-01-02T00:00:00Z',
    },
  ],
  page: 1,
  page_size: 12,
  total: 1,
};

describe('ProjectsView (component)', () => {
  beforeEach(() => {
    listMock.mockReset();
    teamListMock.mockReset();
    teamCreateMock.mockReset();
    envListMock.mockReset();
    envCreateMock.mockReset();
    tagListMock.mockReset();
    tagCreateMock.mockReset();
    userListMock.mockReset();
    listMock.mockResolvedValue(sampleProjects);
    teamListMock.mockResolvedValue([
      {
        id: 't1',
        name: 'Platform',
        slug: 'platform',
        accountable: { id: 'u1', email: 'a@example.com' },
        created_at: '',
        updated_at: '',
      },
    ]);
    envListMock.mockResolvedValue([
      { id: 'e1', title: 'dev', created_at: '', updated_at: '' },
    ]);
    tagListMock.mockResolvedValue([
      { id: 'g1', title: 'backend', color: '#000000', created_at: '', updated_at: '' },
    ]);
    userListMock.mockResolvedValue([
      { id: 'u1', email: 'a@example.com', role: 'super_admin', token: null },
    ]);
  });

  it('renders projects with create/edit and delete for admin', async () => {
    const { wrapper } = await mountComponent(ProjectsView, {
      session: adminSession,
      routes: projectRoutes,
    });
    await flushPromises();

    expect(wrapper.get('h1').text()).toBe('Projetos');
    expect(wrapper.text()).toContain('Payment');
    expect(wrapper.text()).toContain('Platform');
    expect(wrapper.text()).toContain('Criar projeto');
    expect(wrapper.get('button[aria-label="Editar"]')).toBeTruthy();
    expect(wrapper.get('button[aria-label="Excluir"]')).toBeTruthy();
    expectComponentSnapshot(wrapper);
  });

  it('hides delete for standard user but keeps create/edit', async () => {
    const { wrapper } = await mountComponent(ProjectsView, {
      session: userSession,
      routes: projectRoutes,
    });
    await flushPromises();

    expect(wrapper.text()).toContain('Payment');
    expect(wrapper.text()).toContain('Criar projeto');
    expect(wrapper.get('button[aria-label="Editar"]')).toBeTruthy();
    expect(wrapper.find('button[aria-label="Excluir"]').exists()).toBe(false);
    expectComponentSnapshot(wrapper);
  });

  it('searches by title via API', async () => {
    listMock
      .mockResolvedValueOnce(sampleProjects)
      .mockResolvedValueOnce({ items: [], page: 1, page_size: 12, total: 0 });

    const { wrapper } = await mountComponent(ProjectsView, {
      session: userSession,
      routes: projectRoutes,
    });
    await flushPromises();

    await wrapper.get('input[name="project-title-search"]').setValue('pay');
    await wrapper.get('form').trigger('submit');
    await flushPromises();

    expect(listMock).toHaveBeenLastCalledWith({
      title: 'pay',
      team_ids: undefined,
      tag_ids: undefined,
      page: 1,
      page_size: 12,
    });
    expect(wrapper.text()).toContain('Nenhum projeto encontrado para essa busca.');
    expectComponentSnapshot(wrapper);
  });

  it('shows catalog add buttons by role and creates tag from project form', async () => {
    const { wrapper } = await mountComponent(ProjectsView, {
      session: userSession,
      routes: projectRoutes,
    });
    await flushPromises();

    await wrapper.get('button.btn-primary').trigger('click');
    await flushPromises();

    expect(wrapper.find('button[aria-label="Adicionar time"]').exists()).toBe(false);
    expect(wrapper.find('button[aria-label="Adicionar ambiente"]').exists()).toBe(false);
    expect(wrapper.get('button[aria-label="Adicionar tag"]')).toBeTruthy();

    tagCreateMock.mockResolvedValue({
      id: 'g2',
      title: 'mobile',
      color: '#6366F1',
      created_at: '',
      updated_at: '',
    });
    tagListMock.mockResolvedValueOnce([
      { id: 'g1', title: 'backend', color: '#000000', created_at: '', updated_at: '' },
      { id: 'g2', title: 'mobile', color: '#6366F1', created_at: '', updated_at: '' },
    ]);

    await wrapper.get('button[aria-label="Adicionar tag"]').trigger('click');
    await flushPromises();

    expect(wrapper.text()).toContain('Nova tag');
    await wrapper.get('input[name="quick-tag-title"]').setValue('mobile');
    await wrapper.get('#project-form-modal form').trigger('submit');
    await flushPromises();

    expect(tagCreateMock).toHaveBeenCalledWith({ title: 'mobile', color: '#6366F1' });
    expect(wrapper.text()).toContain('mobile');
    expectComponentSnapshot(wrapper);
  });

  it('lets admin create environment from project form and selects it', async () => {
    const { wrapper } = await mountComponent(ProjectsView, {
      session: adminSession,
      routes: projectRoutes,
    });
    await flushPromises();

    await wrapper.get('button.btn-primary').trigger('click');
    await flushPromises();

    expect(wrapper.get('button[aria-label="Adicionar ambiente"]')).toBeTruthy();

    envCreateMock.mockResolvedValue({
      id: 'e2',
      title: 'staging',
      created_at: '',
      updated_at: '',
    });
    envListMock.mockResolvedValueOnce([
      { id: 'e1', title: 'dev', created_at: '', updated_at: '' },
      { id: 'e2', title: 'staging', created_at: '', updated_at: '' },
    ]);

    await wrapper.get('button[aria-label="Adicionar ambiente"]').trigger('click');
    await flushPromises();

    expect(wrapper.text()).toContain('Novo ambiente');
    await wrapper.get('input[name="quick-environment-title"]').setValue('staging');
    await wrapper.get('#project-form-modal form').trigger('submit');
    await flushPromises();

    expect(envCreateMock).toHaveBeenCalledWith({ title: 'staging' });
    expect(wrapper.text()).toContain('staging');
    expectComponentSnapshot(wrapper);
  });

  it('lets admin create team from project form and selects it', async () => {
    const { wrapper } = await mountComponent(ProjectsView, {
      session: adminSession,
      routes: projectRoutes,
    });
    await flushPromises();

    await wrapper.get('button.btn-primary').trigger('click');
    await flushPromises();

    expect(wrapper.get('button[aria-label="Adicionar time"]')).toBeTruthy();

    teamCreateMock.mockResolvedValue({
      id: 't2',
      name: 'Mobile',
      slug: 'mobile',
      accountable: { id: 'u1', email: 'a@example.com' },
      created_at: '',
      updated_at: '',
    });
    teamListMock.mockResolvedValueOnce([
      {
        id: 't1',
        name: 'Platform',
        slug: 'platform',
        accountable: { id: 'u1', email: 'a@example.com' },
        created_at: '',
        updated_at: '',
      },
      {
        id: 't2',
        name: 'Mobile',
        slug: 'mobile',
        accountable: { id: 'u1', email: 'a@example.com' },
        created_at: '',
        updated_at: '',
      },
    ]);

    await wrapper.get('button[aria-label="Adicionar time"]').trigger('click');
    await flushPromises();

    expect(userListMock).toHaveBeenCalled();
    expect(wrapper.text()).toContain('Novo time');
    await wrapper.get('input[name="quick-team-name"]').setValue('Mobile');
    await wrapper.get('#project-form-modal form').trigger('submit');
    await flushPromises();

    expect(teamCreateMock).toHaveBeenCalledWith({
      name: 'Mobile',
      accountable_id: 'u1',
    });
    expect(wrapper.text()).toContain('Mobile');
    expectComponentSnapshot(wrapper);
  });
});
