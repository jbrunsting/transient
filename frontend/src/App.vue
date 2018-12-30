<template>
  <div id="app">
    <router-view v-on:auth="updateAuth()" :authenticated="!!username && !!sessionId"/>
  </div>
</template>

<script>
export default {
    data() {
        return {
            username: this.$cookie.get(this.$usernameCookie),
            sessionId: this.$cookie.get(this.$sessionIdCookie),
        };
    },
    methods: {
        updateAuth() {
            this.username = this.$cookie.get(this.$usernameCookie);
            this.sessionId = this.$cookie.get(this.$sessionIdCookie);
            if (!this.username || !this.sessionId) {
                this.$router.push('/');
            }
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
