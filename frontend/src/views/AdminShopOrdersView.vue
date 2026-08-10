<template>
    <AdminLayout pageTitle="Bestellingen">
        <template #actions>
            <span class="order-count">{{ totalOrders }} bestellingen totaal</span>
        </template>

        <!-- Filters -->
        <div class="admin-card" style="margin-bottom: 1.5rem;">
            <div class="admin-card-body">
                <div class="admin-filters">
                    <div class="admin-filter-group">
                        <label class="admin-filter-label">Status</label>
                        <select v-model="selectedStatus" @change="fetchOrders" class="admin-select">
                            <option value="">Alle Statussen</option>
                            <option value="pending">In afwachting</option>
                            <option value="paid">Betaald</option>
                            <option value="confirmed">Bevestigd</option>
                            <option value="shipped">Verzonden</option>
                            <option value="fulfilled">Vervuld</option>
                            <option value="cancelled">Geannuleerd</option>
                        </select>
                    </div>
                </div>
            </div>
        </div>

        <!-- Loading State -->
        <div v-if="loading" class="admin-card">
            <div class="admin-card-body" style="text-align: center; padding: 3rem;">
                <div class="admin-spinner"></div>
                <p style="margin-top: 1rem; color: var(--admin-muted);">Bestellingen laden...</p>
            </div>
        </div>

        <!-- Empty State -->
        <div v-else-if="orders.length === 0" class="admin-card">
            <div class="admin-empty">
                <div class="admin-empty-icon">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                        <circle cx="9" cy="21" r="1"/>
                        <circle cx="20" cy="21" r="1"/>
                        <path d="M1 1h4l2.68 13.39a2 2 0 0 0 2 1.61h9.72a2 2 0 0 0 2-1.61L23 6H6"/>
                    </svg>
                </div>
                <h3 class="admin-empty-title">Geen Bestellingen Gevonden</h3>
                <p class="admin-empty-description">
                    {{ selectedStatus ? 'Geen bestellingen gevonden met de geselecteerde status.' : 'Nog geen shopbestellingen.' }}
                </p>
            </div>
        </div>

        <!-- Orders Table -->
        <div v-else class="admin-card">
            <div class="admin-table-wrapper">
                <table class="admin-table">
                    <thead>
                        <tr>
                            <th style="width: 5rem;">Order</th>
                            <th>Item</th>
                            <th>Koper</th>
                            <th style="width: 4rem;">Aantal</th>
                            <th style="width: 7rem;">Bedrag</th>
                            <th style="width: 7rem;">Levering</th>
                            <th style="width: 7rem;">Status</th>
                            <th style="width: 10rem;">Datum</th>
                            <th style="width: 6rem; text-align: right;">Acties</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-for="order in orders" :key="order.id">
                            <td><span style="font-weight: 600;">#{{ order.id }}</span></td>
                            <td>
                                <span style="font-weight: 500;">{{ order.item_title || order.item?.title || '—' }}</span>
                            </td>
                            <td>{{ order.buyer_email || '—' }}</td>
                            <td>{{ order.quantity }}</td>
                            <td style="font-weight: 500;">{{ order.amount_display || order.amount || '—' }}</td>
                            <td>
                                <span class="fulfillment-cell">
                                    <span class="fulfillment-icon" v-html="getFulfillmentIcon(order.fulfillment_type)"></span>
                                    {{ getFulfillmentLabel(order.fulfillment_type) }}
                                </span>
                            </td>
                            <td>
                                <span class="admin-badge" :class="getStatusBadgeClass(order.status)">
                                    {{ getStatusLabel(order.status) }}
                                </span>
                            </td>
                            <td>{{ formatDate(order.created_at) }}</td>
                            <td>
                                <div class="admin-table-actions">
                                    <button
                                        class="admin-btn admin-btn-ghost admin-btn-icon admin-btn-sm"
                                        @click="openOrderDetail(order)"
                                        title="Details bekijken"
                                    >
                                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="width: 1rem; height: 1rem;">
                                            <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/>
                                            <circle cx="12" cy="12" r="3"/>
                                        </svg>
                                    </button>
                                    <button
                                        class="admin-btn admin-btn-ghost admin-btn-icon admin-btn-sm"
                                        @click="openOrderDetail(order)"
                                        title="Status wijzigen"
                                    >
                                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="width: 1rem; height: 1rem;">
                                            <path d="M21 12a9 9 0 1 1-3-6.7L21 8"/>
                                            <path d="M21 3v5h-5"/>
                                        </svg>
                                    </button>
                                </div>
                            </td>
                        </tr>
                    </tbody>
                </table>
            </div>
        </div>

        <!-- Order Detail Modal -->
        <Teleport to="body">
            <div v-if="showModal" class="admin-modal-overlay" @click.self="closeModal">
                <div class="admin-modal admin-modal-lg">
                    <div class="admin-modal-header">
                        <h2 class="admin-modal-title">
                            Bestelling <span v-if="selectedOrder">#{{ selectedOrder.id }}</span>
                        </h2>
                        <button class="admin-modal-close" @click="closeModal" aria-label="Sluiten">
                            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                                <line x1="18" y1="6" x2="6" y2="18"/>
                                <line x1="6" y1="6" x2="18" y2="18"/>
                            </svg>
                        </button>
                    </div>

                    <div class="admin-modal-body">
                        <!-- Loading Details -->
                        <div v-if="loadingDetails" style="text-align: center; padding: 3rem;">
                            <div class="admin-spinner"></div>
                        </div>

                        <template v-else-if="orderDetails">
                            <!-- Overview -->
                            <div class="detail-section">
                                <h4 class="detail-section-title">Overzicht</h4>
                                <dl class="detail-grid">
                                    <div class="detail-row">
                                        <dt>Item</dt>
                                        <dd>{{ orderDetails.item_title || orderDetails.item?.title || '—' }}</dd>
                                    </div>
                                    <div class="detail-row">
                                        <dt>Koper</dt>
                                        <dd>{{ orderDetails.buyer_email || '—' }}</dd>
                                    </div>
                                    <div class="detail-row">
                                        <dt>Aantal</dt>
                                        <dd>{{ orderDetails.quantity }}</dd>
                                    </div>
                                    <div class="detail-row">
                                        <dt>Bedrag</dt>
                                        <dd style="font-weight: 600;">{{ orderDetails.amount_display || orderDetails.amount || '—' }}</dd>
                                    </div>
                                    <div class="detail-row">
                                        <dt>Levering</dt>
                                        <dd>
                                            <span class="fulfillment-cell">
                                                <span class="fulfillment-icon" v-html="getFulfillmentIcon(orderDetails.fulfillment_type)"></span>
                                                {{ getFulfillmentLabel(orderDetails.fulfillment_type) }}
                                            </span>
                                        </dd>
                                    </div>
                                    <div class="detail-row">
                                        <dt>Datum</dt>
                                        <dd>{{ formatDate(orderDetails.created_at) }}</dd>
                                    </div>
                                </dl>
                            </div>

                            <!-- Shipping Address -->
                            <div v-if="orderDetails.fulfillment_type === 'shipping' && hasShippingAddress(orderDetails)" class="detail-section">
                                <h4 class="detail-section-title">Verzendadres</h4>
                                <div class="address-block">
                                    <template v-if="typeof orderDetails.shipping_address === 'string'">
                                        <p class="address-line">{{ orderDetails.shipping_address }}</p>
                                    </template>
                                    <template v-else>
                                        <p v-if="orderDetails.shipping_address?.name" class="address-line">{{ orderDetails.shipping_address.name }}</p>
                                        <p v-if="orderDetails.shipping_address?.line1" class="address-line">{{ orderDetails.shipping_address.line1 }}</p>
                                        <p v-if="orderDetails.shipping_address?.line2" class="address-line">{{ orderDetails.shipping_address.line2 }}</p>
                                        <p class="address-line">
                                            <span v-if="orderDetails.shipping_address?.postal_code">{{ orderDetails.shipping_address.postal_code }} </span>
                                            <span v-if="orderDetails.shipping_address?.city">{{ orderDetails.shipping_address.city }}</span>
                                        </p>
                                        <p v-if="orderDetails.shipping_address?.country" class="address-line">{{ orderDetails.shipping_address.country }}</p>
                                    </template>
                                </div>
                            </div>

                            <!-- Status -->
                            <div class="detail-section">
                                <h4 class="detail-section-title">Status</h4>
                                <div class="status-current">
                                    <span class="admin-badge" :class="getStatusBadgeClass(orderDetails.status)">
                                        {{ getStatusLabel(orderDetails.status) }}
                                    </span>
                                </div>
                                <div class="status-buttons">
                                    <button
                                        v-for="status in statusFlow"
                                        :key="status.value"
                                        class="status-btn"
                                        :class="{ active: orderDetails.status === status.value }"
                                        :data-status="status.value"
                                        :disabled="saving"
                                        @click="changeStatus(status.value)"
                                    >
                                        {{ status.label }}
                                    </button>
                                </div>
                                <div class="status-buttons" style="margin-top: 0.5rem;">
                                    <button
                                        class="status-btn status-btn-cancel"
                                        :class="{ active: orderDetails.status === 'cancelled' }"
                                        :disabled="saving || orderDetails.status === 'cancelled'"
                                        @click="changeStatus('cancelled')"
                                    >
                                        Annuleren
                                    </button>
                                </div>
                            </div>

                            <!-- Notes -->
                            <div class="detail-section">
                                <h4 class="detail-section-title">Interne Notities</h4>
                                <textarea
                                    v-model="notesDraft"
                                    class="admin-textarea"
                                    placeholder="Interne notities voor deze bestelling..."
                                    rows="3"
                                ></textarea>
                                <button
                                    class="admin-btn admin-btn-secondary"
                                    style="margin-top: 0.5rem;"
                                    :disabled="saving"
                                    @click="saveNotes"
                                >
                                    {{ saving ? 'Opslaan...' : 'Notities Opslaan' }}
                                </button>
                            </div>

                            <!-- Payment -->
                            <div class="detail-section">
                                <h4 class="detail-section-title">Betaling</h4>
                                <dl class="detail-grid">
                                    <div class="detail-row">
                                        <dt>Stripe Session ID</dt>
                                        <dd class="mono">{{ orderDetails.stripe_session_id || '—' }}</dd>
                                    </div>
                                    <div class="detail-row">
                                        <dt>Stripe Payment ID</dt>
                                        <dd class="mono">{{ orderDetails.stripe_payment_id || '—' }}</dd>
                                    </div>
                                </dl>
                            </div>
                        </template>
                    </div>

                    <div class="admin-modal-footer">
                        <button class="admin-btn admin-btn-secondary" @click="closeModal" :disabled="saving">
                            Sluiten
                        </button>
                    </div>
                </div>
            </div>
        </Teleport>
    </AdminLayout>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import AdminLayout from '@/components/admin/AdminLayout.vue'
