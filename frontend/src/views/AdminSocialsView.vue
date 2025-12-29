<script setup>
import { ref, computed, onMounted } from 'vue';
import AdminLayout from '@/components/admin/AdminLayout.vue';
import config from '@/data/config.js';

// State
const loading = ref(true);
const socials = ref([]);
const search = ref('');

// Modal state
const showModal = ref(false);
const modalMode = ref('create');
const saving = ref(false);
const editingSocial = ref(null);

// File upload
const fileInput = ref(null);
const imagePreview = ref('');

// Form data
const formData = ref({
    platform: '',
    url: '',
    icon: '',
    active: true,
    sort_order: 0
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
async function fetchSocials() {
    loading.value = true;
    try {
        const response = await apiRequest('admin/socials');
        if (response?.ok) {
            const data = await response.json();
            socials.value = data.data || data || [];
        }
    } catch (err) {
        console.error('Failed to fetch socials:', err);
    }
    loading.value = false;
}

// Computed
const filteredSocials = computed(() => {
    if (!search.value) return socials.value;
    const q = search.value.toLowerCase();
    return socials.value.filter(s => 
        s.platform?.toLowerCase().includes(q)
    );
});

// Modal handlers
function openCreateModal() {
    modalMode.value = 'create';
    editingSocial.value = null;
    imagePreview.value = '';
    
    formData.value = {
        platform: '',
        url: '',
        icon: '',
        active: true,
        sort_order: 0
    };
    
    showModal.value = true;
}

function openEditModal(social) {
    modalMode.value = 'edit';
    editingSocial.value = social;
    
    formData.value = {
        platform: social.platform || '',
        url: social.url || '',
        icon: social.icon || '',
        active: social.active ?? true,
        sort_order: social.sort_order || 0
    };
    
    imagePreview.value = social.icon ? `${config.apiUrl}images/${social.icon}` : '';
    showModal.value = true;
}

function closeModal() {
    showModal.value = false;
    editingSocial.value = null;
    if (fileInput.value) fileInput.value.value = '';
}

function handleImageChange(e) {
    const file = e.target.files?.[0];
    if (file) {
        const reader = new FileReader();
        reader.onload = (e) => {
            imagePreview.value = e.target.result;
        };
        reader.readAsDataURL(file);
    }
}

async function handleSave() {
    if (!formData.value.platform || !formData.value.url) {
        window.$toast?.error('Platform name and URL are required');
        return;
    }
    
    saving.value = true;
    
    try {
        // Build JSON payload (backend expects JSON, not FormData)
        const payload = {
            platform: formData.value.platform,
            url: formData.value.url,
            icon: formData.value.icon || '',
            active: formData.value.active,
            sort_order: formData.value.sort_order || 0
        };
        
        const endpoint = modalMode.value === 'create' 
            ? 'admin/socials'
            : `admin/socials/${editingSocial.value.id}`;
        
        const method = modalMode.value === 'create' ? 'POST' : 'PUT';
        
        const response = await apiRequest(endpoint, {
            method,
            body: JSON.stringify(payload)
        });
        
        if (response?.ok) {
            window.$toast?.success(modalMode.value === 'create' ? 'Sociale link aangemaakt!' : 'Sociale link bijgewerkt!');
            closeModal();
            fetchSocials();
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
    if (!editingSocial.value) return;
    
    if (!confirm(`Weet je zeker dat je "${editingSocial.value.platform}" wilt verwijderen?`)) return;
    
    saving.value = true;
    try {
        const response = await apiRequest(`admin/socials/${editingSocial.value.id}`, {
            method: 'DELETE'
        });
        
        if (response?.ok) {
            window.$toast?.success('Sociale link verwijderd');
            closeModal();
            fetchSocials();
        } else {
            window.$toast?.error('Verwijderen mislukt');
        }
    } catch (err) {
        window.$toast?.error('Er is een fout opgetreden');
    }
    saving.value = false;
}

// Initialize
onMounted(fetchSocials);
</script>

<template>
    <AdminLayout pageTitle="Sociale Links">
        <template #actions>
            <button class="admin-btn admin-btn-primary" @click="openCreateModal">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="width: 1rem; height: 1rem;">
                    <line x1="12" y1="5" x2="12" y2="19"/>
                    <line x1="5" y1="12" x2="19" y2="12"/>
                </svg>
                Nieuwe Link
            </button>
        </template>

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

        <!-- Socials Grid -->
        <div class="socials-grid">
            <div v-if="loading" class="loading-state">
                <div class="admin-spinner"></div>
            </div>
            <template v-else-if="filteredSocials.length === 0">
                <div class="admin-card" style="grid-column: 1 / -1;">
                    <div class="admin-empty">
                        <div class="admin-empty-icon">
                            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                                <path d="M18 2h-3a5 5 0 0 0-5 5v3H7v4h3v8h4v-8h3l1-4h-4V7a1 1 0 0 1 1-1h3z"/>
                            </svg>
                        </div>
                        <p class="admin-empty-title">Geen sociale links gevonden</p>
                        <p class="admin-empty-description">Voeg je sociale media profielen toe</p>
                        <button class="admin-btn admin-btn-primary" @click="openCreateModal">
                            Nieuwe Link
                        </button>
                    </div>
                </div>
            </template>
            <div 
                v-else 
                v-for="social in filteredSocials" 
                :key="social.id" 
                class="social-card admin-card"
            >
                <div class="social-card-body">
                    <div class="social-icon">
                        <img v-if="social.icon" :src="`${config.apiUrl}images/${social.icon}`" :alt="social.platform">
                        <span v-else class="icon-placeholder">
                            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                                <path d="M18 2h-3a5 5 0 0 0-5 5v3H7v4h3v8h4v-8h3l1-4h-4V7a1 1 0 0 1 1-1h3z"/>
                            </svg>
                        </span>
                    </div>
                    <div class="social-info">
                        <h3 class="social-name">{{ social.platform }}</h3>
                        <a :href="social.url" target="_blank" class="social-url">{{ social.url }}</a>
                    </div>
                </div>
                <div class="social-card-actions">
                    <a :href="social.url" target="_blank" class="admin-btn admin-btn-ghost admin-btn-sm">
                        Bezoeken
                    </a>
                    <button class="admin-btn admin-btn-ghost admin-btn-sm" @click="openEditModal(social)">
                        Bewerken
                    </button>
                </div>
            </div>
        </div>

        <!-- Modal -->
        <Teleport to="body">
            <div v-if="showModal" class="admin-modal-overlay" @click.self="closeModal">
                <div class="admin-modal admin-modal-md">
                    <div class="admin-modal-header">
                        <h2 class="admin-modal-title">
                            {{ modalMode === 'create' ? 'Nieuwe Sociale Link' : 'Sociale Link Bewerken' }}
                        </h2>
                        <button class="admin-modal-close" @click="closeModal">
                            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                                <line x1="18" y1="6" x2="6" y2="18"/>
                                <line x1="6" y1="6" x2="18" y2="18"/>
                            </svg>
                        </button>
                    </div>
                    <div class="admin-modal-body">
                        <div class="modal-grid">
                            <div class="form-section">
                                <div class="admin-form-group">
                                    <label class="admin-label">Platform Naam *</label>
                                    <input v-model="formData.platform" type="text" class="admin-input" placeholder="Facebook, Instagram..." required>
                                </div>

                                <div class="admin-form-group">
                                    <label class="admin-label">URL *</label>
                                    <input v-model="formData.url" type="url" class="admin-input" placeholder="https://..." required>
                                    <span class="admin-form-hint">Volledige URL naar je sociale media profiel</span>
                                </div>

                                <div class="admin-form-group">
                                    <label class="admin-checkbox">
                                        <input type="checkbox" v-model="formData.active">
                                        <span>Actief</span>
                                    </label>
                                    <span class="admin-form-hint">Alleen actieve links worden getoond op de website</span>
                                </div>
                            </div>

                            <div class="image-section">
                                <div class="admin-form-group">
                                    <label class="admin-label">Icoon URL</label>
                                    <input v-model="formData.icon" type="text" class="admin-input" placeholder="icon-name.png">
                                    <span class="admin-form-hint">Icoon bestandsnaam (opgeslagen in images folder)</span>
                                </div>
                            </div>
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
.socials-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
    gap: 1rem;
}

.loading-state {
    grid-column: 1 / -1;
    display: flex;
    justify-content: center;
    padding: 3rem;
}

.social-card {
    padding: 0;
    overflow: hidden;
    transition: transform 0.2s, box-shadow 0.2s;
}

.social-card:hover {
    transform: translateY(-2px);
    box-shadow: 0 8px 16px rgba(0, 0, 0, 0.1);
}

.social-card-body {
    display: flex;
    align-items: center;
    gap: 1rem;
    padding: 1.25rem;
}

.social-icon {
    width: 3rem;
    height: 3rem;
    border-radius: var(--admin-radius);
    overflow: hidden;
    flex-shrink: 0;
    background: var(--admin-bg);
    display: flex;
    align-items: center;
    justify-content: center;
}

.social-icon img {
    width: 100%;
    height: 100%;
    object-fit: contain;
}

.icon-placeholder {
    color: var(--admin-text-muted);
}

.icon-placeholder svg {
    width: 1.5rem;
    height: 1.5rem;
}

.social-info {
    flex: 1;
    min-width: 0;
}

.social-name {
    font-size: 1rem;
    font-weight: 600;
    color: var(--admin-text);
    margin: 0;
}

.social-url {
    font-size: 0.75rem;
    color: var(--admin-text-muted);
    text-decoration: none;
    display: block;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
}

.social-url:hover {
    color: var(--admin-primary);
}

.social-card-actions {
    display: flex;
    gap: 0.5rem;
    padding: 0.75rem 1.25rem;
    border-top: 1px solid var(--admin-border-light);
    background: var(--admin-bg);
}

.modal-grid {
    display: grid;
    grid-template-columns: 1fr auto;
    gap: 2rem;
}

@media (max-width: 600px) {
    .modal-grid {
        grid-template-columns: 1fr;
    }
}

.form-section {
    display: flex;
    flex-direction: column;
}

.image-section {
    width: 10rem;
}

.image-upload-area {
    border: 2px dashed var(--admin-border);
    border-radius: var(--admin-radius-lg);
    overflow: hidden;
    background: var(--admin-bg);
}

.image-preview {
    width: 100%;
    height: 8rem;
    object-fit: contain;
    background: var(--admin-surface);
}

.image-placeholder {
    height: 8rem;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 0.5rem;
    color: var(--admin-text-muted);
    font-size: 0.75rem;
}

.image-placeholder svg {
    width: 2rem;
    height: 2rem;
}
</style>