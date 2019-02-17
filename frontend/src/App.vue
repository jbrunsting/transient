<template>
  <div id="app">
    <router-view
        v-if="authenticated != null"
        v-on:auth="updateAuth()"
        :authenticated="authenticated"/>
  </div>
</template>

<script>
import router from './router'

export default {
    data() {
        return {
            authenticated: null,
        };
    },
    methods: {
        updateAuth() {
            this.$http.get('/api/authenticated')
                .then(() => {
                    if (!this.authenticated) {
                        this.authenticated = true;
                    }
                }).catch(() => {
                    if (this.authenticated || this.authenticated == null) {
                        this.authenticated = false;
                        router.push('/');
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
  background-color: $base2;
}

#app {
  font-family: $font-stack;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: $text0;
}

textarea,
input {
  font-family: $font-stack;
  color: $text0;
  padding: $margin0 $margin1;
  margin: $margin0;
  border-radius: $margin0;
  border-color: $accent0;
  width: 150px;
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
  text-decoration: none;

  &:visited {
    color: $accent0;
  }
}

.content {
  max-width: $page-width;
  margin: $margin2 auto;
  padding: $margin2;
  background-color: $base0;
  border-radius: $margin0;
}
</style>
