import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '@/views/HomeView.vue'
import { StaticContentProvider } from '@/services/StaticContentService'
import NotFoundView from '@/views/NotFoundView.vue';
import { LanguageProvider } from '@/services/LanguageService';

await initStaticContentProvider();

async function initStaticContentProvider() {
  const staticContentProvider = new StaticContentProvider()
  await staticContentProvider.init()
}

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      redirect: () => {
        return { path: `/${LanguageProvider.CURR_LANG.toLocaleLowerCase()}` }
      }
    },
    {
      path: '/:lang([A-Za-z]{2})',
      redirect: to => {
        return `/${to.params.lang}/${StaticContentProvider.DICTIONARY.NavHome[to.params.lang.toLocaleUpperCase()]}`
      }
    },
    {
      path: StaticContentProvider.ROUTES.navHome.path,
      name: "home",
      props: true,
      component: HomeView,
      alias: StaticContentProvider.ROUTES.navHome.aliases
    },
    {
      path: StaticContentProvider.ROUTES.navEvents.path,
      name: "events",
      props: true,
      component: HomeView,
      alias: StaticContentProvider.ROUTES.navEvents.aliases
    },
    {
      path: StaticContentProvider.ROUTES.navGeocaches.path,
      name: "geocaches",
      props: true,
      component: HomeView,
      alias: StaticContentProvider.ROUTES.navGeocaches.aliases
    },
    {
      path: StaticContentProvider.ROUTES.navShop.path,
      name: "shop",
      props: true,
      component: HomeView,
      alias: StaticContentProvider.ROUTES.navShop.aliases
    },
    {
      path: '/:pathMatch(.*)',
      name: "NotFound",
      component: NotFoundView,
    }
  ]
})

export default router
