<script setup>
import { ref, computed } from 'vue';
import { useRoute, useRouter } from 'vue-router';

const props = defineProps({
    userProfile: {
        type: Object,
        required: true
    },
    contactCount: {
        type: Number,
        default: 0
    }
});

const emit = defineEmits(['logout']);

const route = useRoute();
const router = useRouter();
const collapsed = ref(false);

const menuItems = [
    { 
        name: 'Overzicht', 
        icon: 'dashboard',
        route: 'adminDashboard'
    },
    {
        name: 'Evenementen',
        icon: 'calendar',
        route: 'adminEvents'
    },
    {
        name: 'Geocaches',
        icon: 'location',
        route: 'adminGeocaches'
    },
    {
        name: 'Berichtgeving',
        icon: 'megaphone',
        route: 'adminMessages'
    },
    {
        name: 'Contactberichten',
        icon: 'inbox',
        route: 'adminContacts',
        badge: true
    },
    { type: 'divider' },
    {
        name: 'Vertalingen',
        icon: 'globe',
        route: 'adminStatic'
    },
    {
        name: 'Talen',
        icon: 'language',
        route: 'adminLanguages'
    },
    {
        name: 'Sociale media',
        icon: 'share',
        route: 'adminSocials'
    },
    { type: 'divider' },
    {
        name: 'Gebruikers',
        icon: 'users',
        route: 'adminUsers'
    }
];

const isActive = (routeName) => route.name === routeName;

const navigateTo = (routeName) => {
    router.push({ name: routeName });
};

const userInitials = computed(() => {
    const name = props.userProfile?.name || '';
    return name.split(' ').map(n => n[0]).join('').toUpperCase().slice(0, 2) || '?';
});
</script>

