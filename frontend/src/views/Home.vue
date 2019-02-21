<template>
  <div class="home">
    <Nav v-on:auth="updateAuth()" :authenticated="authenticated"/>
    <div v-if="!authenticated" class="content">
      <Landing v-on:auth="updateAuth()"/>
    </div>
    <Home v-if="authenticated"/>
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

<style>
.content {
  margin-top: 0;
}
</style>
