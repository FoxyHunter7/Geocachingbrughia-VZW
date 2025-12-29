<script setup>
import { ref, computed, onMounted } from 'vue';
import AdminLayout from '@/components/admin/AdminLayout.vue';
import config from '@/data/config.js';

// State
const loading = ref(true);
const languages = ref([]);
const search = ref('');

// Modal state
const showModal = ref(false);
const modalMode = ref('create');
const saving = ref(false);
const editingLanguage = ref(null);

// Form data
const formData = ref({
    code: '',
    name: '',
    flag_url: '',
    is_active: true
});

// API helpers
function getToken() {
    return localStorage.getItem('admin_token');
}

async function apiRequest(endpoint, options = {}) {
    const token = getToken();
    const headers = {
        'Accept': 'application/json',
        ...(!(options.body instanceof FormData) && { 'Content-Type': 'application/json' }),
        ...(token && { 'Authorization': `Bearer ${token}` }),
        ...options.headers
    };

    try {
        const response = await fetch(`${config.apiUrl}${endpoint}`, {
            ...options,
            headers
        });
        return response;
    } catch (err) {
        console.error('API request failed:', err);
        return null;
    }
}

// Fetch data
async function fetchLanguages() {
    loading.value = true;
    try {
        const response = await apiRequest('admin/languages');
        if (response?.ok) {
            const data = await response.json();
            languages.value = data.data || data || [];
        }
    } catch (err) {
        console.error('Failed to fetch languages:', err);
    }
    loading.value = false;
}

// Computed
const filteredLanguages = computed(() => {
    if (!search.value) return languages.value;
    const q = search.value.toLowerCase();
    return languages.value.filter(l => 
        l.code?.toLowerCase().includes(q) ||
        l.name?.toLowerCase().includes(q)
    );
});

// Modal handlers
function openCreateModal() {
    modalMode.value = 'create';
    editingLanguage.value = null;
    
    formData.value = {
        code: '',
        name: '',
        flag_url: '',
        is_active: true
    };
    
    showModal.value = true;
}

function openEditModal(language) {
    modalMode.value = 'edit';
    editingLanguage.value = language;
    
    formData.value = {
        code: language.code || '',
        name: language.name || '',
        flag_url: language.flag_url || '',
        is_active: language.active ?? true
    };
    
    showModal.value = true;
}

function closeModal() {
    showModal.value = false;
    editingLanguage.value = null;
}

async function handleSave() {
    if (!formData.value.code || !formData.value.name) {
        window.$toast?.error('Code en naam zijn verplicht');
        return;
    }
    
    saving.value = true;
    
    try {
        const data = {
            code: formData.value.code.toUpperCase(),
            name: formData.value.name,
            flag_url: formData.value.flag_url || '',
            active: formData.value.is_active
        };
        
        const endpoint = modalMode.value === 'create' 
            ? 'admin/languages'
            : `admin/languages/${editingLanguage.value.code}`;
        
        const method = modalMode.value === 'create' ? 'POST' : 'PUT';
        
        const response = await apiRequest(endpoint, {
            method,
            body: JSON.stringify(data)
        });
        
        if (response?.ok) {
            window.$toast?.success(modalMode.value === 'create' ? 'Taal aangemaakt!' : 'Taal bijgewerkt!');
            closeModal();
            fetchLanguages();
        } else {
            const err = await response?.json();
            window.$toast?.error(err?.message || 'Opslaan mislukt');
        }
    } catch (err) {
        console.error('Save failed:', err);
        window.$toast?.error('Er is een fout opgetreden bij het opslaan');
    }
    
    saving.value = false;
}

async function handleDelete() {
    if (!editingLanguage.value) return;
    
    if (!confirm(`Weet je zeker dat je "${editingLanguage.value.name}" wilt verwijderen? Dit verwijdert ook alle vertalingen voor deze taal.`)) return;
    
    saving.value = true;
    try {
        const response = await apiRequest(`admin/languages/${editingLanguage.value.code}`, {
            method: 'DELETE'
        });
        
        if (response?.ok) {
            window.$toast?.success('Taal verwijderd');
            closeModal();
            fetchLanguages();
        } else {
            window.$toast?.error('Verwijderen mislukt');
        }
    } catch (err) {
        window.$toast?.error('Er is een fout opgetreden');
    }
    saving.value = false;
}

