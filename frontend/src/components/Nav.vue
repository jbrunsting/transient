<template>
  <div class="nav">
    <template v-if="authenticated">
      <router-link to="/">Home</router-link>
      <router-link to="/profile">Profile</router-link>
      <router-link to="/following">Following</router-link>
      <router-link to="/settings">Settings</router-link>
      <Logout v-on:logout="updateAuth()"/>
    </template>
    <template v-if="!authenticated">
      <router-link to="/">Home</router-link>
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
@import "../styles/settings.scss";

.nav {
  height: $nav-height - 2 * $margin1;
  display: flex;
  align-items: center;
  background: $accent1;
  padding: $margin1 $margin2;

  a {
    font-size: $fontsize2;
    padding: $margin1;
    font-weight: bold;
    color: $text1;
  }

  * {
    display: inline-block;
  }

  div {
    margin-left: auto;
  }
}
</style>
