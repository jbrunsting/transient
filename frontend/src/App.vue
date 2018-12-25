<template>
  <div id="app">
    <Nav/>
    <Signup/>
    <router-view/>
  </div>
</template>

<script>
import Nav from '@/components/Nav.vue';
import Signup from '@/components/Signup.vue';

export default {
    components: {
        Nav,
        Signup,
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
