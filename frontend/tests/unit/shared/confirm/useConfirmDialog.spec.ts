import { describe, expect, it } from 'vitest';
import { useConfirmDialog } from '@/shared/confirm/useConfirmDialog';

describe('useConfirmDialog', () => {
  it('resolves true on accept and false on cancel', async () => {
    const dialog = useConfirmDialog();

    const pending = dialog.confirm({
      title: 'Title',
      message: 'Message',
      confirmLabel: 'OK',
      cancelLabel: 'Cancel',
      variant: 'danger',
    });
    expect(dialog.open.value).toBe(true);
    expect(dialog.title.value).toBe('Title');

    dialog.accept();
    await expect(pending).resolves.toBe(true);
    expect(dialog.open.value).toBe(false);

    const cancelled = dialog.confirm({
      title: 'T2',
      message: 'M2',
      confirmLabel: 'OK',
      cancelLabel: 'Cancel',
    });
    dialog.cancel();
    await expect(cancelled).resolves.toBe(false);
  });
});
