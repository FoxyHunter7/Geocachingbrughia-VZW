<script setup>
    import { computed } from 'vue';
    const props = defineProps({
        error: String,
        date: String
    });

    const formattedDate = computed(() => {
        if (!props.date) return '';

        const date = new Date(props.date);
        const options = {
            year: 'numeric',
            month: 'long',
            day: 'numeric',
            hour: '2-digit',
            minute: '2-digit',
            timeZoneName: 'short'
        };
        return date.toLocaleString(navigator.language, options);
    });
</script>

<template>
    <div>
        <div class="warn-icon"></div>
        <p>{{ error + formattedDate }}</p>
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
        border-radius: 0.5rem;
    }

    div.warn-icon {
        height: 80%;
        aspect-ratio: 1 / 1;
        background-color: var(--color-text2);
        mask: url(../assets/media/alert-triangle.svg);
        mask-size: contain;
    }

    p {
        color: var(--color-text2);
    }
</style>