async function toggleActive(language) {
    try {
        const response = await apiRequest(`admin/languages/${language.code}`, {
            method: 'PUT',
            body: JSON.stringify({
                ...language,
                active: !language.active
            })
        });
        
        if (response?.ok) {
            window.$toast?.success(`Taal ${language.active ? 'uitgeschakeld' : 'ingeschakeld'}`);
            fetchLanguages();
        }
    } catch (err) {
        window.$toast?.error('Bijwerken mislukt');
    }
}

// Initialize
onMounted(fetchLanguages);
</script>

<template>
    <AdminLayout pageTitle="Talen">
        <template #actions>
            <button class="admin-btn admin-btn-primary" @click="openCreateModal">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="width: 1rem; height: 1rem;">
                    <line x1="12" y1="5" x2="12" y2="19"/>
                    <line x1="5" y1="12" x2="19" y2="12"/>
                </svg>
                Nieuwe Taal
            </button>
        </template>

        <!-- Info Card -->
        <div class="admin-alert admin-alert-info" style="margin-bottom: 1.5rem;">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="width: 1.25rem; height: 1.25rem; flex-shrink: 0;">
                <circle cx="12" cy="12" r="10"/>
                <line x1="12" y1="16" x2="12" y2="12"/>
                <line x1="12" y1="8" x2="12.01" y2="8"/>
            </svg>
            <div>
                <strong>Beheer beschikbare talen</strong>
                <p style="margin: 0.25rem 0 0; font-size: 0.875rem; opacity: 0.9;">
                    Talen bepalen welke vertalingen beschikbaar zijn voor inhoud. Een taal uitschakelen verbergt deze voor bezoekers maar behoudt de vertalingen.
                </p>
            </div>
        </div>

        <!-- Filters -->
        <div class="admin-card" style="margin-bottom: 1.5rem;">
            <div class="admin-card-body">
                <div class="admin-filters">
                    <div class="admin-search" style="flex: 1; max-width: 20rem;">
                        <span class="admin-search-icon">
                            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="width: 1rem; height: 1rem;">
                                <circle cx="11" cy="11" r="8"/>
                                <line x1="21" y1="21" x2="16.65" y2="16.65"/>
                            </svg>
                        </span>
                        <input v-model="search" type="search" class="admin-input" placeholder="Zoeken...">
                    </div>
                </div>
            </div>
        </div>

        <!-- Languages Grid -->
        <div class="languages-grid">
            <div v-if="loading" class="loading-state">
                <div class="admin-spinner"></div>
            </div>
            <template v-else-if="filteredLanguages.length === 0">
                <div class="admin-card" style="grid-column: 1 / -1;">
                    <div class="admin-empty">
                        <div class="admin-empty-icon">
                            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                                <path d="M3 5h12M9 3v2m1.048 9.5A18.022 18.022 0 016.412 9m6.088 9h7M11 21l5-10 5 10M12.751 5C11.783 10.77 8.07 15.61 3 18.129"/>
                            </svg>
                        </div>
                        <p class="admin-empty-title">Geen talen gevonden</p>
                        <p class="admin-empty-description">Voeg je eerste taal toe om vertalingen in te schakelen</p>
                        <button class="admin-btn admin-btn-primary" @click="openCreateModal">
                            Nieuwe Taal
                        </button>
                    </div>
                </div>
            </template>
            <div 
                v-else 
                v-for="language in filteredLanguages" 
                :key="language.code" 
                class="language-card admin-card"
                :class="{ 'is-inactive': !language.active }"
            >
                <div class="language-card-body">
                    <div class="language-flag">
                        <img v-if="language.flag_url" :src="language.flag_url" :alt="language.code">
                        <span v-else class="flag-placeholder">{{ language.code }}</span>
                    </div>
                    <div class="language-info">
                        <h3 class="language-name">{{ language.name }}</h3>
                        <span class="language-code">{{ language.code }}</span>
                    </div>
                    <div class="language-status">
                        <span :class="['admin-badge', language.active ? 'admin-badge-success' : 'admin-badge-neutral']">
                            {{ language.active ? 'Actief' : 'Inactief' }}
                        </span>
                    </div>
                </div>
                <div class="language-card-actions">
                    <button class="admin-btn admin-btn-ghost admin-btn-sm" @click="toggleActive(language)">
                        {{ language.active ? 'Uitschakelen' : 'Inschakelen' }}
                    </button>
                    <button class="admin-btn admin-btn-ghost admin-btn-sm" @click="openEditModal(language)">
                        Bewerken
                    </button>
                </div>
            </div>
        </div>

        <!-- Modal -->
        <Teleport to="body">
            <div v-if="showModal" class="admin-modal-overlay" @click.self="closeModal">
                <div class="admin-modal admin-modal-sm">
                    <div class="admin-modal-header">
                        <h2 class="admin-modal-title">
                            {{ modalMode === 'create' ? 'Nieuwe Taal' : 'Taal Bewerken' }}
                        </h2>
                        <button class="admin-modal-close" @click="closeModal">
                            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                                <line x1="18" y1="6" x2="6" y2="18"/>
                                <line x1="6" y1="6" x2="18" y2="18"/>
                            </svg>
                        </button>
                    </div>
                    <div class="admin-modal-body">
                        <div class="admin-form-group">
                            <label class="admin-label">Taalcode *</label>
                            <input 
                                v-model="formData.code" 
                                type="text" 
                                class="admin-input" 
                                placeholder="EN, NL, FR..."
                                maxlength="5"
                                :disabled="modalMode === 'edit'"
                                required
                            >
                            <span class="admin-form-hint">Korte code (2-5 karakters, bijv. EN, NL, FR)</span>
                        </div>

                        <div class="admin-form-group">
                            <label class="admin-label">Taalnaam *</label>
                            <input v-model="formData.name" type="text" class="admin-input" placeholder="English, Nederlands..." required>
                        </div>

                        <div class="admin-form-group">
                            <label class="admin-label">Vlag URL</label>
                            <input v-model="formData.flag_url" type="url" class="admin-input" placeholder="https://...">
                            <span class="admin-form-hint">URL naar een vlagafbeelding (optioneel)</span>
                        </div>

                        <div class="admin-form-group">
                            <label class="admin-checkbox">
                                <input type="checkbox" v-model="formData.is_active">
                                <span>Actief</span>
                            </label>
                        </div>
                    </div>
                    <div class="admin-modal-footer">
                        <button v-if="modalMode === 'edit'" class="admin-btn admin-btn-danger" @click="handleDelete" :disabled="saving">
                            Verwijderen
                        </button>
                        <div style="flex: 1;"></div>
                        <button class="admin-btn admin-btn-secondary" @click="closeModal" :disabled="saving">
                            Annuleren
                        </button>
                        <button class="admin-btn admin-btn-primary" @click="handleSave" :disabled="saving">
                            {{ saving ? 'Opslaan...' : 'Opslaan' }}
                        </button>
                    </div>
                </div>
            </div>
        </Teleport>
    </AdminLayout>
