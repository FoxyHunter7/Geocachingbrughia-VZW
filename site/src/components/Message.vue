<script setup>
    import { computed } from 'vue';

    const props = defineProps({
        message: Object
    });

    const messageP = computed(() => (props.message.body) ? props.message.body : props.message.title);
    const messageH = computed(() => (props.message.body) ? props.message.title : "");
</script>

<template>
    <div>
        <h2 v-if="messageH">{{ messageH }}</h2>
        <p>{{ messageP }}</p>
    </div>
</template>

<style scoped>
    div {
        color: var(--color-text);
        background-color: var(--color-secondary);
        padding: 1rem;
        border-radius: 0.5rem;
        width: fit-content;
    }

    h2 {
        font-size: 1.1rem;
        font-weight: bold;
        text-transform: capitalize;
    }
</style>