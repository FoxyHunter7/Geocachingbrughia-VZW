<script setup>
    import { useEditor, EditorContent, generateHTML } from '@tiptap/vue-3';
    import { watch } from 'vue';
    import Document from '@tiptap/extension-document';
    import Paragraph from '@tiptap/extension-paragraph';
    import Text from '@tiptap/extension-text';
    import Bold from '@tiptap/extension-bold';
    import Italic from '@tiptap/extension-italic';
    import Underline from '@tiptap/extension-underline';
    import Link from '@tiptap/extension-link';
    import BulletList from '@tiptap/extension-bullet-list';
    import ListItem from '@tiptap/extension-list-item';
    import Strike from '@tiptap/extension-strike';
    import Gapcursor from '@tiptap/extension-gapcursor';
    import Table from '@tiptap/extension-table';
    import TableCell from '@tiptap/extension-table-cell';
    import TableHeader from '@tiptap/extension-table-header';
    import TableRow from '@tiptap/extension-table-row';

    const extensions = [
        Document,
        Paragraph,
        Text,
        Bold,
        Italic,
        Underline,
        Link.configure({
            openOnClick: false,
        }),
        BulletList,
        ListItem,
        Strike,
        Gapcursor,
        Table.configure({
            resizable: true,
        }),
        TableRow,
        TableHeader,
        TableCell
    ];

    const props = defineProps({
        content: String,
        editable: Boolean,
        langCode: String
    });

    function parseContent(content) {
        if (!content) return '<p></p>';
        
        try {
            const parsed = JSON.parse(content);
            return generateHTML(parsed, extensions);
        } catch (e) {
            // If JSON parsing fails, return content as-is or empty
            return content || '<p></p>';
        }
    }

    const editor = useEditor({
        extensions,
        editable: props.editable,
        content: parseContent(props.content)
    });

    // Watch for content changes from parent
    watch(() => props.content, (newContent) => {
        if (editor.value && newContent !== undefined) {
            const currentContent = JSON.stringify(editor.value.getJSON());
            // Only update if content is different to avoid cursor reset
            if (newContent !== currentContent) {
                editor.value.commands.setContent(parseContent(newContent));
            }
        }
    });

    const getContent = () => {
        if (!editor.value) return { lang_code: props.langCode, description: '' };
        return {
            lang_code: props.langCode,
            description: JSON.stringify(editor.value.getJSON())
        };
    };

    defineExpose({
        getContent
    });

    const setLink = () => {
        const previousUrl = editor.value.getAttributes('link').href;
        const url = window.prompt('URL', previousUrl);

        // If cancelled or no URL provided, unset the link
        if (url === null || url === '') {
            editor.value.chain().focus().extendMarkRange('link').unsetLink().run();
        } else {
            editor.value.chain().focus().extendMarkRange('link').setLink({ href: url }).run();
        }
    };
</script>

<template>
    <article>
        <div v-if="editor" class="menu">
            <div class="button-group">
                <button type="button" @click="editor.chain().focus().toggleBold().run()" :class="{ 'is-active': editor.isActive('bold') }">
                    Vet
                </button>
                <button type="button" @click="editor.chain().focus().toggleItalic().run()" :class="{ 'is-active': editor.isActive('italic') }">
                    Cursief
                </button>
                <button type="button" @click="editor.chain().focus().toggleStrike().run()" :class="{ 'is-active': editor.isActive('strike') }">
                    Doorhalen
                </button>
                <button type="button" @click="setLink" :class="{ 'is-active': editor.isActive('link') }">
                    Link toevoegen
                </button>
                <button type="button" @click="editor.chain().focus().unsetLink().run()" :disabled="!editor.isActive('link')">
                    Link verwijderen
                </button>
            </div>
            <div class="button-group">
                <button @click="editor.chain().focus().insertTable({ rows: 3, cols: 3, withHeaderRow: true }).run()">
                    Tabel invoegen
                </button>
                <button @click="editor.chain().focus().deleteTable().run()" class="btn-red">
                    Tabel verwijderen
                </button>
                <button @click="editor.chain().focus().addColumnBefore().run()">
                    Kolom voor
                </button>
                <button @click="editor.chain().focus().addColumnAfter().run()">
                    Kolom na
                </button>
                <button @click="editor.chain().focus().deleteColumn().run()" class="btn-red">
                    Kolom weg
                </button>
            </div>
            <div class="button-group">
                <button @click="editor.chain().focus().addRowBefore().run()">
                    Rij boven
                </button>
                <button @click="editor.chain().focus().addRowAfter().run()">
                    Rij onder
                </button>
                <button @click="editor.chain().focus().deleteRow().run()" class="btn-red">
                    Rij weg
                </button>
                <button @click="editor.chain().focus().mergeCells().run()">
                    Cellen samenvoegen
                </button>
                <button @click="editor.chain().focus().splitCell().run()">
                    Cel splitsen
                </button>
            </div>
        </div>
        <editor-content class="editor" :editor="editor"/>
    </article>
</template>

<style scoped>
    article {
        display: flex;
        flex-direction: column;
        gap: 1rem;
        width: 100%;
    }

    .editor {
        border-radius: 0.3rem;
        border: solid 0.1rem var(--color-text);
        padding: 0.75rem 1rem;
        width: 100%;
        min-height: 200px;
        max-height: 500px;
        overflow-y: auto;
    }

    .editor :deep(.ProseMirror) {
        min-height: 180px;
        outline: none;
    }

    .editor :deep(.ProseMirror p) {
        margin: 0.5em 0;
    }

    .menu {
        display: flex;
        flex-direction: column;
        gap: 0.5rem;
        width: 100%;
        padding: 0.5rem;
        background: var(--color-background-2);
        border-radius: 0.3rem;
    }

    .menu button {
        color: var(--color-text);
        background-color: var(--color-primary);
        border: none;
        height: 1.85rem;
        font-family: inherit;
        border-radius: 0.4rem;
        box-shadow: var(--color-background-2) 0.5rem 0.5rem;
        text-transform: capitalize;
        font-weight: bold;
        scale: 100%;
        transition: scale 0.15s;
        padding: 0 0.5rem;
    }

    .menu div.button-group {
        display: flex;
        padding: 0.5rem 0;
        gap: 1rem;
        flex-wrap: wrap;
    }

    .menu button:hover {
        cursor: pointer;
        scale: 103%;
        transition: scale 0.25s;
    }

    .menu button:disabled {
        display: none;
    }

    button.btn-red {
        background-color: rgb(222, 85, 61);
        color: var(--color-text);
        box-shadow: #360f01 0.5rem 0.5rem;
    }
</style>