<template>
  <div class="settings">
    <Nav v-on:auth="updateAuth()" :authenticated="authenticated"/>
    <div class="content">
      <form @submit.prevent="invalidateSessions">
        <button type="submit">Logout of all other sessions</button>
      </form>
      <Login submitText="Delete account" apiPath="/api/user/delete" v-on:login="updateAuth()"/>
    </div>
  </div>
</template>
<script>
import Nav from '@/components/Nav.vue';
import Login from '@/components/Login.vue';

export default {
    name: 'about',
    props: {
        authenticated: Boolean,
    },
    components: {
        Nav,
        Login,
    },
    data() {
        return {
            username: '',
            password: '',
        };
    },
    methods: {
        updateAuth() {
            this.$emit('auth');
        },
        invalidateSessions() {
            this.$http.post('/api/user/invalidate')
                .then(() => {
                    alert('Successfully logged out of all other sessions'); // eslint-disable-line no-alert
                }).catch((e) => {
                    console.log(`Error ${JSON.stringify(e)}`);
                });
        },
    },
};
</script>
