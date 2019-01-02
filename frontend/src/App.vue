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
                    if (this.authenticated || this.authenticated == null) {
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
@import "styles/settings.scss";

body {
  margin: 0;
}

#app {
  font-family: $font-stack;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: $text0;
}

input {
  font-family: $font-stack;
  color: $text0;
}

button {
  font-weight: bold;
  font-size: $fontsize1;
  background-color: $accent0;
  color: $text1;
  padding: $margin0 $margin1;
  border-radius: $margin0;
  border: none;
  transition: background-color $duration0;
  display: inline-block;

  &:hover {
    cursor: pointer;
    background-color: lighten($accent0, 5%);
  }
}

a {
  color: $accent0;
  padding: $margin0 $margin1;
  text-decoration: none;

  &:visited {
    color: $accent0;
  }
}
</style>
