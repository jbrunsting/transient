<template>
  <div id="app">
    <Nav/>
    <Signup v-on:signup="updateSessionId()"/>
    <Signin v-on:signin="updateSessionId()"/>
    <router-view :sessionId="sessionId"/>
  </div>
</template>

<script>
import Nav from '@/components/Nav.vue';
import Signup from '@/components/Signup.vue';
import Signin from '@/components/Signin.vue';

export default {
    data() {
        console.log(`Session ID is '${this.$cookie.get('sessionId')}'`);
        return {
            sessionId: this.$cookie.get('sessionId'),
        };
    },
    components: {
        Nav,
        Signup,
        Signin,
    },
    methods: {
        updateSessionId() {
            this.sessionId = this.$cookie.get('sessionId');
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
