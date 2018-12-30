import Vue from 'vue';
import axios from 'axios';
import cookie from 'vue-cookie';
import App from './App.vue';
import router from './router';

Vue.config.productionTip = false;
Vue.prototype.$http = axios;

Vue.use(cookie);
Vue.prototype.$usernameCookie = 'username';
Vue.prototype.$sessionIdCookie = 'sessionId';

new Vue({
    router,
    render: h => h(App),
}).$mount('#app');
