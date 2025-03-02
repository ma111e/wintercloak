import * as VueRouter from "vue-router";
import {createMemoryHistory, createRouter} from "vue-router";
import Home from "@/views/Home.vue"
import ThemeSettings from "@/views/ThemeSettings.vue";
import EnginesContainer from "@/views/EnginesContainer.vue";

const routes: VueRouter.RouteRecordRaw[] = [
    {
        path: "/",
        name: 'home',
        component: Home,
        children: []
    },
    {
        path: "/settings/theme",
        name: 'themeSettings',
        component: ThemeSettings,
        children: []
    },
    {
        path: '/engines/:engineName?',
        name: 'engines',
        component: EnginesContainer, // Default component
        // beforeEnter: async (to, _, next) => {
            // if (!to.query.mode) {
            //     to.query.mode = "playground"
            //     next({query: to.query});
            // }
        // },
        children: [],
    },
]

const router = createRouter({
    // history: createWebHistory(),
    history: createMemoryHistory(),
    routes,
})

export default router;