<template>
    <aside class="sidebar" :class="{ collapsed }">
        <!-- Logo -->
        <div class="sidebar-brand">
            <div class="brand-icon">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0 1 18 0z"/>
                    <circle cx="12" cy="10" r="3"/>
                </svg>
            </div>
            <span class="brand-text" v-if="!collapsed">Admin Panel</span>
        </div>

        <!-- Navigation -->
        <nav class="sidebar-nav">
            <template v-for="(item, index) in menuItems" :key="index">
                <div v-if="item.type === 'divider'" class="nav-divider"></div>
                <button
                    v-else
                    class="nav-item"
                    :class="{ active: isActive(item.route) }"
                    @click="navigateTo(item.route)"
                    :title="collapsed ? item.name : ''"
                >
                    <span class="nav-icon">
                        <!-- Dashboard -->
                        <svg v-if="item.icon === 'dashboard'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                            <rect x="3" y="3" width="7" height="9" rx="1"/>
                            <rect x="14" y="3" width="7" height="5" rx="1"/>
                            <rect x="14" y="12" width="7" height="9" rx="1"/>
                            <rect x="3" y="16" width="7" height="5" rx="1"/>
                        </svg>
                        <!-- Calendar -->
                        <svg v-else-if="item.icon === 'calendar'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                            <rect x="3" y="4" width="18" height="18" rx="2"/>
                            <line x1="16" y1="2" x2="16" y2="6"/>
                            <line x1="8" y1="2" x2="8" y2="6"/>
                            <line x1="3" y1="10" x2="21" y2="10"/>
                        </svg>
                        <!-- Location -->
                        <svg v-else-if="item.icon === 'location'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                            <path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0 1 18 0z"/>
                            <circle cx="12" cy="10" r="3"/>
                        </svg>
                        <!-- Megaphone -->
                        <svg v-else-if="item.icon === 'megaphone'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                            <path d="M3 11v2a4 4 0 0 0 4 4h1l3 5h2v-5h2a4 4 0 0 0 4-4V9a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2z"/>
                            <path d="M21 11l2-1v4l-2-1"/>
                        </svg>
                        <!-- Inbox -->
                        <svg v-else-if="item.icon === 'inbox'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                            <path d="M22 12h-6l-2 3h-4l-2-3H2"/>
                            <path d="M5.45 5.11L2 12v6a2 2 0 0 0 2 2h16a2 2 0 0 0 2-2v-6l-3.45-6.89A2 2 0 0 0 16.76 4H7.24a2 2 0 0 0-1.79 1.11z"/>
                        </svg>
                        <!-- Globe -->
                        <svg v-else-if="item.icon === 'globe'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                            <circle cx="12" cy="12" r="10"/>
                            <line x1="2" y1="12" x2="22" y2="12"/>
                            <path d="M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1-4 10 15.3 15.3 0 0 1-4-10 15.3 15.3 0 0 1 4-10z"/>
                        </svg>
                        <!-- Language -->
                        <svg v-else-if="item.icon === 'language'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                            <path d="M5 8l6 6"/>
                            <path d="M4 14l6-6 2-3"/>
                            <path d="M2 5h12"/>
                            <path d="M7 2v3"/>
                            <path d="M22 22l-5-10-5 10"/>
                            <path d="M14 18h6"/>
                        </svg>
                        <!-- Share -->
                        <svg v-else-if="item.icon === 'share'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                            <circle cx="18" cy="5" r="3"/>
                            <circle cx="6" cy="12" r="3"/>
                            <circle cx="18" cy="19" r="3"/>
                            <line x1="8.59" y1="13.51" x2="15.42" y2="17.49"/>
                            <line x1="15.41" y1="6.51" x2="8.59" y2="10.49"/>
                        </svg>
                        <!-- Users -->
                        <svg v-else-if="item.icon === 'users'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                            <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
                            <circle cx="9" cy="7" r="4"/>
                            <path d="M23 21v-2a4 4 0 0 0-3-3.87"/>
                            <path d="M16 3.13a4 4 0 0 1 0 7.75"/>
                        </svg>
                    </span>
                    <span class="nav-label" v-if="!collapsed">{{ item.name }}</span>
                    <span v-if="item.badge && contactCount > 0 && !collapsed" class="nav-badge">
                        {{ contactCount > 99 ? '99+' : contactCount }}
                    </span>
                    <span v-if="item.badge && contactCount > 0 && collapsed" class="nav-badge-dot"></span>
                </button>
            </template>
        </nav>

        <!-- Collapse Toggle -->
        <button class="collapse-btn" @click="collapsed = !collapsed" :title="collapsed ? 'Expand' : 'Collapse'">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <polyline v-if="collapsed" points="9 18 15 12 9 6"/>
                <polyline v-else points="15 18 9 12 15 6"/>
            </svg>
        </button>

        <!-- User Section -->
        <div class="sidebar-user">
            <div class="user-avatar">{{ userInitials }}</div>
            <div class="user-info" v-if="!collapsed">
                <span class="user-name">{{ userProfile.name }}</span>
                <span class="user-email">{{ userProfile.email }}</span>
            </div>
            <button class="logout-btn" @click="$emit('logout')" title="Logout">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"/>
                    <polyline points="16 17 21 12 16 7"/>
                    <line x1="21" y1="12" x2="9" y2="12"/>
                </svg>
            </button>
        </div>
    </aside>
</template>

<style scoped>
.sidebar {
    width: var(--sidebar-width);
    min-width: var(--sidebar-width);
    height: 100vh;
    background: var(--admin-surface);
    border-right: 1px solid var(--admin-border);
    display: flex;
    flex-direction: column;
    position: fixed;
    left: 0;
    top: 0;
    z-index: 100;
    transition: width 0.2s ease, min-width 0.2s ease;
}

.sidebar.collapsed {
    width: var(--sidebar-collapsed-width);
    min-width: var(--sidebar-collapsed-width);
}

/* Brand */
.sidebar-brand {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    padding: 1.25rem;
    border-bottom: 1px solid var(--admin-border-light);
}

.brand-icon {
    width: 2.25rem;
    height: 2.25rem;
    background: var(--admin-primary);
    border-radius: var(--admin-radius);
    display: flex;
    align-items: center;
    justify-content: center;
    color: white;
    flex-shrink: 0;
}

.brand-icon svg {
    width: 1.25rem;
    height: 1.25rem;
}

.brand-text {
    font-weight: 600;
    font-size: 1rem;
    color: var(--admin-text);
    white-space: nowrap;
}

