<script setup>
import { ref } from 'vue';

const toasts = ref([]);
let toastId = 0;

function addToast(message, type = 'info', duration = 4000) {
    const id = ++toastId;
    toasts.value.push({ id, message, type, visible: true });
    
    if (duration > 0) {
        setTimeout(() => removeToast(id), duration);
    }
    
    return id;
}

function removeToast(id) {
    const index = toasts.value.findIndex(t => t.id === id);
    if (index > -1) {
        toasts.value[index].visible = false;
        setTimeout(() => {
            toasts.value = toasts.value.filter(t => t.id !== id);
        }, 300);
    }
}

// Expose methods globally
if (typeof window !== 'undefined') {
    window.$toast = {
        success: (msg, duration) => addToast(msg, 'success', duration),
        error: (msg, duration) => addToast(msg, 'error', duration),
        warning: (msg, duration) => addToast(msg, 'warning', duration),
        info: (msg, duration) => addToast(msg, 'info', duration),
    };
}

defineExpose({ addToast, removeToast });
</script>

<template>
    <Teleport to="body">
        <div class="toast-container">
            <TransitionGroup name="toast">
                <div 
                    v-for="toast in toasts" 
                    :key="toast.id"
                    :class="['toast', `toast-${toast.type}`, { 'toast-hiding': !toast.visible }]"
                >
                    <span class="toast-icon">
                        <!-- Success -->
                        <svg v-if="toast.type === 'success'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                            <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/>
                            <polyline points="22 4 12 14.01 9 11.01"/>
                        </svg>
                        <!-- Error -->
                        <svg v-else-if="toast.type === 'error'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                            <circle cx="12" cy="12" r="10"/>
                            <line x1="15" y1="9" x2="9" y2="15"/>
                            <line x1="9" y1="9" x2="15" y2="15"/>
                        </svg>
                        <!-- Warning -->
                        <svg v-else-if="toast.type === 'warning'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                            <path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"/>
                            <line x1="12" y1="9" x2="12" y2="13"/>
                            <line x1="12" y1="17" x2="12.01" y2="17"/>
                        </svg>
                        <!-- Info -->
                        <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                            <circle cx="12" cy="12" r="10"/>
                            <line x1="12" y1="16" x2="12" y2="12"/>
                            <line x1="12" y1="8" x2="12.01" y2="8"/>
                        </svg>
                    </span>
                    <span class="toast-message">{{ toast.message }}</span>
                    <button class="toast-close" @click="removeToast(toast.id)" aria-label="Close">
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                            <line x1="18" y1="6" x2="6" y2="18"/>
                            <line x1="6" y1="6" x2="18" y2="18"/>
                        </svg>
                    </button>
                </div>
            </TransitionGroup>
        </div>
    </Teleport>
</template>

<style scoped>
.toast-container {
    position: fixed;
    bottom: 1.5rem;
    right: 1.5rem;
    z-index: 9999;
    display: flex;
    flex-direction: column-reverse;
    gap: 0.75rem;
    max-width: 24rem;
}

.toast {
    display: flex;
    align-items: flex-start;
    gap: 0.75rem;
    padding: 1rem 1.25rem;
    border-radius: 0.75rem;
    background: #1e293b;
    color: white;
    box-shadow: 0 10px 25px -5px rgba(0, 0, 0, 0.25), 0 8px 10px -6px rgba(0, 0, 0, 0.1);
    transition: all 0.3s ease;
}

.toast-hiding {
    opacity: 0;
    transform: translateX(100%) scale(0.95);
}

.toast-icon {
    flex-shrink: 0;
    width: 1.25rem;
    height: 1.25rem;
    margin-top: 0.0625rem;
}

.toast-icon svg {
    width: 100%;
    height: 100%;
}

.toast-message {
    flex: 1;
    font-size: 0.875rem;
    line-height: 1.5;
}

.toast-close {
    flex-shrink: 0;
    background: none;
    border: none;
    color: rgba(255, 255, 255, 0.5);
    cursor: pointer;
    padding: 0.125rem;
    border-radius: 0.25rem;
    transition: all 0.15s ease;
}

.toast-close:hover {
    color: white;
    background: rgba(255, 255, 255, 0.1);
}

.toast-close svg {
    width: 1rem;
    height: 1rem;
    display: block;
}

.toast-success .toast-icon {
    color: #4ade80;
}

.toast-error .toast-icon {
    color: #f87171;
}

.toast-warning .toast-icon {
    color: #fbbf24;
}

.toast-info .toast-icon {
    color: #60a5fa;
}

/* Transitions */
.toast-enter-active {
    animation: toastSlideIn 0.3s ease;
}

.toast-leave-active {
    animation: toastSlideOut 0.3s ease;
}

@keyframes toastSlideIn {
    from {
        opacity: 0;
        transform: translateX(100%) scale(0.95);
    }
    to {
        opacity: 1;
        transform: translateX(0) scale(1);
    }
}

@keyframes toastSlideOut {
    from {
        opacity: 1;
        transform: translateX(0) scale(1);
    }
    to {
        opacity: 0;
        transform: translateX(100%) scale(0.95);
    }
}

@media (max-width: 640px) {
    .toast-container {
        left: 1rem;
        right: 1rem;
        bottom: 1rem;
        max-width: none;
    }
}
</style>
