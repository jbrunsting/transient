<template>
  <div class="nav">
    <router-link to="/">Home</router-link> |
    <router-link to="/about">About</router-link>
    <Signin v-on:signin="updateAuth()"/>
    <Signup v-if="!authenticated" v-on:signup="updateAuth()"/>
  </div>
</template>

<script>
import Signup from '@/components/Signup.vue';
import Signin from '@/components/Signin.vue';

export default {
    name: 'Nav',
    props: {
        authenticated: Boolean,
    },
    components: {
        Signup,
        Signin,
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
