<template>
  <div class="about">
    <Nav v-on:auth="updateAuth()" :authenticated="authenticated"/>
    <div class="content">
      <h1>This is an about page</h1>
      <Error>
        There was a problem! Hi <a href="www.google.com">google</a>
      </Error>
    </div>
  </div>
</template>

<script>
import Nav from '@/components/Nav.vue';
import Error from '@/components/Error.vue';

export default {
    name: 'about',
    props: {
        authenticated: Boolean,
    },
    components: {
        Nav,
        Error,
    },
    methods: {
        updateAuth() {
            this.$emit('auth');
        },
    },
};
</script>