import config from '@/data/config.js'

const router = useRouter()

// Auth & API helpers
const verifyLogin = () => {
    const token = localStorage.getItem('admin_token')
    if (!token) {
        router.push({ name: 'admin' })
        return false
    }
    return true
}

const apiRequest = async (endpoint, options = {}) => {
    const token = localStorage.getItem('admin_token')
    const response = await fetch(`${config.apiUrl}${endpoint}`, {
        ...options,
        headers: {
            'Authorization': `Bearer ${token}`,
            'Content-Type': 'application/json',
            ...options.headers
        }
    })

    if (response.status === 401) {
        localStorage.removeItem('admin_token')
        router.push({ name: 'admin' })
        throw new Error('Session expired')
    }

    if (!response.ok) {
        let msg = `Request failed (${response.status})`
        try {
            const err = await response.json()
            msg = err.error || err.message || msg
        } catch { /* response had no JSON body */ }
        throw new Error(msg)
    }

    return response
}

// State
const orders = ref([])
const loading = ref(true)
const selectedStatus = ref('')
const totalOrders = ref(0)

const showModal = ref(false)
const selectedOrder = ref(null)
const orderDetails = ref(null)
const loadingDetails = ref(false)
const saving = ref(false)
const notesDraft = ref('')

const statusFlow = [
    { value: 'pending', label: 'In afwachting' },
    { value: 'paid', label: 'Betaald' },
    { value: 'confirmed', label: 'Bevestigd' },
    { value: 'shipped', label: 'Verzonden' },
    { value: 'fulfilled', label: 'Vervuld' }
]

