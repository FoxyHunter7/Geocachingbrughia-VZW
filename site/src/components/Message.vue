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
        padding: 0.75rem 1rem;
        width: 100%;
    }

    h2 {
        font-size: 1rem;
        font-weight: bold;
        text-transform: capitalize;
    }
</style>