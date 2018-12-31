<template>
  <div class="nav">
    <router-link v-if="!authenticated" to="/">Home</router-link>
    <router-link v-if="authenticated" to="/">Profile</router-link>
    <template v-if="!authenticated">
    |
    <router-link to="/about">About</router-link>
    </template>
    <Login v-if="!authenticated" v-on:login="updateAuth()"/>
    <Logout v-if="authenticated" v-on:logout="updateAuth()"/>
  </div>
</template>

<script>
import Login from '@/components/Login.vue';
import Logout from '@/components/Logout.vue';

export default {
    name: 'Nav',
    props: {
        authenticated: Boolean,
    },
    components: {
        Login,
        Logout,
    },
    methods: {
        updateAuth() {
            this.$emit('auth');
        },
    },
};
</script>

<style scoped lang="scss">
.nav {
  a {
    padding: 30px;
    font-weight: bold;

    &.router-link-exact-active {
      color: red;
    }
  }
}
</style>
