import Vue from 'vue';
import Router from 'vue-router';
import Home from './views/Home.vue';
import cookie from 'vue-cookie';

Vue.use(Router);

let router = new Router({
    mode: 'history',
    base: process.env.BASE_URL,
    routes: [{
            path: '/',
            name: 'home',
            component: Home,
        },
        {
            path: '/about',
            name: 'about',
            // route level code-splitting
            // this generates a separate chunk (about.[hash].js) for this route
            // which is lazy-loaded when the route is visited.
            component: () =>
                import( /* webpackChunkName: "about" */ './views/About.vue'),
            meta: {
                noAuth: true
            },
        },
    ],
});

router.beforeEach((to, from, next) => {
    const authenticated = (Vue.prototype.$cookie.get(Vue.prototype.$sessionIdCookie) &&
        Vue.prototype.$cookie.get(Vue.prototype.$usernameCookie));
    if (to.meta.auth && !authenticated || to.meta.noAuth && authenticated) {
        next({
            name: 'home'
        })
        return
    }
    next()
});

export default router
