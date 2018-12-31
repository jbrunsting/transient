<template>
  <div class="nav">
    <template v-if="authenticated">
      <router-link to="/">Profile</router-link>|
      <router-link to="/settings">Settings</router-link>
      <Logout v-on:logout="updateAuth()"/>
    </template>
    <template v-if="!authenticated">
      <router-link to="/">Home</router-link>|
      <router-link to="/about">About</router-link>
      <Login v-on:login="updateAuth()"/>
    </template>
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
