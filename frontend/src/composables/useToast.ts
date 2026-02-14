import { ref } from "vue";

export type ToastVariant = "success" | "error" | "info" | "warning";

export interface Toast {
  id: number;
  message: string;
  variant: ToastVariant;
}

const toasts = ref<Toast[]>([]);
let toastId = 0;

export function useToast() {
  const addToast = (message: string, variant: ToastVariant = "info") => {
    const id = ++toastId;
    toasts.value.push({ id, message, variant });

    // Auto-remove after 3 seconds
    setTimeout(() => {
      removeToast(id);
    }, 3000);
  };

  const removeToast = (id: number) => {
    toasts.value = toasts.value.filter((t) => t.id !== id);
  };

  return {
    toasts,
    toast: addToast,
    removeToast,
  };
}