// Status helpers
const getStatusLabel = (status) => {
    const labels = {
        pending: 'In afwachting',
        paid: 'Betaald',
        confirmed: 'Bevestigd',
        shipped: 'Verzonden',
        fulfilled: 'Vervuld',
        cancelled: 'Geannuleerd'
    }
    return labels[status] || status
}

const getStatusBadgeClass = (status) => {
    const classes = {
        pending: 'admin-badge-danger',
        paid: 'admin-badge-warning',
        confirmed: 'admin-badge-info',
        shipped: 'admin-badge-success',
        fulfilled: 'admin-badge-neutral',
        cancelled: 'admin-badge-neutral'
    }
    return classes[status] || 'admin-badge-neutral'
}

// Fulfillment helpers
const getFulfillmentLabel = (type) => {
    const labels = {
        pickup: 'Afhalen',
        shipping: 'Verzending'
    }
    return labels[type] || type || '—'
}

const getFulfillmentIcon = (type) => {
    if (type === 'pickup') {
        return '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="width: 0.875rem; height: 0.875rem;"><path d="M3 9l1-5h16l1 5"/><path d="M4 9h16v11H4z"/><path d="M9 13h6"/></svg>'
    }
    if (type === 'shipping') {
        return '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="width: 0.875rem; height: 0.875rem;"><rect x="1" y="3" width="15" height="13"/><path d="M16 8h4l3 3v5h-7"/><circle cx="5.5" cy="18.5" r="2.5"/><circle cx="18.5" cy="18.5" r="2.5"/></svg>'
    }
    return ''
}

