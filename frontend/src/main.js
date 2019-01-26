import Vue from 'vue';
import axios from 'axios';
import App from './App.vue';
import router from './router';

Vue.config.productionTip = false;
Vue.prototype.$http = axios;
Vue.prototype.$http.getProtected = function(...args) {
    return axios.get(...args)
        .catch((e) => {
            if (e.response.status === 401) {
                this.$router.push('/');
                return;
            }
            throw e;
        });
};

Vue.prototype.CONNECTION = 'connection';
Vue.prototype.NOT_FOUND = 'not_found';
Vue.prototype.DATA_VIOLATION = 'data_volation';
Vue.prototype.UNIQUENESS_VIOLATION = 'uniqueness_violation';
Vue.prototype.UNEXPECTED = 'unexpected';

new Vue({
    router,
    render: h => h(App),
}).$mount('#app');
