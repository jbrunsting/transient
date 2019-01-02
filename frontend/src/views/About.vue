<template>
  <div class="about">
    <Nav v-on:auth="updateAuth()" :authenticated="authenticated"/>
    <div class="content">
      <h1>This is an about page</h1>
    </div>
  </div>
</template>

<script>
import Nav from '@/components/Nav.vue';

export default {
    name: 'about',
    props: {
        authenticated: Boolean,
    },
    components: {
        Nav,
    },
    methods: {
        updateAuth() {
            this.$emit('auth');
        },
    },
};
</script>
