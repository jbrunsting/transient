<template>
  <div class="nav">
    <router-link v-if="!authenticated" to="/">Home</router-link>
    <router-link v-if="authenticated" to="/">Profile</router-link>
    <template v-if="!authenticated">
    |
    <router-link to="/about">About</router-link>
    </template>
    <Signin v-if="!authenticated" v-on:signin="updateAuth()"/>
    <Signout v-if="authenticated" v-on:signout="updateAuth()"/>
  </div>
</template>

<script>
import Signin from '@/components/Signin.vue';
import Signout from '@/components/Signout.vue';

export default {
    name: 'Nav',
    props: {
        authenticated: Boolean,
    },
    components: {
        Signin,
        Signout,
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