const hasShippingAddress = (order) => {
    const addr = order.shipping_address
    if (!addr) return false
    if (typeof addr === 'string') return addr.trim().length > 0
    return true
}

// Date helper
const formatDate = (dateStr) => {
    if (!dateStr) return '—'
    return new Date(dateStr).toLocaleDateString('en-US', {
        day: 'numeric',
        month: 'short',
        year: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
    })
}

// Fetch orders
const fetchOrders = async () => {
    if (!verifyLogin()) return
    loading.value = true

    try {
        let url = 'admin/shop/orders'
        if (selectedStatus.value) {
            url += `?status=${selectedStatus.value}`
        }

        const response = await apiRequest(url)
        const data = await response.json()

        orders.value = data.data || data || []
        totalOrders.value = data.total || orders.value.length
    } catch (error) {
        console.error('Error fetching orders:', error)
        window.$toast?.error('Bestellingen laden mislukt')
    } finally {
        loading.value = false
    }
}

// Open order detail
const openOrderDetail = async (order) => {
    selectedOrder.value = order
    orderDetails.value = null
    notesDraft.value = ''
    showModal.value = true
    await fetchOrderDetails(order.id)
}

// Fetch order details
const fetchOrderDetails = async (id) => {
    loadingDetails.value = true

    try {
        const response = await apiRequest(`admin/shop/orders/${id}`)
        orderDetails.value = await response.json()
        notesDraft.value = orderDetails.value?.notes || ''

        // Sync list item status
        const index = orders.value.findIndex(o => o.id === id)
        if (index !== -1) {
            orders.value[index].status = orderDetails.value.status
        }
    } catch (error) {
        console.error('Error fetching order details:', error)
        window.$toast?.error('Bestelgegevens laden mislukt')
    } finally {
        loadingDetails.value = false
    }
}

// Close modal
const closeModal = () => {
    if (saving.value) return
    showModal.value = false
    selectedOrder.value = null
    orderDetails.value = null
    notesDraft.value = ''
}

// Save status (and notes) via PUT admin/shop/orders/{id}/status
const saveOrderStatus = async (status, notes) => {
    if (!orderDetails.value) return

    saving.value = true

    try {
        await apiRequest(`admin/shop/orders/${orderDetails.value.id}/status`, {
            method: 'PUT',
            body: JSON.stringify({ status, notes })
        })

        orderDetails.value.status = status
        if (notes !== undefined) {
            orderDetails.value.notes = notes
        }

        // Sync list item status
        const index = orders.value.findIndex(o => o.id === orderDetails.value.id)
        if (index !== -1) {
            orders.value[index].status = status
        }

        window.$toast?.success('Bestelling bijgewerkt')
    } catch (error) {
        console.error('Error updating order:', error)
        window.$toast?.error('Bijwerken mislukt')
    } finally {
        saving.value = false
    }
}

// Change status (keep current notes)
const changeStatus = (status) => {
    if (!orderDetails.value || orderDetails.value.status === status) return
    saveOrderStatus(status, notesDraft.value)
}

