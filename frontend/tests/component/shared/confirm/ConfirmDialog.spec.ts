import { flushPromises } from '@vue/test-utils';
import { afterEach, beforeEach, describe, expect, it } from 'vitest';
import ConfirmDialog from '@/shared/confirm/ConfirmDialog.vue';
import { useConfirmDialog } from '@/shared/confirm/useConfirmDialog';
import { expectComponentSnapshot, mountComponent } from '../../helpers/mount';

describe('ConfirmDialog (component)', () => {
  beforeEach(() => {
    const dialog = useConfirmDialog();
    if (dialog.open.value) {
      dialog.cancel();
    }
  });

  afterEach(() => {
    const dialog = useConfirmDialog();
    if (dialog.open.value) {
      dialog.cancel();
    }
  });

  it('renders closed by default', async () => {
    const { wrapper } = await mountComponent(ConfirmDialog);

    expect(wrapper.get('#app-confirm-dialog').attributes('open')).toBeUndefined();
    expectComponentSnapshot(wrapper);
  });

  it('shows title, message and danger confirm button when opened', async () => {
    const { wrapper } = await mountComponent(ConfirmDialog);
    const dialog = useConfirmDialog();

    void dialog.confirm({
      title: 'Revogar chave?',
      message: 'Esta ação não pode ser desfeita.',
      confirmLabel: 'Revogar',
      cancelLabel: 'Cancelar',
      variant: 'danger',
    });
    await flushPromises();

    expect(wrapper.get('h3').text()).toBe('Revogar chave?');
    expect(wrapper.get('p').text()).toBe('Esta ação não pode ser desfeita.');
    expect(wrapper.get('[data-testid="confirm-cancel"]').text()).toBe('Cancelar');
    const accept = wrapper.get('[data-testid="confirm-accept"]');
    expect(accept.text()).toBe('Revogar');
    expect(accept.classes()).toContain('btn-error');
    expectComponentSnapshot(wrapper);
  });

  it('uses primary variant when requested', async () => {
    const { wrapper } = await mountComponent(ConfirmDialog);
    const dialog = useConfirmDialog();

    void dialog.confirm({
      title: 'Confirmar',
      message: 'Prosseguir?',
      confirmLabel: 'OK',
      cancelLabel: 'Voltar',
      variant: 'primary',
    });
    await flushPromises();

    expect(wrapper.get('[data-testid="confirm-accept"]').classes()).toContain('btn-primary');
    expectComponentSnapshot(wrapper);
  });

  it('resolves true when accept is clicked', async () => {
    const { wrapper } = await mountComponent(ConfirmDialog);
    const dialog = useConfirmDialog();

    const pending = dialog.confirm({
      title: 'Title',
      message: 'Message',
      confirmLabel: 'OK',
      cancelLabel: 'Cancel',
    });
    await flushPromises();

    await wrapper.get('[data-testid="confirm-accept"]').trigger('click');
    await expect(pending).resolves.toBe(true);
    expect(dialog.open.value).toBe(false);
  });

  it('resolves false when cancel is clicked', async () => {
    const { wrapper } = await mountComponent(ConfirmDialog);
    const dialog = useConfirmDialog();

    const pending = dialog.confirm({
      title: 'Title',
      message: 'Message',
      confirmLabel: 'OK',
      cancelLabel: 'Cancel',
    });
    await flushPromises();

    await wrapper.get('[data-testid="confirm-cancel"]').trigger('click');
    await expect(pending).resolves.toBe(false);
    expect(dialog.open.value).toBe(false);
  });
});
