<script setup>
import { ref, computed, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import AdminLayout from '@/components/admin/AdminLayout.vue';
import config from '@/data/config.js';
import { getAdminGoldenKeyMonths } from '@/services/GoldenKeyMonthService';

const router = useRouter();
const loading = ref(true);

const stats = ref({
    events: { total: 0, published: 0, draft: 0 },
    geocaches: { total: 0, active: 0 },
    contacts: { newCount: 0, inProgress: 0 },
    messages: { total: 0, published: 0 },
    languages: { total: 0, active: 0 },
    users: { total: 0 }
});

const goldenKey = ref({
    isActive: false,
    activationTime: null,
    months: []
});

const gkActiveMonths = computed(() => goldenKey.value.months.filter(m => m.state === 'active'));
const gkFoundMonths = computed(() => goldenKey.value.months.filter(m => m.state === 'found'));
const gkLockedMonths = computed(() => goldenKey.value.months.filter(m => m.state === 'locked'));

const recentContacts = ref([]);

function getToken() {
    return localStorage.getItem('admin_token');
}

async function apiRequest(endpoint) {
    const token = getToken();
    const headers = {
        'Content-Type': 'application/json',
        'Accept': 'application/json',
        ...(token && { 'Authorization': `Bearer ${token}` })
    };

    try {
        const response = await fetch(`${config.apiUrl}${endpoint}`, { headers });
        if (!response.ok) return null;
        return await response.json();
    } catch (err) {
        console.error('API request failed:', err);
        return null;
    }
}

async function fetchDashboardData() {
    try {
        const contactsData = await apiRequest('admin/contacts?status=new');
        if (contactsData && contactsData.data) {
            recentContacts.value = contactsData.data.slice(0, 5);
            stats.value.contacts.newCount = contactsData.total || 0;
        }

        const inProgressData = await apiRequest('admin/contacts?status=in_progress');
        if (inProgressData) {
            stats.value.contacts.inProgress = inProgressData.total || 0;
        }

        const eventsData = await apiRequest('admin/events');
        const events = eventsData?.data || (Array.isArray(eventsData) ? eventsData : []);
        stats.value.events.total = events.length;
        stats.value.events.published = events.filter(e => e.state === 'published' || e.state === 'ONLINE').length;
        stats.value.events.draft = events.filter(e => e.state === 'draft').length;

        const geocachesData = await apiRequest('admin/geocaches');
        const geocaches = geocachesData?.data || (Array.isArray(geocachesData) ? geocachesData : []);
        stats.value.geocaches.total = geocaches.length;
        stats.value.geocaches.active = geocaches.filter(g => g.status === 'active').length;

        const messagesData = await apiRequest('admin/messages');
        const messages = messagesData?.data || (Array.isArray(messagesData) ? messagesData : []);
        stats.value.messages.total = messages.length;
        stats.value.messages.published = messages.filter(m => m.state === 'published').length;

        const languagesData = await apiRequest('admin/languages');
        const languages = languagesData?.data || (Array.isArray(languagesData) ? languagesData : []);
        stats.value.languages.total = languages.length;
        stats.value.languages.active = languages.filter(l => l.active).length;

        const usersData = await apiRequest('admin/users');
        if (usersData && usersData.data) {
            stats.value.users.total = usersData.data.length;
        }

        const gkSettings = await apiRequest('admin/golden-key');
        if (gkSettings) {
            goldenKey.value.isActive = new Date() >= new Date(gkSettings.activation_time);
            goldenKey.value.activationTime = gkSettings.activation_time;
        }

        const gkMonths = await getAdminGoldenKeyMonths();
        goldenKey.value.months = Array.isArray(gkMonths) ? gkMonths : [];
    } catch (err) {
        console.error('Failed to fetch dashboard data:', err);
    }
    loading.value = false;
}

function formatDate(dateString) {
    if (!dateString) return '-';
    return new Date(dateString).toLocaleDateString('nl-BE', { day: 'numeric', month: 'short', year: 'numeric' });
}

function getTimeAgo(dateString) {
    if (!dateString) return '';
    const diffMs = new Date() - new Date(dateString);
    const diffMins = Math.floor(diffMs / 60000);
    const diffHours = Math.floor(diffMs / 3600000);
    const diffDays = Math.floor(diffMs / 86400000);
    if (diffMins < 1) return 'zojuist';
    if (diffMins < 60) return `${diffMins} min geleden`;
    if (diffHours < 24) return `${diffHours} uur geleden`;
    if (diffDays === 1) return 'gisteren';
    if (diffDays < 7) return `${diffDays} dagen geleden`;
    return formatDate(dateString);
}

function getStatusColor(status) {
    return { new: 'danger', in_progress: 'warning', resolved: 'success', closed: 'neutral' }[status] || 'neutral';
}

onMounted(fetchDashboardData);
</script>

<template>
    <AdminLayout pageTitle="Overzicht">
        <div class="dashboard">
            <div class="stats-grid">
                <div class="admin-stat-card clickable" @click="router.push({ name: 'adminContacts' })">
                    <div class="admin-stat-icon admin-stat-icon-warning">
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                            <path d="M22 12h-6l-2 3h-4l-2-3H2"/><path d="M5.45 5.11L2 12v6a2 2 0 0 0 2 2h16a2 2 0 0 0 2-2v-6l-3.45-6.89A2 2 0 0 0 16.76 4H7.24a2 2 0 0 0-1.79 1.11z"/>
                        </svg>
                    </div>
                    <div class="admin-stat-content">
                        <div class="admin-stat-value">{{ stats.contacts.newCount }}</div>
                        <div class="admin-stat-label">Nieuwe berichten</div>
                    </div>
                </div>
                <div class="admin-stat-card clickable" @click="router.push({ name: 'adminEvents' })">
                    <div class="admin-stat-icon admin-stat-icon-info">
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                            <rect x="3" y="4" width="18" height="18" rx="2"/><line x1="16" y1="2" x2="16" y2="6"/><line x1="8" y1="2" x2="8" y2="6"/><line x1="3" y1="10" x2="21" y2="10"/>
                        </svg>
                    </div>
                    <div class="admin-stat-content">
                        <div class="admin-stat-value">{{ stats.events.published }}</div>
                        <div class="admin-stat-label">Gepubliceerde evenementen</div>
                    </div>
                </div>
                <div class="admin-stat-card clickable" @click="router.push({ name: 'adminGeocaches' })">
                    <div class="admin-stat-icon admin-stat-icon-success">
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                            <path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0 1 18 0z"/><circle cx="12" cy="10" r="3"/>
                        </svg>
                    </div>
                    <div class="admin-stat-content">
                        <div class="admin-stat-value">{{ stats.geocaches.total }}</div>
                        <div class="admin-stat-label">Geocaches</div>
                    </div>
                </div>
                <div class="admin-stat-card clickable" @click="router.push({ name: 'adminLanguages' })">
                    <div class="admin-stat-icon admin-stat-icon-primary">
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                            <circle cx="12" cy="12" r="10"/><line x1="2" y1="12" x2="22" y2="12"/><path d="M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1-4 10 15.3 15.3 0 0 1-4-10 15.3 15.3 0 0 1 4-10z"/>
                        </svg>
                    </div>
                    <div class="admin-stat-content">
                        <div class="admin-stat-value">{{ stats.languages.active }}</div>
                        <div class="admin-stat-label">Actieve talen</div>
                    </div>
                </div>
            </div>

            <div class="content-grid">
                <div class="admin-card">
                    <div class="admin-card-header">
                        <h3 class="admin-card-title">Recente contactberichten</h3>
                        <router-link :to="{ name: 'adminContacts' }" class="admin-btn admin-btn-ghost admin-btn-sm">Alles bekijken</router-link>
                    </div>
                    <div class="admin-card-body" style="padding: 0;">
                        <div v-if="loading" class="loading-state"><div class="admin-spinner"></div></div>
                        <div v-else-if="recentContacts.length === 0" class="admin-empty" style="padding: 2rem;">
                            <div class="admin-empty-icon">
                                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                                    <path d="M22 12h-6l-2 3h-4l-2-3H2"/><path d="M5.45 5.11L2 12v6a2 2 0 0 0 2 2h16a2 2 0 0 0 2-2v-6l-3.45-6.89A2 2 0 0 0 16.76 4H7.24a2 2 0 0 0-1.79 1.11z"/>
                                </svg>
                            </div>
                            <p class="admin-empty-title">Geen nieuwe berichten</p>
                            <p class="admin-empty-description">Contactformulier berichten verschijnen hier</p>
                        </div>
                        <div v-else class="contact-list">
                            <div v-for="contact in recentContacts" :key="contact.id" class="contact-item" @click="router.push({ name: 'adminContacts' })">
                                <div class="contact-indicator" :class="`indicator-${getStatusColor(contact.status)}`"></div>
                                <div class="contact-content">
                                    <div class="contact-header">
                                        <span class="contact-email">{{ contact.email }}</span>
                                        <span class="contact-time">{{ getTimeAgo(contact.created_at) }}</span>
                                    </div>
                                    <p class="contact-subject">{{ contact.subject }}</p>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>

                <!-- Golden Key status card -->
                <div class="admin-card">
                    <div class="admin-card-header">
                        <h3 class="admin-card-title">Golden Key</h3>
                        <router-link :to="{ name: 'adminGoldenKey' }" class="admin-btn admin-btn-ghost admin-btn-sm">Beheren</router-link>
                    </div>
                    <div class="admin-card-body">
                        <div v-if="loading" class="loading-state"><div class="admin-spinner"></div></div>
                        <template v-else>
                            <div class="gk-status-row">
                                <span class="gk-status-badge" :class="goldenKey.isActive ? 'gk-active' : 'gk-soon'">
                                    {{ goldenKey.isActive ? 'Actief' : 'Binnenkort' }}
                                </span>
                                <span v-if="goldenKey.activationTime" class="gk-activation-hint">
                                    {{ goldenKey.isActive ? 'Geactiveerd op' : 'Activeert op' }}
                                    {{ new Date(goldenKey.activationTime).toLocaleString('nl-BE', { dateStyle: 'medium', timeStyle: 'short' }) }}
                                </span>
                            </div>
                            <div class="gk-month-stats">
                                <div class="gk-month-stat">
                                    <span class="gk-month-count gk-count-active">{{ gkActiveMonths.length }}</span>
                                    <span class="gk-month-label">Actief</span>
                                </div>
                                <div class="gk-month-stat">
                                    <span class="gk-month-count gk-count-found">{{ gkFoundMonths.length }}</span>
                                    <span class="gk-month-label">Gevonden</span>
                                </div>
                                <div class="gk-month-stat">
                                    <span class="gk-month-count gk-count-locked">{{ gkLockedMonths.length }}</span>
                                    <span class="gk-month-label">Vergrendeld</span>
                                </div>
                            </div>
                            <div v-if="gkActiveMonths.length > 0" class="gk-active-list">
                                <p class="gk-active-list-title">Lopende zoektochten</p>
                                <div v-for="month in gkActiveMonths" :key="month.id" class="gk-active-item">
                                    <span class="gk-active-name">{{ month.month_name }}</span>
                                    <router-link :to="`/admin/golden-key/months/${month.id}`" class="admin-btn admin-btn-ghost admin-btn-sm">Bekijken</router-link>
                                </div>
                            </div>
                            <p v-else-if="!goldenKey.isActive" class="gk-inactive-note">Golden Key is nog niet geactiveerd.</p>
                            <p v-else class="gk-inactive-note">Geen actieve zoektochten.</p>
                        </template>
                    </div>
                </div>
            </div>
        </div>
    </AdminLayout>
</template>

<style scoped>
.dashboard { display: flex; flex-direction: column; gap: 1.5rem; }
.stats-grid { display: grid; grid-template-columns: repeat(4, 1fr); gap: 1rem; }
.clickable { cursor: pointer; transition: transform 0.15s ease, box-shadow 0.15s ease; }
.clickable:hover { transform: translateY(-2px); box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1); }
@media (max-width: 1200px) { .stats-grid { grid-template-columns: repeat(2, 1fr); } }
@media (max-width: 640px) { .stats-grid { grid-template-columns: 1fr; } }
.content-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 1.5rem; }
@media (max-width: 900px) { .content-grid { grid-template-columns: 1fr; } }
.gk-status-row { display: flex; align-items: center; gap: 0.75rem; margin-bottom: 1.25rem; flex-wrap: wrap; }
.gk-status-badge { display: inline-flex; align-items: center; padding: 0.25rem 0.75rem; border-radius: 9999px; font-size: 0.75rem; font-weight: 700; letter-spacing: 0.05em; text-transform: uppercase; }
.gk-active { background: var(--admin-primary-bg); color: var(--admin-primary); }
.gk-soon { background: var(--admin-warning-bg); color: var(--admin-warning); }
.gk-activation-hint { font-size: 0.8rem; color: var(--admin-text-muted); }
.gk-month-stats { display: flex; gap: 1.5rem; margin-bottom: 1.25rem; }
.gk-month-stat { display: flex; flex-direction: column; align-items: center; gap: 0.2rem; }
.gk-month-count { font-size: 1.75rem; font-weight: 700; line-height: 1; }
.gk-count-active { color: var(--admin-primary); }
.gk-count-found { color: var(--admin-success); }
.gk-count-locked { color: var(--admin-text-muted); }
.gk-month-label { font-size: 0.725rem; color: var(--admin-text-muted); text-transform: uppercase; letter-spacing: 0.04em; }
.gk-active-list { border-top: 1px solid var(--admin-border-light); padding-top: 1rem; display: flex; flex-direction: column; gap: 0.5rem; }
.gk-active-list-title { margin: 0 0 0.5rem; font-size: 0.8rem; font-weight: 600; color: var(--admin-text-secondary); text-transform: uppercase; letter-spacing: 0.04em; }
.gk-active-item { display: flex; align-items: center; justify-content: space-between; gap: 0.5rem; }
.gk-active-name { font-size: 0.875rem; font-weight: 500; color: var(--admin-text); }
.gk-inactive-note { font-size: 0.875rem; color: var(--admin-text-muted); font-style: italic; margin: 0; }
.loading-state { display: flex; align-items: center; justify-content: center; padding: 3rem; }
.contact-list { display: flex; flex-direction: column; }
.contact-item { display: flex; align-items: flex-start; gap: 0.875rem; padding: 1rem 1.25rem; cursor: pointer; transition: background 0.15s ease; border-bottom: 1px solid var(--admin-border-light); }
.contact-item:last-child { border-bottom: none; }
.contact-item:hover { background: var(--admin-surface-hover); }
.contact-indicator { width: 0.5rem; height: 0.5rem; border-radius: 50%; margin-top: 0.375rem; flex-shrink: 0; }
.indicator-danger { background: var(--admin-danger); }
.indicator-warning { background: var(--admin-warning); }
.indicator-success { background: var(--admin-success); }
.indicator-neutral { background: var(--admin-text-muted); }
.contact-content { flex: 1; min-width: 0; }
.contact-header { display: flex; align-items: center; justify-content: space-between; gap: 0.5rem; margin-bottom: 0.25rem; }
.contact-email { font-size: 0.875rem; font-weight: 500; color: var(--admin-text); white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.contact-time { font-size: 0.75rem; color: var(--admin-text-muted); white-space: nowrap; }
.contact-subject { font-size: 0.8125rem; color: var(--admin-text-secondary); margin: 0; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }

</style>
