<template>
  <div class="settings">
    <Nav v-on:auth="updateAuth()" :authenticated="authenticated"/>
    <form @submit.prevent="invalidateSessions">
      <button type="submit">Logout of all other sessions</button>
    </form>
    <form @submit.prevent="deleteAccount">
      <input placeholder="username" v-model="username">
      <input type="password" placeholder="password" v-model="password">
      <button type="deleteAccount">Delete account</button>
    </form>
  </div>
</template>
<script>
import Nav from '@/components/Nav.vue';

export default {
    name: 'about',
    props: {
        authenticated: Boolean,
    },
    components: {
        Nav,
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
        deleteAccount() {
            const identification = {
                username: this.username,
                password: this.password,
            };
            this.$http.post('/api/user/delete', identification)
                .then(() => {
                    this.$emit('auth');
                }).catch((e) => {
                    console.log(`Error ${JSON.stringify(e)}`);
                });
        },
    },
};
</script>
