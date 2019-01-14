<template>
  <div class="home">
    <Nav v-on:auth="updateAuth()" :authenticated="authenticated"/>
    <div class="content">
      <Landing v-if="!authenticated" v-on:auth="updateAuth()"/>
      <Home v-if="authenticated"/>
    </div>
  </div>
</template>

<script>
import Nav from '@/components/Nav.vue';
import Landing from '@/components/Landing.vue';
import Home from '@/components/Home.vue';

export default {
    name: 'home',
    props: {
        authenticated: Boolean,
    },
    components: {
        Nav,
        Landing,
        Home,
    },
    methods: {
        updateAuth() {
            this.$emit('auth');
        },
    },
};
</script>
