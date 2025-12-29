<script setup>
import { ref, onMounted, provide, computed } from 'vue';
import { useRouter } from 'vue-router';
import AdminSidebar from '@/components/admin/AdminSidebar.vue';
import ToastNotification from '@/components/admin/ToastNotification.vue';
import '@/css/admin.css';

const props = defineProps({
    pageTitle: {
        type: String,
        default: 'Dashboard'
    }
});

const router = useRouter();
const toastRef = ref(null);

// Auth state
const loggedIn = ref(false);
const doneChecking = ref(false);
const userProfile = ref({});
const needsSetup = ref(false);

// Login/Register form
const formMode = ref('login');
const name = ref('');
const email = ref('');
const password = ref('');
const formError = ref('');
const isSubmitting = ref(false);

// Contact count for badge
const contactCount = ref(0);

// Token management
function getToken() {
    return localStorage.getItem('admin_token');
}

function setToken(token) {
    localStorage.setItem('admin_token', token);
}

function removeToken() {
    localStorage.removeItem('admin_token');
}

// API helper
async function apiRequest(endpoint, options = {}) {
    const token = getToken();
    const headers = {
        'Content-Type': 'application/json',
        'Accept': 'application/json',
        ...(token && { 'Authorization': `Bearer ${token}` }),
        ...options.headers
    };

    const response = await fetch(`${import.meta.env.VITE_API_URL}${endpoint}`, {
        ...options,
        headers
    });

    return response;
}

// Check if setup is needed
async function checkSetupStatus() {
    try {
        const response = await fetch(`${import.meta.env.VITE_API_URL}setup-status`);
        const data = await response.json();
        needsSetup.value = data.needsSetup === true;
        if (needsSetup.value) {
            formMode.value = 'register';
        }
    } catch (err) {
        console.error('Setup status check failed:', err);
    }
}

// Auth functions
async function checkAuth() {
    await checkSetupStatus();
    
    const token = getToken();
    if (!token) {
        doneChecking.value = true;
        return;
    }

    try {
        const response = await apiRequest('admin/profile');
        if (response.ok) {
            const data = await response.json();
            if (data.status && data.data) {
                userProfile.value = data.data;
                loggedIn.value = true;
                fetchContactCount();
            }
        } else {
            removeToken();
        }
    } catch (err) {
        console.error('Auth check failed:', err);
        removeToken();
    }
    
    doneChecking.value = true;
}

async function fetchContactCount() {
    try {
        const response = await apiRequest('admin/contacts?status=new');
        if (response.ok) {
            const data = await response.json();
            contactCount.value = data.total || 0;
        }
    } catch (err) {
        console.error('Failed to fetch contact count:', err);
    }
}

async function handleSubmit() {
    if (isSubmitting.value) return;
    
    formError.value = '';
    isSubmitting.value = true;

    try {
        const endpoint = formMode.value === 'register' ? 'register' : 'login';
        const body = formMode.value === 'register' 
            ? { name: name.value, email: email.value, password: password.value }
            : { email: email.value, password: password.value };

        const response = await apiRequest(endpoint, {
            method: 'POST',
            body: JSON.stringify(body)
        });

        const data = await response.json();

        if (data.status && data.token) {
            setToken(data.token);
            userProfile.value = data.user;
            loggedIn.value = true;
            needsSetup.value = false;
            fetchContactCount();
            window.$toast?.success(formMode.value === 'register' ? 'Account created successfully!' : 'Welcome back!');
        } else if (data.errors) {
            const errors = Object.values(data.errors).flat();
            formError.value = errors.join(', ');
        } else {
            formError.value = data.message || 'Authentication failed';
        }
    } catch (err) {
        console.error('Auth failed:', err);
        formError.value = 'Something went wrong. Please try again.';
    }

    isSubmitting.value = false;
}

async function handleLogout() {
    try {
        await apiRequest('admin/logout', { method: 'POST' });
    } catch (err) {
        // Ignore logout errors
    }
    
    removeToken();
    userProfile.value = {};
    loggedIn.value = false;
    contactCount.value = 0;
    router.push({ name: 'admin' });
    window.$toast?.info('You have been logged out');
}

// Provide auth state and helpers to child components
provide('auth', {
    userProfile,
    apiRequest,
    getToken
});

provide('apiRequest', apiRequest);
provide('refreshContactCount', fetchContactCount);

const showToast = (message, type = 'info') => {
    window.$toast?.[type]?.(message);
};
provide('showToast', showToast);

onMounted(checkAuth);
</script>

