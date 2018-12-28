<template>
  <div id="app">
    <Nav/>
    <Signup v-on:signup="updateSession()"/>
    <Signin v-on:signin="updateSession()"/>
    <router-view :session="session"/>
  </div>
</template>

<script>
import Nav from '@/components/Nav.vue';
import Signup from '@/components/Signup.vue';
import Signin from '@/components/Signin.vue';

export default {
    data() {
        console.log(`Session is '${this.$cookie.get('session')}'`);
        return {
            session: this.$cookie.get('session'),
        };
    },
    components: {
        Nav,
        Signup,
        Signin,
    },
    methods: {
        updateSession() {
            this.session = this.$cookie.get('session');
        },
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
