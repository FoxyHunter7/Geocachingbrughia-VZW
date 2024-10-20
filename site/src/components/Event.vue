<script setup>
    import { computed, toRef, watch } from 'vue';
    import { useEditor, EditorContent, generateHTML } from '@tiptap/vue-3';
    import config from '../data/config.js';
    import LanguageProvider from '@/services/LanguageService';
    import Document from '@tiptap/extension-document';
    import Paragraph from '@tiptap/extension-paragraph';
    import Text from '@tiptap/extension-text';
    import Bold from '@tiptap/extension-bold';
    import Italic from '@tiptap/extension-italic';
    import Underline from '@tiptap/extension-underline';
    import Link from '@tiptap/extension-link';
    import BulletList from '@tiptap/extension-bullet-list';
    import ListItem from '@tiptap/extension-list-item';
    import { StaticContentProvider as SCP } from '@/services/StaticContentService.js';

    const props = defineProps({
        event: Object
    });

    const dictionary = SCP.DICTIONARY;
    const lang = computed(() => LanguageProvider.CURR_LANG.value);

    const event = toRef(props, 'event');

    const parseEventDescription = (event) => {
        return generateHTML(JSON.parse(event.translations[0].description), [
            Document,
            Paragraph,
            Text,
            Bold,
            Italic,
            Underline,
            Link,
            BulletList,
            ListItem
        ]);
    };

    const eventDescription = computed(() => {
        return parseEventDescription(event.value);
    });

    const editor = useEditor({
        extensions: [
            Document,
            Paragraph,
            Text,
            Bold,
            Italic,
            Underline,
            Link,
            BulletList,
            ListItem
        ],
        editable: false,
        content: eventDescription.value
    });

    watch(event, (newEvent) => {
        if (editor.value && newEvent) {
            const newContent = parseEventDescription(newEvent);
            editor.value.commands.setContent(newContent);
        }
    }, { immediate: true });

    const {formattedStartDate, formattedEndDate} = computed(() => {
        const startDate = new Date(event.value.start_date);
        const endDate = new Date(event.value.end_date);

        const options = {
            year: 'numeric',
            month: 'short',
            day: 'numeric',
            hour: '2-digit',
            minute: '2-digit'
        }

        return {
            formattedStartDate: startDate.toLocaleString(LanguageProvider.CURR_LANG.value, options),
            formattedEndDate: endDate.toLocaleDateString(LanguageProvider.CURR_LANG.value, options)
        };
    }).value;
</script>

<template>
    <section>
        <div>
            <img :src="`${config.apiUrl}images/${event.imageUrl}`">
        </div>
        <div>
            <h2>{{ event.title }}</h2>
            <editor-content :editor="editor" />
            <ul>
                <li><img :src="`/assets/media/eventtypes/${event.type}.png`"><p>{{ event.type }}</p></li>
                <li><p><span>{{ formattedStartDate }}</span> - <span>{{ formattedEndDate }}</span></p></li>
                <li v-if="event.ticket_purchase_url"><a class="capitalize" :href="event.ticket_purchase_url" target="_blank">{{ dictionary.ButtonPurchase[lang] }}</a></li>
                <li><a :href="event.geolink" target="_blank">Geocaching.com</a></li>
            </ul>
        </div>
    </section>
</template>

<style scoped>
    section {
        display: grid;
        grid-template-columns: 1fr 2fr;
        gap: 3rem;
        padding: 3rem;
        max-height: 40rem;
    }

    section > div:first-child {
        height: 100%;
        max-height: 40rem;
        display: flex;
        justify-content: center;
        align-items: flex-start;
        overflow: hidden;
    }

    section > div > img {
        max-width: 100%;
        max-height: 100%;
        object-fit: contain;
        border-radius: 0.5rem;
    }

    h2 {
        font-weight: bold;
    }

    h2::first-letter {
        text-transform: uppercase;
    }

    ul {
        margin-top: 1rem;
        padding-inline-start: 0;
        display: flex;
        justify-content: flex-start;
        align-items: center;
        height: 3rem;
        gap: 2rem;
        font-size: 0.8rem;
    }

    ul li {
        margin: none;
        padding: none;
        list-style: none;
        display: block;
    }

    ul li:first-child {
        display: flex;
        gap: 0.5rem;
        align-items: center;
        justify-content: center;
        height: 100%;
    }

    ul li:first-child img {
        height: 100%;
        object-fit: contain;
    }

    a {
        display: block;
        background-color: var(--color-accent-light);
        color: var(--color-text);
        padding: 0.5rem 1rem;
        text-decoration: none;
        border-radius: 0.4rem;
        transform: scale(100%);
        transition: transform 0.15s;
    }

    a.capitalize {
        text-transform: capitalize;
    }

    a:hover, a.capitalize:hover {
        cursor: pointer;
        transform: scale(102%);
        transition: transform 0.15s;
    }

    @media screen and (max-width: 1000px) {
        section {
            height: max-content !important;
            max-height: max-content !important;
            grid-template-columns: 1fr;
            grid-template-rows: auto;
            justify-content: center;
            align-items: center;
            padding: 1rem;
        }

        section > div:first-child {
            height: 70vh;
            max-height: 40rem;
            display: flex;
            justify-content: center;
            align-items: flex-start;
            overflow: hidden;
        }

        ul {
            height: auto;
            flex-direction: column;
        }

        ul li:first-child {
            display: flex;
            gap: 0.5rem;
            align-items: center;
            height: 2.5rem;
        }

        li a {
            box-sizing: border-box
        }

        h2 {
            margin-top: auto;
        }

        section > div:first-child {
            align-items: center;
        }
    }
</style>
