import { flushPromises } from '@vue/test-utils';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { tagService } from '@/modules/tag/services/TagService';
import TagsView from '@/modules/tag/views/TagsView.vue';
import { expectComponentSnapshot, mountComponent } from '../../../helpers/mount';

vi.mock('@/modules/tag/services/TagService', () => ({
  tagService: {
    list: vi.fn(),
    create: vi.fn(),
    update: vi.fn(),
    remove: vi.fn(),
  },
}));

vi.mock('@/modules/tag/types', async (importOriginal) => {
  const actual = await importOriginal<typeof import('@/modules/tag/types')>();
  return {
    ...actual,
    randomTagColor: () => '#6366F1',
  };
});

const listMock = vi.mocked(tagService.list);

const userSession = {
  token: 'jwt-token',
  user: { id: '2', email: 'user@example.com', role: 'user' as const },
};

const sampleTags = [
  {
    id: 'tag1',
    title: 'Backend',
    color: '#6366F1',
    created_at: '2026-01-01T00:00:00Z',
    updated_at: '2026-01-02T00:00:00Z',
  },
];

describe('TagsView (component)', () => {
  beforeEach(() => {
    listMock.mockReset();
    listMock.mockResolvedValue(sampleTags);
  });

  it('renders tags and mutate actions for standard user', async () => {
    const { wrapper } = await mountComponent(TagsView, { session: userSession });
    await flushPromises();

    expect(wrapper.get('h1').text()).toBe('Tags');
    expect(wrapper.text()).toContain('Backend');
    expect(wrapper.text()).toContain('#6366F1');
    expect(wrapper.text()).toContain('Criar tag');
    expect(wrapper.get('button[aria-label="Editar"]')).toBeTruthy();
    expect(wrapper.get('button[aria-label="Excluir"]')).toBeTruthy();
    expectComponentSnapshot(wrapper);
  });

  it('searches by title via API', async () => {
    listMock.mockResolvedValueOnce(sampleTags).mockResolvedValueOnce([]);

    const { wrapper } = await mountComponent(TagsView, { session: userSession });
    await flushPromises();

    await wrapper.get('input[type="search"]').setValue('back');
    await wrapper.get('form').trigger('submit');
    await flushPromises();

    expect(listMock).toHaveBeenLastCalledWith({ title: 'back' });
    expect(wrapper.text()).toContain('Nenhuma tag encontrada para essa busca.');
    expectComponentSnapshot(wrapper);
  });
});
