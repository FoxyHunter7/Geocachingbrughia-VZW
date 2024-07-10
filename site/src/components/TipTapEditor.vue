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
            Strike
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
    <div v-if="editor" class="menu">
        <button @click="editor.chain().focus().toggleBold().run()" :class="{ 'is-active': editor.isActive('bold') }">
          Bold
        </button>
        <button @click="editor.chain().focus().toggleItalic().run()" :class="{ 'is-active': editor.isActive('italic') }">
          Italic
        </button>
        <button @click="editor.chain().focus().toggleStrike().run()" :class="{ 'is-active': editor.isActive('strike') }">
          Strike
        </button>
        <button @click="setLink" :class="{ 'is-active': editor.isActive('link') }">
          Set link
        </button>
        <button @click="editor.chain().focus().unsetLink().run()" :disabled="!editor.isActive('link')">
          Unset link
        </button>
      </div>
      <editor-content :editor="editor"/>
</template>