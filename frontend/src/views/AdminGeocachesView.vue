<script setup>
import { ref, computed, onMounted, watch } from 'vue';
import AdminLayout from '@/components/admin/AdminLayout.vue';
import config from '@/data/config.js';

// State
const loading = ref(true);
const geocaches = ref([]);
const currPage = ref(1);
const lastPage = ref(1);
const totalItems = ref(0);
const search = ref('');
const statusFilter = ref('');

// Modal state
const showModal = ref(false);
const modalMode = ref('create');
const saving = ref(false);
const editingGeocache = ref(null);

// Form data - matches backend Geocache struct
const formData = ref({
    geolink: '',
    name: '',
    difficulty: 1,
    terrain: 1,
    size: '',
    type: 'traditional',
    placed_date: '',
    status: 'active'
});

// Size options
const sizeOptions = [
    { value: 'micro', label: 'Micro' },
    { value: 'small', label: 'Klein' },
    { value: 'regular', label: 'Normaal' },
    { value: 'large', label: 'Groot' },
    { value: 'other', label: 'Anders' },
    { value: 'virtual', label: 'Virtueel' },
    { value: 'not_chosen', label: 'Niet gekozen' }
];

// Type options
const typeOptions = [
    { value: 'traditional', label: 'Traditional' },
    { value: 'multi', label: 'Multi-cache' },
    { value: 'mystery', label: 'Mystery/Puzzle' },
    { value: 'letterbox', label: 'Letterbox Hybrid' },
    { value: 'wherigo', label: 'Wherigo' },
    { value: 'earthcache', label: 'EarthCache' },
    { value: 'virtual', label: 'Virtual' },
    { value: 'event', label: 'Event' },
    { value: 'cito', label: 'CITO' },
    { value: 'mega', label: 'Mega-Event' },
    { value: 'giga', label: 'Giga-Event' },
    { value: 'lab', label: 'Lab Cache' }
];

// API helpers
function getToken() {
    return localStorage.getItem('admin_token');
}

