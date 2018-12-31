<template>
  <div id="app">
    <router-view
        v-if="authenticated != null"
        v-on:auth="updateAuth()"
        :authenticated="authenticated"/>
  </div>
</template>

<script>
export default {
    data() {
        return {
            authenticated: null,
        };
    },
    methods: {
        updateAuth() {
            this.$http.get('/api/user/authenticated')
                .then(() => {
                    if (!this.authenticated) {
                        this.authenticated = true;
                        this.$router.push('/');
                    }
                }).catch(() => {
                    if (this.authenticated) {
                        this.authenticated = false;
                        this.$router.push('/');
                    }
                });
        },
    },
    created() {
        this.updateAuth();
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