</template>

<style scoped>
.languages-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
    gap: 1rem;
}

.loading-state {
    grid-column: 1 / -1;
    display: flex;
    justify-content: center;
    padding: 3rem;
}

.language-card {
    padding: 0;
    overflow: hidden;
    transition: transform 0.2s, box-shadow 0.2s;
}

.language-card:hover {
    transform: translateY(-2px);
    box-shadow: 0 8px 16px rgba(0, 0, 0, 0.1);
}

.language-card.is-inactive {
    opacity: 0.6;
}

.language-card-body {
    display: flex;
    align-items: center;
    gap: 1rem;
    padding: 1.25rem;
}

.language-flag {
    width: 3rem;
    height: 2rem;
    border-radius: var(--admin-radius-sm);
    overflow: hidden;
    flex-shrink: 0;
}

.language-flag img {
    width: 100%;
    height: 100%;
    object-fit: cover;
}

.flag-placeholder {
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
    background: var(--admin-bg);
    font-weight: 600;
    font-size: 0.75rem;
    color: var(--admin-text-secondary);
}

.language-info {
    flex: 1;
    min-width: 0;
}

.language-name {
    font-size: 1rem;
    font-weight: 600;
    color: var(--admin-text);
    margin: 0;
}

.language-code {
    font-size: 0.75rem;
    color: var(--admin-text-muted);
    text-transform: uppercase;
}

.language-card-actions {
    display: flex;
    gap: 0.5rem;
    padding: 0.75rem 1.25rem;
    border-top: 1px solid var(--admin-border-light);
    background: var(--admin-bg);
}
</style>