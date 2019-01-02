<template>
  <div class="home">
    <Nav v-on:auth="updateAuth()" :authenticated="authenticated"/>
    <div class="content">
      <Landing v-if="!authenticated" v-on:auth="updateAuth()"/>
      <Profile v-if="authenticated"/>
    </div>
  </div>
</template>

<script>
import Nav from '@/components/Nav.vue';
import Landing from '@/components/Landing.vue';
import Profile from '@/components/Profile.vue';

export default {
    name: 'home',
    props: {
        authenticated: Boolean,
    },
    components: {
        Nav,
        Landing,
        Profile,
    },
    methods: {
        updateAuth() {
            this.$emit('auth');
        },
    },
};
</script>
