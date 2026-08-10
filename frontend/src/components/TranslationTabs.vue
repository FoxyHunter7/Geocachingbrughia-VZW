<script setup>
import { ref, watch, computed } from 'vue';

const props = defineProps({
    languages: {
        type: Array,
        required: true,
    },
    translations: {
        type: Array,
        required: true,
    },
    label: {
        type: String,
        default: 'Beschrijving',
    },
    useTextarea: {
        type: Boolean,
        default: true,
    },
    useSlot: {
        type: Boolean,
        default: false,
    },
});

const emit = defineEmits(['update:translations', 'active-lang-change']);

const activeLang = ref('');

const availableLangs = computed(() => {
    if (!props.languages || props.languages.length === 0) return [];
    return props.languages;
});

watch(availableLangs, (langs) => {
    if (langs.length > 0 && !langs.find(l => l.code === activeLang.value)) {
        activeLang.value = langs[0].code;
        emit('active-lang-change', activeLang.value);
    }
}, { immediate: true });

function selectLang(code) {
    activeLang.value = code;
    emit('active-lang-change', code);
}

function getTranslation(langCode) {
    const t = props.translations.find(t => t.lang_code === langCode);
    return t ? t.description : '';
}

function updateTranslation(langCode, value) {
    const idx = props.translations.findIndex(t => t.lang_code === langCode);
    const updated = [...props.translations];
    if (idx > -1) {
        updated[idx] = { ...updated[idx], description: value };
    } else {
        updated.push({ lang_code: langCode, description: value });
    }
    emit('update:translations', updated);
}

function getFlagUrl(lang) {
    if (lang.flag_url) {
        if (lang.flag_url.includes('/assets/') || lang.flag_url.includes('fallbackLangFlags')) {
            return lang.flag_url;
        }
        return `/api/images/${lang.flag_url}`;
    }
    if (lang.imageUrl) {
        return `/assets/media/${lang.imageUrl}`;
    }
    return `/assets/media/fallbackLangFlags/${lang.code}.svg`;
}

function isActive(langCode) {
    const t = props.translations.find(t => t.lang_code === langCode);
    return t && t.description && t.description.trim() !== '';
}

defineExpose({ activeLang, updateTranslation });
</script>

<template>
    <div class="translation-tabs" v-if="availableLangs.length > 0">
        <label class="translation-label">{{ label }}</label>
        <div class="translation-tabs-body" :class="{ 'with-slot': useSlot }">
            <div class="translation-tab-list" role="tablist" :aria-label="label + ' vertalingen'">
                <button
                    v-for="lang in availableLangs"
                    :key="lang.code"
                    role="tab"
                    :aria-selected="activeLang === lang.code"
                    :aria-controls="`translation-panel-${lang.code}`"
                    :id="`translation-tab-${lang.code}`"
                    class="translation-tab"
                    :class="{ active: activeLang === lang.code, 'has-content': isActive(lang.code) }"
                    @click="selectLang(lang.code)"
                    :title="lang.name"
                >
                    <img :src="getFlagUrl(lang)" :alt="lang.name" class="translation-flag" />
                    <span class="translation-code">{{ lang.code }}</span>
                </button>
            </div>
            <div class="translation-panel-container">
                <div
                    v-for="lang in availableLangs"
                    :key="lang.code"
                    v-show="activeLang === lang.code"
                    role="tabpanel"
                    :id="`translation-panel-${lang.code}`"
                    :aria-labelledby="`translation-tab-${lang.code}`"
                >
                    <slot
                        v-if="useSlot"
                        :lang="lang"
                        :value="getTranslation(lang.code)"
                        :update="(val) => updateTranslation(lang.code, val)"
                    ></slot>
                    <textarea
                        v-else-if="useTextarea"
                        :value="getTranslation(lang.code)"
                        @input="updateTranslation(lang.code, $event.target.value)"
                        class="admin-textarea translation-input"
                        :placeholder="`${label} in ${lang.name}...`"
                        :aria-label="`${label} in ${lang.name}`"
                        rows="4"
                    ></textarea>
                    <input
                        v-else
                        type="text"
                        :value="getTranslation(lang.code)"
                        @input="updateTranslation(lang.code, $event.target.value)"
                        class="admin-input translation-input"
                        :placeholder="`${label} in ${lang.name}...`"
                        :aria-label="`${label} in ${lang.name}`"
                    />
                </div>
            </div>
        </div>
    </div>
</template>

<style scoped>
.translation-tabs {
    display: flex;
    flex-direction: column;
    gap: 0.375rem;
}

.translation-label {
    font-size: 0.875rem;
    font-weight: 500;
    color: var(--admin-text);
    margin-bottom: 0.125rem;
}

.translation-tabs-body {
    display: flex;
    gap: 0;
    border: 1px solid var(--admin-border);
    border-radius: var(--admin-radius);
    overflow: hidden;
}

.translation-tab-list {
    display: flex;
    flex-direction: column;
    gap: 0;
    border-right: 1px solid var(--admin-border);
    background: var(--admin-bg);
    overflow-y: auto;
    max-height: 200px;
    flex-shrink: 0;
    width: 52px;
}

.translation-tab {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 0.25rem;
    padding: 0.5rem 0.375rem;
    border: none;
    background: transparent;
    cursor: pointer;
    transition: background 0.15s;
    border-bottom: 1px solid var(--admin-border-light);
    position: relative;
}

.translation-tab:hover {
    background: var(--admin-surface-hover);
}

.translation-tab.active {
    background: var(--admin-primary-bg);
}

.translation-tab.active::after {
    content: '';
    position: absolute;
    left: 0;
    top: 0;
    bottom: 0;
    width: 3px;
    background: var(--admin-primary);
}

.translation-tab.has-content::before {
    content: '';
    position: absolute;
    right: 0.25rem;
    top: 0.25rem;
    width: 6px;
    height: 6px;
    border-radius: 50%;
    background: var(--admin-success);
}

.translation-flag {
    width: 1.5rem;
    height: 1rem;
    object-fit: cover;
    border-radius: 2px;
}

.translation-code {
    font-size: 0.625rem;
    font-weight: 600;
    color: var(--admin-text-muted);
    text-transform: uppercase;
}

.translation-tab.active .translation-code {
    color: var(--admin-primary);
}

.translation-panel-container {
    flex: 1;
    min-width: 0;
    display: flex;
    flex-direction: column;
}

.translation-tabs-body.with-slot {
    min-height: 250px;
}

.translation-tabs-body.with-slot .translation-panel-container {
    padding: 0.75rem;
}

.translation-input {
    border: none;
    border-radius: 0;
    height: 100%;
    min-height: 120px;
}

textarea.translation-input {
    resize: vertical;
    min-height: 120px;
}
</style>
