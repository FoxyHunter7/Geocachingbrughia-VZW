<script setup>
    import { useEditor, EditorContent, generateHTML } from '@tiptap/vue-3';
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

    const props = defineProps({
        content: String,
        editable: Boolean,
        langCode: String
    });

    function parseContent(content) {
        return generateHTML(JSON.parse(content), [
            Document,
            Paragraph,
            Text,
            Bold,
            Italic,
            Underline,
            Link,
            BulletList,
            ListItem,
            Strike,
            Gapcursor,
            Table.configure({
                resizable: false,
            }),
            TableRow,
            TableHeader,
            TableCell
        ]);
    };

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
            ListItem,
            Strike,
            Gapcursor,
            Table.configure({
                resizable: true,
            }),
            TableRow,
            TableHeader,
            TableCell
        ],
        editable: props.editable,
        content: parseContent(props.content)
    });

    const getContent = () => {
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
                    Bold
                </button>
                <button type="button" @click="editor.chain().focus().toggleItalic().run()" :class="{ 'is-active': editor.isActive('italic') }">
                    Italic
                </button>
                <button type="button" @click="editor.chain().focus().toggleStrike().run()" :class="{ 'is-active': editor.isActive('strike') }">
                    Strike
                </button>
                <button type="button" @click="setLink" :class="{ 'is-active': editor.isActive('link') }">
                    Set link
                </button>
                <button type="button" @click="editor.chain().focus().unsetLink().run()" :disabled="!editor.isActive('link')">
                    Unset link
                </button>
            </div>
            <br>
            <div>
                <div class="button-group">
                    <button @click="editor.chain().focus().insertTable({ rows: 3, cols: 3, withHeaderRow: true }).run()">
                        Insert table
                    </button>
                    <button @click="editor.chain().focus().deleteTable().run()" class="btn-red">
                        Delete table
                    </button>
                </div>
                <div class="button-group">
                    <button @click="editor.chain().focus().addColumnBefore().run()">
                      Add column before
                    </button>
                    <button @click="editor.chain().focus().addColumnAfter().run()">
                      Add column after
                    </button>
                    <button @click="editor.chain().focus().deleteColumn().run()" class="btn-red">
                      Delete column
                    </button>
                </div>
                <div class="button-group">
                    <button @click="editor.chain().focus().addRowBefore().run()">
                      Add row before
                    </button>
                    <button @click="editor.chain().focus().addRowAfter().run()">
                      Add row after
                    </button>
                    <button @click="editor.chain().focus().deleteRow().run()" class="btn-red">
                      Delete row
                    </button>
                </div>
                <br>
                <div class="button-group">
                    <button @click="editor.chain().focus().mergeCells().run()">
                      Merge cells
                    </button>
                    <button @click="editor.chain().focus().splitCell().run()">
                      Split cell
                    </button>
                    <button @click="editor.chain().focus().toggleHeaderColumn().run()">
                      Toggle header column
                    </button>
                    <button @click="editor.chain().focus().toggleHeaderRow().run()">
                      Toggle header row
                    </button>
                    <button @click="editor.chain().focus().toggleHeaderCell().run()">
                      Toggle header cell
                    </button>
                    <button @click="editor.chain().focus().mergeOrSplit().run()">
                      Merge or split
                    </button>
                    <button @click="editor.chain().focus().setCellAttribute('colspan', 2).run()">
                      Set cell attribute
                    </button>
                    <button @click="editor.chain().focus().fixTables().run()">
                      Fix tables
                    </button>
                    <button @click="editor.chain().focus().goToNextCell().run()">
                      Go to next cell
                    </button>
                    <button @click="editor.chain().focus().goToPreviousCell().run()">
                      Go to previous cell
                    </button>
                </div>
            </div>
        </div>
        <editor-content class="editor" :editor="editor"/>
    </article>
</template>

<style scoped>
    article {
        display: flex;
        flex-direction: row-reverse;
        gap: 2rem;
    }

    .editor {
        border-radius: 0.3rem;
        border: solid 0.1rem var(--color-text);
        padding: 0.3rem 0.5rem;
        width: 50%;
    }

    .menu {
        display: flex;
        flex-direction: column;
        gap: 1rem;
        width: 50%;
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