import './assets/css/main.css'

import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import { StaticContentProvider } from '@/services/StaticContentService'

async function init() {
    const staticContentProvider = new StaticContentProvider()
    await staticContentProvider.init()

    const app = createApp(App)
    app.use(router())
    app.mount('#app')
}

init();