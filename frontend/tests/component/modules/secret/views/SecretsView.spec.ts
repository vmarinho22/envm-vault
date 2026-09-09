import { flushPromises } from '@vue/test-utils';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { secretService } from '@/modules/secret/services/SecretService';
import SecretsView from '@/modules/secret/views/SecretsView.vue';
import { expectComponentSnapshot, mountComponent } from '../../../helpers/mount';

vi.mock('@/modules/secret/services/SecretService', () => ({
  secretService: {
    listCatalog: vi.fn(),
  },
}));

const listCatalogMock = vi.mocked(secretService.listCatalog);

const session = {
  token: 'jwt-token',
  user: { id: '1', email: 'admin@example.com', role: 'super_admin' as const },
};

describe('SecretsView', () => {
  beforeEach(() => {
    listCatalogMock.mockReset();
    listCatalogMock.mockResolvedValue({
      items: [
        {
          id: 's1',
          name: 'DATABASE_URL',
          project_count: 2,
          created_at: '2026-01-01T00:00:00Z',
          updated_at: '2026-01-02T00:00:00Z',
        },
      ],
      page: 1,
      page_size: 12,
      total: 1,
    });
  });

  it('renders paginated catalog list with clickable rows to detail', async () => {
    const { wrapper, router } = await mountComponent(SecretsView, {
      session,
      routes: [
        { path: '/secrets', name: 'secrets', component: SecretsView },
        { path: '/secrets/:id', name: 'secret-detail', component: { template: '<div />' } },
      ],
      initialRoute: '/secrets',
    });
    await flushPromises();

    expect(wrapper.text()).toContain('Secrets');
    expect(wrapper.text()).toContain('DATABASE_URL');
    expect(wrapper.text()).toContain('2');
    const row = wrapper.get('tbody tr[role="link"]');
    expect(row.text()).toContain('DATABASE_URL');
    await row.trigger('click');
    await flushPromises();
    expect(router.currentRoute.value.name).toBe('secret-detail');
    expect(router.currentRoute.value.params.id).toBe('s1');
    expect(listCatalogMock).toHaveBeenCalledWith({
      title: undefined,
      page: 1,
      pageSize: 12,
    });
    expectComponentSnapshot(wrapper);
  });

  it('applies title search', async () => {
    const { wrapper } = await mountComponent(SecretsView, {
      session,
      routes: [
        { path: '/secrets', name: 'secrets', component: SecretsView },
        { path: '/secrets/:id', name: 'secret-detail', component: { template: '<div />' } },
      ],
      initialRoute: '/secrets',
    });
    await flushPromises();

    await wrapper.get('input[name="secret-catalog-search"]').setValue('DB');
    await wrapper.get('form').trigger('submit');
    await flushPromises();

    expect(listCatalogMock).toHaveBeenLastCalledWith({
      title: 'DB',
      page: 1,
      pageSize: 12,
    });
    expectComponentSnapshot(wrapper);
  });
});