<template>
    <ToastNotification ref="toastRef" />
    
    <!-- Loading state -->
    <div v-if="!doneChecking" class="auth-loading admin-panel">
        <div class="admin-spinner"></div>
        <p>Loading...</p>
    </div>

    <!-- Login/Register form -->
    <div v-else-if="!loggedIn" class="login-page admin-panel">
        <div class="login-wrapper">
            <div class="login-card">
                <!-- Header -->
                <div class="login-header">
                    <div class="login-icon">
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                            <path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0 1 18 0z"/>
                            <circle cx="12" cy="10" r="3"/>
                        </svg>
                    </div>
                    <h1 v-if="needsSetup">Welcome!</h1>
                    <h1 v-else>Admin Panel</h1>
                    <p v-if="needsSetup">Create your administrator account to get started</p>
                    <p v-else>Sign in to manage your website</p>
                </div>
                
                <!-- Form -->
                <form @submit.prevent="handleSubmit" class="login-form">
                    <div v-if="formMode === 'register'" class="admin-form-group">
                        <label class="admin-label" for="name">Full Name</label>
                        <input 
                            v-model="name" 
                            type="text" 
                            id="name" 
                            class="admin-input"
                            required 
                            autocomplete="name"
                            placeholder="Your full name"
                        >
                    </div>

                    <div class="admin-form-group">
                        <label class="admin-label" for="email">Email Address</label>
                        <input 
                            v-model="email" 
                            type="email" 
                            id="email" 
                            class="admin-input"
                            required 
                            autocomplete="email"
                            placeholder="you@example.com"
                        >
                    </div>
                    
                    <div class="admin-form-group">
                        <label class="admin-label" for="password">Password</label>
                        <input 
                            v-model="password" 
                            type="password" 
                            id="password" 
                            class="admin-input"
                            required 
                            autocomplete="current-password"
                            placeholder="••••••••"
                            minlength="8"
                        >
                    </div>

                    <div v-if="formError" class="admin-alert admin-alert-danger">
                        <svg class="admin-alert-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                            <circle cx="12" cy="12" r="10"/>
                            <line x1="12" y1="8" x2="12" y2="12"/>
                            <line x1="12" y1="16" x2="12.01" y2="16"/>
                        </svg>
                        <span>{{ formError }}</span>
                    </div>
                    
                    <button type="submit" class="admin-btn admin-btn-primary admin-btn-lg login-submit" :disabled="isSubmitting">
                        <div v-if="isSubmitting" class="admin-spinner" style="width: 1rem; height: 1rem;"></div>
                        <span v-else>{{ formMode === 'register' ? 'Create Account' : 'Sign In' }}</span>
                    </button>
                </form>
            </div>
            
            <p class="login-footer">
                Geocaching Brughia VZW &copy; {{ new Date().getFullYear() }}
            </p>
        </div>
    </div>

    <!-- Admin layout -->
    <div v-else class="admin-panel admin-layout">
        <AdminSidebar 
            :userProfile="userProfile" 
            :contactCount="contactCount"
            @logout="handleLogout" 
        />
        
        <main class="admin-main">
            <header class="admin-header">
                <div class="header-title">
                    <h1>{{ pageTitle }}</h1>
                </div>
                <div class="header-actions">
                    <slot name="actions"></slot>
                </div>
            </header>
            
            <div class="admin-content">
                <slot></slot>
            </div>
        </main>
    </div>
</template>

<style scoped>
/* Loading state */
.auth-loading {
    height: 100vh;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 1rem;
    color: var(--admin-text-secondary);
}

/* Login Page */
.login-page {
    min-height: 100vh;
    display: flex;
    align-items: center;
    justify-content: center;
    background: var(--admin-bg);
    padding: 1.5rem;
}

.login-wrapper {
    width: 100%;
    max-width: 400px;
}

.login-card {
    background: var(--admin-surface);
    border-radius: var(--admin-radius-xl);
    border: 1px solid var(--admin-border);
    box-shadow: var(--admin-shadow-lg);
    padding: 2.5rem;
}

.login-header {
    text-align: center;
    margin-bottom: 2rem;
}

.login-icon {
    width: 3.5rem;
    height: 3.5rem;
    background: var(--admin-primary);
    border-radius: var(--admin-radius-lg);
    display: flex;
    align-items: center;
    justify-content: center;
    color: white;
    margin: 0 auto 1.25rem;
}

.login-icon svg {
    width: 1.75rem;
    height: 1.75rem;
}

.login-header h1 {
    font-size: 1.5rem;
    font-weight: 700;
    color: var(--admin-text);
    margin: 0 0 0.5rem;
}

.login-header p {
    font-size: 0.875rem;
    color: var(--admin-text-secondary);
    margin: 0;
}

.login-form {
    display: flex;
    flex-direction: column;
    gap: 1.25rem;
}

.login-submit {
    width: 100%;
    margin-top: 0.5rem;
}

.login-footer {
    text-align: center;
    font-size: 0.75rem;
    color: var(--admin-text-muted);
    margin-top: 1.5rem;
}

/* Admin Layout */
.admin-layout {
    display: flex;
    min-height: 100vh;
    max-height: 100vh;
    overflow: hidden;
}

.admin-main {
    flex: 1;
    margin-left: var(--sidebar-width);
    display: flex;
    flex-direction: column;
    min-height: 100vh;
    max-height: 100vh;
    overflow: hidden;
    transition: margin-left 0.2s ease;
}

.admin-header {
    background: var(--admin-surface);
    padding: 1.25rem 2rem;
    border-bottom: 1px solid var(--admin-border);
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 1rem;
    position: sticky;
    top: 0;
    z-index: 50;
    flex-shrink: 0;
}

.header-title h1 {
    margin: 0;
    font-size: 1.5rem;
    font-weight: 600;
    color: var(--admin-text);
}

.header-actions {
    display: flex;
    align-items: center;
    gap: 0.75rem;
}

.admin-content {
    flex: 1;
    padding: 1.5rem 2rem 2rem;
    overflow-y: auto;
    overflow-x: hidden;
}
</style>
