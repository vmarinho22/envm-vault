import { flushPromises } from '@vue/test-utils';
import { describe, expect, it } from 'vitest';
import ProjectDownloadSnippet from '@/modules/secret/components/ProjectDownloadSnippet.vue';
import { expectComponentSnapshot, mountComponent } from '../../../helpers/mount';

describe('ProjectDownloadSnippet', () => {
  it('renders linux curl snippet with project slug and env', async () => {
    const { wrapper } = await mountComponent(ProjectDownloadSnippet, {
      props: {
        projectSlug: 'envm-api',
        environments: [
          { id: 'e1', title: 'dev' },
          { id: 'e2', title: 'production' },
        ],
      },
      session: {
        token: 'jwt',
        user: { id: 'u1', email: 'user@example.com', role: 'user' },
      },
    });
    await flushPromises();

    expect(wrapper.text()).toContain('Baixar .env no projeto');
    expect(wrapper.get('[data-testid="download-snippet"]').text()).toContain('curl -fsSL');
    expect(wrapper.get('[data-testid="download-snippet"]').text()).toContain(
      '/v1/projects/envm-api/env?environment=dev&format=env',
    );
    expect(wrapper.get('[data-testid="production-download-note"]').text()).toContain(
      'super admins',
    );
    expectComponentSnapshot(wrapper);
  });

  it('includes production for super_admin', async () => {
    const { wrapper } = await mountComponent(ProjectDownloadSnippet, {
      props: {
        projectSlug: 'envm-api',
        environments: [
          { id: 'e1', title: 'dev' },
          { id: 'e2', title: 'production' },
        ],
      },
      session: {
        token: 'jwt',
        user: { id: 'u1', email: 'admin@example.com', role: 'super_admin' },
      },
    });
    await flushPromises();

    const options = wrapper.findAll('option').map((o) => o.text());
    expect(options).toContain('production');
    expect(wrapper.find('[data-testid="production-download-note"]').exists()).toBe(false);
    expectComponentSnapshot(wrapper);
  });
});
