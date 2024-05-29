<script setup>
import { ref } from 'vue'
import { fetchAllEvents } from '../services/EventService.js'
import { defineProps } from 'vue'

const props = defineProps({
    lang: String
})

const events = ref(null)
events.value = (await fetchAllEvents(props.lang)).data
</script>

<template>
    <p>All events:</p>
    <ul v-for="event in events">
        <li>{{ event.title }}</li>
        <li>{{ event.geolink }}</li>
        <li><p style="color: gray">{{ event?.translations[0]?.description }}</p></li>
    </ul>
</template>