/* Navigation */
.sidebar-nav {
    flex: 1;
    padding: 1rem 0.75rem;
    overflow-y: auto;
    overflow-x: hidden;
}

.nav-divider {
    height: 1px;
    background: var(--admin-border-light);
    margin: 0.75rem 0;
}

.nav-item {
    width: 100%;
    display: flex;
    align-items: center;
    gap: 0.75rem;
    padding: 0.625rem 0.75rem;
    border: none;
    background: transparent;
    color: var(--admin-text-secondary);
    font-size: 0.875rem;
    font-weight: 500;
    border-radius: var(--admin-radius);
    cursor: pointer;
    transition: all 0.15s ease;
    position: relative;
    text-align: left;
}

.nav-item:hover {
    background: var(--admin-surface-hover);
    color: var(--admin-text);
}

.nav-item.active {
    background: var(--admin-primary-bg);
    color: var(--admin-primary);
}

.nav-icon {
    width: 1.25rem;
    height: 1.25rem;
    flex-shrink: 0;
}

.nav-icon svg {
    width: 100%;
    height: 100%;
}

.nav-label {
    white-space: nowrap;
    overflow: hidden;
}

.nav-badge {
    margin-left: auto;
    background: var(--admin-danger);
    color: white;
    font-size: 0.6875rem;
    font-weight: 600;
    padding: 0.125rem 0.375rem;
    border-radius: 9999px;
    min-width: 1.25rem;
    text-align: center;
}

.nav-badge-dot {
    position: absolute;
    top: 0.375rem;
    right: 0.375rem;
    width: 0.5rem;
    height: 0.5rem;
    background: var(--admin-danger);
    border-radius: 50%;
}

/* Collapse Button */
.collapse-btn {
    margin: 0 0.75rem;
    padding: 0.5rem;
    border: 1px solid var(--admin-border);
    background: var(--admin-surface);
    border-radius: var(--admin-radius);
    cursor: pointer;
    color: var(--admin-text-muted);
    transition: all 0.15s ease;
}

.collapse-btn:hover {
    background: var(--admin-surface-hover);
    color: var(--admin-text);
}

.collapse-btn svg {
    width: 1rem;
    height: 1rem;
    display: block;
}

/* User Section */
.sidebar-user {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    padding: 1rem;
    border-top: 1px solid var(--admin-border-light);
    margin-top: 0.5rem;
}

.user-avatar {
    width: 2.25rem;
    height: 2.25rem;
    background: var(--admin-primary-bg);
    color: var(--admin-primary);
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    font-weight: 600;
    font-size: 0.8125rem;
    flex-shrink: 0;
}

.user-info {
    flex: 1;
    min-width: 0;
    overflow: hidden;
}

.user-name {
    display: block;
    font-weight: 500;
    font-size: 0.875rem;
    color: var(--admin-text);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
}

.user-email {
    display: block;
    font-size: 0.75rem;
    color: var(--admin-text-muted);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
}

.logout-btn {
    width: 2rem;
    height: 2rem;
    display: flex;
    align-items: center;
    justify-content: center;
    border: none;
    background: transparent;
    color: var(--admin-text-muted);
    border-radius: var(--admin-radius);
    cursor: pointer;
    transition: all 0.15s ease;
    flex-shrink: 0;
}

.logout-btn:hover {
    background: var(--admin-danger-bg);
    color: var(--admin-danger);
}

.logout-btn svg {
    width: 1.125rem;
    height: 1.125rem;
}

/* Collapsed State */
.sidebar.collapsed .sidebar-brand {
    justify-content: center;
    padding: 1.25rem 0.5rem;
}

.sidebar.collapsed .sidebar-nav {
    padding: 1rem 0.5rem;
}

.sidebar.collapsed .nav-item {
    justify-content: center;
    padding: 0.625rem;
}

.sidebar.collapsed .collapse-btn {
    margin: 0 auto;
}

.sidebar.collapsed .sidebar-user {
    flex-direction: column;
    gap: 0.5rem;
    padding: 0.75rem 0.5rem;
}

.sidebar.collapsed .logout-btn {
    margin-top: 0.25rem;
}
</style>
