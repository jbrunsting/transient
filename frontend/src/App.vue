<template>
  <div id="app">
    <Nav/>
    <Signup/>
    <Signin/>
    <router-view/>
  </div>
</template>

<script>
import Nav from '@/components/Nav.vue';
import Signup from '@/components/Signup.vue';
import Signin from '@/components/Signin.vue';

export default {
    components: {
        Nav,
        Signup,
        Signin,
    },
};
</script>

<style lang="scss">
#app {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
}
</style>
