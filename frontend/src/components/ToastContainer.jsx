import { useToast } from "../context/ToastContext";

export default function ToastContainer() {
  const { toasts, removeToast } = useToast();

  return (
    <div className="toast-container" aria-live="polite" aria-atomic="true">
      {toasts.map((toast) => (
        <div key={toast.id} className={`toast toast-${toast.type}`}>
          <span>{toast.message}</span>
          <button type="button" className="toast-close" onClick={() => removeToast(toast.id)}>
            x
          </button>
        </div>
      ))}
    </div>
  );
}

