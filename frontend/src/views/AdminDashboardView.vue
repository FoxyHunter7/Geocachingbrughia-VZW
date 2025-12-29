<script setup>
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import AdminLayout from '@/components/admin/AdminLayout.vue';
import config from '@/data/config.js';

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
            </div>

            <div class="summary-grid">
                <div class="summary-card">
                    <h4>Evenementen</h4>
                    <div class="summary-stats">
                        <div class="summary-stat"><span class="stat-value">{{ stats.events.published }}</span><span class="stat-label">Gepubliceerd</span></div>
                        <div class="summary-stat"><span class="stat-value">{{ stats.events.draft }}</span><span class="stat-label">Concept</span></div>
                    </div>
                </div>
                <div class="summary-card">
                    <h4>Contactberichten</h4>
                    <div class="summary-stats">
                        <div class="summary-stat"><span class="stat-value text-danger">{{ stats.contacts.newCount }}</span><span class="stat-label">Nieuw</span></div>
                        <div class="summary-stat"><span class="stat-value text-warning">{{ stats.contacts.inProgress }}</span><span class="stat-label">In behandeling</span></div>
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
.content-grid { display: grid; grid-template-columns: 1fr; gap: 1.5rem; max-width: 800px; }
.summary-grid { display: grid; grid-template-columns: repeat(2, 1fr); gap: 1rem; }
@media (max-width: 640px) { .summary-grid { grid-template-columns: 1fr; } }
.summary-card { background: var(--admin-surface); border: 1px solid var(--admin-border); border-radius: var(--admin-radius); padding: 1.25rem; }
.summary-card h4 { margin: 0 0 1rem 0; font-size: 0.875rem; font-weight: 600; color: var(--admin-text-secondary); text-transform: uppercase; }
.summary-stats { display: flex; gap: 2rem; }
.summary-stat { display: flex; flex-direction: column; gap: 0.25rem; }
.stat-value { font-size: 1.5rem; font-weight: 700; color: var(--admin-text); }
.stat-value.text-danger { color: var(--admin-danger); }
.stat-value.text-warning { color: var(--admin-warning); }
.stat-label { font-size: 0.75rem; color: var(--admin-text-muted); }
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
