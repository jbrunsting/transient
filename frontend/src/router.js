import Vue from 'vue';
import Router from 'vue-router';
import Home from './views/Home.vue';

Vue.use(Router);

export default new Router({
    mode: 'history',
    base: process.env.BASE_URL,
    routes: [
        {
            path: '/',
            name: 'home',
            component: Home,
        },
        {
            path: '/settings',
            name: 'settings',
            component: () => import(/* webpackChunkName: "settings" */
                './views/Settings',
            ),
        },
        {
            path: '/about',
            name: 'about',
            component: () => import(/* webpackChunkName: "about" */
                './views/About.vue',
            ),
        },
        {
            path: '/profile/:username',
            name: 'profile',
            component: () => import(/* webpackChunkName: "profile" */
                './views/Profile.vue',
            ),
        },
        // TODO: Don't allow unless authenticated
        {
            path: '/following',
            name: 'following',
            component: () => import(/* webpackChunkName: "following" */
                './views/Following.vue',
            ),
        },
    ],
});
