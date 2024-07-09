<script setup>
    import { fetchEvents, getProfileData } from '@/services/AdminService';
    import { onMounted, ref, computed, watch } from 'vue';
    import { useRouter } from 'vue-router';

    const router = useRouter();

    async function verifyLogin() {
        const response = await getProfileData();

        if (response.status) {
            fetchData();
        } else {
            router.push({ name: "admin" });
        }
    }

    const events = ref([]);
    const currPage = ref(1);
    const lastPage = ref(1);
    const search = ref("");

    async function fetchData() {
        const pagedResponse = await fetchEvents();

        if (pagedResponse.data) {
            events.value = pagedResponse.data;
            currPage.value = pagedResponse.current_page;
            lastPage.value = pagedResponse.lastPage;
        } else if (pagedResponse.access_denied) {
            window.alert(response.access_denied);
            router.push({ name: "admin" });
        } else {
            router.push({ name: "admin" });
        }
    }

    const filteredEvents = computed(() => {
        if (!search.value) {
            return events.value;
        }

        return events.value.filter(event =>
            event.title.toLowerCase().includes(search.value.toLowerCase())
        );
    });

    function nextPage() {
        if (currPage.value < lastPage.value) {
            currPage.value++;
        }
    }

    function prevPage() {
        if (currPage > 1) {
            currPage.value--;
        }
    }

    function dateTimeFormatter(dateTimeString) {
        if (!dateTimeString) {
            "n/a";
        }

        const date = new Date(dateTimeString);
        if (isNaN(date)) {
            return "parsing error";
        }

        return date.toLocaleString("nl-be", {
            year: "numeric",
            month: "short",
            day: "2-digit",
            hour: "2-digit",
            minute: "2-digit"
        });
    }

    onMounted(verifyLogin);
    watch(currPage, async () => await fetchData());
</script>

<template>
    <header>
        <button type="button">Terug</button>
        <form v-show="events.length > 0 && events[0] !== 'loading'" method="post" @submit.prevent="" >
            <div>
                <input v-model="search" type="search" id="search" name="search" autocomplete="search" required placeholder="Zoeken">
            </div>
        </form>
    </header>
    <main>
        <table>
            <thead>
                <th>ID</th>
                <th>Type</th>
                <th>Status</th>
                <th>Op Start</th>
                <th>Titel</th>
                <th>Start</th>
                <th>Eind</th>
                <th></th>
            </thead>
            <tbody>
                <tr v-for="event in filteredEvents">
                    <td>{{ event.id }}</td>
                    <td><div><img :src="`/src/assets/media/eventtypes/${event.type}.png`"></div></td>
                    <td :class="event.state">{{ event.state }}</td>
                    <td :class="event.on_home">{{ event.on_home }}</td>
                    <td>{{ event.title }}</td>
                    <td>{{ dateTimeFormatter(event.start_date) }}</td>
                    <td>{{ dateTimeFormatter(event.end_date) }}</td>
                    <td><div class="icon-edit"></div></td>
                </tr>
            </tbody>
        </table>
    </main>
</template>

<style scoped>
    header {
        display: flex;
        justify-content: space-between;
    }

    button {
        color: var(--color-text);
        background-color: var(--color-secondary);
        border: none;
        width: 9rem;
        height: 2rem;
        font-family: inherit;
        border-radius: 0.4rem;
        box-shadow: var(--color-background2) 0.5rem 0.5rem;
        text-transform: capitalize;
        font-weight: bold;
        scale: 100%;
        transition: scale 0.15s;
        margin: 1rem;
    }

    form {
        display: flex;
        justify-content: flex-end;
        padding: 1rem;
        gap: 1rem;
    }

    form > div input {
        width: 20rem;
        max-width: 95vw;
    }

    form > div {
        position: relative;
    }

    form > div::before {
        content: '';
        height: 1.3rem;
        width: 1.3rem;
        position: absolute;
        left: 0.5rem;
        top: 0.3rem;
        background-color: var(--color-text);
        mask: url(../assets/media/search.svg);
        mask-size: contain;
    }

    form input, form textarea {
        height: 2rem;
        border-radius: 0.3rem;
        border: solid 0.1rem var(--color-text);
        outline: none;
        font-family: inherit;
        font-size: 1rem;
        box-sizing: border-box;
        padding: 0 0.5rem 0 2rem;
        line-height: 2rem;
        text-transform: capitalize;
        background-color: var(--color-background);
        color: var(--color-text)
    }

    table {
        width: 90vw;
        max-width: 80rem;
        margin: 3rem auto;
    }

    thead {
        background-color: var(--color-primary);
        color: var(--color-text);
    }

    th, td, tr {
        text-align: center;
        padding: 0 0.5rem;
    }

    tbody tr:nth-child(even) {
        background-color: var(--color-background2);
    }

    th:nth-child(5), td:nth-child(5) {
        text-align: start;
    }

    td {
        height: 1.8rem;
    }

    td div {
        display: flex;
        justify-content: center;
    }

    td img {
        height: 1.5rem;
    }

    th:nth-child(1) {
        width: 3rem;
    }

    th:nth-child(2) {
        width: 5rem;
    }

    th:nth-child(3), th:nth-child(4) {
        width: 6rem;
    }

    td:nth-child(4) {
        text-transform: uppercase;
    }

    th:nth-child(6), th:nth-child(7) {
        width: 12rem;
    }

    th:nth-child(8) {
        width: 3rem;
    }

    td.ONLINE, td.true {
        background-color: rgb(125, 209, 0);
    }

    td.DRAFT {
        background-color: darkorange;
        color: var(--color-text3);
    }

    td.ARCHIVED, td.false {
        background-color: red;
        color: var(--color-text4);
    }

    div.icon-edit {
        margin: auto;
        height: 60%;
        aspect-ratio: 1 / 1;
        mask: url(../assets/media/edit.svg);
        mask-size: contain;
        background-color: var(--color-text);
    }

    div.icon-edit:hover {
        cursor: pointer;
        background-color: var(--color-primary);
    }

    @media (prefers-color-scheme: dark) {
        thead {
            background-color: var(--color-secondary);
        }

        td.ONLINE, td.true {
            background-color: green;
        }

        td.ARCHIVED, td.false {
            background-color: darkred;
        }
    }
</style>