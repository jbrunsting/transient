<template>
  <div class="following">
    <Nav v-on:auth="updateAuth()" :authenticated="authenticated"/>
    <div class="content">
      <Following/>
    </div>
  </div>
</template>

<script>
import Nav from '@/components/Nav.vue';
import Following from '@/components/Following.vue';

export default {
    name: 'following',
    props: {
        authenticated: Boolean,
    },
    components: {
        Nav,
        Following,
    },
    methods: {
        updateAuth() {
            this.$emit('auth');
        },
    },
};
</script>
