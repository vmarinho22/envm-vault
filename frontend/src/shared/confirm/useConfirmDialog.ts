import { ref } from 'vue';

export type ConfirmVariant = 'primary' | 'danger';

export interface ConfirmOptions {
  title: string;
  message: string;
  confirmLabel: string;
  cancelLabel: string;
  variant?: ConfirmVariant;
}

const open = ref(false);
const title = ref('');
const message = ref('');
const confirmLabel = ref('');
const cancelLabel = ref('');
const variant = ref<ConfirmVariant>('danger');

let resolver: ((confirmed: boolean) => void) | null = null;

function settle(confirmed: boolean): void {
  open.value = false;
  const resolve = resolver;
  resolver = null;
  resolve?.(confirmed);
}

/**
 * App-wide confirm dialog. Mount `<ConfirmDialog />` once (e.g. DefaultLayout).
 * Call `confirm({...})` from any composable/view — returns a Promise<boolean>.
 */
export function useConfirmDialog() {
  function confirm(options: ConfirmOptions): Promise<boolean> {
    if (resolver) {
      settle(false);
    }

    title.value = options.title;
    message.value = options.message;
    confirmLabel.value = options.confirmLabel;
    cancelLabel.value = options.cancelLabel;
    variant.value = options.variant ?? 'danger';
    open.value = true;

    return new Promise<boolean>((resolve) => {
      resolver = resolve;
    });
  }

  function accept(): void {
    settle(true);
  }

  function cancel(): void {
    settle(false);
  }

  return {
    open,
    title,
    message,
    confirmLabel,
    cancelLabel,
    variant,
    confirm,
    accept,
    cancel,
  };
}
