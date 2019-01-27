<template>
  <div class="following">
    <Nav v-on:auth="updateAuth()" :authenticated="authenticated"/>
    <div class="content">
      <Followings/>
      <Follow/>
    </div>
  </div>
</template>

<script>
import Nav from '@/components/Nav.vue';
import Followings from '@/components/Followings.vue';
import Follow from '@/components/Follow.vue';

export default {
    name: 'following',
    props: {
        authenticated: Boolean,
    },
    components: {
        Nav,
        Followings,
        Follow,
    },
    methods: {
        updateAuth() {
            this.$emit('auth');
        },
    },
};
</script>