// Save notes (keep current status)
const saveNotes = () => {
    if (!orderDetails.value) return
    saveOrderStatus(orderDetails.value.status, notesDraft.value)
}

onMounted(fetchOrders)
</script>

<style scoped>
.order-count {
    color: var(--admin-muted);
    font-size: 0.875rem;
}

/* Fulfillment cell */
.fulfillment-cell {
    display: inline-flex;
    align-items: center;
    gap: 0.375rem;
    font-size: 0.8125rem;
    color: var(--admin-text);
}

.fulfillment-icon {
    display: inline-flex;
    align-items: center;
    color: var(--admin-primary);
}

/* Detail sections */
.detail-section {
    margin-bottom: 1.5rem;
}

.detail-section:last-child {
    margin-bottom: 0;
}

.detail-section-title {
    font-size: 0.75rem;
    font-weight: 600;
    color: var(--admin-muted);
    text-transform: uppercase;
    letter-spacing: 0.05em;
    margin: 0 0 0.75rem 0;
}

/* Detail grid (dl) */
.detail-grid {
    display: grid;
    grid-template-columns: 1fr;
    gap: 0.625rem;
    margin: 0;
}

.detail-row {
    display: grid;
    grid-template-columns: 9rem 1fr;
    gap: 1rem;
    align-items: baseline;
}

.detail-row dt {
    font-size: 0.8125rem;
    color: var(--admin-muted);
    font-weight: 500;
}

.detail-row dd {
    margin: 0;
    font-size: 0.875rem;
    color: var(--admin-text);
    word-break: break-word;
}

.detail-row dd.mono {
    font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
    font-size: 0.8125rem;
    color: var(--admin-text-secondary);
}

/* Address block */
.address-block {
    background: var(--admin-bg);
    border: 1px solid var(--admin-border-light);
    border-radius: var(--admin-radius);
    padding: 1rem 1.25rem;
}

.address-line {
    margin: 0;
    font-size: 0.875rem;
    color: var(--admin-text);
    line-height: 1.5;
}

/* Status section */
.status-current {
    margin-bottom: 0.75rem;
}

.status-buttons {
    display: flex;
    flex-wrap: wrap;
    gap: 0.5rem;
}

.status-btn {
    padding: 0.375rem 0.875rem;
    border: 2px solid var(--admin-border);
    border-radius: 9999px;
    background: var(--admin-surface);
    font-size: 0.8125rem;
    cursor: pointer;
    transition: all 0.2s;
    color: var(--admin-muted);
}

.status-btn:hover:not(:disabled) {
    border-color: var(--admin-primary);
}

.status-btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
}

.status-btn.active[data-status="pending"] {
    border-color: var(--admin-danger);
    background: var(--admin-danger-bg);
    color: var(--admin-danger);
}

.status-btn.active[data-status="paid"] {
    border-color: var(--admin-warning);
    background: var(--admin-warning-bg);
    color: var(--admin-warning);
}

.status-btn.active[data-status="confirmed"] {
    border-color: var(--admin-info);
    background: var(--admin-info-bg);
    color: var(--admin-info);
}

.status-btn.active[data-status="shipped"] {
    border-color: var(--admin-success);
    background: var(--admin-success-bg);
    color: var(--admin-success);
}

.status-btn.active[data-status="fulfilled"] {
    border-color: var(--admin-text-muted);
    background: var(--admin-surface-hover);
    color: var(--admin-text-secondary);
}

.status-btn-cancel,
.status-btn-cancel.active {
    border-color: var(--admin-border);
    background: var(--admin-surface);
    color: var(--admin-muted);
}

.status-btn-cancel:hover:not(:disabled) {
    border-color: var(--admin-danger);
    color: var(--admin-danger);
}

.status-btn-cancel.active {
    border-color: var(--admin-danger);
    background: var(--admin-danger-bg);
    color: var(--admin-danger);
}

/* Notes textarea */
.detail-section .admin-textarea {
    resize: vertical;
    min-height: 80px;
    width: 100%;
}

/* Responsive */
@media (max-width: 768px) {
    .detail-row {
        grid-template-columns: 1fr;
        gap: 0.125rem;
    }

    .detail-row dt {
        font-size: 0.75rem;
    }
}
</style>
