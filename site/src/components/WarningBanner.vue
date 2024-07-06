<script setup>
    import LanguageProvider from '@/services/LanguageService';
    import { computed } from 'vue';

    const props = defineProps({
        error: Object
    });
    const errorMessage = computed(() => props.error[LanguageProvider.CURR_LANG.value]);
</script>

<template>
    <div>
        <div class="warn-icon"></div>
        <p>{{ errorMessage }}</p>
    </div>
</template>

<style scoped>
    div {
        background-color: var(--color-alert);
        display: flex;
        justify-content: flex-start;
        align-items: center;
        gap: 2rem;
        padding: 1rem;
    }

    div.warn-icon {
        height: 80%;
        aspect-ratio: 1 / 1;
        background-color: var(--color-text);
        mask: url(../assets/media/alert-triangle.svg);
        mask-size: contain;
    }

    p {
        color: var(--color-text);
    }
</style>