async function apiRequest(endpoint, options = {}) {
    const token = getToken();
    const headers = {
        'Accept': 'application/json',
        'Content-Type': 'application/json',
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
async function fetchGeocaches() {
    loading.value = true;
    try {
        const params = new URLSearchParams({ page: currPage.value });
        if (statusFilter.value) params.append('status', statusFilter.value);
        
        const response = await apiRequest(`admin/geocaches?${params}`);
        if (response?.ok) {
            const data = await response.json();
            if (data.data) {
                geocaches.value = data.data;
                currPage.value = data.current_page || 1;
                lastPage.value = data.last_page || 1;
                totalItems.value = data.total || data.data.length;
            } else if (Array.isArray(data)) {
                geocaches.value = data;
                totalItems.value = data.length;
            }
        }
    } catch (err) {
        console.error('Failed to fetch geocaches:', err);
    }
    loading.value = false;
}

// Computed
const filteredGeocaches = computed(() => {
    if (!search.value) return geocaches.value;
    const q = search.value.toLowerCase();
    return geocaches.value.filter(g => 
        g.name?.toLowerCase().includes(q) ||
        g.geolink?.toLowerCase().includes(q)
    );
});

// Modal handlers
function openCreateModal() {
    modalMode.value = 'create';
    editingGeocache.value = null;
    
    formData.value = {
        geolink: '',
        name: '',
        difficulty: 1,
        terrain: 1,
        size: 'regular',
        type: 'traditional',
        placed_date: '',
        status: 'active'
    };
    
    showModal.value = true;
}

function openEditModal(geocache) {
    modalMode.value = 'edit';
    editingGeocache.value = geocache;
    
    formData.value = {
        geolink: geocache.geolink || '',
        name: geocache.name || '',
        difficulty: geocache.difficulty || 1,
        terrain: geocache.terrain || 1,
        size: geocache.size || 'regular',
        type: geocache.type || 'traditional',
        placed_date: geocache.placed_date || '',
        status: geocache.status || 'active'
    };
    
    showModal.value = true;
}

function closeModal() {
    showModal.value = false;
    editingGeocache.value = null;
}

async function handleSave() {
    if (!formData.value.geolink || !formData.value.name) {
        window.$toast?.error('Link en Naam zijn verplicht');
        return;
    }
    
    saving.value = true;
    
    try {
        const payload = {
            geolink: formData.value.geolink,
            name: formData.value.name,
            difficulty: parseFloat(formData.value.difficulty) || 1,
            terrain: parseFloat(formData.value.terrain) || 1,
            size: formData.value.size || '',
            type: formData.value.type || 'traditional',
            placed_date: formData.value.placed_date || '',
            status: formData.value.status || 'active'
        };
        
        const endpoint = modalMode.value === 'create' 
            ? 'admin/geocaches'
            : `admin/geocaches/${editingGeocache.value.id}`;
        
        const method = modalMode.value === 'create' ? 'POST' : 'PUT';
        
        const response = await apiRequest(endpoint, {
            method,
            body: JSON.stringify(payload)
        });
        
        if (response?.ok) {
            window.$toast?.success(modalMode.value === 'create' ? 'Geocache aangemaakt!' : 'Geocache bijgewerkt!');
            closeModal();
            fetchGeocaches();
        } else {
            const err = await response?.json();
            window.$toast?.error(err?.message || err?.error || 'Opslaan mislukt');
        }
    } catch (err) {
        console.error('Save failed:', err);
        window.$toast?.error('Er is een fout opgetreden bij het opslaan');
    }
    
    saving.value = false;
}

async function handleDelete() {
    if (!editingGeocache.value) return;
    
    if (!confirm(`Weet je zeker dat je "${editingGeocache.value.name}" wilt verwijderen?`)) return;
    
    saving.value = true;
    try {
        const response = await apiRequest(`admin/geocaches/${editingGeocache.value.id}`, {
            method: 'DELETE'
        });
        
        if (response?.ok) {
            window.$toast?.success('Geocache verwijderd');
            closeModal();
            fetchGeocaches();
        } else {
            window.$toast?.error('Verwijderen mislukt');
        }
    } catch (err) {
        window.$toast?.error('Er is een fout opgetreden');
    }
    saving.value = false;
}

// Formatting
function getStatusBadge(status) {
    const badges = {
        active: 'success',
        disabled: 'warning',
        archived: 'neutral'
    };
    return badges[status] || 'neutral';
}

// Initialize
onMounted(fetchGeocaches);

watch(currPage, fetchGeocaches);
watch(statusFilter, () => {
    currPage.value = 1;
    fetchGeocaches();
});
</script>

<template>
    <AdminLayout pageTitle="Geocaches">
        <template #actions>
            <button class="admin-btn admin-btn-primary" @click="openCreateModal">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="width: 1rem; height: 1rem;">
                    <line x1="12" y1="5" x2="12" y2="19"/>
                    <line x1="5" y1="12" x2="19" y2="12"/>
                </svg>
                Geocache Toevoegen
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
                        <input v-model="search" type="search" class="admin-input" placeholder="Zoeken op naam of GC code...">
                    </div>
                    <div class="admin-filter-group">
                        <label class="admin-filter-label">Status:</label>
                        <select v-model="statusFilter" class="admin-select" style="width: auto;">
                            <option value="">Alles</option>
                            <option value="active">Actief</option>
                            <option value="disabled">Uitgeschakeld</option>
                            <option value="archived">Gearchiveerd</option>
                        </select>
                    </div>
                </div>
            </div>
        </div>

        <!-- Table -->
        <div class="admin-card">
            <div class="admin-table-wrapper">
                <table class="admin-table">
                    <thead>
                        <tr>
                            <th>Link</th>
                            <th>Naam</th>
                            <th style="width: 5rem;">D/T</th>
                            <th style="width: 6rem;">Grootte</th>
                            <th style="width: 6rem;">Status</th>
                            <th style="width: 5rem;"></th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-if="loading">
                            <td colspan="6" style="text-align: center; padding: 3rem;">
                                <div class="admin-spinner" style="margin: 0 auto;"></div>
                            </td>
                        </tr>
                        <tr v-else-if="filteredGeocaches.length === 0">
                            <td colspan="6">
                                <div class="admin-empty">
                                    <div class="admin-empty-icon">
                                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                                            <path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0 1 18 0z"/>
                                            <circle cx="12" cy="10" r="3"/>
                                        </svg>
                                    </div>
                                    <p class="admin-empty-title">Geen geocaches gevonden</p>
                                    <p class="admin-empty-description">Maak uw eerste geocache aan om te beginnen</p>
                                    <button class="admin-btn admin-btn-primary" @click="openCreateModal">
                                        Geocache Toevoegen
                                    </button>
                                </div>
                            </td>
                        </tr>
                        <tr v-else v-for="geocache in filteredGeocaches" :key="geocache.id">
                            <td>
                                <a v-if="geocache.geolink" :href="geocache.geolink" target="_blank" style="color: var(--admin-primary); text-decoration: none;">
                                    {{ geocache.geolink }}
                                </a>
                                <span v-else>-</span>
                            </td>
                            <td>
                                <span style="font-weight: 500;">{{ geocache.name }}</span>
                            </td>
                            <td>
                                <span class="dt-rating">{{ geocache.difficulty || '-' }} / {{ geocache.terrain || '-' }}</span>
                            </td>
                            <td>{{ geocache.size || '-' }}</td>
                            <td>
                                <span :class="['admin-badge', `admin-badge-${getStatusBadge(geocache.status)}`]">
                                    {{ geocache.status }}
                                </span>
                            </td>
                            <td>
                                <div class="admin-table-actions">
                                    <button class="admin-btn admin-btn-ghost admin-btn-icon admin-btn-sm" @click="openEditModal(geocache)" title="Edit">
                                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="width: 1rem; height: 1rem;">
                                            <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
                                            <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
                                        </svg>
                                    </button>
                                </div>
                            </td>
                        </tr>
                    </tbody>
                </table>
            </div>

            <!-- Pagination -->
            <div v-if="lastPage > 1" class="admin-card-footer">
                <div class="admin-pagination">
                    <button 
                        class="admin-pagination-btn" 
                        :disabled="currPage === 1"
                        @click="currPage--"
                    >
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="width: 1rem; height: 1rem;">
                            <polyline points="15 18 9 12 15 6"/>
                        </svg>
                    </button>
                    <span style="font-size: 0.875rem; color: var(--admin-text-secondary); padding: 0 0.5rem;">
                        Pagina {{ currPage }} van {{ lastPage }}
                    </span>
                    <button 
                        class="admin-pagination-btn" 
                        :disabled="currPage === lastPage"
                        @click="currPage++"
                    >
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="width: 1rem; height: 1rem;">
                            <polyline points="9 18 15 12 9 6"/>
                        </svg>
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
                            {{ modalMode === 'create' ? 'Geocache Toevoegen' : 'Geocache Bewerken' }}
                        </h2>
                        <button class="admin-modal-close" @click="closeModal">
                            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                                <line x1="18" y1="6" x2="6" y2="18"/>
                                <line x1="6" y1="6" x2="18" y2="18"/>
                            </svg>
                        </button>
                    </div>
                    <div class="admin-modal-body">
                        <div class="form-row">
                            <div class="admin-form-group">
                                <label class="admin-label">Link *</label>
                                <input v-model="formData.geolink" type="url" class="admin-input" placeholder="https://coord.info/GC12345" required>
                                <span class="admin-form-hint">De volledige geocaching.com link</span>
                            </div>
                            <div class="admin-form-group">
                                <label class="admin-label">Status</label>
                                <select v-model="formData.status" class="admin-select">
                                    <option value="active">Actief</option>
                                    <option value="disabled">Uitgeschakeld</option>
                                    <option value="archived">Gearchiveerd</option>
                                </select>
                            </div>
                        </div>

                        <div class="admin-form-group">
                            <label class="admin-label">Cache Naam *</label>
                            <input v-model="formData.name" type="text" class="admin-input" required>
                        </div>

                        <div class="form-row">
                            <div class="admin-form-group">
                                <label class="admin-label">Moeilijkheid</label>
                                <select v-model="formData.difficulty" class="admin-select">
                                    <option v-for="n in [1, 1.5, 2, 2.5, 3, 3.5, 4, 4.5, 5]" :key="n" :value="n">
                                        {{ n }}
                                    </option>
                                </select>
                            </div>
                            <div class="admin-form-group">
                                <label class="admin-label">Terrein</label>
                                <select v-model="formData.terrain" class="admin-select">
                                    <option v-for="n in [1, 1.5, 2, 2.5, 3, 3.5, 4, 4.5, 5]" :key="n" :value="n">
                                        {{ n }}
                                    </option>
                                </select>
                            </div>
                        </div>

                        <div class="form-row">
                            <div class="admin-form-group">
                                <label class="admin-label">Grootte</label>
                                <select v-model="formData.size" class="admin-select">
                                    <option v-for="size in sizeOptions" :key="size.value" :value="size.value">
                                        {{ size.label }}
                                    </option>
                                </select>
                            </div>
                            <div class="admin-form-group">
                                <label class="admin-label">Type</label>
                                <select v-model="formData.type" class="admin-select">
                                    <option v-for="t in typeOptions" :key="t.value" :value="t.value">
                                        {{ t.label }}
                                    </option>
                                </select>
                            </div>
                        </div>

                        <div class="admin-form-group">
                            <label class="admin-label">Datum Geplaatst</label>
                            <input v-model="formData.placed_date" type="date" class="admin-input">
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
.form-row {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 1rem;
}

.dt-rating {
    font-weight: 500;
    color: var(--admin-text);
    font-size: 0.875rem;
}
</style>
