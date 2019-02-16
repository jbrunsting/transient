<template>
  <div class="home">
    <Nav v-on:auth="updateAuth()" :authenticated="authenticated"/>
    <div class="content">
      <HomeProfile />
    </div>
  </div>
</template>

<script>
import Nav from '@/components/Nav.vue';
import HomeProfile from '@/components/HomeProfile.vue';

export default {
    name: 'home',
    props: {
        authenticated: Boolean,
    },
    components: {
        Nav,
        HomeProfile,
    },
    methods: {
        updateAuth() {
            this.$emit('auth');
        },
    },
};
</script>
