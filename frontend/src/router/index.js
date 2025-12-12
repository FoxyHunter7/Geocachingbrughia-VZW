import { computed } from 'vue';
import { createRouter, createWebHistory } from 'vue-router';
import HomeView from '@/views/HomeView.vue';
import NotFoundView from '@/views/NotFoundView.vue';
import EventsView from '@/views/EventsView.vue';
import GeocachesView from '@/views/GeocachesView.vue';
import ShopView from '@/views/ShopView.vue';
import AdminView from '@/views/AdminView.vue';
import AdminEventsView from "@/views/AdminEventsView.vue";
import AdminGeocachesView from "@/views/AdminGeocachesView.vue";
import AdminSocialsView from "@/views/AdminSocialsView.vue";
import AdminMessagesView from "@/views/AdminMessagesView.vue";
import AdminStaticView from "@/views/AdminStaticView.vue";
import AdminLanguagesView from "@/views/AdminLanguagesView.vue";
import AdminContactFormView from "@/views/AdminContactFormView.vue";
import { StaticContentProvider } from '@/services/StaticContentService';
import { LanguageProvider } from '@/services/LanguageService';

export default function setupRouter() {
  const lang = computed(() => LanguageProvider.CURR_LANG.value);

  return createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    routes: [
      {
        path: '/',
        redirect: () => {
          return { path: `/${lang.value.toLocaleLowerCase()}` }
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
        component: EventsView,
        alias: StaticContentProvider.ROUTES.navEvents.aliases
      },
      {
        path: StaticContentProvider.ROUTES.navGeocaches.path,
        name: "geocaches",
        props: true,
        component: GeocachesView,
        alias: StaticContentProvider.ROUTES.navGeocaches.aliases
      },
      {
        path: StaticContentProvider.ROUTES.navShop.path,
        name: "shop",
        props: true,
        component: ShopView,
        alias: StaticContentProvider.ROUTES.navShop.aliases
      },
      {
        path: '/admin',
        name: "admin",
        props: false,
        component: AdminView
      },
      {
        path: '/admin/events',
        name: "adminEvents",
        props: false,
        component: AdminEventsView
      },
      {
        path: '/admin/geocaches',
        name: "adminGeocaches",
        props: false,
        component: AdminGeocachesView
      },
      {
        path: '/admin/socials',
        name: "adminSocials",
        props: false,
        component: AdminSocialsView
      },
      {
        path: '/admin/messages',
        name: "adminMessages",
        props: false,
        component: AdminMessagesView
      },
      {
        path: '/admin/static',
        name: "adminStatic",
        props: false,
        component: AdminStaticView
      },
      {
        path: '/admin/languages',
        name: "adminLanguages",
        props: false,
        component: AdminLanguagesView
      },
      {
        path: '/admin/contact-form',
        name: "adminContactForm",
        props: false,
        component: AdminContactFormView
      },
      {
        path: '/:pathMatch(.*)',
        name: "NotFound",
        component: NotFoundView,
      }
    ]
  })
}
